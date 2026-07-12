package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/auth"
	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/database"
	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/models"
	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/validators"
)

func (cfg *Config) HandlerRegister(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(
			w,
			400,
			fmt.Sprintf("Invalid JSON: %v", err),
		)
		return
	}

	// Basic Validation
	
	if req.Name == "" || req.Email == "" || req.Password == "" {
		respondWithError(
			w,
			400,
			"All fields are required",
		)

		return
	}
	if err := validators.ValidateEmail(req.Email); err != nil {
		respondWithError(
			w,
			400,
			err.Error(),
		)
		return
	}

	if err := validators.ValidatePassword(req.Password); err != nil {
		respondWithError(
			w,
			400,
			err.Error(),
		)
		return
	}
	// Removing already registered email 
	_, err = cfg.DB.GetUserByEmail(
		r.Context(),
		req.Email,
	)

	if err == nil {
		respondWithError(
			w,
			http.StatusConflict,
			"email already registered",
		)
		return
	}

	if !errors.Is(err, sql.ErrNoRows) {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"failed to check existing user",
		)
		return
	}

	// Hash Password
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		respondWithError(
			w,
			500,
			"Failed to Hash Password",
		)
		return
	}

	user, err := cfg.DB.CreateUser(
		r.Context(),
		database.CreateUserParams{
			Name:  req.Name,
			Email: req.Email,

			PasswordHash: sql.NullString{
				String: passwordHash,
				Valid:  true,
			},

			Provider:      "local",
			ProviderID:    sql.NullString{},
			AvatarUrl:     sql.NullString{},
			EmailVerified: false,
		},
	)
	if err != nil {
		respondWithError(
			w,
			500,
			fmt.Sprintf(
				"Couldn't Create user: %v",
				err,
			),
		)
		return
	}

	respondWithJSON(
		w,
		201,
		models.UserToResponse(user),
	)
}
