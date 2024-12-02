package model

import "gorm.io/gorm"

type News struct {
	gorm.Model
	Author      string `gorm:"not null"`
	Title       string `gorm:"not null"`
	Description string `gorm:"not null"`
	Topic       string `gorm:"not null"`
}

type Topic struct {
	gorm.Model
	Name string `gorm:"unique"`
}

type Subscription struct {
	gorm.Model
	UserID  uint
	TopicID uint
}
