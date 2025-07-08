package routes

import (
	"regexp"
	"reschedule-program/models"
	"reschedule-program/services"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	userService := services.NewUserService()

	r.POST("/register", func(c *gin.Context) {
		var user struct {
			UserID   string `json:"userID"`
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&user); err != nil || user.UserID == "" || user.Username == "" || user.Password == "" {
			c.JSON(400, gin.H{"msg": "Invalid input"})
			return
		}

		// 验证UserID是否为10位数字
		matched, _ := regexp.MatchString(`^\d{10}$`, user.UserID)
		if !matched {
			c.JSON(400, gin.H{"msg": "UserID must be 10 digits"})
			return
		}

		// 检查UserID是否已存在
		if userService.UserIDExists(user.UserID) {
			c.JSON(400, gin.H{"msg": "UserID already exists"})
			return
		}

		if userService.UserExists(user.Username) {
			c.JSON(400, gin.H{"msg": "Username already exists"})
			return
		}

		newUser := &models.User{
			UserID:   user.UserID,
			Username: user.Username,
			Password: user.Password,
			UserType: "viewer", // 默认为观察者用户
		}

		if err := userService.CreateUser(newUser); err != nil {
			c.JSON(500, gin.H{"msg": "Failed to create user"})
			return
		}

		c.JSON(200, gin.H{"msg": "Register success"})
	})

	r.POST("/login", func(c *gin.Context) {
		var user struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(400, gin.H{"msg": "Invalid input"})
			return
		}

		// 优先检查admin用户
		if user.Username == "Admin" && user.Password == "88888888" {
			c.JSON(200, gin.H{
				"msg":      "Login success",
				"userType": "admin",
				"username": "Admin",
			})
			return
		}

		// 普通用户登录
		dbUser, err := userService.GetUserByUsername(user.Username)
		if err != nil || dbUser.Password != user.Password {
			c.JSON(401, gin.H{"msg": "Wrong username or password"})
			return
		}

		c.JSON(200, gin.H{
			"msg":      "Login success",
			"userType": dbUser.UserType,
			"username": dbUser.Username,
		})
	})
}
