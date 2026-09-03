package preset_test

import (
	"testing"

	"github.com/Jeff1507/go-to-save-api/internal/common/pickup"
	"github.com/Jeff1507/go-to-save-api/internal/core/preset"
	"github.com/stretchr/testify/require"
)

func mustTime(t *testing.T, hr, min int) pickup.Time {
	t.Helper()

	pickupTime, err := pickup.NewTime(hr, min)
	require.NoError(t, err)

	return pickupTime
}

func TestPreset_New(t *testing.T) {
	tests := []struct {
		name    string
		start   pickup.Time
		end     pickup.Time
		wantErr bool
	}{
		{
			name:    "overnight interval equal to 30 minutes",
			start:   mustTime(t, 23, 45),
			end:     mustTime(t, 0, 15),
			wantErr: false,
		},
		{
			name:    "overnight interval equal to 6 hours",
			start:   mustTime(t, 22, 0),
			end:     mustTime(t, 4, 0),
			wantErr: false,
		},
		{
			name:    "overnight interval less than 30 minutes",
			start:   mustTime(t, 23, 45),
			end:     mustTime(t, 0, 14),
			wantErr: true,
		},
		{
			name:    "overnight interval greater than 6 hours",
			start:   mustTime(t, 22, 0),
			end:     mustTime(t, 4, 1),
			wantErr: true,
		},
		{
			name:    "interval equal to 30 minutes",
			start:   mustTime(t, 0, 0),
			end:     mustTime(t, 0, 30),
			wantErr: false,
		},
		{
			name:    "interval equal to 6 hours",
			start:   mustTime(t, 0, 0),
			end:     mustTime(t, 6, 0),
			wantErr: false,
		},
		{
			name:    "start equal to end",
			start:   mustTime(t, 0, 0),
			end:     mustTime(t, 0, 0),
			wantErr: true,
		},
		{
			name:    "interval less than 30 minutes",
			start:   mustTime(t, 0, 0),
			end:     mustTime(t, 0, 29),
			wantErr: true,
		},
		{
			name:    "interval greater than 6 hours",
			start:   mustTime(t, 0, 0),
			end:     mustTime(t, 6, 1),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := preset.New("id", "store-id", "Fechamento", tt.start, tt.end)

			if gotErr := err != nil; gotErr != tt.wantErr {
				t.Errorf("expected error: %t, got: %v", tt.wantErr, err)
			}
		})
	}
}
