package email

import (
	"github.com/Pr3c10us/boilerplate/internals/domains/email"
	"github.com/Pr3c10us/boilerplate/packages/configs"
	"github.com/go-mail/mail"
)

type GoMailEmailRepository struct {
	environmentVariables *configs.EnvironmentVariables
}

func NewGoMailEmailRepository(environmentVariables *configs.EnvironmentVariables) email.Repository {
	return &GoMailEmailRepository{
		environmentVariables: environmentVariables,
	}
}

func (repo *GoMailEmailRepository) SendEmail(params *email.MessageEmailParams) error {
	mailer := mail.NewMessage()
	mailer.SetHeader("From", repo.environmentVariables.SMTP.FromAddress)

	mailer.SetHeader("To", params.Email)

	mailer.SetHeader("Subject", params.Subject)

	mailer.SetBody("text/html", params.Message)

	dialer := mail.NewDialer(
		repo.environmentVariables.SMTP.Host,
		repo.environmentVariables.SMTP.Port,
		repo.environmentVariables.SMTP.Username,
		repo.environmentVariables.SMTP.Password,
	)

	// Send the email to Kate, Noah and Oliver.

	if err := dialer.DialAndSend(mailer); err != nil {
		return err
	}
	return nil
}
