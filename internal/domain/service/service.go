package service

import (
	"context"
	"fmt"
	"log/slog"
	"schedule-app/internal/config"
	"schedule-app/internal/domain/entity"
	"time"
)

type repository interface {
	CreateSchedule(context.Context, *entity.Schedule) (*entity.Schedule, error)
	GetSchedulesByUserId(context.Context, int64) ([]entity.Schedule, error)
	GetScheduleByIdAndUserId(context.Context, int64, int64) (*entity.Schedule, error)
}

type Service struct {
	r       repository
	configs config.MedPeriodConfig
}

func New(r repository, configs *config.MedPeriodConfig) *Service {
	return &Service{
		r:       r,
		configs: *configs,
	}
}

func (s *Service) CreateSchedule(ctx context.Context, schedule *entity.Schedule) (*entity.Schedule, error) {
	createSchedule, err := s.r.CreateSchedule(ctx, schedule)
	if err != nil {
		return nil, fmt.Errorf("%w repository.CreateSchedule", err)
	}
	return createSchedule, nil
}

func (s *Service) GetUsersSchedules(ctx context.Context, userId int64) ([]int64, error) {
	usersSchedules, err := s.r.GetSchedulesByUserId(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("%w repository.GetSchedulesByUserId", err)
	}

	usersSchedulesId := make([]int64, 0)
	for _, schdl := range usersSchedules {
		usersSchedulesId = append(usersSchedulesId, schdl.ID)
	}
	slog.InfoContext(ctx, "GetUsersSchedules OK")
	return usersSchedulesId, nil
}

func (s *Service) NextTaking(ctx context.Context, userId int64) ([]entity.ScheduleTo, error) {
	usersSchedules, err := s.r.GetSchedulesByUserId(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("%w repository.GetSchedulesByUserId", err)
	}

	periodC, _ := time.ParseDuration(s.configs.Period)

	nowTime := time.Now()
	t1 := time.Date(0000, 01, 01, nowTime.Hour(), nowTime.Minute(), nowTime.Second(), 0, time.UTC)
	t2 := t1.Add(periodC)

	usersSchedulesId := make([]string, 0)
	yy := make([]entity.ScheduleTo, 0)

	for _, schdl := range usersSchedules {
		timePoints := s.CountTimeForMedicament(schdl.MedicationPerDay)
		for _, timePoint := range timePoints {
			if t1.Before(timePoint) && t2.After(timePoint) || t1.Equal(timePoint) || t2.Equal(timePoint) {
				usersSchedulesId = append(usersSchedulesId, timePoint.Format("15:04"))
			}
		}
		if len(usersSchedulesId) != 0 {
			yy = append(yy, entity.ScheduleTo{
				ID:                 schdl.ID,
				NameMedication:     schdl.NameMedication,
				MedicationPerDay:   schdl.MedicationPerDay,
				ScheduleMedication: usersSchedulesId,
			})
			usersSchedulesId = make([]string, 0)
		}
	}

	return yy, nil
}

func (s *Service) GetScheduleByScheduleId(ctx context.Context, scheduleId, userId int64) (*entity.ScheduleTo, error) {
	schedule, err := s.r.GetScheduleByIdAndUserId(ctx, scheduleId, userId)
	if err != nil {
		return nil, fmt.Errorf("%w repository.GetScheduleByIdAndUserId", err)
	}

	timeForMedication := s.CountTimeForMedicament(schedule.MedicationPerDay)
	timeForMedicationInString := fromTimeToString(timeForMedication)

	usersSchedulesId := entity.ScheduleTo{
		ID:                 scheduleId,
		NameMedication:     schedule.NameMedication,
		MedicationPerDay:   schedule.MedicationPerDay,
		ScheduleMedication: timeForMedicationInString,
	}
	return &usersSchedulesId, nil
}

func (s *Service) GetDurationToTakePills() time.Duration {
	startTime, _ := time.Parse(time.TimeOnly, s.configs.Start)
	finishTime, _ := time.Parse(time.TimeOnly, s.configs.End)
	return finishTime.Sub(startTime)
}

func (s *Service) CountTimeForMedicament(medicationPerDay int) []time.Time {
	if medicationPerDay < 1 {
		return nil
	}

	period := time.Duration(medicationPerDay - 1)
	if medicationPerDay > 1 {
		period = s.GetDurationToTakePills() / time.Duration(medicationPerDay-1)
	}

	var medicationPeriod []time.Time
	d := time.Duration(15 * time.Minute)

	timePoint, _ := time.Parse(time.TimeOnly, s.configs.Start)
	for i := 0; i < medicationPerDay; i++ {
		medicationPeriod = append(medicationPeriod, timePoint.Round(d))
		timePoint = timePoint.Add(period)
	}

	return medicationPeriod
}

func fromTimeToString(timeInTime []time.Time) []string {
	stringTime := make([]string, 0)
	for _, tm := range timeInTime {
		stringTime = append(stringTime, tm.Format("15:04"))
	}
	return stringTime
}
