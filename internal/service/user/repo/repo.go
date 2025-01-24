package repo

import (
	"errors"
	"github.com/jmoiron/sqlx"
)

var (
	ErrOnQuery           = errors.New("query error")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
)

type UserRepo interface {
	CreateUser(userData CreateUserData) error
	GetUser(name, password string) (User, error)
}

type Repo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) CreateUser(userData CreateUserData) error {
	query := `
	INSERT INTO users (name, email, password_hash, profile_picture_url)
	VALUES ($1, $2, $3, $4) ON CONFLICT (msisdn, subscriber_id) DO NOTHING;
	`

	res, err := r.db.Exec(query, userData.Name, userData.Email, userData.Password, userData.ProfilePicUrl)
	if err != nil {
		return errors.Join(ErrOnQuery, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.Join(ErrOnQuery, err)
	}

	if rowsAffected == 0 {
		return ErrUserAlreadyExists
	}
	return nil
}

func (r *Repo) GetUser(name, password string) (User, error) {
	var user User
	q := `SELECT id, name, email, profile_picture_url, created_at FROM users WHERE email = $1 AND password_hash = $2;`
	err := r.db.Get(&user, q, name, password)
	if err != nil {
		return User{}, errors.Join(ErrUserNotFound, err)
	}
	return user, nil
}
