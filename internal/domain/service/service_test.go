package service

import (
	"reflect"
	"testing"
	"time"
)

func TestCountTimeForMedicament(t *testing.T) {
	type args struct {
		medicationPerDay int
		startTime        time.Time
		finishTime       time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name:    "borderline: period 08:00-22:00",
			args:    args{medicationPerDay: 1, startTime: time.Date(0, 0, 0, 8, 0, 0, 0, time.UTC), finishTime: time.Date(0, 0, 0, 22, 0, 0, 0, time.UTC)},
			want:    []string{"08:00"},
			wantErr: false,
		},
		{
			name:    "success: 2 times",
			args:    args{medicationPerDay: 2, startTime: time.Date(0, 0, 0, 8, 0, 0, 0, time.UTC), finishTime: time.Date(0, 0, 0, 22, 0, 0, 0, time.UTC)},
			want:    []string{"08:00", "22:00"},
			wantErr: false,
		},
		{
			name:    "success: 3 times",
			args:    args{medicationPerDay: 3, startTime: time.Date(0, 0, 0, 8, 0, 0, 0, time.UTC), finishTime: time.Date(0, 0, 0, 22, 0, 0, 0, time.UTC)},
			want:    []string{"08:00", "15:00", "22:00"},
			wantErr: false,
		},
		{
			name:    "borderline: period 08:00-22:00",
			args:    args{medicationPerDay: 15, startTime: time.Date(0, 0, 0, 8, 0, 0, 0, time.UTC), finishTime: time.Date(0, 0, 0, 22, 0, 0, 0, time.UTC)},
			want:    []string{"08:00", "09:00", "10:00", "11:00", "12:00", "13:00", "14:00", "15:00", "16:00", "17:00", "18:00", "19:00", "20:00", "21:00", "22:00"},
			wantErr: false,
		},
		{
			name:    "negative: big number",
			args:    args{medicationPerDay: 25, startTime: time.Date(0, 0, 0, 8, 0, 0, 0, time.UTC), finishTime: time.Date(0, 0, 0, 22, 0, 0, 0, time.UTC)},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "zero",
			args:    args{medicationPerDay: 0, startTime: time.Date(0, 0, 0, 8, 0, 0, 0, time.UTC), finishTime: time.Date(0, 0, 0, 22, 0, 0, 0, time.UTC)},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "negative: negative number",
			args:    args{medicationPerDay: -12, startTime: time.Date(0, 0, 0, 8, 0, 0, 0, time.UTC), finishTime: time.Date(0, 0, 0, 22, 0, 0, 0, time.UTC)},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "success: 4 period 10:00-22:00",
			args:    args{medicationPerDay: 4, startTime: time.Date(0, 0, 0, 10, 0, 0, 0, time.UTC), finishTime: time.Date(0, 0, 0, 22, 0, 0, 0, time.UTC)},
			want:    []string{"10:00", "14:00", "18:00", "22:00"},
			wantErr: false,
		},
		{
			name:    "negative: period 22:00-20:00",
			args:    args{medicationPerDay: 7, startTime: time.Date(0, 0, 0, 22, 0, 0, 0, time.UTC), finishTime: time.Date(0, 0, 0, 20, 0, 0, 0, time.UTC)},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "success: 6 period 10:00-22:00",
			args:    args{medicationPerDay: 6, startTime: time.Date(0, 0, 0, 10, 0, 0, 0, time.UTC), finishTime: time.Date(0, 0, 0, 22, 0, 0, 0, time.UTC)},
			want:    []string{"10:00", "12:30", "14:45", "17:15", "19:30", "22:00"},
			wantErr: false,
		},
		{
			name:    "border line: 13 period 10:00-22:00",
			args:    args{medicationPerDay: 13, startTime: time.Date(0, 0, 0, 10, 0, 0, 0, time.UTC), finishTime: time.Date(0, 0, 0, 22, 0, 0, 0, time.UTC)},
			want:    []string{"10:00", "11:00", "12:00", "13:00", "14:00", "15:00", "16:00", "17:00", "18:00", "19:00", "20:00", "21:00", "22:00"},
			wantErr: false,
		},
		{
			name:    "negative: 14 period 10:00-22:00",
			args:    args{medicationPerDay: 14, startTime: time.Date(0, 0, 0, 10, 0, 0, 0, time.UTC), finishTime: time.Date(0, 0, 0, 22, 0, 0, 0, time.UTC)},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "success: 12 period 10:00-22:00",
			args:    args{medicationPerDay: 12, startTime: time.Date(0, 0, 0, 10, 0, 0, 0, time.UTC), finishTime: time.Date(0, 0, 0, 22, 0, 0, 0, time.UTC)},
			want:    []string{"10:00", "11:00", "12:15", "13:15", "14:15", "15:30", "16:30", "17:45", "18:45", "19:45", "21:00", "22:00"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CountTimeForMedicament(tt.args.medicationPerDay, tt.args.startTime, tt.args.finishTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("CountTimeForMedicament() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CountTimeForMedicament() = %v, want %v", got, tt.want)
			}
		})
	}
}
