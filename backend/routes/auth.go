package routes

import (
    "github.com/gin-gonic/gin"
    "reschedule-program/models"
    "reschedule-program/store"
)

func AuthRoutes(r *gin.Engine) {
    r.POST("/register", func(c *gin.Context) {
        var user models.User
        if err := c.ShouldBindJSON(&user); err != nil || user.Username == "" || user.Password == "" {
            c.JSON(400, gin.H{"msg": "Invalid input"})
            return
        }

        store.UserMutex.Lock()
        defer store.UserMutex.Unlock()

        if _, exists := store.Users[user.Username]; exists {
            c.JSON(400, gin.H{"msg": "User already exists"})
            return
        }

        store.Users[user.Username] = user.Password
        c.JSON(200, gin.H{"msg": "Register success"})
    })

    r.POST("/login", func(c *gin.Context) {
        var user models.User
        if err := c.ShouldBindJSON(&user); err != nil {
            c.JSON(400, gin.H{"msg": "Invalid input"})
            return
        }

        store.UserMutex.RLock()
        defer store.UserMutex.RUnlock()

        if pwd, ok := store.Users[user.Username]; !ok || pwd != user.Password {
            c.JSON(401, gin.H{"msg": "Wrong username or password"})
            return
        }

        c.JSON(200, gin.H{"msg": "Login success"})
    })
}