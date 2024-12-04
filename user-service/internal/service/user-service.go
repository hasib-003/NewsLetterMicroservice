package service

import (
	"context"
	"fmt"
	models "github.com/hasib-003/newsLetterMicroservice/user-service/internal/model"
	"github.com/hasib-003/newsLetterMicroservice/user-service/internal/repository"
	"github.com/hasib-003/newsLetterMicroservice/user-service/proto/email"
	"github.com/hasib-003/newsLetterMicroservice/user-service/proto/subscription"

	"log"
)

type UserService struct {
	repository  *repository.UserRepository
	newsClient  subscription.NewsServiceClient
	emailClient email.EmailServiceClient
}

func NewUserService(repository *repository.UserRepository, newsClient subscription.NewsServiceClient, emailClient email.EmailServiceClient) *UserService {
	return &UserService{
		repository:  repository,
		newsClient:  newsClient,
		emailClient: emailClient,
	}
}
func (s *UserService) CreateUser(email, name, password string) (*models.User, error) {
	user := &models.User{
		Email:    email,
		Name:     name,
		Password: password,
	}
	return s.repository.CreateUser(user)
}
func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	var user *models.User
	user, err := s.repository.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (s *UserService) SubscribeToTopic(email string, topic string) error {
	user, err := s.GetUserByEmail(email)
	if err != nil {
		log.Println("no email found")
		return err
	}
	log.Printf("Userid>>>>>>>>%v", user.ID)
	req := &subscription.SubscribeRequest{
		UserId:    uint32(user.ID),
		TopicName: topic,
	}
	log.Printf("request userid %v %v", req.GetUserId(), req)
	res, err := s.newsClient.SubscribeToTopic(context.Background(), req)

	if err != nil {
		log.Println(err)
		return fmt.Errorf("subscribe to topic error: %v", err)
	}
	if !res.Success {
		return fmt.Errorf("subscription failed: %s", res.Message)
	}
	return nil
}
func (s *UserService) GetSubscribedTopics(userID uint) ([]string, error) {
	req := &subscription.GetTopicRequest{UserId: uint32(userID)}
	res, err := s.newsClient.GetSubscribedTopics(context.Background(), req)
	if err != nil {
		log.Printf("get topic error: %v", err)
		return nil, err
	}
	return res.Topics, nil
}
func (s *UserService) GetSubscribedNews(userID uint) ([]*subscription.NewsItem, error) {
	req := &subscription.GetSubscribedNewsRequest{
		UserId: uint32(userID),
	}
	res, err := s.newsClient.GetSubscribedNews(context.Background(), req)
	if err != nil {
		log.Printf("get news error: %v", err)
		return nil, err
	}
	return res.NewsItems, nil

}
func (s *UserService) GetAllUserEmails() ([]string, error) {
	emails, err := s.repository.GetAllUserEmails()
	if err != nil {
		return nil, err
	}
	return emails, nil
}

func (s *UserService) GetUserWithNews() ([]*email.UserWithNews, error) {
	userEmails, err := s.repository.GetAllUserEmails()
	if err != nil {
		return nil, fmt.Errorf("get user emails error: %v", err)
	}
	var userWithNews []*email.UserWithNews

	for _, Email := range userEmails {
		user, err := s.repository.GetUserByEmail(Email)
		if err != nil {
			log.Printf("get user email error: %v", err)
			continue
		}
		newsItems, err := s.GetSubscribedNews(user.ID)
		if err != nil {
			log.Printf("get news error: %v", err)
			continue
		}
		var newsList []*email.News
		for _, news := range newsItems {
			newsList = append(newsList, &email.News{
				Title:       news.Title,
				Description: news.Description,
				TopicName:   news.TopicName,
			})
		}
		userWithNews = append(userWithNews, &email.UserWithNews{
			Email:    Email,
			NewsList: newsList,
		})
	}
	return userWithNews, nil
}
func (s *UserService) SendEmailsToAllUsers() error {
	userWithNews, err := s.GetUserWithNews()
	if err != nil {
		return fmt.Errorf("failed to get users with news: %v", err)
	}
	req := &email.SendEmailsRequest{
		UsersWithNews: userWithNews,
	}
	_, err = s.emailClient.SendEmails(context.Background(), req)
	if err != nil {
		return fmt.Errorf("failed to send emails: %v", err)
	}
	return nil
}
