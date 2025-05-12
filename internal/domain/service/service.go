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
		usersSchedulesId = append(usersSchedulesId, schdl.Id)
	}
	slog.InfoContext(ctx, "GetUsersSchedules OK")
	return usersSchedulesId, nil
}

func (s *Service) NextTaking(ctx context.Context, userId int64) ([]entity.ScheduleWithTime, error) {
	usersSchedules, err := s.r.GetSchedulesByUserId(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("%w repository.GetSchedulesByUserId", err)
	}

	periodC, _ := time.ParseDuration(s.configs.Period)

	nowTime := time.Now()
	t1 := time.Date(0000, 01, 01, nowTime.Hour(), nowTime.Minute(), nowTime.Second(), 0, time.UTC)
	t2 := t1.Add(periodC)

	usersSchedulesId := make([]string, 0)
	scheduletime := make([]entity.ScheduleWithTime, 0)

	for _, schdl := range usersSchedules {
		timePoints, err := s.GetDurationToTakePills(schdl.MedicationPerDay)
		if err != nil {
			return nil, err
		}
		for _, timePoin := range timePoints {
			timePoint, _ := time.Parse("15:04", timePoin)
			if t1.Before(timePoint) && t2.After(timePoint) || t1.Equal(timePoint) || t2.Equal(timePoint) {
				usersSchedulesId = append(usersSchedulesId, timePoint.Format("15:04"))
			}
		}
		if len(usersSchedulesId) != 0 {
			scheduletime = append(scheduletime, entity.ScheduleWithTime{
				Id:                 schdl.Id,
				NameMedication:     schdl.NameMedication,
				MedicationPerDay:   schdl.MedicationPerDay,
				ScheduleMedication: usersSchedulesId,
			})
			usersSchedulesId = make([]string, 0)
		}
	}

	return scheduletime, nil
}

func (s *Service) GetScheduleByScheduleId(ctx context.Context, scheduleId, userId int64) (*entity.ScheduleWithTime, error) {
	schedule, err := s.r.GetScheduleByIdAndUserId(ctx, scheduleId, userId)
	if err != nil {
		return nil, fmt.Errorf("%w repository.GetScheduleByIdAndUserId", err)
	}

	timeForMedication, err := s.GetDurationToTakePills(schedule.MedicationPerDay)
	if err != nil {
		return nil, err
	}

	usersSchedulesId := entity.ScheduleWithTime{
		Id:                 scheduleId,
		NameMedication:     schedule.NameMedication,
		MedicationPerDay:   schedule.MedicationPerDay,
		ScheduleMedication: timeForMedication,
	}
	return &usersSchedulesId, nil
}

func (s *Service) GetDurationToTakePills(medicationPerDay int) ([]string, error) {
	startTime, _ := time.Parse(time.TimeOnly, s.configs.Start)
	finishTime, _ := time.Parse(time.TimeOnly, s.configs.End)
	timeSchedule, err := CountTimeForMedicament(medicationPerDay, startTime, finishTime)
	return timeSchedule, err
}

func CountTimeForMedicament(medicationPerDay int, startTime time.Time, finishTime time.Time) ([]string, error) {
	if medicationPerDay < 1 || medicationPerDay > 24 {
		return nil, fmt.Errorf("%s CountTimeMedication", "invalid medicationPerDay")
	}

	durationToTakePills := finishTime.Sub(startTime)
	if int(durationToTakePills.Hours()+1) < medicationPerDay {
		return nil, fmt.Errorf("%s CountTimeMedication", "medicationPerDay more that 1 time an hour")
	}

	period := time.Duration(medicationPerDay - 1)
	if medicationPerDay > 1 {
		period = durationToTakePills / time.Duration(medicationPerDay-1)
	}

	var medicationPeriod []string
	d := time.Duration(15 * time.Minute)

	for i := 0; i < medicationPerDay; i++ {
		medicationPeriod = append(medicationPeriod, startTime.Round(d).Format("15:04"))
		startTime = startTime.Add(period)
	}

	return medicationPeriod, nil
}
