package tweet

import (
	"github.com/Pr3c10us/boilerplate/internals/domains/embedding"
	"github.com/Pr3c10us/boilerplate/internals/domains/llm"
	"github.com/Pr3c10us/boilerplate/internals/services/tweet/command"
)

type Services struct {
	Commands
	Queries
}

type Commands struct {
	Tweet *command.Tweet
}

type Queries struct {
}

func NewTweetService(llm llm.Repository, embedding embedding.Repository) Services {
	return Services{
		Commands: Commands{
			Tweet: command.NewTweet(llm, embedding),
		},
		Queries: Queries{},
	}
}
