package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `json:"email" gorm:"unique"`
	Name     string `json:"name"`
	Password string `json:"-"`
}
type Subscription struct {
	ID     string `gorm:"primaryKey"`
	UserID string
	Topic  string
}
