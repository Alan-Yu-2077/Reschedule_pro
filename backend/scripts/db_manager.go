package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"reschedule-program/database"
	"reschedule-program/models"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("=== 数据库管理工具 ===")
	fmt.Println("输入 'help' 查看可用命令")
	
	// 初始化数据库
	database.InitDB()
	
	scanner := bufio.NewScanner(os.Stdin)
	
	for {
		fmt.Print("\n> ")
		if !scanner.Scan() {
			break
		}
		
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}
		
		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}
		
		command := parts[0]
		args := parts[1:]
		
		switch command {
		case "help":
			showHelp()
		case "users":
			showUsers()
		case "schedules":
			showSchedules()
		case "courses":
			showCourses()
		case "logs":
			showLogs()
		case "add-user":
			if len(args) >= 2 {
				addUser(args[0], args[1])
			} else {
				fmt.Println("用法: add-user <username> <password>")
			}
		case "add-schedule":
			if len(args) >= 1 {
				addSchedule(strings.Join(args, " "))
			} else {
				fmt.Println("用法: add-schedule <class_name>")
			}
		case "add-course":
			if len(args) >= 5 {
				addCourse(args)
			} else {
				fmt.Println("用法: add-course <schedule_id> <name> <day> <slot> <week_from> <week_to>")
			}
		case "delete-user":
			if len(args) >= 1 {
				deleteUser(args[0])
			} else {
				fmt.Println("用法: delete-user <username>")
			}
		case "delete-schedule":
			if len(args) >= 1 {
				deleteSchedule(args[0])
			} else {
				fmt.Println("用法: delete-schedule <class_name>")
			}
		case "stats":
			showStats()
		case "clear":
			clearScreen()
		case "exit", "quit":
			fmt.Println("再见！")
			return
		default:
			fmt.Printf("未知命令: %s\n", command)
			fmt.Println("输入 'help' 查看可用命令")
		}
	}
}

func showHelp() {
	fmt.Println(`
可用命令:
  help                    - 显示此帮助信息
  users                   - 显示所有用户
  schedules               - 显示所有课程表
  courses                 - 显示所有课程
  logs                    - 显示最近日志
  add-user <u> <p>        - 添加用户
  add-schedule <name>     - 添加课程表
  add-course <args>       - 添加课程 (schedule_id name day slot week_from week_to)
  delete-user <username>  - 删除用户
  delete-schedule <name>  - 删除课程表
  stats                   - 显示统计信息
  clear                   - 清屏
  exit/quit               - 退出程序
`)
}

func showUsers() {
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		log.Printf("查询用户失败: %v", err)
		return
	}
	
	fmt.Printf("\n=== 用户列表 (%d 个) ===\n", len(users))
	for _, user := range users {
		fmt.Printf("ID: %d | 用户名: %s | 创建时间: %s\n", 
			user.ID, user.Username, user.CreatedAt.Format("2006-01-02 15:04:05"))
	}
}

func showSchedules() {
	var schedules []models.ClassSchedule
	if err := database.DB.Preload("Courses").Find(&schedules).Error; err != nil {
		log.Printf("查询课程表失败: %v", err)
		return
	}
	
	fmt.Printf("\n=== 课程表列表 (%d 个) ===\n", len(schedules))
	for _, schedule := range schedules {
		fmt.Printf("ID: %d | 班级: %s | 课程数: %d | 创建时间: %s\n", 
			schedule.ID, schedule.ClassName, len(schedule.Courses), 
			schedule.CreatedAt.Format("2006-01-02 15:04:05"))
	}
}

func showCourses() {
	var courses []models.Course
	if err := database.DB.Find(&courses).Error; err != nil {
		log.Printf("查询课程失败: %v", err)
		return
	}
	
	fmt.Printf("\n=== 课程列表 (%d 门) ===\n", len(courses))
	for _, course := range courses {
		fmt.Printf("ID: %d | 名称: %s | 星期: %d | 时间段: %d | 周数: %d-%d | 课程表ID: %d\n", 
			course.ID, course.Name, course.Day, course.Slot, course.WeekFrom, course.WeekTo, course.ScheduleID)
	}
}

func showLogs() {
	var logs []models.ActivityLog
	if err := database.DB.Order("created_at desc").Limit(10).Find(&logs).Error; err != nil {
		log.Printf("查询日志失败: %v", err)
		return
	}
	
	fmt.Printf("\n=== 最近日志 (%d 条) ===\n", len(logs))
	for _, log := range logs {
		fmt.Printf("[%s] %s\n", log.CreatedAt.Format("2006-01-02 15:04:05"), log.Message)
	}
}

func addUser(username, password string) {
	user := &models.User{
		Username: username,
		Password: password,
	}
	
	if err := database.DB.Create(user).Error; err != nil {
		log.Printf("创建用户失败: %v", err)
		return
	}
	
	fmt.Printf("✅ 用户创建成功: %s (ID: %d)\n", username, user.ID)
	
	// 记录日志
	logEntry := &models.ActivityLog{
		Message: fmt.Sprintf("创建用户: %s", username),
	}
	database.DB.Create(logEntry)
}

func addSchedule(className string) {
	schedule := &models.ClassSchedule{
		ClassName: className,
	}
	
	if err := database.DB.Create(schedule).Error; err != nil {
		log.Printf("创建课程表失败: %v", err)
		return
	}
	
	fmt.Printf("✅ 课程表创建成功: %s (ID: %d)\n", className, schedule.ID)
	
	// 记录日志
	logEntry := &models.ActivityLog{
		Message: fmt.Sprintf("创建课程表: %s", className),
	}
	database.DB.Create(logEntry)
}

func addCourse(args []string) {
	if len(args) < 6 {
		fmt.Println("参数不足")
		return
	}
	
	scheduleID, err := strconv.ParseUint(args[0], 10, 32)
	if err != nil {
		fmt.Println("课程表ID格式错误")
		return
	}
	
	day, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println("星期格式错误")
		return
	}
	
	slot, err := strconv.Atoi(args[3])
	if err != nil {
		fmt.Println("时间段格式错误")
		return
	}
	
	weekFrom, err := strconv.Atoi(args[4])
	if err != nil {
		fmt.Println("起始周数格式错误")
		return
	}
	
	weekTo, err := strconv.Atoi(args[5])
	if err != nil {
		fmt.Println("结束周数格式错误")
		return
	}
	
	course := &models.Course{
		Name:       args[1],
		Day:        day,
		Slot:       slot,
		WeekFrom:   weekFrom,
		WeekTo:     weekTo,
		ScheduleID: uint(scheduleID),
	}
	
	if err := database.DB.Create(course).Error; err != nil {
		log.Printf("创建课程失败: %v", err)
		return
	}
	
	fmt.Printf("✅ 课程创建成功: %s (ID: %d)\n", course.Name, course.ID)
	
	// 记录日志
	logEntry := &models.ActivityLog{
		Message: fmt.Sprintf("创建课程: %s", course.Name),
	}
	database.DB.Create(logEntry)
}

func deleteUser(username string) {
	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		fmt.Printf("用户不存在: %s\n", username)
		return
	}
	
	if err := database.DB.Delete(&user).Error; err != nil {
		log.Printf("删除用户失败: %v", err)
		return
	}
	
	fmt.Printf("✅ 用户删除成功: %s\n", username)
	
	// 记录日志
	logEntry := &models.ActivityLog{
		Message: fmt.Sprintf("删除用户: %s", username),
	}
	database.DB.Create(logEntry)
}

func deleteSchedule(className string) {
	var schedule models.ClassSchedule
	if err := database.DB.Where("class_name = ?", className).First(&schedule).Error; err != nil {
		fmt.Printf("课程表不存在: %s\n", className)
		return
	}
	
	if err := database.DB.Delete(&schedule).Error; err != nil {
		log.Printf("删除课程表失败: %v", err)
		return
	}
	
	fmt.Printf("✅ 课程表删除成功: %s\n", className)
	
	// 记录日志
	logEntry := &models.ActivityLog{
		Message: fmt.Sprintf("删除课程表: %s", className),
	}
	database.DB.Create(logEntry)
}

func showStats() {
	var userCount, scheduleCount, courseCount, logCount int64
	
	database.DB.Model(&models.User{}).Count(&userCount)
	database.DB.Model(&models.ClassSchedule{}).Count(&scheduleCount)
	database.DB.Model(&models.Course{}).Count(&courseCount)
	database.DB.Model(&models.ActivityLog{}).Count(&logCount)
	
	fmt.Printf("\n=== 数据库统计 ===\n")
	fmt.Printf("用户总数: %d\n", userCount)
	fmt.Printf("课程表总数: %d\n", scheduleCount)
	fmt.Printf("课程总数: %d\n", courseCount)
	fmt.Printf("日志总数: %d\n", logCount)
	
	// 计算平均课程数
	if scheduleCount > 0 {
		avgCourses := float64(courseCount) / float64(scheduleCount)
		fmt.Printf("平均每课程表课程数: %.2f\n", avgCourses)
	}
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
} 