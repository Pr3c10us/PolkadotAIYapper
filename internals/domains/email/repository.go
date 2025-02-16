package email

type Repository interface {
	SendEmail(params *MessageEmailParams) error
}
