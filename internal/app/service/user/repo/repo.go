package repo

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	ErrOnQuery           = errors.New("query error")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
)

type UserRepo interface {
	CreateUser(ctx context.Context, userData CreateUserData) (uuid.UUID, error)
	GetUser(ctx context.Context, email string) (User, error)
}

type Repo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) CreateUser(ctx context.Context, userData CreateUserData) (uuid.UUID, error) {
	query := `
	INSERT INTO users (name, email, password_hash, profile_picture_url)
	VALUES ($1, $2, $3, $4) ON CONFLICT (email) DO NOTHING RETURNING id;
	`

	res := r.db.QueryRowxContext(ctx, query, userData.Name, userData.Email, userData.Password, userData.ProfilePicUrl)

	var id uuid.UUID
	err := res.Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, ErrUserAlreadyExists
		}
		return uuid.Nil, errors.Join(ErrOnQuery, err)
	}

	if id == uuid.Nil {
		return uuid.Nil, ErrUserAlreadyExists
	}
	return id, nil
}

func (r *Repo) GetUser(ctx context.Context, email string) (User, error) {
	var user User
	q := `SELECT id, name, email, profile_picture_url, created_at, password_hash FROM users WHERE email = $1;`
	err := r.db.GetContext(ctx, &user, q, email)
	if err != nil {
		return User{}, errors.Join(ErrUserNotFound, err)
	}
	return user, nil
}
