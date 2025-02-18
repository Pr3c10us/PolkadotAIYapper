package command

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Pr3c10us/boilerplate/internals/domains/embedding"
	"github.com/Pr3c10us/boilerplate/internals/domains/llm"
	"github.com/Pr3c10us/boilerplate/internals/domains/xdotcom"
	"math/rand"
	"strings"
)

type Tweet struct {
	llm       llm.Repository
	embedding embedding.Repository
	xdotcom   xdotcom.Repository
}

func NewTweet(llm llm.Repository, embedding embedding.Repository, xdotcom xdotcom.Repository) *Tweet {
	return &Tweet{llm: llm, embedding: embedding, xdotcom: xdotcom}
}

func (service *Tweet) Tweets() ([]string, bool, error) {
	topicType := service.RandomTopicType()
	var topic, context string
	var err error

	switch topicType {
	case PRODUCT:
		count := 0
		for count < 5 {
			count++
			topic, err = service.GetProductTopics()
			if err != nil {
				continue
			}
			break
		}
		if count >= 5 {
			return nil, true, errors.New("error getting appropriate response from model")
		}

	case JAM:
	default:
		count := 0
		for count < 5 {
			count++
			topic, err = service.GetStandardTopics()
			if err != nil {
				continue
			}
			break
		}
		if count >= 5 {
			return nil, true, errors.New("error getting appropriate response from model")
		}
	}

	if topic == "" {
		return nil, true, errors.New("no topic")
	}
	topicTweeted, embeddingStr, err := service.TopicAlreadyTweeted(topic)
	if err != nil {
		return nil, false, err
	}
	if *topicTweeted {
		return nil, true, errors.New("topic tweeted")
	}

	err = service.embedding.AddEmbedding(embeddingStr, topic)
	if err != nil {
		return nil, false, err
	}

	tweets, err := service.GetTweet(topic, context)
	if err != nil {
		return nil, false, err
	}

	return tweets, false, nil
}

func (service *Tweet) SendTweet(tweets []string) error {
	prevTweetID := ""
	for _, tweet := range tweets {
		id, err := service.xdotcom.Tweet(xdotcom.Tweet{
			Text:            tweet,
			PreviousTweetID: prevTweetID,
		})
		if err != nil {
			return err
		}
		prevTweetID = id
	}
	return nil
}

const (
	PRODUCT  = "product"
	STANDARD = "standard"
	JAM      = "jam"
)

func (service *Tweet) RandomTopicType() string {
	list := []string{
		PRODUCT, STANDARD, JAM,
	}

	randIndex := rand.Intn(len(list) - 0)

	return list[randIndex]
}

func (service *Tweet) convertToArray(input string) ([]string, error) {
	if strings.HasPrefix(input, "```json") && strings.HasSuffix(input, "```") {
		// Remove the code block markers
		trimmedInput := strings.TrimPrefix(input, "```json")
		trimmedInput = strings.TrimSuffix(trimmedInput, "```")

		// Trim any extra spaces, newlines, or tabs around the JSON content
		trimmedInput = strings.TrimSpace(trimmedInput)

		// Remove newlines and extra spaces inside the JSON string (if required)
		// Only if you're sure that the newlines should not be part of the actual data
		// this will collapse all the lines into a single line
		trimmedInput = strings.ReplaceAll(trimmedInput, "\n", "")

		// Parse the JSON string into a slice of strings
		var data []string
		err := json.Unmarshal([]byte(trimmedInput), &data)
		if err != nil {
			// Handle JSON parsing error
			fmt.Println("Error parsing JSON:", err)
			return nil, err
		}

		// Return the parsed data
		return data, nil
	}
	var result []string
	err := json.Unmarshal([]byte(input), &result)
	if err != nil {
		return nil, fmt.Errorf("invalid format: %v", err)
	}
	return result, nil
}

func (service *Tweet) GetProductTopics() (string, error) {
	response, err := service.llm.Prompt(service.ProductListPrompt())
	if err != nil {
		println(err)
		return "", err
	}

	products, err := service.convertToArray(response)
	if err != nil {
		println(err)
		return "", err
	}

	topicsPrompt := service.ProductTopicPrompt(products[rand.Intn(len(products)-0)])
	topicResponse, err := service.llm.Prompt(topicsPrompt)
	if err != nil {
		println(err)
		return "", err
	}

	topics, err := service.convertToArray(topicResponse)
	if err != nil {
		println(err)
		return "", err
	}

	return topics[rand.Intn(len(topics)-0)], nil
}

func (service *Tweet) GetStandardTopics() (string, error) {
	topicsPrompt := service.RandomStandardPrompt()
	topicResponse, err := service.llm.Prompt(topicsPrompt)
	if err != nil {
		println(err)
		return "", err
	}

	topics, err := service.convertToArray(topicResponse)
	if err != nil {
		println(err)
		return "", err
	}

	return topics[rand.Intn(len(topics)-0)], nil
}

func (service *Tweet) TopicAlreadyTweeted(topic string) (*bool, []float32, error) {
	tweetEmbedding, err := service.llm.Embed(topic)
	if err != nil {
		return nil, nil, err
	}
	used, err := service.embedding.SimilarValuesExist(tweetEmbedding)
	if err != nil {
		return nil, nil, err
	}

	return used, tweetEmbedding, nil
}

func (service *Tweet) GetTweet(topic, context string) ([]string, error) {
	tweetTypes := []string{"short", "thread"}
	tweetType := tweetTypes[rand.Intn(len(tweetTypes)-0)]
	//tweetType := tweetTypes[0]

	var prompt string
	if tweetType == "short" {
		prompt = service.ShortTweetPrompt(topic, context)
	} else {
		prompt = service.TweetThreadPrompt(topic, context)
	}

	switch tweetType {
	case "short":
		response, err := service.llm.Prompt(prompt)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		if strings.HasPrefix(response, "\"") && strings.HasSuffix(response, "\"") {
			// Remove the code block markers
			trimmedInput := strings.TrimPrefix(response, "\"")
			trimmedInput = strings.TrimSuffix(trimmedInput, "\"")

			// Trim any extra spaces, newlines, or tabs around the JSON content
			trimmedInput = strings.TrimSpace(trimmedInput)
			return []string{trimmedInput}, nil
		}
		return []string{response}, nil
	default:
		response, err := service.llm.Prompt(prompt)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		tweets, err := service.convertToArray(response)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		return tweets, nil
	}
}
