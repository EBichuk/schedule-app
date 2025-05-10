package storage

import (
	"context"
	"fmt"
	"schedule-app/internal/domain/entity"

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

func (r *Repository) CreateSchedule(ctx context.Context, schedule *entity.Schedule) (*entity.Schedule, error) {
	err := r.db.Create(schedule).Error
	if err != nil {
		return nil, err
	}
	return schedule, nil
}

func (r *Repository) GetSchedulesByUserId(ctx context.Context, userId int64) ([]entity.Schedule, error) {
	var usersShedules []entity.Schedule

	err := r.db.Where("user_id = ?", userId).Find(&usersShedules).Error
	if err != nil {
		return nil, err
	}

	return usersShedules, nil
}

func (r *Repository) GetScheduleByIdAndUserId(ctx context.Context, scheduleId, userId int64) (*entity.Schedule, error) {
	var usersShedule entity.Schedule

	err := r.db.First(&usersShedule, "id = ?", scheduleId).Error
	if err != nil {
		return nil, err
	}

	if usersShedule.UserID != userId {
		return nil, fmt.Errorf("schedule with schedule_id %d and user_id %d not find in db", scheduleId, userId)
	}

	return &usersShedule, nil
}
