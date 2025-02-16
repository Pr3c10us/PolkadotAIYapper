package commands

import (
	"errors"

	"github.com/Pr3c10us/boilerplate/internals/domains/authentication"
	"github.com/Pr3c10us/boilerplate/packages/appError"
	"github.com/markbates/goth"
)

type Authenticate interface {
	Handle(user *goth.User) (*authentication.User, error)
}

type authenticate struct {
	authenticationRepository authentication.Repository
}

func (service *authenticate) Handle(user *goth.User) (*authentication.User, error) {
	var organizationNotFoundErr = errors.New("user does not exit")
	switch organizationFetched, err := service.authenticationRepository.GetUserDetails(
		&authentication.GetUserParams{Email: user.Email},
	); {
	case errors.Is(err, appError.NotFound(organizationNotFoundErr)) || organizationFetched == nil:
		organizationFetched, err = service.authenticationRepository.AddUserOAuth(user)
		if err != nil {
			return nil, err
		}
		return organizationFetched, nil
	case err == nil:
		return organizationFetched, nil
	default:
		return nil, err
	}
}

func NewAuthenticate(repository authentication.Repository) Authenticate {
	return &authenticate{
		repository,
	}
}
