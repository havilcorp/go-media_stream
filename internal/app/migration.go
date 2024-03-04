package app

import (
	"context"
	"database/sql"
)

func Bootstrap(ctx context.Context, db *sql.DB) error {
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INT PRIMARY KEY AUTO_INCREMENT,
			login varchar(300) NOT NULL UNIQUE,
			password varchar(300) NOT NULL
		);
	`)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS video (
			id INT PRIMARY KEY AUTO_INCREMENT,
			name varchar(300) UNIQUE NOT NULL, 
			tag varchar(100),
			select_audio_id INT,
			user_id INT NOT NULL,
			time SMALLINT NOT NULL DEFAULT 0,
			updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
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
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS queue (
			id INT PRIMARY KEY AUTO_INCREMENT,
			user_id INT NOT NULL,
			video_id INT,
			folder varchar(300) NOT NULL,
			title varchar(300) NOT NULL DEFAULT "",
			idx INT,
			type varchar(100) NOT NULL,
			is_work boolean NOT NULL DEFAULT 0,
			is_done boolean NOT NULL DEFAULT 0
		);
	`)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
		ALTER TABLE audio
		ADD FOREIGN KEY (video_id) REFERENCES video(id);
	`)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
		ALTER TABLE video
		ADD FOREIGN KEY (select_audio_id) REFERENCES audio(id);
	`)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
		ALTER TABLE video
		ADD FOREIGN KEY (user_id) REFERENCES users(id);
	`)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
		ALTER TABLE queue
		ADD FOREIGN KEY (user_id) REFERENCES users(id);
	`)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
		ALTER TABLE queue
		ADD FOREIGN KEY (video_id) REFERENCES video(id);
	`)
	if err != nil {
		return err
	}
	err = tx.Commit()
	return err
}
