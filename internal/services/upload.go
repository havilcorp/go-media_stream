package services

import (
	"mime/multipart"

	"go-media-stream/internal/domain"

	"golang.org/x/net/context"
)

type (
	Queuer interface {
		AddVideoToQueue(ctx context.Context, userId int, folder string, t string) error
	}
	Uploader interface {
		IsValidVideoName(ctx context.Context, name string) (bool, error)
	}
	Storager interface {
		SaveFile(file *multipart.File, fileName string) (string, error)
	}
)

type UploadService struct {
	upload Uploader
	store  Storager
	queue  Queuer
}

func NewUploadServices(upload Uploader, store Storager, queue Queuer) *UploadService {
	return &UploadService{
		upload: upload,
		store:  store,
		queue:  queue,
	}
}

func (s *UploadService) Upload(ctx context.Context, userId int, name string, videoFile *multipart.File) error {
	ok, err := s.upload.IsValidVideoName(ctx, name)
	if err != nil {
		return err
	}
	if !ok {
		return domain.ErrInvalidVideoName
	}
	path, err := s.store.SaveFile(videoFile, name)
	if err != nil {
		return err
	}
	err = s.queue.AddVideoToQueue(ctx, userId, path, "video")
	if err != nil {
		return err
	}
	return nil
}

//////////
