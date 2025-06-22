package store

import (
    "reschedule-program/models"
    "sync"
)

var (
    Users       = make(map[string]string)
    Schedules   = make(map[string]models.ClassSchedule)
    Logs        = make([]string, 0)
    UserMutex   = sync.RWMutex{}
    ScheduleMux = sync.RWMutex{}
)