package storage

import (
	"fmt"
	"schedule-app/internal/app/model"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateSchedule(schedule *model.Schedule) (*model.Schedule, error) {
	err := r.db.Create(schedule).Error
	if err != nil {
		return nil, err
	}
	return schedule, nil
}

func (r *Repository) GetSchedulesByUserId(userId uint64) ([]model.Schedule, error) {
	var usersShedules []model.Schedule

	err := r.db.Where("user_id = ?", userId).Find(&usersShedules).Error
	if err != nil {
		return nil, err
	}

	return usersShedules, nil
}

func (r *Repository) GetScheduleByIdAndUserId(scheduleId, userId uint64) (*model.Schedule, error) {
	var usersShedule model.Schedule

	err := r.db.First(&usersShedule, "id = ?", scheduleId).Error
	if err != nil {
		return nil, err
	}

	if usersShedule.UserID != userId {
		return nil, fmt.Errorf("schedule with schedule_id %d and user_id %d not find in db", scheduleId, userId)
	}

	return &usersShedule, nil
}
