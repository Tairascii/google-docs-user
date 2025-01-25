package repo

import (
	"github.com/google/uuid"
	"time"
)

type CreateUserData struct {
	Name          string `db:"name"`
	Email         string `db:"email"`
	Password      string `db:"password_hash"`
	ProfilePicUrl string `db:"profile_picture_url"`
}

type User struct {
	ID            uuid.UUID `db:"id"`
	Name          string    `db:"name"`
	Email         string    `db:"email"`
	ProfilePicUrl string    `db:"profile_picture_url"`
	CreatedAt     time.Time `db:"created_at"`
	PasswordHash  string    `db:"password_hash"`
}
