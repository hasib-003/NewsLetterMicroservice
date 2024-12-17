package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	email "github.com/hasib-003/newsLetterMicroservice/email-service/proto"
	"github.com/hasib-003/newsLetterMicroservice/email-service/service"
	"log"
	"net/http"
)

type EmailHandler struct {
	email.UnimplementedEmailServiceServer
	service *service.EmailService
}

func NewEmailHandler(service *service.EmailService) *EmailHandler {
	return &EmailHandler{
		service: service,
	}
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
		err := h.service.SendEmail(userEmail, "your Weekly NewsLetter", body)
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
	err := h.service.SendEmail(userEmail, "verification code", body)
	if err != nil {
		log.Printf("send email error: %v", err)
	}
	return &email.SendIndividualEmailResponse{
		Message: "email sent to the individual user",
		Success: true,
	}, nil
}
func (h *EmailHandler) StartListening(c *gin.Context) {
	userWithNewsChan, err := h.service.Repository.ConsumeMessage()
	if err != nil {
		log.Fatalf("fail to consume message, err:%v", err)
	}
	for msg := range userWithNewsChan {
		userEmail := msg.Email
		body := "Here are your subscribed news :\n\n"
		for _, news := range msg.NewsList {
			body += "Title: " + news.Title + "\n"
			body += "Description: " + news.Description + "\n"
			body += "Topic: " + news.TopicName + "\n\n"
		}
		err := h.service.SendEmail(userEmail, "Your weekly newsLetter ", body)
		if err != nil {
			log.Fatalf("fail to send email to %s: %v", userWithNewsChan, err)
		}
		log.Printf("Email sent successfully to %v with body:%v ", msg.Email, body)
	}
	c.JSON(http.StatusOK, gin.H{"message :": "email service started listening"})
}
