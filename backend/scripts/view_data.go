package main

import (
	"fmt"
	"log"
	"reschedule-program/database"
	"reschedule-program/models"
)

func main() {
	// Initialize database
	database.InitDB()

	// View users
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		log.Printf("Failed to get users: %v", err)
	} else {
		fmt.Println("=== Users ===")
		for _, user := range users {
			fmt.Printf("ID: %d, Username: %s, Created: %s\n", user.ID, user.Username, user.CreatedAt.Format("2006-01-02 15:04:05"))
		}
	}

	// View schedules
	var schedules []models.ClassSchedule
	if err := database.DB.Preload("Courses").Find(&schedules).Error; err != nil {
		log.Printf("Failed to get schedules: %v", err)
	} else {
		fmt.Println("\n=== Schedules ===")
		for _, schedule := range schedules {
			fmt.Printf("ID: %d, Class: %s, Created: %s\n", schedule.ID, schedule.ClassName, schedule.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Println("  Courses:")
			for _, course := range schedule.Courses {
				fmt.Printf("    - %s (Day: %d, Slot: %d, Weeks: %d-%d)\n", course.Name, course.Day, course.Slot, course.WeekFrom, course.WeekTo)
			}
		}
	}

	// View logs
	var logs []models.ActivityLog
	if err := database.DB.Order("created_at desc").Limit(10).Find(&logs).Error; err != nil {
		log.Printf("Failed to get logs: %v", err)
	} else {
		fmt.Println("\n=== Recent Logs ===")
		for _, log := range logs {
			fmt.Printf("[%s] %s\n", log.CreatedAt.Format("2006-01-02 15:04:05"), log.Message)
		}
	}
} 