package queries

import (
	"errors"

	"github.com/Pr3c10us/boilerplate/internals/domains/authentication"
	"github.com/Pr3c10us/boilerplate/packages/appError"
	"github.com/google/uuid"
)

type GetUserDetails interface {
	Handle(params *authentication.GetUserParams) (*authentication.User, error)
}

type getUserDetails struct {
	repository authentication.Repository
}

func NewGetUserDetails(repository authentication.Repository) GetUserDetails {
	return &getUserDetails{
		repository,
	}
}

func (service *getUserDetails) Handle(params *authentication.GetUserParams) (*authentication.User, error) {
	if params.Email == "" && params.ID == uuid.Nil {
		return nil, appError.BadRequest(errors.New("provide either 'email' or 'id"))
	}

	rider, err := service.repository.GetUserDetails(params)
	if err != nil {
		return nil, err
	}

	return rider, nil
}
