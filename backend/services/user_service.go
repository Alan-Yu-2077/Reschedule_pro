package services

import (
	"reschedule-program/database"
	"reschedule-program/models"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) CreateUser(user *models.User) error {
	return database.DB.Create(user).Error
}

func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := database.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) UserExists(username string) bool {
	var count int64
	database.DB.Model(&models.User{}).Where("username = ?", username).Count(&count)
	return count > 0
}

func (s *UserService) UserIDExists(userID string) bool {
	var count int64
	database.DB.Model(&models.User{}).Where("user_id = ?", userID).Count(&count)
	return count > 0
}
