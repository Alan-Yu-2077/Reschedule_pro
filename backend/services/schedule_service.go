package services

import (
	"reschedule-program/database"
	"reschedule-program/models"
	"strconv"
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

// SaveSchedule 保存课程表数据（批量导入不记录操作者场景）
func SaveSchedule(data ScheduleData) error {
	var class models.Class
	result := database.DB.Where("name = ?", data.ClassName).First(&class)
	if result.Error != nil {
		class = models.Class{Name: data.ClassName}
		if err := database.DB.Create(&class).Error; err != nil {
			return err
		}
	}
	for row := 0; row < len(data.Schedule); row++ {
		for col := 0; col < len(data.Schedule[row]); col++ {
			courseData := data.Schedule[row][col]
			if courseData == nil {
				continue
			}
			var course models.Course
			result := database.DB.Where("name = ?", courseData.Name).First(&course)
			if result.Error != nil {
				course = models.Course{Name: courseData.Name}
				if err := database.DB.Create(&course).Error; err != nil {
					return err
				}
			}
			var weeks []int
			if courseData.WeekType == "continuous" {
				for w := courseData.StartWeek; w <= courseData.EndWeek; w++ {
					weeks = append(weeks, w)
				}
			} else {
				weeks = courseData.SelectedWeeks
			}
			for _, week := range weeks {
				weeklySchedule := models.WeeklySchedule{ClassID: class.ID, CourseID: course.ID, WeekNumber: week, TimeSlotRow: row, TimeSlotCol: col}
				if err := database.DB.Create(&weeklySchedule).Error; err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func GetScheduleByClass(className string, weekNumber int) ([]models.WeeklySchedule, error) {
	var schedules []models.WeeklySchedule
	err := database.DB.Preload("Class").Preload("Course").Joins("JOIN classes ON classes.id = weekly_schedules.class_id").Where("classes.name = ? AND weekly_schedules.week_number = ?", className, weekNumber).Find(&schedules).Error
	return schedules, err
}

func GetAllClasses() ([]models.Class, error) {
	var classes []models.Class
	err := database.DB.Find(&classes).Error
	return classes, err
}

// DeleteSchedule 删除指定时间槽的课程（记录操作者）
func DeleteSchedule(className string, weekNumber int, timeSlotRow int, timeSlotCol int, operator string) error {
	var class models.Class
	if err := database.DB.Where("name = ?", className).First(&class).Error; err != nil {
		return err
	}
	result := database.DB.Where("class_id = ? AND week_number = ? AND time_slot_row = ? AND time_slot_col = ?", class.ID, weekNumber, timeSlotRow, timeSlotCol).Delete(&models.WeeklySchedule{})
	if result.Error != nil {
		return result.Error
	}
	_ = NewLogService().AddClassLogWithOperator(class.ID, operator, "Deleted course at week "+strconv.Itoa(weekNumber)+", ("+strconv.Itoa(timeSlotRow)+","+strconv.Itoa(timeSlotCol)+")")
	return nil
}

// MoveSchedule 移动课程（记录操作者）
func MoveSchedule(className string, sourceWeek int, sourceRow int, sourceCol int, targetWeek int, targetRow int, targetCol int, operator string) error {
	var class models.Class
	if err := database.DB.Where("name = ?", className).First(&class).Error; err != nil {
		return err
	}
	var sourceSchedule models.WeeklySchedule
	if err := database.DB.Where("class_id = ? AND week_number = ? AND time_slot_row = ? AND time_slot_col = ?", class.ID, sourceWeek, sourceRow, sourceCol).First(&sourceSchedule).Error; err != nil {
		return err
	}
	var targetSchedule models.WeeklySchedule
	if err := database.DB.Where("class_id = ? AND week_number = ? AND time_slot_row = ? AND time_slot_col = ?", class.ID, targetWeek, targetRow, targetCol).First(&targetSchedule).Error; err == nil {
		return err
	}
	targetSchedule = models.WeeklySchedule{ClassID: sourceSchedule.ClassID, CourseID: sourceSchedule.CourseID, WeekNumber: targetWeek, TimeSlotRow: targetRow, TimeSlotCol: targetCol}
	if err := database.DB.Create(&targetSchedule).Error; err != nil {
		return err
	}
	if err := database.DB.Delete(&sourceSchedule).Error; err != nil {
		return err
	}
	// 同步教师在该班该周的授课占位：若源位置存在教师占位，则随课程移动到目标位置
	if err := syncTeacherSlotsOnMove(class.ID, sourceWeek, sourceRow, sourceCol, targetWeek, targetRow, targetCol); err != nil {
		_ = NewLogService().AddClassLogWithOperator(class.ID, operator, "Teacher plan sync failed: "+err.Error())
	}
	_ = NewLogService().AddClassLogWithOperator(class.ID, operator, "Moved course from week "+strconv.Itoa(sourceWeek)+" ("+strconv.Itoa(sourceRow)+","+strconv.Itoa(sourceCol)+") to week "+strconv.Itoa(targetWeek)+" ("+strconv.Itoa(targetRow)+","+strconv.Itoa(targetCol)+")")
	return nil
}

// AddSchedule 添加新课程（记录操作者）
func AddSchedule(className string, courseName string, weekNumber int, timeSlotRow int, timeSlotCol int, operator string) error {
	var class models.Class
	result := database.DB.Where("name = ?", className).First(&class)
	if result.Error != nil {
		class = models.Class{Name: className}
		if err := database.DB.Create(&class).Error; err != nil {
			NewLogService().AddLog("Failed to create class: " + className + " - " + err.Error())
			return err
		}
		// 记录创建新班级的日志（没有操作者，因为这是系统自动创建的）
		NewLogService().AddLog("Created new class: " + className)
	}
	var course models.Course
	result = database.DB.Where("name = ?", courseName).First(&course)
	if result.Error != nil {
		course = models.Course{Name: courseName}
		if err := database.DB.Create(&course).Error; err != nil {
			NewLogService().AddLog("Failed to create course: " + courseName + " - " + err.Error())
			return err
		}
		// 记录创建新课程的日志（没有操作者，因为这是系统自动创建的）
		NewLogService().AddLog("Created new course: " + courseName)
	}
	var existingSchedule models.WeeklySchedule
	if err := database.DB.Where("class_id = ? AND week_number = ? AND time_slot_row = ? AND time_slot_col = ?", class.ID, weekNumber, timeSlotRow, timeSlotCol).First(&existingSchedule).Error; err == nil {
		_ = NewLogService().AddClassLogWithOperator(class.ID, operator, "Failed to add course: slot already occupied at week "+strconv.Itoa(weekNumber)+" ("+strconv.Itoa(timeSlotRow)+","+strconv.Itoa(timeSlotCol)+")")
		return err
	}
	weeklySchedule := models.WeeklySchedule{ClassID: class.ID, CourseID: course.ID, WeekNumber: weekNumber, TimeSlotRow: timeSlotRow, TimeSlotCol: timeSlotCol}
	if err := database.DB.Create(&weeklySchedule).Error; err != nil {
		_ = NewLogService().AddClassLogWithOperator(class.ID, operator, "Failed to create weekly schedule: "+err.Error())
		return err
	}
	_ = NewLogService().AddClassLogWithOperator(class.ID, operator, "Added course '"+courseName+"' week "+strconv.Itoa(weekNumber)+" at ("+strconv.Itoa(timeSlotRow)+","+strconv.Itoa(timeSlotCol)+")")
	return nil
}

// syncTeacherSlotsOnMove 将教师在某班的授课占位从源位置移动到目标位置（若存在）
// - 查找 teacher_teach_slots 中 class_id=classID 且 (week,row,col) 为源位置的记录
// - 若目标位置已存在相同教师同班同周的占位：删除源位置记录
// - 否则：更新为目标位置
func syncTeacherSlotsOnMove(classID uint, sourceWeek int, sourceRow int, sourceCol int, targetWeek int, targetRow int, targetCol int) error {
	var slots []models.TeacherTeachSlot
	if err := database.DB.Where("class_id = ? AND week = ? AND row = ? AND col = ?", classID, sourceWeek, sourceRow, sourceCol).Find(&slots).Error; err != nil {
		return err
	}
	if len(slots) == 0 {
		return nil
	}
	for _, s := range slots {
		var count int64
		if err := database.DB.Model(&models.TeacherTeachSlot{}).Where("teacher_id = ? AND class_id = ? AND week = ? AND row = ? AND col = ?", s.TeacherID, classID, targetWeek, targetRow, targetCol).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			if err := database.DB.Delete(&s).Error; err != nil {
				return err
			}
			continue
		}
		s.Week = targetWeek
		s.Row = targetRow
		s.Col = targetCol
		if err := database.DB.Save(&s).Error; err != nil {
			return err
		}
	}
	return nil
}
