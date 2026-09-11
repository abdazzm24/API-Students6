package repository

import (
	"context"
	"errors"
	"fmt"

	"api-students/app/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, user model.User) (model.User, error)
	FindByID(ctx context.Context, id int) (model.User, error)
	FindByUsername(ctx context.Context, username string) (model.User, error)
}

var (
	ErrNotFound  = errors.New("data tidak ditemukan")
	ErrDuplicate = errors.New("data sudah ada")
)

type userPostgresRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(
	pool *pgxpool.Pool,
) UserRepository {
	return &userPostgresRepository{
		pool: pool,
	}
}

func (r *userPostgresRepository) Create(
	ctx context.Context,
	user model.User,
) (model.User, error) {

	var result model.User

	err := r.pool.QueryRow(
		ctx,
		`
		INSERT INTO users (
			username,
			email,
			password,
			role,
			is_active
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING
			id,
			username,
			email,
			password,
			role,
			is_active,
			created_at
		`,
		user.Username,
		user.Email,
		user.Password,
		user.Role,
		user.IsActive,
	).Scan(
		&result.ID,
		&result.Username,
		&result.Email,
		&result.Password,
		&result.Role,
		&result.IsActive,
		&result.CreatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) &&
			pgErr.Code == "23505" {
			return model.User{}, ErrDuplicate
		}

		return model.User{}, fmt.Errorf(
			"membuat user: %w",
			err,
		)
	}

	return result, nil
}

func (r *userPostgresRepository) FindByID(
	ctx context.Context,
	id int,
) (model.User, error) {

	var user model.User

	err := r.pool.QueryRow(
		ctx,
		`
		SELECT
			id,
			username,
			email,
			password,
			role,
			is_active,
			created_at
		FROM users
		WHERE id = $1
		`,
		id,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, ErrNotFound
		}

		return model.User{}, fmt.Errorf(
			"mencari user berdasarkan id: %w",
			err,
		)
	}

	return user, nil
}

func (r *userPostgresRepository) FindByUsername(
	ctx context.Context,
	username string,
) (model.User, error) {

	var user model.User

	err := r.pool.QueryRow(
		ctx,
		`
		SELECT
			id,
			username,
			email,
			password,
			role,
			is_active,
			created_at
		FROM users
		WHERE LOWER(username) = LOWER($1)
		`,
		username,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, ErrNotFound
		}

		return model.User{}, fmt.Errorf(
			"mencari user berdasarkan username: %w",
			err,
		)
	}

	return user, nil
}