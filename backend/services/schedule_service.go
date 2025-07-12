package services

import (
	"reschedule-program/database"
	"reschedule-program/models"
)

// ScheduleData 前端传来的课程表数据
type ScheduleData struct {
	ClassName string                    `json:"className"`
	Schedule  [][]*CourseAssignmentData `json:"schedule"`
}

// CourseAssignmentData 课程分配数据
type CourseAssignmentData struct {
	Name          string `json:"name"`
	WeekType      string `json:"weekType"` // "continuous" or "discrete"
	StartWeek     int    `json:"startWeek"`
	EndWeek       int    `json:"endWeek"`
	SelectedWeeks []int  `json:"selectedWeeks"`
}

// SaveSchedule 保存课程表数据
func SaveSchedule(data ScheduleData) error {
	// 1. 创建或获取班级
	var class models.Class
	result := database.DB.Where("name = ?", data.ClassName).First(&class)
	if result.Error != nil {
		// 班级不存在，创建新班级
		class = models.Class{Name: data.ClassName}
		if err := database.DB.Create(&class).Error; err != nil {
			return err
		}
	}

	// 2. 处理课程表数据
	for row := 0; row < len(data.Schedule); row++ {
		for col := 0; col < len(data.Schedule[row]); col++ {
			courseData := data.Schedule[row][col]
			if courseData == nil {
				continue
			}

			// 3. 创建或获取课程
			var course models.Course
			result := database.DB.Where("name = ?", courseData.Name).First(&course)
			if result.Error != nil {
				// 课程不存在，创建新课程
				course = models.Course{Name: courseData.Name}
				if err := database.DB.Create(&course).Error; err != nil {
					return err
				}
			}

			// 4. 生成周记录
			var weeks []int
			if courseData.WeekType == "continuous" {
				// 连续周
				for week := courseData.StartWeek; week <= courseData.EndWeek; week++ {
					weeks = append(weeks, week)
				}
			} else {
				// 离散周
				weeks = courseData.SelectedWeeks
			}

			// 5. 为每一周创建记录
			for _, week := range weeks {
				weeklySchedule := models.WeeklySchedule{
					ClassID:     class.ID,
					CourseID:    course.ID,
					WeekNumber:  week,
					TimeSlotRow: row,
					TimeSlotCol: col,
				}

				if err := database.DB.Create(&weeklySchedule).Error; err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// GetScheduleByClass 根据班级名获取课程表
func GetScheduleByClass(className string, weekNumber int) ([]models.WeeklySchedule, error) {
	var schedules []models.WeeklySchedule

	err := database.DB.Preload("Class").Preload("Course").
		Joins("JOIN classes ON classes.id = weekly_schedules.class_id").
		Where("classes.name = ? AND weekly_schedules.week_number = ?", className, weekNumber).
		Find(&schedules).Error

	return schedules, err
}

// GetAllClasses 获取所有班级
func GetAllClasses() ([]models.Class, error) {
	var classes []models.Class
	err := database.DB.Find(&classes).Error
	return classes, err
}

// DeleteSchedule 删除指定时间槽的课程
func DeleteSchedule(className string, weekNumber int, timeSlotRow int, timeSlotCol int) error {
	// 1. 获取班级ID
	var class models.Class
	if err := database.DB.Where("name = ?", className).First(&class).Error; err != nil {
		return err
	}

	// 2. 删除指定时间槽的课程记录
	result := database.DB.Where("class_id = ? AND week_number = ? AND time_slot_row = ? AND time_slot_col = ?",
		class.ID, weekNumber, timeSlotRow, timeSlotCol).Delete(&models.WeeklySchedule{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// MoveSchedule 移动课程从源位置到目标位置（支持跨周）
func MoveSchedule(className string, sourceWeek int, sourceRow int, sourceCol int, targetWeek int, targetRow int, targetCol int) error {
	// 1. 获取班级ID
	var class models.Class
	if err := database.DB.Where("name = ?", className).First(&class).Error; err != nil {
		return err
	}

	// 2. 检查源位置是否有课程
	var sourceSchedule models.WeeklySchedule
	if err := database.DB.Where("class_id = ? AND week_number = ? AND time_slot_row = ? AND time_slot_col = ?",
		class.ID, sourceWeek, sourceRow, sourceCol).First(&sourceSchedule).Error; err != nil {
		return err
	}

	// 3. 检查目标位置是否为空
	var targetSchedule models.WeeklySchedule
	if err := database.DB.Where("class_id = ? AND week_number = ? AND time_slot_row = ? AND time_slot_col = ?",
		class.ID, targetWeek, targetRow, targetCol).First(&targetSchedule).Error; err == nil {
		// 目标位置已有课程，返回错误
		return err
	}

	// 4. 创建新的目标记录
	targetSchedule = models.WeeklySchedule{
		ClassID:     sourceSchedule.ClassID,
		CourseID:    sourceSchedule.CourseID,
		WeekNumber:  targetWeek,
		TimeSlotRow: targetRow,
		TimeSlotCol: targetCol,
	}

	if err := database.DB.Create(&targetSchedule).Error; err != nil {
		return err
	}

	// 5. 删除源记录
	if err := database.DB.Delete(&sourceSchedule).Error; err != nil {
		return err
	}

	return nil
}
