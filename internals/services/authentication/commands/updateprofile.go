package commands

import (
	"errors"

	"github.com/Pr3c10us/boilerplate/internals/domains/authentication"
	"github.com/Pr3c10us/boilerplate/packages/appError"
	"github.com/google/uuid"
)

type UpdateProfile interface {
	Handle(params *authentication.UserProfileParams) error
}

type updateProfile struct {
	repository authentication.Repository
}

func NewUpdateProfile(repository authentication.Repository) UpdateProfile {
	return &updateProfile{
		repository,
	}
}

func (service *updateProfile) Handle(params *authentication.UserProfileParams) error {
	if params.ID == uuid.Nil {
		return appError.Unauthorized(errors.New("not permitted"))
	}
	return service.repository.UpdateProfile(params)
}
