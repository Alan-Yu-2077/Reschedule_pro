package main

import (
	"net/http"
	"reschedule-program/database"
	"reschedule-program/middleware"
	"reschedule-program/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	database.InitDB()

	r := gin.Default()
	r.Use(middleware.CORS())

	// 添加静态文件服务
	// 静态资源文件（CSS、JS、图片等）
	r.Static("/static", "./static")

	// 加载HTML模板
	r.LoadHTMLGlob("./static/*.html")

	// 前端路由 - 所有前端页面都返回index.html
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// 处理前端路由（SPA应用需要）
	r.NoRoute(func(c *gin.Context) {
		// 检查是否是API请求
		if c.Request.URL.Path[:4] == "/api" ||
			c.Request.URL.Path[:6] == "/admin" ||
			c.Request.URL.Path[:8] == "/teacher" ||
			c.Request.URL.Path == "/login" ||
			c.Request.URL.Path == "/register" {
			c.JSON(http.StatusNotFound, gin.H{"error": "API endpoint not found"})
			return
		}

		// 前端路由，返回index.html
		c.HTML(http.StatusOK, "index.html", nil)
	})

	routes.AuthRoutes(r)
	routes.SetupScheduleRoutes(r)
	routes.AdminRoutes(r)
	routes.TeacherRoutes(r)

	r.Run("0.0.0.0:80")
}
