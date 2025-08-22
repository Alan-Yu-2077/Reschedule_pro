package routes

import (
	"net/http"
	"reschedule-program/database"
	"reschedule-program/models"
	"reschedule-program/services"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func SetupScheduleRoutes(router *gin.Engine) {
	scheduleGroup := router.Group("/api/schedule")
	{
		scheduleGroup.POST("/save", saveSchedule)
		scheduleGroup.POST("/add", addSchedule)
		scheduleGroup.GET("/class/:className/week/:weekNumber", getScheduleByClass)
		scheduleGroup.GET("/classes", getAllClasses)
		scheduleGroup.GET("/logs", getClassLogs)
		scheduleGroup.DELETE("/delete", deleteSchedule)
		scheduleGroup.POST("/move", moveSchedule)
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
	if err := services.SaveSchedule(scheduleData); err != nil {
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
	weekNumber, err := strconv.Atoi(weekNumberStr)
	if err != nil || weekNumber < 1 || weekNumber > 20 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid week number"})
		return
	}
	resolvedName := className
	if id, err := strconv.Atoi(className); err == nil && id > 0 {
		var cls models.Class
		if err2 := database.DB.Where("id = ?", uint(id)).First(&cls).Error; err2 == nil {
			resolvedName = cls.Name
		}
	}
	schedules, err := services.GetScheduleByClass(resolvedName, weekNumber)
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

// getClassLogs 按班级获取日志
func getClassLogs(c *gin.Context) {
	className := strings.TrimSpace(c.Query("className"))
	limitStr := c.Query("limit")
	var limit int
	if limitStr != "" {
		if v, err := strconv.Atoi(limitStr); err == nil {
			limit = v
		}
	}
	if className == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "className is required"})
		return
	}
	var classID uint
	var cls models.Class
	if err := database.DB.Where("name = ?", className).First(&cls).Error; err == nil {
		classID = cls.ID
	} else if id, err2 := strconv.Atoi(className); err2 == nil && id > 0 {
		var byID models.Class
		if err3 := database.DB.First(&byID, uint(id)).Error; err3 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Class not found"})
			return
		}
		classID = byID.ID
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Class not found"})
		return
	}
	logs, err := services.NewLogService().GetLogsByClass(classID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get logs"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"logs": logs})
}

// deleteSchedule 删除课程
func deleteSchedule(c *gin.Context) {
	var request struct {
		ClassName   string `json:"className"`
		WeekNumber  int    `json:"weekNumber"`
		TimeSlotRow int    `json:"timeSlotRow"`
		TimeSlotCol int    `json:"timeSlotCol"`
		Operator    string `json:"operator"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if request.ClassName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Class name is required"})
		return
	}
	if request.WeekNumber < 1 || request.WeekNumber > 20 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid week number"})
		return
	}
	if request.TimeSlotRow < 0 || request.TimeSlotRow > 4 || request.TimeSlotCol < 0 || request.TimeSlotCol > 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time slot"})
		return
	}
	resolvedName := request.ClassName
	if id, err := strconv.Atoi(request.ClassName); err == nil && id > 0 {
		var cls models.Class
		if err2 := database.DB.Where("id = ?", uint(id)).First(&cls).Error; err2 == nil {
			resolvedName = cls.Name
		}
	}
	if err := services.DeleteSchedule(resolvedName, request.WeekNumber, request.TimeSlotRow, request.TimeSlotCol, request.Operator); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete schedule: " + err.Error()})
		return
	}
	// 追加记录操作者（通过服务层日志统一写的已存在；这里可补充operator）
	c.JSON(http.StatusOK, gin.H{"message": "Schedule deleted successfully"})
}

// moveSchedule 移动课程
func moveSchedule(c *gin.Context) {
	var request struct {
		ClassName  string `json:"className"`
		SourceWeek int    `json:"sourceWeek"`
		SourceRow  int    `json:"sourceRow"`
		SourceCol  int    `json:"sourceCol"`
		TargetWeek int    `json:"targetWeek"`
		TargetRow  int    `json:"targetRow"`
		TargetCol  int    `json:"targetCol"`
		Operator   string `json:"operator"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if request.ClassName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Class name is required"})
		return
	}
	if request.SourceWeek < 1 || request.SourceWeek > 20 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid source week number"})
		return
	}
	if request.TargetWeek < 1 || request.TargetWeek > 20 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target week number"})
		return
	}
	if request.SourceRow < 0 || request.SourceRow > 4 || request.SourceCol < 0 || request.SourceCol > 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid source time slot"})
		return
	}
	if request.TargetRow < 0 || request.TargetRow > 4 || request.TargetCol < 0 || request.TargetCol > 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target time slot"})
		return
	}
	resolvedName := request.ClassName
	if id, err := strconv.Atoi(request.ClassName); err == nil && id > 0 {
		var cls models.Class
		if err2 := database.DB.Where("id = ?", uint(id)).First(&cls).Error; err2 == nil {
			resolvedName = cls.Name
		}
	}
	if err := services.MoveSchedule(resolvedName, request.SourceWeek, request.SourceRow, request.SourceCol, request.TargetWeek, request.TargetRow, request.TargetCol, request.Operator); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to move schedule: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Schedule moved successfully"})
}

// addSchedule 添加新课程
func addSchedule(c *gin.Context) {
	var request struct {
		ClassName   string `json:"className"`
		CourseName  string `json:"courseName"`
		WeekNumber  int    `json:"weekNumber"`
		TimeSlotRow int    `json:"timeSlotRow"`
		TimeSlotCol int    `json:"timeSlotCol"`
		Operator    string `json:"operator"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if request.ClassName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Class name is required"})
		return
	}
	if request.CourseName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Course name is required"})
		return
	}
	if request.WeekNumber < 1 || request.WeekNumber > 20 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid week number"})
		return
	}
	if request.TimeSlotRow < 0 || request.TimeSlotRow > 4 || request.TimeSlotCol < 0 || request.TimeSlotCol > 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time slot"})
		return
	}
	resolvedName := request.ClassName
	if id, err := strconv.Atoi(request.ClassName); err == nil && id > 0 {
		var cls models.Class
		if err2 := database.DB.Where("id = ?", uint(id)).First(&cls).Error; err2 == nil {
			resolvedName = cls.Name
		}
	}
	if err := services.AddSchedule(resolvedName, request.CourseName, request.WeekNumber, request.TimeSlotRow, request.TimeSlotCol, request.Operator); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add schedule: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Schedule added successfully"})
}
