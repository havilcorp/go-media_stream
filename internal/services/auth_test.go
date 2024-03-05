package services

import (
	"context"
	"testing"

	"go-media-stream/internal/domain"
	"go-media-stream/internal/services/mocks"
	"go-media-stream/internal/utils"

	"github.com/stretchr/testify/mock"
)

func TestAuthService_SignUp(t *testing.T) {
	authProvider := mocks.NewAuthProvider(t)

	authProvider.On("GetUserByLogin", mock.Anything, "login", utils.GetMD5("password")).Return(
		&domain.User{
			ID:       1,
			Login:    "login",
			Password: utils.GetMD5("password"),
		}, nil,
	)

	authProvider.On("GetUserByLogin", mock.Anything, "none", utils.GetMD5("password")).Return(nil, domain.ErrUserNotFound)

	authProvider.On("GetUserByLogin", mock.Anything, "login", utils.GetMD5("wrong")).Return(nil, domain.ErrWrongPassword)

	authProvider.On("CreateUser", mock.Anything, &domain.User{
		Login:    "none",
		Password: utils.GetMD5("password"),
	}).Return(int64(2), nil)

	authProvider.On("GetUserById", mock.Anything, int64(2)).Return(
		&domain.User{
			ID:       2,
			Login:    "none",
			Password: utils.GetMD5("password"),
		}, nil,
	)

	type args struct {
		login    string
		password string
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
				login:    "login",
				password: "password",
			},
			want:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxIn0.lqqfs37UK7VsrQKS27uL5b2iZb-7abQhjY56N7wDBIk",
			wantErr: false,
		},
		{
			name: "Good2",
			args: args{
				login:    "none",
				password: "password",
			},
			want:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIyIn0.7epdXB_aXtN0dKO4HhqYRpwhvFd_A1KPCxTHPtWsdhw",
			wantErr: false,
		},
		{
			name: "wrong password",
			args: args{
				login:    "login",
				password: "wrong",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AuthService{
				authProvider: authProvider,
			}
			got, err := s.SignUp(context.Background(), tt.args.login, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthService.SignUp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AuthService.SignUp() = %v, want %v", got, tt.want)
			}
		})
	}
}
