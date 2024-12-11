package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Email             string    `json:"email" gorm:"unique"`
	Name              string    `json:"name"`
	Password          string    `json:"-"`
	Role              string    `json:"role"`
	Verified          bool      `json:"-" gorm:"default:false"`
	VerificationToken string    `json:"-"`
	TokenExpiresAt    time.Time `json:"-"`
}
type Subscription struct {
	ID     string `gorm:"primaryKey"`
	UserID string
	Topic  string
}
