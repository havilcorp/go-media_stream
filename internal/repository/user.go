package repository

import (
	"context"
	"database/sql"

	"go-media-stream/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (repo *UserRepository) GetUserById(ctx context.Context, id int64) (*domain.User, error) {
	row := repo.db.QueryRowContext(ctx, "SELECT id, login FROM users WHERE id = ?", id)
	user := domain.User{}
	if err := row.Scan(&user.ID, &user.Login); err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) UserIsExists(ctx context.Context, login string) (bool, error) {
	rows, err := repo.db.QueryContext(ctx, `Select * from users WHERE login = ?`, login)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return rows.Next(), nil
}

func (repo *UserRepository) CreateUser(ctx context.Context, user *domain.User) (int64, error) {
	res, err := repo.db.ExecContext(ctx, `
		INSERT INTO users (id, login, password) 
		VALUES (NULL, ?, ?);
	`, user.Login, user.Password)
	if err != nil {
		return 0, err
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastId, nil
}

func (repo *UserRepository) GetUserByLogin(ctx context.Context, login string, password string) (*domain.User, error) {
	row := repo.db.QueryRowContext(ctx, "SELECT id, login, password FROM users WHERE login = ?", login)
	user := domain.User{}
	if err := row.Scan(&user.ID, &user.Login, &user.Password); err != nil {
		return nil, domain.ErrUserNotFound
	}
	if user.Password != password {
		return nil, domain.ErrWrongPassword
	}
	return &user, nil
}
