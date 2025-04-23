package server

import (
	"context"
	"schedule-app/internal/app/model"
	grpc_service "schedule-app/proto/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type Service interface {
	CreateSchedule(*model.Schedule) (*model.Schedule, error)
	GetUsersSchedules(uint64) ([]uint64, error)
	GetScheduleByScheduleId(uint64, uint64) (*model.ScheduleTo, error)
	NextTaking(uint64) ([]model.ScheduleTo, error)
}

type serverAPI struct {
	grpc_service.UnimplementedUserServiceServer
	s Service
}

func RegisterServerAPI(gRPC *grpc.Server, s Service) {
	grpc_service.RegisterUserServiceServer(gRPC, &serverAPI{s: s})
}

func (s *serverAPI) GetSchedulesByUser(ctx context.Context, req *grpc_service.UserRequest) (*grpc_service.SchedulesIDPesponce, error) {
	schedulesByUser, err := s.s.GetUsersSchedules(req.GetUserId())

	if err != nil {
		return nil, status.Error(500, "invalid user id")
	}

	return &grpc_service.SchedulesIDPesponce{
		Schedules: schedulesByUser,
	}, nil
}

func (s *serverAPI) CreateSchedule(ctx context.Context, req *grpc_service.CreateScheduleRequest) (*grpc_service.ScheduleResponce, error) {
	schedule := model.Schedule{
		UserID:             req.GetUserId(),
		NameMedication:     req.NameMedication,
		MedicationPerDay:   int(req.MedicationPerDay),
		DurationMedication: int(req.DurationMedication),
	}

	createdSchedule, err := s.s.CreateSchedule(&schedule)
	if err != nil {
		return nil, status.Error(500, "invalid user id")
	}

	return &grpc_service.ScheduleResponce{
		ScheduleId:         createdSchedule.ID,
		UserId:             createdSchedule.UserID,
		NameMedication:     createdSchedule.NameMedication,
		MedicationPerDay:   int64(createdSchedule.MedicationPerDay),
		DurationMedication: int64(createdSchedule.DurationMedication),
	}, nil
}

func (s *serverAPI) GetScheduleById(ctx context.Context, req *grpc_service.ScheduleRequest) (*grpc_service.ScheduleTimeResponce, error) {
	scheduleById, err := s.s.GetScheduleByScheduleId(req.GetScheduleId(), req.GetUserId())

	if err != nil {
		return nil, status.Error(500, "invalid user id")
	}

	return &grpc_service.ScheduleTimeResponce{
		ScheduleId:       scheduleById.ID,
		NameMedication:   scheduleById.NameMedication,
		MedicationPerDay: int64(scheduleById.MedicationPerDay),
		Schedule:         scheduleById.ScheduleMedication,
	}, nil
}

func (s *serverAPI) NextTaking(ctx context.Context, req *grpc_service.UserRequest) (*grpc_service.SchedulesResponce, error) {
	schedulesByUser, err := s.s.NextTaking(req.GetUserId())

	if err != nil {
		return nil, status.Error(500, "invalid user id")
	}

	var fromScheduleToProto []*grpc_service.ScheduleTimeResponce
	for _, i := range schedulesByUser {
		fromScheduleToProto = append(fromScheduleToProto, &grpc_service.ScheduleTimeResponce{
			ScheduleId:       i.ID,
			NameMedication:   i.NameMedication,
			MedicationPerDay: int64(i.MedicationPerDay),
			Schedule:         i.ScheduleMedication,
		})
	}

	return &grpc_service.SchedulesResponce{
		NextSchedule: fromScheduleToProto,
	}, nil
}
