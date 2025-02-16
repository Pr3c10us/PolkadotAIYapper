package authentication

import (
	"github.com/Pr3c10us/boilerplate/internals/domains/authentication"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (repo *MockRepository) UpdateProfile(params *authentication.UserProfileParams) error {
	args := repo.Called(params)
	return args.Error(0)
}

func (repo *MockRepository) CreateUser(params *authentication.AddUserParams) error {
	args := repo.Called(params)
	return args.Error(0)
}

func (repo *MockRepository) GetUserDetails(params *authentication.GetUserParams) (*authentication.User, error) {
	args := repo.Called(params)
	return args.Get(0).(*authentication.User), args.Error(1)
}
