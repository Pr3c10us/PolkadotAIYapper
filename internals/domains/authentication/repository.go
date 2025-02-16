package authentication

import "github.com/markbates/goth"

type Repository interface {
	CreateUser(user *AddUserParams) error
	GetUserDetails(params *GetUserParams) (*User, error)
	UpdateProfile(params *UserProfileParams) error
	AddUserOAuth(user *goth.User) (*User, error)
}
