package cache

import (
	"github.com/stretchr/testify/mock"
	"time"
)

type MockRepository struct {
	mock.Mock
}

func (repo *MockRepository) Set(key string, value string, expiration time.Duration) error {
	args := repo.Called(key, value, expiration)
	return args.Error(0)
}

func (repo *MockRepository) Get(key string) (string, error) {
	args := repo.Called(key)
	return args.String(0), args.Error(1)
}
