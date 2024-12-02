package repository

import (
	"github.com/hasib-003/newsLetterMicroservice/news-service/internal/model"
	"gorm.io/gorm"
)

type NewsRepository struct {
	DB *gorm.DB
}

func NewNewsRepository(db *gorm.DB) *NewsRepository {
	return &NewsRepository{
		DB: db,
	}
}
func (repo *NewsRepository) SaveNews(news []map[string]interface{}, topic string) error {
	for _, data := range news {
		author, authorOk := data["author"].(string)
		if !authorOk {
			author = "Unknown"
		}

		title, titleOk := data["title"].(string)
		if !titleOk {
			title = "Untitled"
		}

		description, descriptionOk := data["description"].(string)
		if !descriptionOk {
			description = "No description available"
		}

		newData := &model.News{
			Author:      author,
			Title:       title,
			Description: description,
			Topic:       topic,
		}
		if err := repo.DB.Create(newData).Error; err != nil {
			return err
		}
	}
	return nil
}

func (repo *NewsRepository) FindTopicByName(name string) ([]model.News, error) {
	var topics []model.News
	if err := repo.DB.Where("topic = ?", name).Find(&topics).Error; err != nil {
		return nil, err
	}
	return topics, nil
}

func (repo *NewsRepository) CreateSubscription(userID uint, topicID uint) error {
	subscription := &model.Subscription{
		UserID:  userID,
		TopicID: topicID,
	}
	if err := repo.DB.Create(subscription).Error; err != nil {
		return err
	}
	return nil
}
