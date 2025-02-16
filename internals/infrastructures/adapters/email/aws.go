package email

import (
	"errors"
	"github.com/Pr3c10us/boilerplate/internals/domains/email"
	"github.com/Pr3c10us/boilerplate/packages/configs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ses"
)

type AWSEmailRepository struct {
	client               *ses.SES
	environmentVariables *configs.EnvironmentVariables
}

func (repository *AWSEmailRepository) SendEmail(params *email.MessageEmailParams) error {
	var input *ses.SendEmailInput
	if params.Type == "text" {
		input = &ses.SendEmailInput{
			Destination: &ses.Destination{
				CcAddresses: []*string{},
				ToAddresses: []*string{
					aws.String(params.Email),
				},
			},
			Message: &ses.Message{
				Body: &ses.Body{
					Text: &ses.Content{
						Charset: aws.String("UTF-8"),
						Data:    aws.String(params.Message),
					},
				},
				Subject: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(params.Subject),
				},
			},
			Source: aws.String(repository.environmentVariables.AWSKeys.AWSFromMail),
		}
	} else if params.Type == "html" {
		input = &ses.SendEmailInput{
			Destination: &ses.Destination{
				CcAddresses: []*string{},
				ToAddresses: []*string{
					aws.String(params.Email),
				},
			},
			Message: &ses.Message{
				Body: &ses.Body{
					Html: &ses.Content{
						Charset: aws.String("UTF-8"),
						Data:    aws.String(params.Message),
					},
				},
				Subject: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(params.Subject),
				},
			},
			Source: aws.String(repository.environmentVariables.AWSKeys.AWSFromMail),
		}
	} else {
		return errors.New("invalid email type")
	}

	_, err := repository.client.SendEmail(input)
	if err != nil {
		var awsErr awserr.Error
		if errors.As(err, &awsErr) {
			switch awsErr.Code() {
			case ses.ErrCodeMessageRejected:
				return errors.New(ses.ErrCodeMessageRejected)
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				return errors.New(ses.ErrCodeMailFromDomainNotVerifiedException)
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				return errors.New(ses.ErrCodeConfigurationSetDoesNotExistException)
			default:
				return errors.New(awsErr.Error())
			}
		}
		return err
	}

	return nil
}

func NewAWSEmailRepository(environmentVariables *configs.EnvironmentVariables, sesClient *ses.SES) email.Repository {
	return &AWSEmailRepository{
		client:               sesClient,
		environmentVariables: environmentVariables,
	}
}
