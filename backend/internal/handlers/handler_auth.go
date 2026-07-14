package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

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
	// Trimming the whitespaces and making the emails lowercase for the consistency
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)
	req.Email = strings.ToLower(req.Email)

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

func (cfg *Config) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			"invalid json",
		)
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	req.Email = strings.ToLower(req.Email)

	// Validate
	if req.Email == "" || req.Password == "" {
		respondWithError(
			w,
			http.StatusBadRequest,
			"email and password are required",
		)
		return
	}

	// Find User
	user, err := cfg.DB.GetUserByEmail(
		r.Context(),
		req.Email,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(
				w,
				http.StatusUnauthorized,
				"invalid email or password",
			)
		}
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(
				w,
				http.StatusUnauthorized,
				"Invalid email or password",
			)
			return
		}

		respondWithError(
			w,
			http.StatusInternalServerError,
			"Database error",
		)
		return
	}

	err = auth.CheckPassword(req.Password, user.PasswordHash.String)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"Invalid email or password",
		)
		return
	}

	//Success
	token, err := auth.GenerateAccessToken(user.ID)

	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"failed to generate access token",
		)
		return
	}

	respondWithJSON(
		w,
		http.StatusOK,
		models.LoginResponse{
			AccessToken: token,
			TokenType:   "Bearer",
			ExpiresIn:   900,
		},
	)
}
