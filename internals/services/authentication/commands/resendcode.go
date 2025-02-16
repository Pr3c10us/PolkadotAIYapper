package commands

import (
	"errors"
	"fmt"
	"github.com/Pr3c10us/boilerplate/internals/domains/authentication"
	"github.com/Pr3c10us/boilerplate/internals/domains/cache"
	emailDomain "github.com/Pr3c10us/boilerplate/internals/domains/email"
	"github.com/Pr3c10us/boilerplate/packages/appError"
	"github.com/Pr3c10us/boilerplate/packages/configs"
	"github.com/Pr3c10us/boilerplate/packages/utils"
	"strconv"
	"time"
)

type ResendCode interface {
	Handle(email string) error
}

type resendCode struct {
	repository           authentication.Repository
	emailRepository      emailDomain.Repository
	cacheRepository      cache.Repository
	environmentVariables *configs.EnvironmentVariables
}

func NewResendCodeService(repository authentication.Repository, emailRepository emailDomain.Repository, cacheRepository cache.Repository, environmentVariables *configs.EnvironmentVariables) ResendCode {
	return &resendCode{repository, emailRepository, cacheRepository, environmentVariables}
}

func (service *resendCode) Handle(email string) error {
	user, err := service.repository.GetUserDetails(&authentication.GetUserParams{Email: email})
	if err != nil {
		newError := errors.New("user does not exist")
		return appError.BadRequest(newError)
	}

	redisKey := fmt.Sprintf("%v:%v", service.environmentVariables.RedisKeys.VerificationCodeKey, user.Email)
	ttl, _ := service.cacheRepository.TTL(redisKey)
	if ttl > ((9 * time.Minute) + (30 * time.Second)) {
		return appError.BadRequest(errors.New("resend still in cool down"))
	}

	code := utils.GenerateRandomNumber(6)
	err = service.cacheRepository.Set(
		redisKey,
		strconv.Itoa(code),
		time.Minute*10,
	)
	if err != nil {
		return err
	}

	var emailParams emailDomain.MessageEmailParams
	emailParams = emailDomain.MessageEmailParams{
		Email:   email,
		Type:    "text",
		Subject: "User Email Verification Code",
		Message: fmt.Sprintf("Here is your verification code '%v'", code),
	}

	return service.emailRepository.SendEmail(&emailParams)
}
