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
	"gopa/internal/core/ports"
)

const userColumns = `id, email, password_hash, display_name, role, avatar_url, created_at, updated_at`

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	const query = `
		INSERT INTO users (id, email, password_hash, display_name, role, avatar_url, created_at, updated_at)
		VALUES (:id, :email, :password_hash, :display_name, :role, :avatar_url, :created_at, :updated_at)`
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
	return r.findOne(ctx, `SELECT `+userColumns+` FROM users WHERE id = $1`, id)
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	return r.findOne(ctx, `SELECT `+userColumns+` FROM users WHERE email = $1`, email)
}

func (r *UserRepository) Update(ctx context.Context, user domain.User) (domain.User, error) {
	result, err := r.db.NamedExecContext(ctx,
		`UPDATE users SET display_name = :display_name, avatar_url = :avatar_url, updated_at = :updated_at WHERE id = :id`,
		user)
	if err != nil {
		return domain.User{}, fmt.Errorf("update user: %w", err)
	}
	if err := requireAffected(result, "update user"); err != nil {
		return domain.User{}, err
	}
	return user, nil
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

var _ ports.UserRepository = (*UserRepository)(nil)
