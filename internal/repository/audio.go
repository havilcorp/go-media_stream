package repository

import (
	"context"
	"database/sql"

	"go-media-stream/internal/domain"

	"github.com/sirupsen/logrus"
)

type AudioRepository struct {
	db *sql.DB
}

func NewAudioRepository(db *sql.DB) *AudioRepository {
	return &AudioRepository{
		db: db,
	}
}

func (repo *AudioRepository) GetAudioById(ctx context.Context, id int) (*domain.Audio, error) {
	audioModel := domain.Audio{}
	sql := `SELECT id, name, idx, video_id FROM audio WHERE id = ?;`
	rows := repo.db.QueryRowContext(ctx, sql, id)
	err := rows.Scan(
		&audioModel.Id,
		&audioModel.Name,
		&audioModel.Idx,
		&audioModel.VideoId,
	)
	if err != nil {
		return &audioModel, err
	}
	return &audioModel, nil
}

func (repo *AudioRepository) GetAudioByVideoId(ctx context.Context, id int) (*[]domain.Audio, error) {
	sql := `SELECT id, name, idx FROM audio WHERE video_id = ?;`
	rows, err := repo.db.QueryContext(ctx, sql, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	audioModels := make([]domain.Audio, 0, len(cols))
	for rows.Next() {
		audioModel := domain.Audio{}
		err := rows.Scan(
			&audioModel.Id,
			&audioModel.Name,
			&audioModel.Idx,
		)
		if err != nil {
			logrus.Error(err)
		}
		audioModels = append(audioModels, audioModel)
	}

	return &audioModels, nil
}
