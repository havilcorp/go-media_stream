package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"go-media-stream/internal/domain"
	"go-media-stream/internal/handlers/mocks"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestVideoHandler_GetVideo(t *testing.T) {
	videoProvider := mocks.NewVideoProvider(t)
	audioProvider := mocks.NewAudioProvider(t)

	videoProvider.On("GetVideoById", mock.Anything, 1).Return(&domain.Video{
		Id:   1,
		Name: "Video1",
	}, nil)
	videoProvider.On("GetVideoById", mock.Anything, 2).Return(nil, errors.New("NOT_FOUND"))
	audioProvider.On("GetAudioByVideoId", mock.Anything, 1).Return(&[]domain.Audio{
		{
			Id:      1,
			Name:    "Audio1",
			Idx:     "1",
			VideoId: 1,
		},
	}, nil)

	type args struct {
		videoId      int
		audioVideoId int
		statusCode   int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "GetVideoById 200",
			args: args{
				videoId:      1,
				audioVideoId: 1,
				statusCode:   200,
			},
		},
		{
			name: "GetVideoById 404",
			args: args{
				videoId:      2,
				audioVideoId: 1,
				statusCode:   404,
			},
		},
	}
	for _, tt := range tests {
		rw := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/video/{id}", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", strconv.Itoa(tt.args.videoId))
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		t.Run(tt.name, func(t *testing.T) {
			h := &VideoHandler{
				video: videoProvider,
				audio: audioProvider,
			}
			h.GetVideo(rw, r)
			res := rw.Result()
			assert.Equal(t, tt.args.statusCode, res.StatusCode)
		})
	}
}

func TestVideoHandler_SetTime(t *testing.T) {
	videoProvider := mocks.NewVideoProvider(t)
	audioProvider := mocks.NewAudioProvider(t)

	var time float32 = 100
	videoProvider.On("SetTime", mock.Anything, 1, time).Return(nil)
	videoProvider.On("SetTime", mock.Anything, 2, time).Return(errors.New("SERVER ERROR"))

	type args struct {
		videoId    int
		time       float32
		statusCode int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "SetTime 200",
			args: args{
				videoId:    1,
				time:       100,
				statusCode: 200,
			},
		},
		{
			name: "SetTime 500",
			args: args{
				videoId:    2,
				time:       100,
				statusCode: 500,
			},
		},
	}
	for _, tt := range tests {
		rw := httptest.NewRecorder()
		r := httptest.NewRequest(
			"GET",
			"/video/{id}/time",
			strings.NewReader(fmt.Sprintf("{\"time\": %f}", tt.args.time)),
		)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", strconv.Itoa(tt.args.videoId))
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		t.Run(tt.name, func(t *testing.T) {
			h := &VideoHandler{
				video: videoProvider,
				audio: audioProvider,
			}
			h.SetTime(rw, r)
			res := rw.Result()
			assert.Equal(t, tt.args.statusCode, res.StatusCode)
		})
	}
}

func TestVideoHandler_SetAudio(t *testing.T) {
	videoProvider := mocks.NewVideoProvider(t)
	audioProvider := mocks.NewAudioProvider(t)

	videoProvider.On("SetAudio", mock.Anything, 1, 1).Return(nil)
	videoProvider.On("SetAudio", mock.Anything, 2, 1).Return(errors.New("SERVER ERROR"))

	type args struct {
		videoId    int
		audioId    int
		statusCode int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "SetAudio 200",
			args: args{
				videoId:    1,
				audioId:    1,
				statusCode: 200,
			},
		},
		{
			name: "SetAudio 500",
			args: args{
				videoId:    2,
				audioId:    1,
				statusCode: 500,
			},
		},
	}
	for _, tt := range tests {
		rw := httptest.NewRecorder()
		r := httptest.NewRequest(
			"GET",
			"/video/{id}/audio",
			strings.NewReader(fmt.Sprintf("{\"audio_id\": %d}", tt.args.audioId)),
		)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", strconv.Itoa(tt.args.videoId))
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		t.Run(tt.name, func(t *testing.T) {
			h := &VideoHandler{
				video: videoProvider,
				audio: audioProvider,
			}
			h.SetAudio(rw, r)
			res := rw.Result()
			assert.Equal(t, tt.args.statusCode, res.StatusCode)
		})
	}
}
