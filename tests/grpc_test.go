package tests

import (
	"context"
	"schedule-app/pkg/dbtest"
	grpc_service "schedule-app/pkg/grpc/gen"

	"google.golang.org/grpc/codes"
)

func (s *Suite) TestGetScheduleGrpc() {
	rq := s.Require()
	ctx := context.Background()

	err := dbtest.MigrateFromFile(s.Db, "testdata/get_schedule.sql")
	rq.NoError(err)

	testCases := []struct {
		name              string
		request           *grpc_service.UserRequest
		expectedResponse  *grpc_service.SchedulesIDPesponce
		expectedErrorCode codes.Code
	}{
		{
			name: "Success. User with 9 schedules",
			request: &grpc_service.UserRequest{
				UserId: 1,
			},
			expectedResponse: &grpc_service.SchedulesIDPesponce{
				Schedules: []int64{1, 2, 4, 5, 6, 7, 8, 9},
			},
			expectedErrorCode: codes.OK,
		},
		{
			name: "Success. User with 1 schedule",
			request: &grpc_service.UserRequest{
				UserId: 3,
			},
			expectedResponse: &grpc_service.SchedulesIDPesponce{
				Schedules: []int64{3},
			},
			expectedErrorCode: codes.OK,
		},
		{
			name: "Success. User withhout schedule",
			request: &grpc_service.UserRequest{
				UserId: 2,
			},
			expectedResponse: &grpc_service.SchedulesIDPesponce{
				Schedules: nil,
			},
			expectedErrorCode: codes.OK,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp, err := s.grpcClient.GetSchedulesByUser(ctx, tc.request)
			rq.Equal(resp.Schedules, tc.expectedResponse.Schedules)
			rq.NoError(err)
		})
	}
}

func (s *Suite) TestGetScheduleByIdGrpc() {
	rq := s.Require()
	ctx := context.Background()

	err := dbtest.MigrateFromFile(s.Db, "testdata/get_schedule.sql")
	rq.NoError(err)
	testCases := []struct {
		name              string
		request           *grpc_service.ScheduleRequest
		expectedResponse  *grpc_service.ScheduleTimeResponce
		expectedErrorCode codes.Code
	}{
		{
			name: "Success. 1 time in day",
			request: &grpc_service.ScheduleRequest{
				ScheduleId: 4,
				UserId:     1,
			},
			expectedResponse: &grpc_service.ScheduleTimeResponce{
				ScheduleId:       4,
				NameMedication:   "Фарингосепт",
				MedicationPerDay: 1,
				Schedule:         []string{"08:00"},
			},
			expectedErrorCode: codes.OK,
		},
		{
			name: "Success. 10 times in day",
			request: &grpc_service.ScheduleRequest{
				ScheduleId: 1,
				UserId:     1,
			},
			expectedResponse: &grpc_service.ScheduleTimeResponce{
				ScheduleId:       1,
				NameMedication:   "Ромашка",
				MedicationPerDay: 10,
				Schedule:         []string{"08:00", "09:30", "11:00", "12:45", "14:15", "15:45", "17:15", "19:00", "20:30", "22:00"},
			},
			expectedErrorCode: codes.OK,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp, err := s.grpcClient.GetScheduleById(ctx, tc.request)
			rq.NoError(err)

			rq.Equal(resp.ScheduleId, tc.expectedResponse.ScheduleId)
			rq.Equal(resp.MedicationPerDay, tc.expectedResponse.MedicationPerDay)
			rq.Equal(resp.Schedule, tc.expectedResponse.Schedule)
			rq.Equal(resp.NameMedication, tc.expectedResponse.NameMedication)
		})
	}
}
