package repository

import (
	"errors"
	"github.com/hasib-003/newsLetterMicroservice/user-service/config"
	models "github.com/hasib-003/newsLetterMicroservice/user-service/internal/model"
	"gorm.io/gorm"
	"log"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}
func (r *UserRepository) CreateUser(user *models.User) (*models.User, error) {
	if err := r.DB.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user *models.User
	if err := config.DB.Select("id,email,name").Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("No user found with email: %s", email)
			return user, nil
		}
		log.Printf("Error fetching user: %v", err)
		return user, err
	}
	return user, nil
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
