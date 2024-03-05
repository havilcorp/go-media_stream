package utils

import (
	"reflect"
	"testing"
)

func TestGenerateJWT(t *testing.T) {
	type args struct {
		userId string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Good",
			args: args{
				userId: "1",
			},
			want: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxIn0.lqqfs37UK7VsrQKS27uL5b2iZb-7abQhjY56N7wDBIk",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateJWT(tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GenerateJWT() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVerifyJWT(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Good",
			args: args{
				token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxIn0.lqqfs37UK7VsrQKS27uL5b2iZb-7abQhjY56N7wDBIk",
			},
			want: "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := VerifyJWT(tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VerifyJWT() = %v, want %v", got, tt.want)
			}
		})
	}
}
