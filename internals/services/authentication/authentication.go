package authentication

import (
	"github.com/Pr3c10us/boilerplate/internals/domains/authentication"
	"github.com/Pr3c10us/boilerplate/internals/domains/cache"
	"github.com/Pr3c10us/boilerplate/internals/domains/email"
	"github.com/Pr3c10us/boilerplate/internals/services/authentication/commands"
	"github.com/Pr3c10us/boilerplate/internals/services/authentication/queries"
	"github.com/Pr3c10us/boilerplate/packages/configs"
)

type Services struct {
	Commands
	Queries
}

type Commands struct {
	CreateUser    commands.CreateUser
	Authenticate  commands.Authenticate
	UpdateProfile commands.UpdateProfile
	ResendCode    commands.ResendCode
}

type Queries struct {
	GetUserDetails queries.GetUserDetails
	Login          queries.Login
	VerifyCode     queries.VerifyCode
}

func NewAuthenticationService(emailRepository email.Repository, cacheRepository cache.Repository, environmentVariables *configs.EnvironmentVariables, repository authentication.Repository) Services {
	return Services{
		Commands: Commands{
			CreateUser:    commands.NewCreateUser(repository, emailRepository, cacheRepository, environmentVariables),
			UpdateProfile: commands.NewUpdateProfile(repository),
			ResendCode:    commands.NewResendCodeService(repository, emailRepository, cacheRepository, environmentVariables),
			Authenticate:  commands.NewAuthenticate(repository),
		},
		Queries: Queries{
			GetUserDetails: queries.NewGetUserDetails(repository),
			Login:          queries.NewLogin(repository, environmentVariables),
			VerifyCode:     queries.NewVerifyCode(repository, cacheRepository, environmentVariables),
		},
	}
}
