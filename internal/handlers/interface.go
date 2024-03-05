package handlers

import (
	"context"

	"go-media-stream/internal/domain"
)

//go:generate mockery --name AuthProvider
type AuthProvider interface {
	SignUp(ctx context.Context, login string, password string) (string, error)
}

//go:generate mockery --name AudioProvider
type AudioProvider interface {
	GetAudioByVideoId(ctx context.Context, id int) (*[]domain.Audio, error)
	GetAudioById(ctx context.Context, id int) (*domain.Audio, error)
}

//go:generate mockery --name VideoProvider
type VideoProvider interface {
	GetVideoById(ctx context.Context, id int) (*domain.Video, error)
	GetVideos(ctx context.Context) (*[]domain.Video, error)
	SetTime(ctx context.Context, videoID int, time float32) error
	SetAudio(ctx context.Context, videoID int, audioID int) error
}
