package services

import (
	"context"

	"go-media-stream/internal/domain"
)

type AudioProcessor interface {
	GetAudioByVideoId(ctx context.Context, id int) (*[]domain.Audio, error)
	GetAudioById(ctx context.Context, id int) (*domain.Audio, error)
}

type AudioService struct {
	repository AudioProcessor
}

func NewAudioServices(repository AudioProcessor) *AudioService {
	return &AudioService{
		repository: repository,
	}
}

func (a *AudioService) GetAudioByVideoId(ctx context.Context, id int) (*[]domain.Audio, error) {
	return a.repository.GetAudioByVideoId(ctx, id)
}

func (a *AudioService) GetAudioById(ctx context.Context, id int) (*domain.Audio, error) {
	return a.repository.GetAudioById(ctx, id)
}
