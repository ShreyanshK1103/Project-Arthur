package handlers

import (
	"context"
	"time"

	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/auth"
	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/database"
	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/models"
)

func (cfg *Config) CreateLoginSession (ctx context.Context, user database.User) (models.LoginResponse, error) {
	accessToken , err := auth.GenerateAccessToken(user.ID)
	if err != nil {
		return models.LoginResponse{}, err
	}

	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return models.LoginResponse{}, err
	}

	hashedRefreshToken := auth.HashRefreshToken(refreshToken)

	_, err = cfg.DB.CreateRefreshToken(
		ctx,
		database.CreateRefreshTokenParams{
			UserID: user.ID,
			TokenHash: hashedRefreshToken,
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		},
	)

	if err != nil {
		return models.LoginResponse{}, err
	}

	return models.LoginResponse{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
		TokenType: "Bearer",
		ExpiresIn: 900,
	}, nil
}