package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"schedule-app/pkg/dbtest"
	api "schedule-app/pkg/types"
)

func (s *Suite) TestGetScheduleById() {
	rq := s.Require()
	//ctx := context.Background()

	err := dbtest.MigrateFromFile(s.Db, "testdata/get_schedule.sql")
	rq.NoError(err)
	testCases := []struct {
		name              string
		userId            int64
		scheduleId        int64
		expectedResponse  api.ScheduleResponce
		expectedErrorCode int
	}{
		{
			name:       "Success",
			userId:     1,
			scheduleId: 4,
			expectedResponse: api.ScheduleResponce{
				Id:               4,
				NameMedication:   "Фарингосепт",
				MedicationPerDay: 1,
				Schedule:         []string{"08:00"},
			},
			expectedErrorCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp, err := s.httpClient.Get(fmt.Sprintf("%s/schedule/%d/%d", s.baseHTTPURL, tc.userId, tc.scheduleId))
			rq.NoError(err)
			//rq.Equal(resp., tc.expectedResponse)
			defer resp.Body.Close()

			if http.StatusOK == resp.StatusCode {
				var dd api.ScheduleResponce
				err = json.NewDecoder(resp.Body).Decode(&dd)
				rq.NoError(err)
				s.Equal(tc.expectedResponse, dd)
			}
		})
	}
}

func (s *Suite) TestGetSchedulesByUser() {
	rq := s.Require()
	//ctx := context.Background()

	err := dbtest.MigrateFromFile(s.Db, "testdata/get_schedule.sql")
	rq.NoError(err)
	testCases := []struct {
		name              string
		userId            int64
		expectedResponse  []int64
		expectedErrorCode int
	}{
		{
			name:              "Success",
			userId:            1,
			expectedResponse:  []int64{1, 2, 4, 5, 6, 7, 8, 9},
			expectedErrorCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp, err := s.httpClient.Get(fmt.Sprintf("%s/schedules/%d/", s.baseHTTPURL, tc.userId))
			rq.NoError(err)
			defer resp.Body.Close()

			if http.StatusOK == resp.StatusCode {
				s.Equal(tc.expectedResponse, resp.Body)
			}
		})
	}
}
