package mysql

import (
	"context"
	"database/sql"

	"go-media-stream/internal/common/ffmpeg"
	"go-media-stream/internal/entity"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

type Mysql struct {
	db *sql.DB
}

func NewMysql(dbConnect string) (*Mysql, error) {
	db, err := sql.Open("mysql", dbConnect)
	if err != nil {
		return nil, err
	}
	mysql := Mysql{
		db: db,
	}
	err = mysql.Bootstrap()
	if err != nil {
		return nil, err
	}
	return &mysql, nil
}

func (m *Mysql) Bootstrap() error {
	tx, err := m.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS video (
			id INT PRIMARY KEY AUTO_INCREMENT,
			name varchar(300) UNIQUE NOT NULL, 
			tag varchar(100),
			select_audio_id INT,
			time SMALLINT NOT NULL DEFAULT 0,
			updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return err
	}
	tx.Exec(`
		CREATE TABLE IF NOT EXISTS audio (
			id INT PRIMARY KEY AUTO_INCREMENT,
			name varchar(300) NOT NULL,
			idx INT NOT NULL,
			video_id INT
		);
	`)
	if err != nil {
		return err
	}
	tx.Exec(`
		ALTER TABLE audio
		ADD FOREIGN KEY (video_id) REFERENCES video(id);
	`)
	if err != nil {
		return err
	}
	tx.Exec(`
		ALTER TABLE video
		ADD FOREIGN KEY (select_audio_id) REFERENCES audio(id);
	`)
	if err != nil {
		return err
	}
	err = tx.Commit()
	return err
}

func (m *Mysql) GetVideos(ctx context.Context) ([]entity.VideoModel, error) {
	videos := make([]entity.VideoModel, 0)
	sql := `SELECT id, name, tag, select_audio_id, time, updated FROM video;`
	rows, err := m.db.QueryContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		videoModel := entity.VideoModel{}
		rows.Scan(
			&videoModel.Id,
			&videoModel.Name,
			&videoModel.Tag,
			&videoModel.Select_audio_id,
			&videoModel.Time,
			&videoModel.Updated,
		)
		videos = append(videos, videoModel)
	}
	return videos, nil
}

func (m *Mysql) GetVideo(ctx context.Context, id int) (entity.VideoModel, error) {
	videoModel := entity.VideoModel{}
	sql := `SELECT id, name, select_audio_id, time FROM video WHERE id = ?;`
	rows := m.db.QueryRowContext(ctx, sql, id)
	err := rows.Scan(
		&videoModel.Id,
		&videoModel.Name,
		&videoModel.Select_audio_id,
		&videoModel.Time,
	)
	if err != nil {
		return videoModel, err
	}
	return videoModel, nil
}

func (m *Mysql) GetAudio(ctx context.Context, id int) (entity.AudioModel, error) {
	audioModel := entity.AudioModel{}
	sql := `SELECT id, name, idx, video_id FROM audio WHERE id = ?;`
	rows := m.db.QueryRowContext(ctx, sql, id)
	err := rows.Scan(
		&audioModel.Id,
		&audioModel.Name,
		&audioModel.Idx,
		&audioModel.VideoId,
	)
	if err != nil {
		return audioModel, err
	}
	return audioModel, nil
}

func (m *Mysql) GetAudioByVideoId(ctx context.Context, id int) ([]entity.AudioModel, error) {
	sql := `SELECT id, name, idx FROM audio WHERE video_id = ?;`
	rows, err := m.db.QueryContext(ctx, sql, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	audioModels := make([]entity.AudioModel, 0, len(cols))
	for rows.Next() {
		audioModel := entity.AudioModel{}
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

	return audioModels, nil
}

func (m *Mysql) SetVideoAudio(ctx context.Context, videoId int, audioId int) error {
	_, err := m.db.QueryContext(ctx, `
		UPDATE video
		SET select_audio_id = ?
		WHERE id = ?;
	`, audioId, videoId)
	if err != nil {
		return err
	}
	return nil
}

func (m *Mysql) SetVideoTime(ctx context.Context, videoID int, time float32) error {
	_, err := m.db.QueryContext(ctx, `
		UPDATE video
		SET time = ?, updated = CURRENT_TIMESTAMP
		WHERE id = ?;
	`, time, videoID)
	if err != nil {
		return err
	}
	return nil
}

func (m *Mysql) IsValidVideoName(ctx context.Context, name string) (bool, error) {
	sql := "SELECT COUNT(*) FROM video WHERE name = ?"
	row := m.db.QueryRowContext(ctx, sql, name)
	var count int
	if err := row.Scan(&count); err != nil {
		return false, err
	}
	return count == 0, nil
}

func (m *Mysql) AddFilm(ctx context.Context, workerRes *ffmpeg.WorkerResult) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	resSql, err := tx.ExecContext(ctx, `
		INSERT INTO video (id, name) 
		VALUES (NULL, ?);
	`, workerRes.Name)
	if err != nil {
		return err
	}
	id, err := resSql.LastInsertId()
	if err != nil {
		return err
	}
	for _, audio := range workerRes.Audios {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO audio (id, name, idx, video_id) 
			VALUES (NULL, ?, ?, ?);
		`, audio.Title, audio.Index, id)
		if err != nil {
			logrus.Error(err)
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
