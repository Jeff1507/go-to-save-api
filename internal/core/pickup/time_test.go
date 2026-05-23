package pickup_test

import (
	"testing"

	"github.com/Jeff1507/go-to-save-api/internal/core/pickup"
)

func TestTime_New(t *testing.T) {
	tests := []struct {
		name    string
		hour    int
		minute  int
		wantErr bool
	}{
		{
			name:    "valid minimum time",
			hour:    0,
			minute:  0,
			wantErr: false,
		},
		{
			name:    "valid maximum time",
			hour:    23,
			minute:  59,
			wantErr: false,
		},
		{
			name:    "negative hour",
			hour:    -1,
			minute:  0,
			wantErr: true,
		},
		{
			name:    "hour greater than 23",
			hour:    24,
			minute:  0,
			wantErr: true,
		},
		{
			name:    "negative minute",
			hour:    0,
			minute:  -1,
			wantErr: true,
		},
		{
			name:    "minute greater than 59",
			hour:    0,
			minute:  60,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := pickup.NewTime(tt.hour, tt.minute)

			if gotErr := err != nil; gotErr != tt.wantErr {
				t.Errorf("expected error: %t, got: %v", tt.wantErr, err)
			}
		})
	}
}
