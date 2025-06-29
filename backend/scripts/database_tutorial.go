package main

import (
	"fmt"
	"log"
	"reschedule-program/database"
	"reschedule-program/models"
)

func main() {
	fmt.Println("=== 数据库操作教学指南 ===\n")
	
	// 初始化数据库
	database.InitDB()
	
	// 1. 创建操作 (Create)
	fmt.Println("1. 创建操作 (Create)")
	createOperations()
	
	// 2. 查询操作 (Read)
	fmt.Println("\n2. 查询操作 (Read)")
	readOperations()
	
	// 3. 更新操作 (Update)
	fmt.Println("\n3. 更新操作 (Update)")
	updateOperations()
	
	// 4. 删除操作 (Delete)
	fmt.Println("\n4. 删除操作 (Delete)")
	deleteOperations()
	
	// 5. 高级查询
	fmt.Println("\n5. 高级查询操作")
	advancedQueries()
	
	fmt.Println("\n=== 教学完成 ===")
}

// 创建操作示例
func createOperations() {
	fmt.Println("--- 创建用户 ---")
	
	// 创建单个用户
	user := &models.User{
		Username: "student1",
		Password: "password123",
	}
	
	if err := database.DB.Create(user).Error; err != nil {
		log.Printf("创建用户失败: %v", err)
	} else {
		fmt.Printf("✅ 创建用户成功: ID=%d, Username=%s\n", user.ID, user.Username)
	}
	
	// 批量创建用户
	users := []models.User{
		{Username: "student2", Password: "password123"},
		{Username: "student3", Password: "password123"},
		{Username: "teacher1", Password: "teacher123"},
	}
	
	if err := database.DB.Create(&users).Error; err != nil {
		log.Printf("批量创建用户失败: %v", err)
	} else {
		fmt.Printf("✅ 批量创建用户成功: %d 个用户\n", len(users))
	}
	
	fmt.Println("--- 创建课程表 ---")
	
	// 创建课程表（包含关联的课程）
	schedule := &models.ClassSchedule{
		ClassName: "计算机科学 2024",
		Courses: []models.Course{
			{Name: "数据结构", Day: 0, Slot: 0, WeekFrom: 1, WeekTo: 8},
			{Name: "算法设计", Day: 1, Slot: 0, WeekFrom: 1, WeekTo: 8},
			{Name: "数据库原理", Day: 2, Slot: 1, WeekFrom: 1, WeekTo: 8},
			{Name: "软件工程", Day: 3, Slot: 1, WeekFrom: 1, WeekTo: 8},
		},
	}
	
	if err := database.DB.Create(schedule).Error; err != nil {
		log.Printf("创建课程表失败: %v", err)
	} else {
		fmt.Printf("✅ 创建课程表成功: %s (包含 %d 门课程)\n", schedule.ClassName, len(schedule.Courses))
	}
}

// 查询操作示例
func readOperations() {
	fmt.Println("--- 基本查询 ---")
	
	// 查询所有用户
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		log.Printf("查询用户失败: %v", err)
	} else {
		fmt.Printf("✅ 查询到 %d 个用户\n", len(users))
		for _, user := range users {
			fmt.Printf("  - ID: %d, Username: %s\n", user.ID, user.Username)
		}
	}
	
	// 根据条件查询
	var specificUser models.User
	if err := database.DB.Where("username = ?", "student1").First(&specificUser).Error; err != nil {
		log.Printf("查询特定用户失败: %v", err)
	} else {
		fmt.Printf("✅ 查询特定用户: %s (ID: %d)\n", specificUser.Username, specificUser.ID)
	}
	
	// 查询课程表（包含关联的课程）
	var schedules []models.ClassSchedule
	if err := database.DB.Preload("Courses").Find(&schedules).Error; err != nil {
		log.Printf("查询课程表失败: %v", err)
	} else {
		fmt.Printf("✅ 查询到 %d 个课程表\n", len(schedules))
		for _, schedule := range schedules {
			fmt.Printf("  - %s (ID: %d, 课程数: %d)\n", schedule.ClassName, schedule.ID, len(schedule.Courses))
		}
	}
	
	// 分页查询
	fmt.Println("--- 分页查询 ---")
	var paginatedUsers []models.User
	limit := 2
	offset := 0
	
	if err := database.DB.Limit(limit).Offset(offset).Find(&paginatedUsers).Error; err != nil {
		log.Printf("分页查询失败: %v", err)
	} else {
		fmt.Printf("✅ 分页查询结果 (限制: %d, 偏移: %d): %d 个用户\n", limit, offset, len(paginatedUsers))
	}
}

// 更新操作示例
func updateOperations() {
	fmt.Println("--- 更新用户信息 ---")
 
	// 更新单个字段
	var user models.User
	if err := database.DB.Where("username = ?", "student1").First(&user).Error; err != nil {
		log.Printf("查找用户失败: %v", err)
		return
	}
	
	// 方法1: 使用 Save
	user.Password = "newpassword123"
	if err := database.DB.Save(&user).Error; err != nil {
		log.Printf("更新用户失败: %v", err)
	} else {
		fmt.Printf("✅ 更新用户密码成功: %s\n", user.Username)
	}
	
	// 方法2: 使用 Update
	if err := database.DB.Model(&models.User{}).Where("username = ?", "student2").Update("password", "updatedpassword123").Error; err != nil {
		log.Printf("批量更新失败: %v", err)
	} else {
		fmt.Printf("✅ 批量更新密码成功\n")
	}
	
	// 更新课程信息
	fmt.Println("--- 更新课程信息 ---")
	var course models.Course
	if err := database.DB.Where("name = ?", "数据结构").First(&course).Error; err != nil {
		log.Printf("查找课程失败: %v", err)
	} else {
		// 更新多个字段
		updates := map[string]interface{}{
			"week_to": 16,
			"slot":    1,
		}
		
		if err := database.DB.Model(&course).Updates(updates).Error; err != nil {
			log.Printf("更新课程失败: %v", err)
		} else {
			fmt.Printf("✅ 更新课程成功: %s (周数: %d-%d, 时间段: %d)\n", course.Name, course.WeekFrom, course.WeekTo, course.Slot)
		}
	}
}

// 删除操作示例
func deleteOperations() {
	fmt.Println("--- 删除操作 ---")
	
	// 软删除（GORM默认使用软删除）
	var userToDelete models.User
	if err := database.DB.Where("username = ?", "student3").First(&userToDelete).Error; err != nil {
		log.Printf("查找要删除的用户失败: %v", err)
	} else {
		if err := database.DB.Delete(&userToDelete).Error; err != nil {
			log.Printf("删除用户失败: %v", err)
		} else {
			fmt.Printf("✅ 软删除用户成功: %s\n", userToDelete.Username)
		}
	}
	
	// 永久删除
	var userToHardDelete models.User
	if err := database.DB.Unscoped().Where("username = ?", "teacher1").First(&userToHardDelete).Error; err != nil {
		log.Printf("查找要永久删除的用户失败: %v", err)
	} else {
		if err := database.DB.Unscoped().Delete(&userToHardDelete).Error; err != nil {
			log.Printf("永久删除用户失败: %v", err)
		} else {
			fmt.Printf("✅ 永久删除用户成功: %s\n", userToHardDelete.Username)
		}
	}
	
	// 批量删除
	if err := database.DB.Where("username LIKE ?", "student%").Delete(&models.User{}).Error; err != nil {
		log.Printf("批量删除失败: %v", err)
	} else {
		fmt.Printf("✅ 批量删除以'student'开头的用户成功\n")
	}
}

// 高级查询示例
func advancedQueries() {
	fmt.Println("--- 高级查询 ---")
	
	// 统计查询
	var userCount int64
	if err := database.DB.Model(&models.User{}).Count(&userCount).Error; err != nil {
		log.Printf("统计用户数量失败: %v", err)
	} else {
		fmt.Printf("✅ 用户总数: %d\n", userCount)
	}
	
	// 聚合查询
	var result struct {
		MaxWeek int
		MinWeek int
		AvgSlot float64
	}
	
	if err := database.DB.Model(&models.Course{}).Select("MAX(week_to) as max_week, MIN(week_from) as min_week, AVG(slot) as avg_slot").Scan(&result).Error; err != nil {
		log.Printf("聚合查询失败: %v", err)
	} else {
		fmt.Printf("✅ 课程统计: 最大周数=%d, 最小周数=%d, 平均时间段=%.2f\n", result.MaxWeek, result.MinWeek, result.AvgSlot)
	}
	
	// 复杂条件查询
	var courses []models.Course
	if err := database.DB.Where("week_from <= ? AND week_to >= ? AND day IN ?", 5, 10, []int{0, 1, 2}).Find(&courses).Error; err != nil {
		log.Printf("复杂条件查询失败: %v", err)
	} else {
		fmt.Printf("✅ 复杂条件查询结果: %d 门课程\n", len(courses))
		for _, course := range courses {
			fmt.Printf("  - %s (周数: %d-%d, 星期: %d, 时间段: %d)\n", course.Name, course.WeekFrom, course.WeekTo, course.Day, course.Slot)
		}
	}
	
	// 关联查询
	var schedulesWithCourses []models.ClassSchedule
	if err := database.DB.Preload("Courses", "week_from <= ?", 8).Find(&schedulesWithCourses).Error; err != nil {
		log.Printf("关联查询失败: %v", err)
	} else {
		fmt.Printf("✅ 关联查询结果: %d 个课程表\n", len(schedulesWithCourses))
		for _, schedule := range schedulesWithCourses {
			fmt.Printf("  - %s: %d 门课程\n", schedule.ClassName, len(schedule.Courses))
		}
	}
	
	// 原生SQL查询
	var courseNames []string
	if err := database.DB.Raw("SELECT DISTINCT name FROM courses WHERE week_to >= ? ORDER BY name", 10).Scan(&courseNames).Error; err != nil {
		log.Printf("原生SQL查询失败: %v", err)
	} else {
		fmt.Printf("✅ 原生SQL查询结果: %d 门课程\n", len(courseNames))
		for _, name := range courseNames {
			fmt.Printf("  - %s\n", name)
		}
	}
} 