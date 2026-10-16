package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

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
	response, err := cfg.CreateLoginSession(
		r.Context(),
		user,
	)

	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Couldn't create a login session",
		)
		return
	}

	SetAuthCookies(
		w,
		response.AccessToken,
		response.RefreshToken,
	)

	respondWithJSON(
		w,
		http.StatusOK,
		map[string]string{
			"message": "Logged in successfully",
		},
	)
}

func (cfg *Config) HandlerMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r.Context())

	if !ok {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"unauthorized",
		)
		return
	}

	user, err := cfg.DB.GetUserByID(
		r.Context(),
		userID,
	)

	if err != nil {
		respondWithError(
			w,
			http.StatusNotFound,
			"user not found",
		)
		return
	}

	respondWithJSON(
		w,
		http.StatusOK,
		models.UserToResponse(user),
	)
}

func (cfg *Config) HandlerRefresh(w http.ResponseWriter, r *http.Request) {

	// Get refresh token from HttpOnly cookie
	cookie, err := r.Cookie(RefreshTokenCookie)

	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"refresh token missing",
		)
		return
	}

	rawRefreshToken := cookie.Value

	if rawRefreshToken == "" {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"refresh token missing",
		)
		return
	}

	// Hash the raw token so we can look it up in the DB
	hashedToken := auth.HashRefreshToken(rawRefreshToken)

	// Find refresh token in DB
	refreshToken, err := cfg.DB.GetRefreshTokenByHash(
		r.Context(),
		hashedToken,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ClearAuthCookies(w)
			respondWithError(
				w,
				http.StatusUnauthorized,
				"invalid refresh token",
			)
			return
		}

		respondWithError(
			w,
			http.StatusInternalServerError,
			"database error",
		)
		return
	}

	// Check expiration
	if refreshToken.ExpiresAt.Before(time.Now()) {

		// Remove expired token
		_ = cfg.DB.DeleteRefreshToken(
			r.Context(),
			hashedToken,
		)

		ClearAuthCookies(w)

		respondWithError(
			w,
			http.StatusUnauthorized,
			"refresh token expired",
		)
		return
	}

	// Delete old refresh token
	err = cfg.DB.DeleteRefreshToken(
		r.Context(),
		hashedToken,
	)

	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"couldn't rotate refresh token",
		)
		return
	}

	// Generate new refresh token
	newRefreshToken, err := auth.GenerateRefreshToken()

	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"couldn't generate refresh token",
		)
		return
	}

	// Hash new refresh token
	newRefreshTokenHash := auth.HashRefreshToken(
		newRefreshToken,
	)

	// Store new refresh token
	_, err = cfg.DB.CreateRefreshToken(
		r.Context(),
		database.CreateRefreshTokenParams{
			UserID:    refreshToken.UserID,
			TokenHash: newRefreshTokenHash,
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		},
	)

	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"couldn't store refresh token",
		)
		return
	}

	// Generate new access token
	accessToken, err := auth.GenerateAccessToken(
		refreshToken.UserID,
	)

	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"couldn't generate access token",
		)
		return
	}

	// Replace cookies
	SetAuthCookies(
		w,
		accessToken,
		newRefreshToken,
	)

	respondWithJSON(
		w,
		http.StatusOK,
		map[string]string{
			"message": "token refreshed successfully",
		},
	)
}

func (cfg *Config) HandlerLogout(w http.ResponseWriter, r *http.Request) {

	// Remove cookies from browser
	ClearAuthCookies(w)

	cookie, err := r.Cookie(RefreshTokenCookie)

	if err == nil && cookie.Value != "" {

		hash := auth.HashRefreshToken(
			cookie.Value,
		)

		err = cfg.DB.DeleteRefreshToken(
			r.Context(),
			hash,
		)

		if err != nil {
			log.Printf("logout: failed to delete refresh token: %v", err)
		}
	}

	respondWithJSON(
		w,
		http.StatusOK,
		map[string]string{
			"message": "Logged out successfully",
		},
	)
}
