package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"gopa/internal/core/domain"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	const query = `
		INSERT INTO users (id, email, password_hash, role, created_at, updated_at)
		VALUES (:id, :email, :password_hash, :role, :created_at, :updated_at)`
	if _, err := r.db.NamedExecContext(ctx, query, user); err != nil {
		var postgresError *pgconn.PgError
		if errors.As(err, &postgresError) && postgresError.Code == "23505" {
			return domain.User{}, domain.ErrConflict
		}
		return domain.User{}, fmt.Errorf("insert user: %w", err)
	}
	return user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	return r.findOne(ctx, `SELECT id, email, password_hash, role, created_at, updated_at FROM users WHERE id = $1`, id)
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	return r.findOne(ctx, `SELECT id, email, password_hash, role, created_at, updated_at FROM users WHERE email = $1`, email)
}

func (r *UserRepository) findOne(ctx context.Context, query string, value any) (domain.User, error) {
	var user domain.User
	if err := r.db.GetContext(ctx, &user, query, value); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, domain.ErrNotFound
		}
		return domain.User{}, fmt.Errorf("find user: %w", err)
	}
	return user, nil
}
