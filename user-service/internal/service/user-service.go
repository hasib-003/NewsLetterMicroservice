package service

import (
	"context"
	"fmt"
	models "github.com/hasib-003/newsLetterMicroservice/user-service/internal/model"
	"github.com/hasib-003/newsLetterMicroservice/user-service/internal/repository"
	subscription "github.com/hasib-003/newsLetterMicroservice/user-service/proto"

	"log"
)

type UserService struct {
	repository *repository.UserRepository
	newsClient subscription.NewsServiceClient
}

func NewUserService(repository *repository.UserRepository, newsClient subscription.NewsServiceClient) *UserService {
	return &UserService{
		repository: repository,
		newsClient: newsClient,
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
func (s *UserService) SubscribeToTopic(userID uint, topic string) error {
	user, err := s.GetUserByEmail("hasibhr17@gmail.com")
	if err != nil {
		log.Println("no email found")
		return err
	}
	req := &subscription.SubscribeRequest{
		UserId:    uint32(user.ID),
		TopicName: topic,
	}
	log.Printf("Sending request with topic: %v", req.TopicName)
	res, err := s.newsClient.SubscribeToTopic(context.Background(), req)
	log.Println("after grpc")

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
