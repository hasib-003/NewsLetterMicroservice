package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hasib-003/newsLetterMicroservice/user-service/config"
	models "github.com/hasib-003/newsLetterMicroservice/user-service/internal/model"
	"github.com/hasib-003/newsLetterMicroservice/user-service/proto/email"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
	"log"
)

type UserRepository struct {
	DB         *gorm.DB
	Connection *amqp.Connection
}

func NewUserRepository(db *gorm.DB, connection *amqp.Connection) *UserRepository {
	return &UserRepository{
		DB:         db,
		Connection: connection,
	}
}
func (r *UserRepository) CreateUser(user *models.User) (*models.User, error) {
	if err := r.DB.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
func (r *UserRepository) PublishUserWithNews(userWithNews []*email.UserWithNews) error {
	if r.Connection == nil || r.Connection.IsClosed() {
		return fmt.Errorf("RabbitMQ connection is not open")
	}
	channel, err := r.Connection.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer func(channel *amqp.Channel) {
		err := channel.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(channel)
	q, err := channel.QueueDeclare(
		"user_with_news",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	for _, userNews := range userWithNews {
		jsonData, err := json.Marshal(userNews)
		if err != nil {
			log.Fatal(err)
		}
		err = channel.Publish(
			"",
			q.Name,
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        jsonData,
			})
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("User with news published successfully %+v", userNews.Email)
	}
	return nil

}
func (r *UserRepository) GetUserById(userId int64) (*models.User, error) {
	var user *models.User
	if err := r.DB.First(&user, userId).Error; err != nil {
		return nil, err
	}
	return user, nil
}
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user *models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("No user found with email: %s", email)
			return user, err
		}
		log.Printf("Error fetching user: %v", err)
		return user, err
	}
	return user, nil
}
func (r *UserRepository) VerifyUserEmail(user *models.User) error {
	if err := r.DB.Save(user).Error; err != nil {
		return fmt.Errorf("error saving user: %v", err)
	}
	return nil
}
func (r *UserRepository) GetAllUserEmails() ([]string, error) {
	var emails []string
	var users []*models.User
	if err := r.DB.Select("email").Find(&users).Error; err != nil {
		return []string{}, err
	}
	for _, user := range users {
		emails = append(emails, user.Email)
	}
	return emails, nil
}

func (r *UserRepository) BuySubscription(id uint64) error {
	result := r.DB.Model(&models.User{}).Where("id=?", id).Update("subscription_limit", 100)
	if result.Error != nil {
		return fmt.Errorf("failed to buy subscription %d", id)
	}
	return nil
}
