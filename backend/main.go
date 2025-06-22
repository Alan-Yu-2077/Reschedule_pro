package main

import (
    "reschedule-program/middleware"
    "reschedule-program/routes"

    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    r.Use(middleware.CORS())

    routes.AuthRoutes(r)
    routes.ScheduleRoutes(r)

    r.Run(":8080")
}