package sms

type MessageSMSParams struct {
	Phone   string `json:"phone"`
	Message string `json:"message"`
}
