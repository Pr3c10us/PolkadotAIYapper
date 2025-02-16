package command

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Pr3c10us/boilerplate/internals/domains/embedding"
	"github.com/Pr3c10us/boilerplate/internals/domains/llm"
	"math/rand"
)

type Tweet struct {
	llm       llm.Repository
	embedding embedding.Repository
}

func NewTweet(llm llm.Repository, embedding embedding.Repository) *Tweet {
	return &Tweet{llm: llm, embedding: embedding}
}

func (service *Tweet) Handle() (bool, error) {
	topicType := service.RandomTopicType()
	var topic string
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
			return true, errors.New("error getting appropriate response from model")
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
			return true, errors.New("error getting appropriate response from model")
		}
	}

	fmt.Println("topic", topic)
	topicTweeted, embeddingStr, err := service.TopicAlreadyTweeted(topic)
	if err != nil {
		return false, err
	}
	if *topicTweeted {
		return true, errors.New("topic tweeted")
	}

	// generate tweet

	//embed topic
	err = service.embedding.AddEmbedding(embeddingStr, topic)
	if err != nil {
		return false, err
	}
	println(topic)

	return false, nil
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
