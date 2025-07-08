package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	UserID    string         `json:"userID" gorm:"primaryKey;size:10;not null"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null"`
	Password  string         `json:"password" gorm:"not null"`
	UserType  string         `json:"userType" gorm:"default:'viewer'"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}
