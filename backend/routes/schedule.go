package routes

import (
	"net/http"
	"reschedule-program/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SetupScheduleRoutes(router *gin.Engine) {
	scheduleGroup := router.Group("/api/schedule")
	{
		scheduleGroup.POST("/save", saveSchedule)
		scheduleGroup.GET("/class/:className/week/:weekNumber", getScheduleByClass)
		scheduleGroup.GET("/classes", getAllClasses)
	}
}

// saveSchedule 保存课程表
func saveSchedule(c *gin.Context) {
	var scheduleData services.ScheduleData

	if err := c.ShouldBindJSON(&scheduleData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if scheduleData.ClassName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Class name is required"})
		return
	}

	err := services.SaveSchedule(scheduleData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save schedule: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Schedule saved successfully"})
}

// getScheduleByClass 根据班级名和周数获取课程表
func getScheduleByClass(c *gin.Context) {
	className := c.Param("className")
	weekNumberStr := c.Param("weekNumber")

	if className == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Class name is required"})
		return
	}

	// 解析周数参数
	weekNumber, err := strconv.Atoi(weekNumberStr)
	if err != nil || weekNumber < 1 || weekNumber > 20 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid week number"})
		return
	}

	schedules, err := services.GetScheduleByClass(className, weekNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get schedule: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"schedules": schedules})
}

// getAllClasses 获取所有班级
func getAllClasses(c *gin.Context) {
	classes, err := services.GetAllClasses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get classes: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"classes": classes})
}
