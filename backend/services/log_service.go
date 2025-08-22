package services

import (
	"reschedule-program/database"
	"reschedule-program/models"
)

type LogService struct{}

func NewLogService() *LogService {
	return &LogService{}
}

func (s *LogService) AddLog(message string) error {
	log := models.ActivityLog{
		Message: message,
	}
	return database.DB.Create(&log).Error
}

func (s *LogService) GetAllLogs() ([]models.ActivityLog, error) {
	var logs []models.ActivityLog
	err := database.DB.Order("created_at desc").Find(&logs).Error
	return logs, err
}

func (s *LogService) GetRecentLogs(limit int) ([]models.ActivityLog, error) {
	var logs []models.ActivityLog
	err := database.DB.Order("created_at desc").Limit(limit).Find(&logs).Error
	return logs, err
}

func (s *LogService) AddClassLog(classID uint, message string) error {
	log := models.ActivityLog{
		Message: message,
		ClassID: classID,
	}
	return database.DB.Create(&log).Error
}

func (s *LogService) AddClassLogWithOperator(classID uint, operator string, message string) error {
	log := models.ActivityLog{
		Message:  message,
		ClassID:  classID,
		Operator: operator,
	}
	return database.DB.Create(&log).Error
}

func (s *LogService) GetLogsByClass(classID uint, limit int) ([]models.ActivityLog, error) {
	var logs []models.ActivityLog
	q := database.DB.Where("class_id = ?", classID).Order("created_at desc")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if err := q.Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}
