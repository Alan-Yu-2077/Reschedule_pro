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

// deleteSchedule 删除课程
func deleteSchedule(c *gin.Context) {
	var request struct {
		ClassName   string `json:"className"`
		WeekNumber  int    `json:"weekNumber"`
		TimeSlotRow int    `json:"timeSlotRow"`
		TimeSlotCol int    `json:"timeSlotCol"`
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

	err := services.DeleteSchedule(request.ClassName, request.WeekNumber, request.TimeSlotRow, request.TimeSlotCol)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete schedule: " + err.Error()})
		return
	}

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

	// 验证时间槽范围
	if request.SourceRow < 0 || request.SourceRow > 4 || request.SourceCol < 0 || request.SourceCol > 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid source time slot"})
		return
	}

	if request.TargetRow < 0 || request.TargetRow > 4 || request.TargetCol < 0 || request.TargetCol > 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target time slot"})
		return
	}

	err := services.MoveSchedule(request.ClassName, request.SourceWeek, request.SourceRow, request.SourceCol, request.TargetWeek, request.TargetRow, request.TargetCol)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to move schedule: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Schedule moved successfully"})
}
