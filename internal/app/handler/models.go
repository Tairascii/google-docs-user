package handler

import (
	"github.com/google/uuid"
	"time"
)

type SignInPayload struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type SignUpPayload struct {
	Email         string `json:"email" binding:"required,email"`
	Password      string `json:"password" binding:"required"`
	Name          string `json:"name," binding:"required"`
	ProfilePicUrl string `json:"profilePicUrl,omitempty"`
}

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type User struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	ProfilePicUrl string    `json:"profilePicUrl"`
	CreatedAt     time.Time `json:"createdAt"`
	PasswordHash  string    `json:"-"`
}
