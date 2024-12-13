package service

import (
	"errors"
	"github.com/hasib-003/newsLetterMicroservice/email-service/repositoty"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"log"
	"os"
)

type EmailService struct {
	Repository *repositoty.RabbitMqRepository
}

func NewEmailService(repository *repositoty.RabbitMqRepository) *EmailService {
	return &EmailService{
		repository,
	}
}
func (s *EmailService) SendEmail(to, subject, body string) error {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	from := os.Getenv("SENDER_EMAIL")
	if apiKey == "" || from == "" {
		return errors.New("SENDER_EMAIL environment variable not set")
	}

	client := sendgrid.NewSendClient(apiKey)
	message := mail.NewSingleEmail(
		mail.NewEmail("Newsletter", from),
		subject,
		mail.NewEmail("User", to),
		body,
		body,
	)

	log.Printf("Sending email to: %s\nSubject: %s\nBody: %s\nFrom: %s\n", to, subject, body, from)
	response, err := client.Send(message)
	if err != nil {
		log.Printf("Failed to send email to %s: %v", to, err)
		return err
	}
	log.Printf("Email sent successfully to %s. Status code: %d", to, response.StatusCode)
	return nil
}
