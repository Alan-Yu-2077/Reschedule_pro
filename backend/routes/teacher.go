package routes

import (
	"net/http"
	"reschedule-program/database"
	"reschedule-program/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func TeacherRoutes(r *gin.Engine) {
	teacherGroup := r.Group("/api/teacher")
	{
		teacherGroup.POST("/teach-schedule/overwrite-batch", overwriteTeacherTeachScheduleBatch)
		teacherGroup.GET("/teach-busy", getTeacherBusySlots)
		teacherGroup.GET("/teach-plan", getTeacherPlan)
		teacherGroup.GET("/classes", getTeacherClasses)
	}
}

// overwriteTeacherTeachScheduleBatch 覆盖式写入：教师在某班的一组周的授课占位
func overwriteTeacherTeachScheduleBatch(c *gin.Context) {
	var req struct {
		TeacherID string `json:"teacherId"`
		ClassName string `json:"className"`
		Weeks     []int  `json:"weeks"`
		Slots     []struct {
			Row int `json:"row"`
			Col int `json:"col"`
		} `json:"slots"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.TeacherID == "" || req.ClassName == "" || len(req.Weeks) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "teacherId, className and weeks are required"})
		return
	}
	// 查找班级
	var class models.Class
	if err := database.DB.Where("name = ?", req.ClassName).First(&class).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Class not found"})
		return
	}

	// 覆盖：删除该教师在该班这些周的所有占位
	if err := database.DB.Where("teacher_id = ? AND class_id = ? AND week IN ?", req.TeacherID, class.ID, req.Weeks).Delete(&models.TeacherTeachSlot{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear previous records"})
		return
	}

	// 批量插入
	var toCreate []models.TeacherTeachSlot
	for _, w := range req.Weeks {
		for _, s := range req.Slots {
			toCreate = append(toCreate, models.TeacherTeachSlot{
				TeacherID: req.TeacherID,
				ClassID:   class.ID,
				Week:      w,
				Row:       s.Row,
				Col:       s.Col,
			})
		}
	}
	if len(toCreate) > 0 {
		if err := database.DB.Create(&toCreate).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save records"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Teacher plan saved"})
}

// getTeacherBusySlots 返回某教师在指定周的忙碌格（可排除一个班）
func getTeacherBusySlots(c *gin.Context) {
	teacherID := c.Query("teacherId")
	weekStr := c.Query("week")
	excludeClassName := c.Query("excludeClassName")

	if teacherID == "" || weekStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "teacherId and week are required"})
		return
	}

	// 解析周
	week, err := strconv.Atoi(weekStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid week"})
		return
	}

	var excludeClassID *uint
	if excludeClassName != "" {
		var cls models.Class
		if err := database.DB.Where("name = ?", excludeClassName).First(&cls).Error; err == nil {
			excludeClassID = &cls.ID
		}
	}

	// 查询并预加载班级名称
	var slots []models.TeacherTeachSlot
	q := database.DB.Preload("Class").Where("teacher_id = ? AND week = ?", teacherID, week)
	if excludeClassID != nil {
		q = q.Where("class_id <> ?", *excludeClassID)
	}
	if err := q.Find(&slots).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load busy slots"})
		return
	}

	// 返回 row/col/classId/className 列表
	type slot struct {
		Row       int    `json:"row"`
		Col       int    `json:"col"`
		ClassID   uint   `json:"classId"`
		ClassName string `json:"className"`
	}
	res := make([]slot, 0, len(slots))
	for _, s := range slots {
		res = append(res, slot{Row: s.Row, Col: s.Col, ClassID: s.ClassID, ClassName: s.Class.Name})
	}
	c.JSON(http.StatusOK, gin.H{"slots": res})
}

// getTeacherPlan 返回教师在某班某周的占位
func getTeacherPlan(c *gin.Context) {
	teacherID := c.Query("teacherId")
	className := c.Query("className")
	weekStr := c.Query("week")
	if teacherID == "" || className == "" || weekStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "teacherId, className and week are required"})
		return
	}
	week, err := strconv.Atoi(weekStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid week"})
		return
	}
	var cls models.Class
	if err := database.DB.Where("name = ?", className).First(&cls).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Class not found"})
		return
	}
	var slots []models.TeacherTeachSlot
	if err := database.DB.Where("teacher_id = ? AND class_id = ? AND week = ?", teacherID, cls.ID, week).Find(&slots).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load teacher plan"})
		return
	}
	type slot struct {
		Row int `json:"row"`
		Col int `json:"col"`
	}
	res := make([]slot, 0, len(slots))
	for _, s := range slots {
		res = append(res, slot{Row: s.Row, Col: s.Col})
	}
	c.JSON(http.StatusOK, gin.H{"slots": res, "className": className, "week": week})
}

// getTeacherClasses 返回教师已录入占位的班级列表
func getTeacherClasses(c *gin.Context) {
	teacherID := c.Query("teacherId")
	if teacherID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "teacherId is required"})
		return
	}
	// 先取 distinct class_id
	var classIDs []uint
	if err := database.DB.Model(&models.TeacherTeachSlot{}).Where("teacher_id = ?", teacherID).Distinct().Pluck("class_id", &classIDs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load classes"})
		return
	}
	if len(classIDs) == 0 {
		c.JSON(http.StatusOK, gin.H{"classes": []models.Class{}})
		return
	}
	var classes []models.Class
	if err := database.DB.Where("id IN ?", classIDs).Find(&classes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load classes"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"classes": classes})
}
