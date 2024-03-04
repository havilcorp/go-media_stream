package repository

import (
	"context"
	"database/sql"

	"go-media-stream/internal/domain"
)

type Videorsitory struct {
	db *sql.DB
}

func NewVideoRepository(db *sql.DB) *Videorsitory {
	return &Videorsitory{
		db: db,
	}
}

func (r *Videorsitory) GetVideos(ctx context.Context) (*[]domain.Video, error) {
	videos := make([]domain.Video, 0)
	sql := `SELECT id, name, tag, select_audio_id, time, updated FROM video;`
	rows, err := r.db.QueryContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		videoModel := domain.Video{}
		rows.Scan(
			&videoModel.Id,
			&videoModel.Name,
			&videoModel.Tag,
			&videoModel.SelectAudioId,
			&videoModel.Time,
			&videoModel.Updated,
		)
		videos = append(videos, videoModel)
	}
	return &videos, nil
}

func (r *Videorsitory) GetVideoById(ctx context.Context, id int) (*domain.Video, error) {
	videoModel := domain.Video{}
	sql := `SELECT id, name, select_audio_id, time FROM video WHERE id = ?;`
	rows := r.db.QueryRowContext(ctx, sql, id)
	err := rows.Scan(
		&videoModel.Id,
		&videoModel.Name,
		&videoModel.SelectAudioId,
		&videoModel.Time,
	)
	if err != nil {
		return nil, err
	}
	return &videoModel, nil
}

func (r *Videorsitory) SetTime(ctx context.Context, videoID int, time float32) error {
	_, err := r.db.QueryContext(ctx, `
	UPDATE video
	SET time = ?, updated = CURRENT_TIMESTAMP
	WHERE id = ?;
`, time, videoID)
	if err != nil {
		return err
	}
	return nil
}

func (r *Videorsitory) SetAudio(ctx context.Context, videoID int, audioID int) error {
	_, err := r.db.QueryContext(ctx, `
	UPDATE video
	SET select_audio_id = ?
	WHERE id = ?;
`, audioID, videoID)
	if err != nil {
		return err
	}
	return nil
}

func (r *Videorsitory) IsValidVideoName(ctx context.Context, name string) (bool, error) {
	sql := "SELECT COUNT(*) FROM video WHERE name = ?"
	row := r.db.QueryRowContext(ctx, sql, name)
	var count int
	if err := row.Scan(&count); err != nil {
		return false, err
	}
	return count == 0, nil
}
