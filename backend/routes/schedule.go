package routes

import (
	"reschedule-program/models"
	"reschedule-program/store"

	"github.com/gin-gonic/gin"
)

func ScheduleRoutes(r *gin.Engine) {
	r.POST("/schedule", func(c *gin.Context) {
		var schedule models.ClassSchedule
		if err := c.ShouldBindJSON(&schedule); err != nil {
			c.JSON(400, gin.H{"msg": "Invalid schedule", "error": err.Error()})
			return
		}

		store.ScheduleMux.Lock()
		store.Schedules[schedule.ClassName] = schedule
		store.Logs = append(store.Logs, "New schedule imported for "+schedule.ClassName)
		store.ScheduleMux.Unlock()

		c.JSON(200, gin.H{"msg": "Schedule saved"})
	})

	r.GET("/schedule/:class", func(c *gin.Context) {
		class := c.Param("class")

		store.ScheduleMux.RLock()
		defer store.ScheduleMux.RUnlock()

		if s, ok := store.Schedules[class]; ok {
			c.JSON(200, s)
		} else {
			c.JSON(404, gin.H{"msg": "Class not found"})
		}
	})

	r.GET("/logs", func(c *gin.Context) {
		store.ScheduleMux.RLock()
		defer store.ScheduleMux.RUnlock()
		c.JSON(200, store.Logs)
	})
}
