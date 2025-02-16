package commands

import (
	"fmt"
	"github.com/Pr3c10us/boilerplate/internals/domains/email"
	"strconv"
	"sync"
	"time"

	"github.com/Pr3c10us/boilerplate/internals/domains/authentication"
	"github.com/Pr3c10us/boilerplate/internals/domains/cache"
	"github.com/Pr3c10us/boilerplate/packages/configs"
	"github.com/Pr3c10us/boilerplate/packages/utils"
	"golang.org/x/crypto/bcrypt"
)

type CreateUser interface {
	Handle(params *authentication.AddUserParams) error
}

type createUser struct {
	identityRepository   authentication.Repository
	emailRepository      email.Repository
	cacheRepository      cache.Repository
	environmentVariables *configs.EnvironmentVariables
}

func NewCreateUser(repository authentication.Repository, emailRepository email.Repository, cacheRepository cache.Repository, environmentVariables *configs.EnvironmentVariables) CreateUser {
	return &createUser{
		repository, emailRepository, cacheRepository, environmentVariables,
	}
}

func (service *createUser) Handle(params *authentication.AddUserParams) error {
	passwordByte, err := bcrypt.GenerateFromPassword([]byte(params.Password), 14)
	if err != nil {
		return err
	}

	params.Password = string(passwordByte)

	if err = service.identityRepository.CreateUser(params); err != nil {
		return err
	}
	code := utils.GenerateRandomNumber(6)

	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	var redisKey string
	redisKey = params.Email
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = service.cacheRepository.Set(
			fmt.Sprintf("%v:%v", service.environmentVariables.RedisKeys.VerificationCodeKey, redisKey),
			strconv.Itoa(code),
			time.Minute*10,
		)
		if err != nil {
			errChan <- err
		}
	}()

	var emailParams email.MessageEmailParams
	emailParams = email.MessageEmailParams{
		Email:   params.Email,
		Type:    "text",
		Subject: "User Email Verification Code",
		Message: fmt.Sprintf("Here is your verification code '%v'", code),
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = service.emailRepository.SendEmail(&emailParams)
		if err != nil {
			errChan <- err
		}
	}()

	wg.Wait()
	close(errChan)

	for err = range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}
