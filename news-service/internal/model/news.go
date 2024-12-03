package model

import "gorm.io/gorm"

type News struct {
	gorm.Model
	Author      string `gorm:"not null"`
	Title       string `gorm:"not null"`
	Description string `gorm:"not null"`
	Topic       string `gorm:"not null"`
}

type UserSubscribedNews struct {
	NewsID      uint   `json:"news_id" gorm:"primaryKey"`
	Title       string `json:"title"`
	Description string `json:"description"`
	TopicName   string `json:"topic_name"`
}

type Subscription struct {
	gorm.Model
	UserID  uint
	TopicID uint
}
