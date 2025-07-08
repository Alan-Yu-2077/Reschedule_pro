package routes

import (
	"fmt"
	"net/http"
	"reschedule-program/database"
	"reschedule-program/models"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.Engine) {
	adminGroup := r.Group("/admin")
	{
		adminGroup.GET("/users", adminGetAllUsers)
		adminGroup.PUT("/users/:id/password", adminUpdateUserPassword)
		adminGroup.PUT("/users/:id", adminUpdateUser)
		adminGroup.POST("/users", adminAddUser)
		adminGroup.DELETE("/users/:id", adminDeleteUser)
		adminGroup.GET("/classes", adminGetAllClasses)
		adminGroup.GET("/courses", adminGetAllCourses)
		adminGroup.GET("/schedules", adminGetAllSchedules)
		adminGroup.GET("/logs", adminGetAllLogs)
	}
}

func adminGetAllUsers(c *gin.Context) {
	var users []models.User
	database.DB.Find(&users)
	c.JSON(http.StatusOK, gin.H{"users": users})
}

func adminGetAllClasses(c *gin.Context) {
	var classes []models.Class
	database.DB.Find(&classes)
	c.JSON(http.StatusOK, gin.H{"classes": classes})
}

func adminGetAllCourses(c *gin.Context) {
	var courses []models.Course
	database.DB.Find(&courses)
	c.JSON(http.StatusOK, gin.H{"courses": courses})
}

func adminGetAllSchedules(c *gin.Context) {
	var schedules []models.WeeklySchedule
	database.DB.Preload("Class").Preload("Course").Find(&schedules)
	c.JSON(http.StatusOK, gin.H{"schedules": schedules})
}

func adminGetAllLogs(c *gin.Context) {
	var logs []models.ActivityLog
	database.DB.Find(&logs)
	c.JSON(http.StatusOK, gin.H{"logs": logs})
}

// adminUpdateUserPassword 管理员修改用户密码
func adminUpdateUserPassword(c *gin.Context) {
	userID := c.Param("id")

	var request struct {
		NewPassword string `json:"newPassword" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "New password is required"})
		return
	}

	var user models.User
	if err := database.DB.Where("user_id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 更新密码
	user.Password = request.NewPassword
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	// 记录日志
	logEntry := &models.ActivityLog{
		Message: "Admin updated user password: " + user.Username,
	}
	database.DB.Create(logEntry)

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

// adminUpdateUser 管理员更新用户信息
func adminUpdateUser(c *gin.Context) {
	oldUserID := c.Param("id")

	var request struct {
		NewUserID   string `json:"newUserID" binding:"required"`
		NewUsername string `json:"newUsername" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required"`
		NewUserType string `json:"newUserType" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields cannot be empty"})
		return
	}

	// 验证新UserID是否为10位数字
	if len(request.NewUserID) != 10 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID must be 10 digits"})
		return
	}

	var user models.User
	if err := database.DB.Where("user_id = ?", oldUserID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 如果UserID发生变化，检查新UserID是否已存在
	if oldUserID != request.NewUserID {
		var existingUser models.User
		if err := database.DB.Where("user_id = ?", request.NewUserID).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "New user ID already exists"})
			return
		}
	}

	// 检查新用户名是否已存在（排除当前用户）
	var existingUserWithUsername models.User
	if err := database.DB.Where("username = ? AND user_id != ?", request.NewUsername, oldUserID).First(&existingUserWithUsername).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "New username already exists"})
		return
	}

	// 更新用户信息
	oldUsername := user.Username
	user.UserID = request.NewUserID
	user.Username = request.NewUsername
	user.Password = request.NewPassword
	user.UserType = request.NewUserType

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	// 记录日志
	logEntry := &models.ActivityLog{
		Message: fmt.Sprintf("Admin updated user information: %s -> %s", oldUsername, user.Username),
	}
	database.DB.Create(logEntry)

	c.JSON(http.StatusOK, gin.H{"message": "User information updated successfully"})
}

// adminAddUser 管理员添加用户
func adminAddUser(c *gin.Context) {
	var req struct {
		UserID   string `json:"userID" binding:"required"`
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		UserType string `json:"userType"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID, username and password cannot be empty"})
		return
	}

	// 验证UserID是否为10位数字
	if len(req.UserID) != 10 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID must be 10 digits"})
		return
	}

	// 检查UserID是否已存在
	var existingUser models.User
	if err := database.DB.Where("user_id = ?", req.UserID).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID already exists"})
		return
	}

	user := models.User{
		UserID:   req.UserID,
		Username: req.Username,
		Password: req.Password,
		UserType: req.UserType,
	}
	if user.UserType == "" {
		user.UserType = "viewer"
	}
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user, username may already exist"})
		return
	}
	// 记录日志
	logEntry := &models.ActivityLog{
		Message: "Admin added user: " + user.Username,
	}
	database.DB.Create(logEntry)
	c.JSON(http.StatusOK, gin.H{"message": "User added successfully", "user": user})
}

// adminDeleteUser 管理员删除用户
func adminDeleteUser(c *gin.Context) {
	userID := c.Param("id")
	var user models.User
	if err := database.DB.Where("user_id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	// 使用Unscoped().Delete()进行硬删除，真正删除记录
	if err := database.DB.Unscoped().Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}
	// 记录日志
	logEntry := &models.ActivityLog{
		Message: "Admin deleted user: " + user.Username,
	}
	database.DB.Create(logEntry)
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
