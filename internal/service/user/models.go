package user

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID            uuid.UUID
	Name          string
	Email         string
	ProfilePicUrl string
	CreatedAt     time.Time
	PasswordHash  string
}

type CreateUserData struct {
	Name          string
	Email         string
	Password      string
	ProfilePicUrl string
}
