package entity

import (
	"gorm.io/gorm"
)

type Schedule struct {
	ID                 int64  `gorm:"primary key;autoIncrement" json:"id"`
	UserID             int64  `json:"user_id"`
	NameMedication     string `json:"name_medication"`
	MedicationPerDay   int    `json:"medication_per_day" validate:"required,gte=0,lte=15"`
	DurationMedication int    `json:"duration_medication"`
}

func MigrationSchedule(db *gorm.DB) error {
	err := db.AutoMigrate(&Schedule{})
	if err != nil {
		return err
	}
	return nil
}
