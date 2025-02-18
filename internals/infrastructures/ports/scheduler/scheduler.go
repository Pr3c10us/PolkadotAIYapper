package main

import (
	"context"
	"fmt"
	"github.com/Pr3c10us/boilerplate/internals/services"
	"github.com/Pr3c10us/boilerplate/packages/configs"
	"log"
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

type DailySchedule struct {
	Windows []PostingWindow
}

// Schedule holds optimal posting times for each day
var WeeklySchedule = map[string]DailySchedule{
	"Monday": {
		Windows: []PostingWindow{
			{StartHour: 10, EndHour: 10},
			{StartHour: 14, EndHour: 16},
		},
	},
	"Tuesday": {
		Windows: []PostingWindow{
			{StartHour: 9, EndHour: 9},
			{StartHour: 13, EndHour: 15},
			{StartHour: 22, EndHour: 22},
		},
	},
	"Wednesday": {
		Windows: []PostingWindow{
			{StartHour: 9, EndHour: 9},
			{StartHour: 13, EndHour: 15},
			{StartHour: 17, EndHour: 19},
		},
	},
	"Thursday": {
		Windows: []PostingWindow{
			{StartHour: 9, EndHour: 9},
			{StartHour: 14, EndHour: 16},
			{StartHour: 20, EndHour: 22},
		},
	},
	"Friday": {
		Windows: []PostingWindow{
			{StartHour: 9, EndHour: 9},
			{StartHour: 14, EndHour: 16},
		},
	},
	"Saturday": {
		Windows: []PostingWindow{
			{StartHour: 13, EndHour: 15},
			{StartHour: 19, EndHour: 21},
		},
	},
	"Sunday": {
		Windows: []PostingWindow{
			{StartHour: 11, EndHour: 16},
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

func NewScheduler(services *services.Services, environment *configs.EnvironmentVariables) (*Scheduler, error) {
	est, err := time.LoadLocation("America/New_York")
	if err != nil {
		return nil, fmt.Errorf("failed to load EST timezone: %v", err)
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
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return &Scheduler{
		rdb:         rdb,
		ctx:         ctx,
		location:    est,
		services:    services,
		environment: environment,
	}, nil
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
		if err := tx.Set(s.ctx, quotaKey, DailyTweetLimit, 24*time.Hour).Err(); err != nil {
			return err
		}

		// Update last reset timestamp
		now := time.Now().In(s.location)
		if err := tx.Set(s.ctx, timestampKey, now.Unix(), 24*time.Hour).Err(); err != nil {
			return err
		}

		return nil
	}

	// Retry transaction if it fails due to WATCH
	for i := 0; i < 3; i++ {
		err := s.rdb.Watch(s.ctx, txf, quotaKey)
		if err == nil {
			return nil
		}
		if err == redis.TxFailedErr {
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
		if err == redis.TxFailedErr {
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

func (s *Scheduler) Run() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		if !s.IsWithinPostingWindow() {
			continue
		}

		// Try to reserve capacity for a single tweet
		tweets, reRun, err := s.services.TweetService.Tweet.Tweets()
		if err != nil || reRun {
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

		// Simulate posting tweet (replace with actual Twitter API call)
		if err := s.services.TweetService.Tweet.SendTweet(tweets); err != nil {
			log.Printf("Error posting tweet: %v", err)
			continue
		}

		// Update usage statistics
		if err := s.UpdateUsageStats(1); err != nil {
			log.Printf("Error updating usage stats: %v", err)
		}
	}
}

func (s *Scheduler) simulatePostTweet() error {
	// Add random delay to simulate API call
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	return nil
}

//func main() {
//	scheduler, err := NewScheduler("localhost:6379")
//	if err != nil {
//		log.Fatalf("Failed to create scheduler: %v", err)
//	}
//
//	if err := scheduler.Initialize(); err != nil {
//		log.Fatalf("Failed to initialize scheduler: %v", err)
//	}
//
//	// Start the scheduler
//	scheduler.Run()
//}
