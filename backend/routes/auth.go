package routes

import (
	"reschedule-program/services"

	"github.com/gin-gonic/gin"
	"reschedule-program/models"
)

func AuthRoutes(r *gin.Engine) {
	userService := services.NewUserService()

	r.POST("/register", func(c *gin.Context) {
		var user struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		
		if err := c.ShouldBindJSON(&user); err != nil || user.Username == "" || user.Password == "" {
			c.JSON(400, gin.H{"msg": "Invalid input"})
			return
		}

		if userService.UserExists(user.Username) {
			c.JSON(400, gin.H{"msg": "User already exists"})
			return
		}

		newUser := &models.User{
			Username: user.Username,
			Password: user.Password,
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

		dbUser, err := userService.GetUserByUsername(user.Username)
		if err != nil || dbUser.Password != user.Password {
			c.JSON(401, gin.H{"msg": "Wrong username or password"})
			return
		}

		c.JSON(200, gin.H{"msg": "Login success"})
	})
}