package queries

import (
	"errors"

	"github.com/Pr3c10us/boilerplate/internals/domains/authentication"
	"github.com/Pr3c10us/boilerplate/packages/appError"
	"github.com/Pr3c10us/boilerplate/packages/configs"
	"github.com/Pr3c10us/boilerplate/packages/utils"
)

type Login interface {
	Handle(params *authentication.LoginParams) (*authentication.User, error)
}

type login struct {
	identityRepository   authentication.Repository
	environmentVariables *configs.EnvironmentVariables
}

func NewLogin(repository authentication.Repository, environmentVariables *configs.EnvironmentVariables) Login {
	return &login{
		repository, environmentVariables,
	}
}

func (service *login) Handle(params *authentication.LoginParams) (*authentication.User, error) {
	user, err := service.identityRepository.GetUserDetails(&authentication.GetUserParams{Email: params.Email})
	if err != nil {
		return nil, err
	}

	if !utils.IsValidPassword(user.Password, params.Password) {
		return nil, appError.BadRequest(errors.New("incorrect password"))
	}

	return user, nil
}
