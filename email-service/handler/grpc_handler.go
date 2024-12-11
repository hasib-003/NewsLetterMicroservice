package handler

import (
	"context"
	"errors"
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
		userEmail := userWithNews.Email
		body := "Here are your subscribed news :\n\n"

		for _, news := range userWithNews.NewsList {
			body += "Title: " + news.Title + "\n"
			body += "Description: " + news.Description + "\n"
			body += "Topic: " + news.TopicName + "\n\n"
		}
		err := SendEmail(userEmail, "your Weekly NewsLetter", body)
		if err != nil {
			log.Printf("send email error: %v", err)
			emailStatus[userEmail] = "Failed"
		} else {
			emailStatus[userEmail] = "Success"
		}

	}
	return &email.SendEmailsResponse{
		EmailStatus: emailStatus,
	}, nil
}

func (h *EmailHandler) SendIndividualEmail(ctx context.Context, req *email.SendIndividualEmailRequest) (*email.SendIndividualEmailResponse, error) {
	userEmail := req.Email
	body := req.VerificationCode
	err := SendEmail(userEmail, "verification code", body)
	if err != nil {
		log.Printf("send email error: %v", err)
	}
	return &email.SendIndividualEmailResponse{
		Message: "email sent to the individual user",
		Success: true,
	}, nil
}
func SendEmail(to, subject, body string) error {
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

	log.Printf("Sending email to: %s\nSubject: %s: %s\nFrom: %s\n", to, subject, body, from)
	response, err := client.Send(message)
	if err != nil {
		log.Printf("Failed to send email to %s: %v", to, err)
		return err
	}
	log.Printf("Email sent successfully to %s. status code %d", to, response.StatusCode)
	return nil
}
