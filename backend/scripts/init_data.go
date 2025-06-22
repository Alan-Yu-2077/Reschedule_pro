package main

import (
	"log"
	"reschedule-program/database"
	"reschedule-program/models"
)

func main() {
	// Initialize database
	database.InitDB()

	// Add sample user
	user := &models.User{
		Username: "admin",
		Password: "admin123",
	}
	
	if err := database.DB.Create(user).Error; err != nil {
		log.Printf("Failed to create sample user: %v", err)
	} else {
		log.Println("Sample user created: admin/admin123")
	}

	// Add sample schedule
	schedule := &models.ClassSchedule{
		ClassName: "Grade 23 - Class 1",
		Courses: []models.Course{
			{Name: "Math", Day: 0, Slot: 0, WeekFrom: 1, WeekTo: 16},
			{Name: "English", Day: 1, Slot: 0, WeekFrom: 1, WeekTo: 16},
			{Name: "Science", Day: 2, Slot: 1, WeekFrom: 1, WeekTo: 16},
			{Name: "History", Day: 3, Slot: 1, WeekFrom: 1, WeekTo: 16},
			{Name: "Art", Day: 4, Slot: 2, WeekFrom: 1, WeekTo: 16},
		},
	}

	if err := database.DB.Create(schedule).Error; err != nil {
		log.Printf("Failed to create sample schedule: %v", err)
	} else {
		log.Println("Sample schedule created: Grade 23 - Class 1")
	}

	// Add sample logs
	logs := []models.ActivityLog{
		{Message: "System initialized"},
		{Message: "Sample data created"},
		{Message: "Database setup completed"},
	}

	for _, logEntry := range logs {
		if err := database.DB.Create(&logEntry).Error; err != nil {
			log.Printf("Failed to create log entry: %v", err)
		}
	}

	log.Println("Database initialization completed!")
} 