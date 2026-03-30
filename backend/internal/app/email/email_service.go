package email

import (
	"context"
	"fmt"

	"github.com/resend/resend-go/v3"
)

type EmailSender interface {
	SendSignUpCode(ctx context.Context, email, code string) error
	SendLoginCode(ctx context.Context, email, code string) error
}

type emailSender struct {
	Client    resend.Client
	EmailFrom string
}

func NewEmailSender(APIKey string, EmailFrom string) EmailSender {
	return &emailSender{
		Client:    *resend.NewClient(APIKey),
		EmailFrom: fmt.Sprintf("Ticketmaster <%v>", EmailFrom),
	}
}

func (e *emailSender) SendSignUpCode(ctx context.Context, email, code string) error {
	params := &resend.SendEmailRequest{
		From:    e.EmailFrom,
		To:      []string{email},
		Subject: "Sign Up Code",
		Html:    fmt.Sprintf("<h1>Welcome!</h1><p>This is your sign up code <code> %v </code></p>", code),
		Text:    fmt.Sprintf("Welcome! This is your sign up code %v", code),
	}

	_, err := e.Client.Emails.Send(params)
	if err != nil {
		return err
	}

	return nil
}
func (e *emailSender) SendLoginCode(ctx context.Context, email, code string) error {
	params := &resend.SendEmailRequest{
		From:    e.EmailFrom,
		To:      []string{email},
		Subject: "Login Code",
		Html:    fmt.Sprintf("<h1>Welcome!</h1><p>This is your login code <code> %v </code></p>", code),
		Text:    fmt.Sprintf("Welcome! This is your login code %v", code),
	}

	_, err := e.Client.Emails.Send(params)
	if err != nil {
		return err
	}

	return nil
}
