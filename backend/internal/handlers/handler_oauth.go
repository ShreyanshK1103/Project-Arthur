package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/database"
	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/oauth"
	githubapi "github.com/google/go-github/v69/github"
)

func (cfg *Config) HandlerGithubLogin(w http.ResponseWriter, r *http.Request) {
	url := oauth.GithubConfig.AuthCodeURL("state")

	http.Redirect(
		w,
		r,
		url,
		http.StatusTemporaryRedirect,
	)
}

func (cfg *Config) HandlerGithubCallback(w http.ResponseWriter, r *http.Request) {
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
			Provider: provider,
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
						UserID: existingUser.ID,
						Provider: provider,
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
						Name: name,
						Email: email,

						PasswordHash: sql.NullString{},

						AvatarUrl: sql.NullString{
							String: githubUser.GetAvatarURL(),
							Valid: true,
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
						UserID: user.ID,
						Provider: provider,
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

	http.Redirect(
		w,
		r,
		fmt.Sprintf(
			"http://localhost:5173/auth/callback?access=%s&refresh=%s",
			response.AccessToken,
			response.RefreshToken,
		),
		http.StatusTemporaryRedirect,
	)
}
