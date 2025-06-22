package models

import "gorm.io/gorm"

type Course struct {
	gorm.Model
	Name       string `json:"name" gorm:"not null"`
	Day        int    `json:"day" gorm:"not null"`
	Slot       int    `json:"slot" gorm:"not null"`
	WeekFrom   int    `json:"weekFrom" gorm:"not null"`
	WeekTo     int    `json:"weekTo" gorm:"not null"`
	ScheduleID uint   `json:"scheduleId" gorm:"not null"`
}

type ClassSchedule struct {
	gorm.Model
	ClassName string   `json:"className" gorm:"uniqueIndex;not null"`
	Courses   []Course `json:"courses" gorm:"foreignKey:ScheduleID;constraint:OnDelete:CASCADE"`
}

type ActivityLog struct {
	gorm.Model
	Message string `json:"message" gorm:"not null"`
}