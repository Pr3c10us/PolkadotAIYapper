package sms

type Repository interface {
	SendSMS(params MessageSMSParams) error
}
