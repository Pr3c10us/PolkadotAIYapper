package queries

import (
	"errors"
	"fmt"

	"github.com/Pr3c10us/boilerplate/internals/domains/authentication"
	"github.com/Pr3c10us/boilerplate/internals/domains/cache"
	"github.com/Pr3c10us/boilerplate/packages/appError"
	"github.com/Pr3c10us/boilerplate/packages/configs"
	"github.com/redis/go-redis/v9"
)

type VerifyCode interface {
	Handle(params *authentication.VerifyCodeParams) error
}

type verifyCode struct {
	identityRepository   authentication.Repository
	cacheRepository      cache.Repository
	environmentVariables *configs.EnvironmentVariables
}

func NewVerifyCode(repository authentication.Repository, cacheRepository cache.Repository, environmentVariables *configs.EnvironmentVariables) VerifyCode {
	return &verifyCode{
		repository,
		cacheRepository,
		environmentVariables,
	}
}

func (service *verifyCode) Handle(params *authentication.VerifyCodeParams) error {
	var redisKey = fmt.Sprintf("%v:%v", service.environmentVariables.RedisKeys.VerificationCodeKey, params.Email)

	code, err := service.cacheRepository.Get(redisKey)
	if errors.Is(err, redis.Nil) {
		newError := errors.New("code has expired")
		return appError.BadRequest(newError)
	} else if err != nil {
		return err
	}

	if code != params.Code {
		newError := errors.New("invalid code")
		return appError.BadRequest(newError)
	}

	return service.identityRepository.UpdateProfile(&authentication.UserProfileParams{
		EmailVerified: true,
	})
}
