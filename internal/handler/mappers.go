package handler

import "github.com/Tairascii/google-docs-user/internal/usecase"

func toTokens(raw usecase.Tokens) Tokens {
	return Tokens{
		AccessToken:  raw.Access,
		RefreshToken: raw.Refresh,
	}
}
