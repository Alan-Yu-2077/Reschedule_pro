package models

import "gorm.io/gorm"

// Class 班级表
type Class struct {
	gorm.Model
	Name string `json:"name" gorm:"uniqueIndex;not null"`
}

// Course 课程表
type Course struct {
	gorm.Model
	Name string `json:"name" gorm:"not null"`
}

// WeeklySchedule 按周存储的课程表
type WeeklySchedule struct {
	gorm.Model
	ClassID     uint   `json:"classId" gorm:"not null"`
	CourseID    uint   `json:"courseId" gorm:"not null"`
	WeekNumber  int    `json:"weekNumber" gorm:"not null"`  // 周数 1-20
	TimeSlotRow int    `json:"timeSlotRow" gorm:"not null"` // 时间段行号 0-4
	TimeSlotCol int    `json:"timeSlotCol" gorm:"not null"` // 时间段列号 0-6
	Class       Class  `json:"class" gorm:"foreignKey:ClassID"`
	Course      Course `json:"course" gorm:"foreignKey:CourseID"`
}

// ActivityLog 活动日志
type ActivityLog struct {
	gorm.Model
	Message  string `json:"message" gorm:"not null"`
	ClassID  uint   `json:"classId" gorm:"index"`
	Operator string `json:"operator" gorm:"size:64"`
}

// TeacherTeachSlot 教师在某班的授课占位（用于跨班冲突判定）
type TeacherTeachSlot struct {
	gorm.Model
	TeacherID string `json:"teacherId" gorm:"size:10;not null;index:idx_teacher_slot,priority:1"`
	ClassID   uint   `json:"classId" gorm:"not null;index:idx_teacher_slot,priority:2"`
	Week      int    `json:"week" gorm:"not null;index:idx_teacher_slot,priority:3"`
	Row       int    `json:"row" gorm:"not null;index:idx_teacher_slot,priority:4"`
	Col       int    `json:"col" gorm:"not null;index:idx_teacher_slot,priority:5"`
	Class     Class  `json:"class" gorm:"foreignKey:ClassID"`
}
