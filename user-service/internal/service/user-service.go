package service

import (
	"context"
	"errors"
	"fmt"
	models "github.com/hasib-003/newsLetterMicroservice/user-service/internal/model"
	"github.com/hasib-003/newsLetterMicroservice/user-service/internal/repository"
	"github.com/hasib-003/newsLetterMicroservice/user-service/proto/email"
	"github.com/hasib-003/newsLetterMicroservice/user-service/proto/subscription"
	"github.com/hasib-003/newsLetterMicroservice/user-service/utils"
	"golang.org/x/crypto/bcrypt"
	"strconv"

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
func (s *UserService) CreateUser(email, name, password, role string) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &models.User{
		Email:    email,
		Name:     name,
		Password: string(hashedPassword),
		Role:     role,
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
func (s *UserService) Login(email, password string) (string, error) {
	user, err := s.repository.GetUserByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password ")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid email or password ")
	}
	token, err := utils.GenerateToken(strconv.Itoa(int(user.ID)), email, user.Role)
	if err != nil {
		return "", errors.New("invalid email or password ")
	}
	return token, nil
}
func (s *UserService) SubscribeToTopic(email string, topic string) error {
	user, err := s.GetUserByEmail(email)
	if err != nil {
		log.Println("no email found")
		return err
	}
	req := &subscription.SubscribeRequest{
		UserId:    uint32(user.ID),
		TopicName: topic,
	}
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

func (s *UserService) GetSubscribedTopics(userID uint32) ([]string, error) {
	req := &subscription.GetTopicRequest{UserId: userID}
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
