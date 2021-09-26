package aoc

import "testing"

func TestDay_Validate(t *testing.T) {
	tests := []struct {
		name    string
		d       Day
		wantErr bool
	}{
		{
			name:    "Valid day",
			d:       13,
			wantErr: false,
		},
		{
			name:    "Day value too small",
			d:       0,
			wantErr: true,
		},
		{
			name:    "Small valid value",
			d:       1,
			wantErr: false,
		},
		{
			name:    "Large valid value",
			d:       25,
			wantErr: false,
		},
		{
			name:    "Too large value",
			d:       26,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.d.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
