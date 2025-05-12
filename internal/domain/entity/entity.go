package entity

type Schedule struct {
	Id                 int64  `gorm:"primary key;autoIncrement" json:"id"`
	UserId             int64  `json:"user_id"`
	NameMedication     string `json:"name_medication"`
	MedicationPerDay   int    `json:"medication_per_day" validate:"required,gte=0,lte=15"`
	DurationMedication int    `json:"duration_medication"`
}

type ScheduleWithTime struct {
	Id                 int64    `json:"id"`
	NameMedication     string   `json:"name_medication"`
	MedicationPerDay   int      `json:"medication_per_day"`
	ScheduleMedication []string `json:"schedule"`
}
