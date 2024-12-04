package handler

import (
	"context"
	email "github.com/hasib-003/newsLetterMicroservice/email-service/proto"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"log"
	"os"
)

type EmailHandler struct {
	email.UnimplementedEmailServiceServer
}

func NewEmailHandler() *EmailHandler {
	return &EmailHandler{}
}
func (h *EmailHandler) SendEmails(ctx context.Context, req *email.SendEmailsRequest) (*email.SendEmailsResponse, error) {
	emailStatus := make(map[string]string)
	for _, userWithNews := range req.UsersWithNews {
		email := userWithNews.Email
		body := "Here are your subscribed news :\n\n"

		for _, news := range userWithNews.NewsList {
			body += "Title: " + news.Title + "\n"
			body += "Description: " + news.Description + "\n"
			body += "Topic: " + news.TopicName + "\n\n"
		}
		err := SendEmail(email, "your Weekly NewsLetter", body)
		if err != nil {
			log.Printf("send email error: %v", err)
			emailStatus[email] = "Failed"
		} else {
			emailStatus[email] = "Success"
		}

	}
	return &email.SendEmailsResponse{
		EmailStatus: emailStatus,
	}, nil
}
func SendEmail(to, subject, body string) error {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	from := os.Getenv("SENDER_EMAIL")
	client := sendgrid.NewSendClient(apiKey)
	message := mail.NewSingleEmail(
		mail.NewEmail("Newsletter", from),
		subject,
		mail.NewEmail("User", to),
		body,
		body,
	)

	log.Printf("Sending email to: %s\nSubject:  %s\nFrom: %s\n", to, subject, from)
	_, err := client.Send(message)
	if err != nil {
		log.Printf("Failed to send email to %s: %v", to, err)
		return err
	}
	return nil
}
