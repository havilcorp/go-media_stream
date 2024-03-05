package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go-media-stream/internal/domain"
	"go-media-stream/internal/handlers/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHomeHandler_MainPage(t *testing.T) {
	videoProvider := mocks.NewVideoProvider(t)
	audioProvider := mocks.NewAudioProvider(t)

	videoProvider.On("GetVideos", mock.Anything).Return(&[]domain.Video{}, nil)

	type args struct {
		statusCode int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "GetVideos",
			args: args{
				statusCode: 200,
			},
		},
	}
	for _, tt := range tests {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		rw := httptest.NewRecorder()
		t.Run(tt.name, func(t *testing.T) {
			h := &HomeHandler{
				video: videoProvider,
				audio: audioProvider,
			}
			h.MainPage(rw, r)
			res := rw.Result()
			assert.Equal(t, tt.args.statusCode, res.StatusCode)
		})
	}
}
