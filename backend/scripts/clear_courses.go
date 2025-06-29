package main

import (
	"fmt"
	"log"
	"reschedule-program/database"
	"reschedule-program/models"
)

func main() {
	fmt.Println("=== 清空课程数据 ===")
	
	// 初始化数据库
	database.InitDB()
	
	// 获取删除前的课程数量
	var courseCount int64
	database.DB.Model(&models.Course{}).Count(&courseCount)
	fmt.Printf("删除前课程数量: %d\n", courseCount)
	
	if courseCount == 0 {
		fmt.Println("✅ 数据库中没有课程数据，无需清空")
		return
	}
	
	// 删除所有课程（使用WHERE条件）
	if err := database.DB.Unscoped().Where("1 = 1").Delete(&models.Course{}).Error; err != nil {
		log.Printf("删除课程失败: %v", err)
		return
	}
	
	// 验证删除结果
	database.DB.Model(&models.Course{}).Count(&courseCount)
	fmt.Printf("✅ 课程数据清空成功！删除后课程数量: %d\n", courseCount)
	
	// 记录日志
	logEntry := &models.ActivityLog{
		Message: "清空所有课程数据",
	}
	database.DB.Create(logEntry)
	
	fmt.Println("✅ 操作完成！")
} 