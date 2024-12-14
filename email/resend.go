package email

import (
	"time"

	"github.com/resendlabs/resend-go"
	"golang.org/x/time/rate"
)

type ResendEmailSender struct {
	client  *resend.Client
	limiter *rate.Limiter
}

func New(apiKey string) *ResendEmailSender {
	client := resend.NewClient(apiKey)
	return &ResendEmailSender{
		client:  client,
		limiter: rate.NewLimiter(rate.Every(time.Second), 1), // 1 request per second
	}
}

// func (r *ResendEmailSender) Send(email emailhandlers.Email) error {
// 	if err := r.limiter.Wait(context.Background()); err != nil {
// 		return err
// 	}

// 	params := &resend.SendEmailRequest{
// 		From:    email.From,
// 		To:      []string{email.To},
// 		Subject: email.Subject,
// 		Html:    email.Body,
// 	}

// 	_, err := r.client.Emails.Send(params)
// 	return err
// }

// func (r *ResendEmailSender) SendMany(manyEmail emailhandlers.ManyEmail) error {
// 	if err := r.limiter.Wait(context.Background()); err != nil {
// 		return err
// 	}

// 	params := &resend.SendEmailRequest{
// 		From:    manyEmail.From,
// 		To:      manyEmail.To,
// 		Subject: manyEmail.Subject,
// 		Html:    manyEmail.Body,
// 	}

// 	_, err := r.client.Emails.Send(params)
// 	return err
// }
