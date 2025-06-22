package routes

import (
	"fmt"
	"reschedule-program/database"
	"reschedule-program/models"
	"reschedule-program/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ScheduleRoutes(r *gin.Engine) {
	scheduleService := services.NewScheduleService()
	logService := services.NewLogService()

	r.POST("/schedule", func(c *gin.Context) {
		var request struct {
			ClassName string `json:"className"`
			Courses   []struct {
				Name      string `json:"name"`
				Day       int    `json:"day"`
				Slot      int    `json:"slot"`
				StartWeek int    `json:"startWeek"`
				EndWeek   int    `json:"endWeek"`
			} `json:"courses"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			fmt.Printf("JSON binding error: %v\n", err)
			c.JSON(400, gin.H{"msg": "Invalid schedule", "error": err.Error()})
			return
		}

		fmt.Printf("Parsed request: className=%s, courses=%+v\n", request.ClassName, request.Courses)

		// 使用事务来确保数据一致性
		err := database.DB.Transaction(func(tx *gorm.DB) error {
			// 创建班级
			schedule := &models.ClassSchedule{
				ClassName: request.ClassName,
			}
			
			if err := tx.Create(schedule).Error; err != nil {
				return err
			}

			// 创建课程并关联到班级
			for _, courseData := range request.Courses {
				course := models.Course{
					Name:       courseData.Name,
					Day:        courseData.Day,
					Slot:       courseData.Slot,
					WeekFrom:   courseData.StartWeek,
					WeekTo:     courseData.EndWeek,
					ScheduleID: schedule.ID,
				}
				
				if err := tx.Create(&course).Error; err != nil {
					return err
				}
			}
			
			return nil
		})

		if err != nil {
			c.JSON(500, gin.H{"msg": "Failed to save schedule", "error": err.Error()})
			return
		}

		logService.AddLog("New schedule imported for " + request.ClassName)
		c.JSON(200, gin.H{"msg": "Schedule saved"})
	})

	r.GET("/schedule/:class", func(c *gin.Context) {
		class := c.Param("class")

		schedule, err := scheduleService.GetScheduleByClassName(class)
		if err != nil {
			c.JSON(404, gin.H{"msg": "Class not found"})
			return
		}

		c.JSON(200, schedule)
	})

	r.GET("/schedules", func(c *gin.Context) {
		schedules, err := scheduleService.GetAllSchedules()
		if err != nil {
			c.JSON(500, gin.H{"msg": "Failed to get schedules"})
			return
		}

		c.JSON(200, schedules)
	})

	r.GET("/logs", func(c *gin.Context) {
		logs, err := logService.GetRecentLogs(50)
		if err != nil {
			c.JSON(500, gin.H{"msg": "Failed to get logs"})
			return
		}

		var messages []string
		for _, log := range logs {
			messages = append(messages, log.Message)
		}

		c.JSON(200, messages)
	})
}
