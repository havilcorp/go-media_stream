package mysql

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
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
		CREATE TABLE IF NOT EXISTS films (
			id INT PRIMARY KEY AUTO_INCREMENT,
			name varchar(300) UNIQUE NOT NULL, 
			tag varchar(100) NOT NULL,
			path varchar(300) UNIQUE NOT NULL,
			select_audio_id INT,
			time float NOT NULL,
			updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return err
	}
	tx.Exec(`
		CREATE TABLE IF NOT EXISTS audio (
			id INT PRIMARY KEY AUTO_INCREMENT,
			name varchar(300) UNIQUE NOT NULL,
			path varchar(300) UNIQUE NOT NULL,
			film_id INT
		);
	`)
	if err != nil {
		return err
	}
	tx.Exec(`
		ALTER TABLE audio
		ADD FOREIGN KEY (film_id) REFERENCES films(id);
	`)
	if err != nil {
		return err
	}
	tx.Exec(`
		ALTER TABLE films
		ADD FOREIGN KEY (select_audio_id) REFERENCES audio(id);
	`)
	if err != nil {
		return err
	}
	err = tx.Commit()
	return err
}

type FilmModel struct {
	Id              int           `db:"id" json:"id"`
	Name            string        `db:"name" json:"name"`
	Tag             string        `db:"tag" json:"tag"`
	Path            string        `db:"path" json:"path"`
	Select_audio_id sql.NullInt32 `db:"select_audio_id" json:"select_audio_id"`
	Time            float32       `db:"time" json:"time"`
	Updated         time.Time     `db:"updated" json:"updated"`
}

func (m *Mysql) GetFilms(ctx context.Context) ([]FilmModel, error) {
	films := make([]FilmModel, 0)
	rows, err := m.db.Query(`SELECT id, name, tag, path, select_audio_id, time, updated FROM films;`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		filmModel := FilmModel{}
		rows.Scan(
			&filmModel.Id,
			&filmModel.Name,
			&filmModel.Tag,
			&filmModel.Path,
			&filmModel.Select_audio_id,
			&filmModel.Time,
			&filmModel.Updated,
		)
		films = append(films, filmModel)
	}
	return films, nil
}

func (m *Mysql) SetFilmTime(ctx context.Context, name string, audio string, time float32) {
	_, err := m.db.QueryContext(ctx, `
		INSERT INTO film_time (id, name, audio, time, updated)
		VALUES (NULL, ?, ?, ?, CURRENT_TIMESTAMP)
		ON DUPLICATE KEY UPDATE audio = ?, time = ?, updated = CURRENT_TIMESTAMP;
	`, name, audio, time, audio, time)
	if err != nil {
		panic(err)
	}
}
