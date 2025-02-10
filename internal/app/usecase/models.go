package usecase

import (
	"github.com/google/uuid"
	"time"
)

type Tokens struct {
	Access  string
	Refresh string
}

type SignUpData struct {
	Name          string
	Email         string
	Password      string
	ProfilePicUrl string
}

type User struct {
	ID            uuid.UUID
	Name          string
	Email         string
	ProfilePicUrl string
	CreatedAt     time.Time
	PasswordHash  string
}
