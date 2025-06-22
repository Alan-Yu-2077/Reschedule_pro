package services

import (
	"reschedule-program/database"
	"reschedule-program/models"
)

type ScheduleService struct{}

func NewScheduleService() *ScheduleService {
	return &ScheduleService{}
}

func (s *ScheduleService) CreateSchedule(schedule *models.ClassSchedule) error {
	return database.DB.Create(schedule).Error
}

func (s *ScheduleService) GetScheduleByClassName(className string) (*models.ClassSchedule, error) {
	var schedule models.ClassSchedule
	err := database.DB.Preload("Courses").Where("class_name = ?", className).First(&schedule).Error
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (s *ScheduleService) GetAllSchedules() ([]models.ClassSchedule, error) {
	var schedules []models.ClassSchedule
	err := database.DB.Preload("Courses").Find(&schedules).Error
	return schedules, err
}

func (s *ScheduleService) UpdateSchedule(schedule *models.ClassSchedule) error {
	return database.DB.Save(schedule).Error
}

func (s *ScheduleService) DeleteSchedule(className string) error {
	return database.DB.Where("class_name = ?", className).Delete(&models.ClassSchedule{}).Error
} 