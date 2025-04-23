package service

import (
	"fmt"
	"schedule-app/config"
	"schedule-app/internal/app/model"
	"schedule-app/internal/app/storage"
	"time"
)

type Service struct {
	r       storage.Repository
	configs config.MedPeriodConfig
}

func New(r *storage.Repository, configs *config.MedPeriodConfig) *Service {
	return &Service{
		r:       *r,
		configs: *configs,
	}
}

func (s *Service) CreateSchedule(schedule *model.Schedule) (*model.Schedule, error) {
	createSchedule, err := s.r.CreateSchedule(schedule)
	if err != nil {
		return nil, fmt.Errorf("%w repository.CreateSchedule", err)
	}
	return createSchedule, nil
}

func (s *Service) GetUsersSchedules(userId uint64) ([]uint64, error) {
	usersSchedules, err := s.r.GetSchedulesByUserId(userId)
	if err != nil {
		return nil, fmt.Errorf("%w repository.GetSchedulesByUserId", err)
	}

	usersSchedulesId := make([]uint64, 0)
	for _, schdl := range usersSchedules {
		usersSchedulesId = append(usersSchedulesId, schdl.ID)
	}

	return usersSchedulesId, nil
}

func (s *Service) NextTaking(userId uint64) ([]model.ScheduleTo, error) {
	usersSchedules, err := s.r.GetSchedulesByUserId(userId)
	if err != nil {
		return nil, fmt.Errorf("%w repository.GetSchedulesByUserId", err)
	}

	periodC, _ := time.ParseDuration(s.configs.Period)

	nowTime := time.Now()
	t1 := time.Date(0000, 01, 01, nowTime.Hour(), nowTime.Minute(), nowTime.Second(), 0, time.UTC)
	t2 := t1.Add(periodC)

	usersSchedulesId := make([]string, 0)
	yy := make([]model.ScheduleTo, 0)

	for _, schdl := range usersSchedules {
		timePoints := s.CountTimeForMedicament(schdl.MedicationPerDay)
		for _, timePoint := range timePoints {
			if t1.Before(timePoint) && t2.After(timePoint) || t1.Equal(timePoint) || t2.Equal(timePoint) {
				usersSchedulesId = append(usersSchedulesId, timePoint.Format("15:04"))
			}
		}
		if len(usersSchedulesId) != 0 {
			yy = append(yy, model.ScheduleTo{
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

func (s *Service) GetScheduleByScheduleId(scheduleId, userId uint64) (*model.ScheduleTo, error) {
	schedule, err := s.r.GetScheduleByIdAndUserId(scheduleId, userId)
	if err != nil {
		return nil, fmt.Errorf("%w repository.GetScheduleByIdAndUserId", err)
	}

	timeForMedication := s.CountTimeForMedicament(schedule.MedicationPerDay)
	timeForMedicationInString := fromTimeToString(timeForMedication)

	usersSchedulesId := model.ScheduleTo{
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
