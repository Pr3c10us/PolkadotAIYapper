package email

type MessageEmailParams struct {
	Email   string `json:"email"`
	Type    string `json:"type"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}
