package email

import "context"

type EmailSender interface {
	SendSignUpCode(ctx context.Context, email, code string) error
	SendLoginCode(ctx context.Context, email, code string) error
}

type emailSender struct{}

func NewEmailSender() EmailSender {
	return &emailSender{}
}

func (e *emailSender) SendSignUpCode(ctx context.Context, email, code string) error
func (e *emailSender) SendLoginCode(ctx context.Context, email, code string) error
