package services

import (
	"github.com/Pr3c10us/boilerplate/internals/infrastructures/adapters"
	"github.com/Pr3c10us/boilerplate/internals/services/authentication"
	"github.com/Pr3c10us/boilerplate/internals/services/tweet"
)

type Services struct {
	AuthenticationServices authentication.Services
	TweetService           tweet.Services
}

func NewServices(adapters *adapters.Adapters) *Services {
	return &Services{
		AuthenticationServices: authentication.NewAuthenticationService(adapters.EmailRepository, adapters.CacheRepository, adapters.EnvironmentVariables, adapters.AuthenticationRepository),
		TweetService:           tweet.NewTweetService(adapters.OpenAiRepository, adapters.EmbeddingRepository),
	}
}
