package scheduler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Pr3c10us/boilerplate/internals/services"
	"github.com/Pr3c10us/boilerplate/packages/configs"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	DailyTweetLimit = 17
	RedisKeyPrefix  = "twitter_scheduler:"
	LockTimeout     = 10 * time.Second
)

type PostingWindow struct {
	StartHour int
	EndHour   int
}

func (w PostingWindow) Duration() int {
	return w.EndHour - w.StartHour + 1
}

type DailySchedule struct {
	Windows []PostingWindow
}

func (d DailySchedule) TotalHours() int {
	total := 0
	for _, window := range d.Windows {
		total += window.Duration()
	}
	return total
}

type TweetDistribution struct {
	Window     PostingWindow
	TweetCount int
	Intervals  []time.Time
}

func (s *Scheduler) calculateDailyDistribution() ([]TweetDistribution, error) {
	now := time.Now().In(s.location)
	daySchedule := WeeklySchedule[now.Weekday().String()]

	// Get remaining tweets for today
	quotaKey := RedisKeyPrefix + "daily_quota"
	remainingTweets, err := s.rdb.Get(s.ctx, quotaKey).Int()
	if err != nil {
		return nil, fmt.Errorf("failed to get remaining tweets: %v", err)
	}

	totalHours := daySchedule.TotalHours()
	var distributions []TweetDistribution

	// Calculate tweets per hour, distributing remaining tweets across available hours
	tweetsPerHour := float64(remainingTweets) / float64(totalHours)

	for _, window := range daySchedule.Windows {
		// Calculate tweets for this window based on its duration
		windowDuration := window.Duration()
		windowTweets := int(math.Round(tweetsPerHour * float64(windowDuration)))

		if windowTweets > 0 {
			// Generate posting times within the window
			intervals := generatePostingTimes(window, windowTweets, now)

			distributions = append(distributions, TweetDistribution{
				Window:     window,
				TweetCount: windowTweets,
				Intervals:  intervals,
			})
		}
	}

	return distributions, nil
}

func generatePostingTimes(window PostingWindow, tweetCount int, today time.Time) []time.Time {
	var times []time.Time

	// Create time range for the window
	start := time.Date(
		today.Year(), today.Month(), today.Day(),
		window.StartHour, 0, 0, 0, today.Location(),
	)
	end := time.Date(
		today.Year(), today.Month(), today.Day(),
		window.EndHour, 59, 59, 0, today.Location(),
	)

	// Calculate duration and divide it into intervals
	duration := end.Sub(start)
	if tweetCount > 1 {
		intervalDuration := duration / time.Duration(tweetCount-1)

		// Add some randomness to each interval (±30% of interval)
		for i := 0; i < tweetCount; i++ {
			baseTime := start.Add(time.Duration(i) * intervalDuration)

			// Add random jitter (±30% of interval)
			maxJitter := int64(intervalDuration) * 30 / 100
			jitter := time.Duration(rand.Int63n(maxJitter*2) - maxJitter)

			postTime := baseTime.Add(jitter)

			// Ensure time stays within window
			if postTime.Before(start) {
				postTime = start
			} else if postTime.After(end) {
				postTime = end
			}

			times = append(times, postTime)
		}
	} else if tweetCount == 1 {
		// For single tweet, pick random time within window
		randomDuration := time.Duration(rand.Int63n(int64(duration)))
		times = append(times, start.Add(randomDuration))
	}

	return times
}

type ScheduledTweet struct {
	PostTime time.Time
	Executed bool
}

var WeeklySchedule = map[string]DailySchedule{
	"Monday": {
		Windows: []PostingWindow{
			{StartHour: 8, EndHour: 12},
			{StartHour: 12, EndHour: 14},
			{StartHour: 14, EndHour: 18},
			{StartHour: 18, EndHour: 22},
		},
	},
	"Tuesday": {
		Windows: []PostingWindow{
			{StartHour: 8, EndHour: 12},
			{StartHour: 12, EndHour: 14},
			{StartHour: 14, EndHour: 18},
			{StartHour: 18, EndHour: 22},
		},
	},
	"Wednesday": {
		Windows: []PostingWindow{
			{StartHour: 8, EndHour: 12},
			{StartHour: 12, EndHour: 14},
			{StartHour: 14, EndHour: 18},
			{StartHour: 18, EndHour: 22},
		},
	},
	"Thursday": {
		Windows: []PostingWindow{
			{StartHour: 8, EndHour: 12},
			{StartHour: 12, EndHour: 14},
			{StartHour: 14, EndHour: 18},
			{StartHour: 18, EndHour: 22},
		},
	},
	"Friday": {
		Windows: []PostingWindow{
			{StartHour: 8, EndHour: 12},
			{StartHour: 12, EndHour: 14},
			{StartHour: 14, EndHour: 18},
			{StartHour: 18, EndHour: 22},
		},
	},
	"Saturday": {
		Windows: []PostingWindow{
			{StartHour: 8, EndHour: 12},
			{StartHour: 12, EndHour: 14},
			{StartHour: 14, EndHour: 18},
			{StartHour: 18, EndHour: 22},
		},
	},
	"Sunday": {
		Windows: []PostingWindow{
			{StartHour: 8, EndHour: 12},
			{StartHour: 12, EndHour: 14},
			{StartHour: 14, EndHour: 18},
			{StartHour: 18, EndHour: 22},
		},
	},
}

type Scheduler struct {
	rdb         *redis.Client
	ctx         context.Context
	mutex       sync.Mutex
	location    *time.Location
	services    *services.Services
	environment *configs.EnvironmentVariables
}

func NewScheduler(services *services.Services, environment *configs.EnvironmentVariables) *Scheduler {
	est, err := time.LoadLocation("America/New_York")
	if err != nil {
		fmt.Println(err)
		panic("failed to load EST timezone: %v")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     environment.RedisCache.Address,
		Password: environment.RedisCache.Password, // Set if required
		Username: environment.RedisCache.Username,
		DB:       0,
	})

	ctx := context.Background()

	// Test Redis connection
	if err := rdb.Ping(ctx).Err(); err != nil {
		fmt.Println(err)
		panic("failed to connect to Redis: %v")
	}

	return &Scheduler{
		rdb:         rdb,
		ctx:         ctx,
		location:    est,
		services:    services,
		environment: environment,
	}
}

func (s *Scheduler) Initialize() error {
	// Set initial daily quota if not exists
	quotaKey := RedisKeyPrefix + "daily_quota"
	exists, err := s.rdb.Exists(s.ctx, quotaKey).Result()
	if err != nil {
		return fmt.Errorf("failed to check quota existence: %v", err)
	}

	if exists == 0 {
		if err := s.resetDailyQuota(); err != nil {
			return fmt.Errorf("failed to initialize daily quota: %v", err)
		}
	}

	return nil
}

func (s *Scheduler) resetDailyQuota() error {
	quotaKey := RedisKeyPrefix + "daily_quota"
	timestampKey := RedisKeyPrefix + "last_reset"

	// Use Redis transaction to update both values atomically
	txf := func(tx *redis.Tx) error {
		// Set new quota
		if err := tx.Set(s.ctx, quotaKey, DailyTweetLimit, 0).Err(); err != nil {
			return err
		}

		// Update last reset timestamp using RFC3339 format
		now := time.Now().In(s.location)
		if err := tx.Set(s.ctx, timestampKey, now.Format(time.RFC3339), 0).Err(); err != nil {
			return err
		}

		return nil
	}

	// Retry transaction if it fails due to WATCH
	for i := 0; i < 3; i++ {
		err := s.rdb.Watch(s.ctx, txf, quotaKey, timestampKey)
		if err == nil {
			return nil
		}
		if errors.Is(err, redis.TxFailedErr) {
			continue
		}
		return fmt.Errorf("failed to reset quota: %v", err)
	}

	return fmt.Errorf("failed to reset quota after retries")
}

func (s *Scheduler) ReserveTweetCapacity(count int) (bool, error) {
	quotaKey := RedisKeyPrefix + "daily_quota"

	// Use Redis transaction to check and update quota atomically
	txf := func(tx *redis.Tx) error {
		// Get current quota
		remaining, err := tx.Get(s.ctx, quotaKey).Int()
		if err != nil {
			return err
		}

		// Check if enough capacity
		if remaining < count {
			return fmt.Errorf("insufficient quota")
		}

		// Decrement quota
		return tx.Set(s.ctx, quotaKey, remaining-count, 0).Err()
	}

	// Retry transaction if it fails due to WATCH
	for i := 0; i < 3; i++ {
		err := s.rdb.Watch(s.ctx, txf, quotaKey)
		if err == nil {
			return true, nil
		}
		if errors.Is(err, redis.TxFailedErr) {
			continue
		}
		if err.Error() == "insufficient quota" {
			return false, nil
		}
		return false, fmt.Errorf("failed to reserve capacity: %v", err)
	}

	return false, fmt.Errorf("failed to reserve capacity after retries")
}

func (s *Scheduler) UpdateUsageStats(tweetCount int) error {
	now := time.Now().In(s.location)
	dateStr := now.Format("2006-01-02")
	statsKey := fmt.Sprintf("%susage_stats:%s", RedisKeyPrefix, dateStr)

	// Update daily usage statistics
	if err := s.rdb.IncrBy(s.ctx, statsKey, int64(tweetCount)).Err(); err != nil {
		return fmt.Errorf("failed to update usage stats: %v", err)
	}

	// Set expiration for stats (keep for 30 days)
	if err := s.rdb.Expire(s.ctx, statsKey, 30*24*time.Hour).Err(); err != nil {
		return fmt.Errorf("failed to set stats expiration: %v", err)
	}

	return nil
}

func (s *Scheduler) IsWithinPostingWindow() bool {
	now := time.Now().In(s.location)
	daySchedule := WeeklySchedule[now.Weekday().String()]
	currentHour := now.Hour()

	for _, window := range daySchedule.Windows {
		if currentHour >= window.StartHour && currentHour <= window.EndHour {
			return true
		}
	}
	return false
}

func (s *Scheduler) AcquireLock(lockKey string) (bool, error) {
	success, err := s.rdb.SetNX(s.ctx, lockKey, "locked", LockTimeout).Result()
	if err != nil {
		return false, fmt.Errorf("failed to acquire lock: %v", err)
	}
	return success, nil
}

func (s *Scheduler) ReleaseLock(lockKey string) error {
	if err := s.rdb.Del(s.ctx, lockKey).Err(); err != nil {
		return fmt.Errorf("failed to release lock: %v", err)
	}
	return nil
}

func (s *Scheduler) checkAndResetQuota() error {
	timestampKey := RedisKeyPrefix + "last_reset"

	// Get last reset timestamp
	lastResetStr, err := s.rdb.Get(s.ctx, timestampKey).Result()
	if err != nil && errors.Is(err, redis.Nil) {
		return fmt.Errorf("failed to get last reset timestamp: %v", err)
	}

	shouldReset := false
	if errors.Is(err, redis.Nil) {
		// No last reset timestamp, need to initialize
		shouldReset = true
	} else {
		lastResetUnix, err := time.Parse(time.RFC3339, lastResetStr)
		if err != nil {
			return fmt.Errorf("failed to parse last reset timestamp: %v", err)
		}

		// Check if 24 hours have passed since last reset
		if time.Since(lastResetUnix) >= 24*time.Hour {
			shouldReset = true
		}
	}

	if shouldReset {
		if err := s.resetDailyQuota(); err != nil {
			return fmt.Errorf("failed to reset quota: %v", err)
		}
	}

	return nil
}

func (s *Scheduler) GetSchedule() ([]ScheduledTweet, error) {
	key := RedisKeyPrefix + "schedule"

	scheduleStr, err := s.rdb.Get(s.ctx, key).Result()
	if err != nil && errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("failed to get last reset timestamp: %v", err)
	}

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	var schedules []ScheduledTweet
	err = json.Unmarshal([]byte(scheduleStr), &schedules)
	if err != nil {
		return nil, err
	}
	return schedules, nil
}

func (s *Scheduler) SetSchedule(schedules []ScheduledTweet) error {
	key := RedisKeyPrefix + "schedule"

	scheduleByte, err := json.Marshal(schedules)
	if err != nil {
		return err
	}

	// Use Redis transaction to update both values atomically
	txf := func(tx *redis.Tx) error {
		// Set new quota
		if err = tx.Set(s.ctx, key, scheduleByte, 0).Err(); err != nil {
			return err
		}

		return nil
	}

	// Retry transaction if it fails due to WATCH
	for i := 0; i < 3; i++ {
		err := s.rdb.Watch(s.ctx, txf, key)
		if err == nil {
			return nil
		}
		if errors.Is(err, redis.TxFailedErr) {
			continue
		}
		return fmt.Errorf("failed to reset quota: %v", err)
	}

	return fmt.Errorf("failed to reset quota after retries")
}

func (s *Scheduler) Run() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	var lastDistributionDate time.Time

	for range ticker.C {
		fmt.Println("Executing Start")

		now := time.Now().In(s.location)

		// Check and reset quota if needed
		if err := s.checkAndResetQuota(); err != nil {
			log.Printf("Error checking/resetting quota: %v", err)
			continue
		}

		// Recalculate distribution if it's a new day or no schedule exists
		if now.Day() != lastDistributionDate.Day() {
			distributions, err := s.calculateDailyDistribution()
			if err != nil {
				log.Printf("Error calculating distribution: %v", err)
				continue
			}

			var scheduledTweets []ScheduledTweet
			for _, dist := range distributions {
				for _, postTime := range dist.Intervals {
					scheduledTweets = append(scheduledTweets, ScheduledTweet{
						PostTime: postTime,
						Executed: false,
					})
				}
			}
			err = s.SetSchedule(scheduledTweets)
			if err != nil {
				log.Printf("Error storing schedules: %v", err)
				continue
			}
			lastDistributionDate = now
		}

		scheduledTweets, err := s.GetSchedule()
		if err != nil {
			log.Printf("Error getting schedule: %v", err)
			continue
		}

		// Check for tweets that should be posted
		for i := range scheduledTweets {
			if scheduledTweets[i].Executed {
				continue
			}

			if math.Abs(now.Sub(scheduledTweets[i].PostTime).Minutes()) <= 5 {
				fmt.Println("Time to Tweet")

				tweets, reRun, err := s.services.TweetService.Tweet.Tweets()
				if err != nil || reRun {
					log.Printf("Error getting tweets: %v", err)
					continue
				}
				reserved, err := s.ReserveTweetCapacity(len(tweets))
				if err != nil {
					log.Printf("Error reserving tweet capacity: %v", err)
					continue
				}
				if !reserved {
					continue
				}

				fmt.Println(tweets)
				//posting tweet
				if err = s.services.TweetService.Tweet.SendTweet(tweets); err != nil {
					log.Printf("Error posting tweet: %v", err)
					continue
				}

				// Update usage statistics
				if err = s.UpdateUsageStats(1); err != nil {
					log.Printf("Error updating usage stats: %v", err)
				}

				scheduledTweets[i].Executed = true
				err = s.SetSchedule(scheduledTweets)
				if err != nil {
					log.Printf("Error storing schedules: %v", err)
					continue
				}
				log.Printf("Posted tweet at %v", now.Format(time.RFC3339))
			}
		}
		fmt.Println("Executing End")
	}
}
