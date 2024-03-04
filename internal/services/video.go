package services

import (
	"context"

	"go-media-stream/internal/domain"
)

type VideoProcessor interface {
	GetVideoById(ctx context.Context, id int) (*domain.Video, error)
	GetVideos(ctx context.Context) (*[]domain.Video, error)
	SetTime(ctx context.Context, videoID int, time float32) error
	SetAudio(ctx context.Context, videoID int, audioID int) error
}

type VideoService struct {
	repository VideoProcessor
}

func NewVideoServices(repository VideoProcessor) *VideoService {
	return &VideoService{
		repository: repository,
	}
}

func (v *VideoService) GetVideoById(ctx context.Context, id int) (*domain.Video, error) {
	return v.repository.GetVideoById(ctx, id)
}

func (v *VideoService) GetVideos(ctx context.Context) (*[]domain.Video, error) {
	return v.repository.GetVideos(ctx)
}

func (v *VideoService) SetTime(ctx context.Context, videoID int, time float32) error {
	return v.repository.SetTime(ctx, videoID, time)
}

func (v *VideoService) SetAudio(ctx context.Context, videoID int, audioID int) error {
	return v.repository.SetAudio(ctx, videoID, audioID)
}
