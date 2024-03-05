package utils

import "testing"

func TestGetProjectRoot(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Good",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetProjectRoot()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProjectRoot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == "" {
				t.Errorf("GetProjectRoot() = %v", got)
			}
		})
	}
}
