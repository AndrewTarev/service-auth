package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"

	"service-auth/internal/app/errs"
	"service-auth/internal/app/models"
)

const (
	DuplicateValue = "23505"
)

func (r *Repository) CreateUser(ctx context.Context, user models.User) (int, error) {
	query := "INSERT INTO users(username, password_hash, email) VALUES ($1, $2, $3) RETURNING id"
	var id int
	err := r.db.QueryRow(ctx, query, user.Username, user.Password, user.Email).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// Обработка ошибок PostgreSQL
			if pgErr.Code == DuplicateValue {
				if pgErr.ConstraintName == "users_username_key" {
					return 0, errs.ErrUserAlreadyExists
				}
				if pgErr.ConstraintName == "users_email_key" {
					return 0, errs.ErrEmailAlreadyUsed
				}
			}
			return 0, fmt.Errorf("database error: %v", pgErr.Message)
		}
		return 0, fmt.Errorf("query error: %v", err)
	}
	return id, nil
}

func (r *Repository) GetUser(ctx context.Context, username string) (models.User, error) {
	var user models.User
	query := "SELECT id, username, password_hash, email FROM users WHERE username = $1"
	row := r.db.QueryRow(ctx, query, username)
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return user, errs.ErrUserNotFound
		}
		return user, err
	}
	return user, nil
}
