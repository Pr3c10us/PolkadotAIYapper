package sms

import (
	"github.com/Pr3c10us/boilerplate/internals/domains/sms"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
)

type AWSSMSRepository struct {
	client *sns.SNS
}

func (repository *AWSSMSRepository) SendSMS(params sms.MessageSMSParams) error {
	input := &sns.PublishInput{
		Message:     aws.String(params.Message),
		PhoneNumber: aws.String(params.Phone),
	}
	_, err := repository.client.Publish(input)
	if err != nil {
		return err
	}
	return nil
}

func NewAWSSMSRepository(client *sns.SNS) sms.Repository {
	return &AWSSMSRepository{
		client: client,
	}
}
