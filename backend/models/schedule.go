package models

type Course struct {
    Name     string `json:"name"`
    Day      int    `json:"day"`
    Slot     int    `json:"slot"`
    WeekFrom int    `json:"weekFrom"`
    WeekTo   int    `json:"weekTo"`
}

type ClassSchedule struct {
    ClassName string   `json:"className"`
    Courses   []Course `json:"courses"`
}