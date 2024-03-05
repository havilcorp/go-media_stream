package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-media-stream/internal/domain"
	"go-media-stream/internal/handlers/mocks"
	"go-media-stream/internal/log"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthHandler_Auth(t *testing.T) {
	logger := log.New()
	authProvider := mocks.NewAuthProvider(t)

	authProvider.On("SignUp", mock.Anything, "login", "password").Return("token", nil)
	authProvider.On("SignUp", mock.Anything, "login", "wrong").Return("", domain.ErrWrongPassword)
	authProvider.On("SignUp", mock.Anything, "wrong", "wrong").Return("", domain.ErrWrongPassword)
	authProvider.On("SignUp", mock.Anything, "servererr", "servererr").Return("", errors.New("OTHER_ERROR"))

	type args struct {
		login      string
		password   string
		statusCode int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "SignUp 200",
			args: args{
				login:      "login",
				password:   "password",
				statusCode: 200,
			},
		},
		{
			name: "SignUp wrong password 401",
			args: args{
				login:      "login",
				password:   "wrong",
				statusCode: 401,
			},
		},
		{
			name: "SignUp wrong all 401",
			args: args{
				login:      "wrong",
				password:   "wrong",
				statusCode: 401,
			},
		},
		{
			name: "SignUp server err 500",
			args: args{
				login:      "servererr",
				password:   "servererr",
				statusCode: 500,
			},
		},
	}
	for _, tt := range tests {
		r := httptest.NewRequest(
			http.MethodPost,
			"/auth",
			strings.NewReader(
				fmt.Sprintf("{\"login\": \"%s\", \"password\": \"%s\"}", tt.args.login, tt.args.password),
			),
		)
		rw := httptest.NewRecorder()
		t.Run(tt.name, func(t *testing.T) {
			h := &AuthHandler{
				logger:       logger,
				authProvider: authProvider,
			}
			h.Auth(rw, r)
			res := rw.Result()
			assert.Equal(t, tt.args.statusCode, res.StatusCode)
		})
	}
}
