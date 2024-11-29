package model

import "gorm.io/gorm"

type News struct {
	gorm.Model
	Author      string `gorm:"not null"`
	Title       string `gorm:"not null"`
	Description string `gorm:"not null"`
	Topic       string `gorm:"not null"`
}
