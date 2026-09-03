package pickup_test

import (
	"testing"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/common/pickup"
)

func makeTime(hr, min int) time.Time {
	return time.Date(2025, time.January, 1, hr, min, 0, 0, time.UTC)
}

func TestWindow_New(t *testing.T) {
	tests := []struct {
		name    string
		start   time.Time
		end     time.Time
		wantErr bool
	}{
		{
			name:    "start equal to end",
			start:   makeTime(22, 0),
			end:     makeTime(22, 0),
			wantErr: true,
		},
		{
			name:    "overnight",
			start:   makeTime(22, 0),
			end:     makeTime(28, 0),
			wantErr: false,
		},
		{
			name:    "start is after end",
			start:   makeTime(28, 0),
			end:     makeTime(22, 0),
			wantErr: true,
		},
		{
			name:    "interval is equal to 30 minutes",
			start:   makeTime(22, 0),
			end:     makeTime(22, 30),
			wantErr: false,
		},
		{
			name:    "interval is less than 30 minutes",
			start:   makeTime(22, 0),
			end:     makeTime(22, 29),
			wantErr: true,
		},
		{
			name:    "interval is equal to 6 hours",
			start:   makeTime(10, 0),
			end:     makeTime(16, 0),
			wantErr: false,
		},
		{
			name:    "interval is greater than 6 hours",
			start:   makeTime(10, 0),
			end:     makeTime(16, 1),
			wantErr: true,
		},
		{
			name:    "overnight interval is equal to 30 minutes",
			start:   makeTime(23, 45),
			end:     makeTime(24, 15),
			wantErr: false,
		},
		{
			name:    "overnight interval is less than 30 minutes",
			start:   makeTime(23, 45),
			end:     makeTime(24, 14),
			wantErr: true,
		},
		{
			name:    "overnight interval is equal to 6 hours",
			start:   makeTime(20, 0),
			end:     makeTime(26, 0),
			wantErr: false,
		},
		{
			name:    "overnight interval is greater than 6 hours",
			start:   makeTime(20, 0),
			end:     makeTime(26, 1),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := pickup.NewWindow(tt.start, tt.end)

			if gotErr := err != nil; gotErr != tt.wantErr {
				t.Errorf("expected error: %t, got: %v", tt.wantErr, err)
			}
		})
	}
}
