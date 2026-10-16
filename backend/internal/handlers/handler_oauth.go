package handlers

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"

	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/config"
	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/database"
	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/oauth"
	githubapi "github.com/google/go-github/v69/github"
)

func (cfg *Config) HandlerGithubLogin(w http.ResponseWriter, r *http.Request) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		respondWithError(w, 500, "could not start login")
		return
	}
	state := hex.EncodeToString(b)

	http.SetCookie(w, &http.Cookie{
		Name: "oauth_state", Value: state, Path: "/v1/auth/github",
		MaxAge: 600, HttpOnly: true, Secure: false, SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, oauth.GithubConfig.AuthCodeURL(state), http.StatusTemporaryRedirect)
}

func (cfg *Config) HandlerGithubCallback(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("oauth_state")
	if err != nil || c.Value == "" || c.Value != r.URL.Query().Get("state") {
		respondWithError(w, http.StatusBadRequest, "invalid oauth state")
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "oauth_state", Value: "", Path: "/v1/auth/github", MaxAge: -1})

	code := r.URL.Query().Get("code")

	if code == "" {
		respondWithError(
			w,
			http.StatusBadRequest,
			"missing authorization code",
		)
		return
	}

	token, err := oauth.GithubConfig.Exchange(
		r.Context(),
		code,
	)

	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"failed to exchange code",
		)
		return
	}

	client := githubapi.NewClient(
		oauth.GithubConfig.Client(
			context.Background(),
			token,
		),
	)

	githubUser, _, err := client.Users.Get(
		context.Background(),
		"",
	)

	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"failed to fetch github user",
		)
		return
	}
	emails, _, err := client.Users.ListEmails(
		context.Background(),
		nil,
	)

	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"failed to fetch github emails",
		)
		return
	}
	var email string

	for _, e := range emails {
		if e.GetPrimary() && e.GetVerified() {
			email = e.GetEmail()
			break
		}
	}

	if email == "" {
		for _, e := range emails {
			if e.GetVerified() {
				email = e.GetEmail()
				break
			}
		}
	}

	if email == "" {
		respondWithError(
			w,
			http.StatusBadRequest,
			"No verified GitHub email found",
		)
		return
	}

	provider := "github"

	providerID := fmt.Sprintf(
		"%d",
		githubUser.GetID(),
	)

	user, err := cfg.DB.GetUserByOAuth(
		r.Context(),
		database.GetUserByOAuthParams{
			Provider:   provider,
			ProviderID: providerID,
		},
	)

	if err != nil {
		// NO GITHUB ACCOUNT LINKED YET
		if errors.Is(err, sql.ErrNoRows) {

			existingUser, emailErr := cfg.DB.GetUserByEmail(
				r.Context(),
				email,
			)

			if emailErr == nil {

				// Existing Arthur account
				// Just link GitHub

				_, err = cfg.DB.CreateOAuthAccount(
					r.Context(),
					database.CreateOAuthAccountParams{
						UserID:     existingUser.ID,
						Provider:   provider,
						ProviderID: providerID,
					},
				)

				if err != nil {
					respondWithError(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)
					return
				}

				user = existingUser

			} else if errors.Is(emailErr, sql.ErrNoRows) {

				// First time user

				name := githubUser.GetName()
				if name == "" {
					name = githubUser.GetLogin()
				}

				user, err = cfg.DB.CreateUser(
					r.Context(),
					database.CreateUserParams{
						Name:  name,
						Email: email,

						PasswordHash: sql.NullString{},

						AvatarUrl: sql.NullString{
							String: githubUser.GetAvatarURL(),
							Valid:  true,
						},

						EmailVerified: true,
					},
				)

				if err != nil {
					respondWithError(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)
					return
				}

				_, err = cfg.DB.CreateOAuthAccount(
					r.Context(),
					database.CreateOAuthAccountParams{
						UserID:     user.ID,
						Provider:   provider,
						ProviderID: providerID,
					},
				)

				if err != nil {
					respondWithError(
						w,
						http.StatusInternalServerError,
						err.Error(),
					)
					return
				}

			} else {

				respondWithError(
					w,
					http.StatusInternalServerError,
					emailErr.Error(),
				)
				return

			}

		} else {

			respondWithError(
				w,
				http.StatusInternalServerError,
				err.Error(),
			)
			return

		}
	}

	response, err := cfg.CreateLoginSession(
		r.Context(),
		user,
	)

	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Couldn't create session",
		)
		return
	}

	SetAuthCookies(
		w,
		response.AccessToken,
		response.RefreshToken,
	)

	http.Redirect(w, r, config.FrontendURL()+"/", http.StatusSeeOther)
}
