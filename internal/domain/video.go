package domain

import (
	"database/sql"
	"time"
)

type Video struct {
	Id            int           `db:"id" json:"id"`
	Name          string        `db:"name" json:"name"`
	Tag           string        `db:"tag" json:"tag"`
	SelectAudioId sql.NullInt32 `db:"select_audio_id" json:"select_audio_id"`
	Time          float32       `db:"time" json:"time"`
	Updated       time.Time     `db:"updated" json:"updated"`
}
