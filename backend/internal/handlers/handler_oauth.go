package handlers

import (
	"context"
	"log"
	"net/http"

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
	if email == "" {
		for _, e := range emails {
			if e.GetPrimary() && e.GetVerified() {
				email = e.GetEmail()
				break
			}
		}
	}

	log.Println("GitHub ID:", githubUser.GetID())
	log.Println("Username:", githubUser.GetLogin())
	log.Println("Name:", githubUser.GetName())
	log.Println("Email:", email)
}
