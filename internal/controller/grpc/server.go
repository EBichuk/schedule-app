package server

import (
	"context"
	"log/slog"
	"schedule-app/internal/domain/entity"
	grpc_service "schedule-app/proto/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service interface {
	CreateSchedule(context.Context, *entity.Schedule) (*entity.Schedule, error)
	GetUsersSchedules(context.Context, int64) ([]int64, error)
	GetScheduleByScheduleId(context.Context, int64, int64) (*entity.ScheduleTo, error)
	NextTaking(context.Context, int64) ([]entity.ScheduleTo, error)
}

type serverAPI struct {
	grpc_service.UnimplementedUserServiceServer
	s      Service
	logger *slog.Logger
}

func RegisterServerAPI(gRPC *grpc.Server, s Service, l *slog.Logger) {
	grpc_service.RegisterUserServiceServer(gRPC, &serverAPI{s: s, logger: l})
}

func (s *serverAPI) GetSchedulesByUser(ctx context.Context, req *grpc_service.UserRequest) (*grpc_service.SchedulesIDPesponce, error) {
	if req.GetUserId() == 0 {
		s.logger.ErrorContext(ctx, "user_id is required")
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	schedulesByUser, err := s.s.GetUsersSchedules(ctx, req.GetUserId())

	if err != nil {
		s.logger.ErrorContext(ctx, err.Error())
		return nil, status.Error(500, "invalid user id")
	}

	return &grpc_service.SchedulesIDPesponce{
		Schedules: schedulesByUser,
	}, nil
}

func (s *serverAPI) CreateSchedule(ctx context.Context, req *grpc_service.CreateScheduleRequest) (*grpc_service.ScheduleResponce, error) {
	if req.GetUserId() == 0 {
		s.logger.ErrorContext(ctx, "user_id is required")
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	schedule := entity.Schedule{
		UserID:             req.GetUserId(),
		NameMedication:     req.NameMedication,
		MedicationPerDay:   int(req.MedicationPerDay),
		DurationMedication: int(req.DurationMedication),
	}

	createdSchedule, err := s.s.CreateSchedule(ctx, &schedule)
	if err != nil {
		s.logger.ErrorContext(ctx, err.Error())
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
	if req.GetUserId() == 0 {
		s.logger.ErrorContext(ctx, "user_id is required")
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.GetScheduleId() == 0 {
		s.logger.ErrorContext(ctx, "schedule_id is required")
		return nil, status.Error(codes.InvalidArgument, "schedule_id is required")
	}

	scheduleById, err := s.s.GetScheduleByScheduleId(ctx, req.GetScheduleId(), req.GetUserId())

	if err != nil {
		s.logger.ErrorContext(ctx, err.Error())
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
	if req.GetUserId() == 0 {
		s.logger.ErrorContext(ctx, "user_id is required")
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	schedulesByUser, err := s.s.NextTaking(ctx, req.GetUserId())

	if err != nil {
		s.logger.ErrorContext(ctx, err.Error())
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
