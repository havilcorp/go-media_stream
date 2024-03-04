package handlers

import (
	"context"

	"go-media-stream/internal/domain"
)

type AudioProvider interface {
	GetAudioByVideoId(ctx context.Context, id int) (*[]domain.Audio, error)
	GetAudioById(ctx context.Context, id int) (*domain.Audio, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name=VideoProvider
type VideoProvider interface {
	GetVideoById(ctx context.Context, id int) (*domain.Video, error)
	GetVideos(ctx context.Context) (*[]domain.Video, error)
	SetTime(ctx context.Context, videoID int, time float32) error
	SetAudio(ctx context.Context, videoID int, audioID int) error
}
