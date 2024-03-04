package repository

import (
	"context"
	"database/sql"
)

type QueueRepository struct {
	db *sql.DB
}

func NewQueueRepository(db *sql.DB) *QueueRepository {
	return &QueueRepository{
		db: db,
	}
}

func (r *QueueRepository) AddVideoToQueue(ctx context.Context, userId int, folder string, t string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO queue (id, user_id, folder, type) 
		VALUES (NULL, ?, ?, ?);
	`, userId, folder, t)
	if err != nil {
		return err
	}
	return nil
}
