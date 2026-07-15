package handlers

import (
	"net/http"

	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/oauth"
)

func (cfg *Config) HandlerGithubLogin(
	w http.ResponseWriter,
	r *http.Request,
) {
	url := oauth.GithubConfig.AuthCodeURL("state")

	http.Redirect(
		w,
		r,
		url,
		http.StatusTemporaryRedirect,
	)
}
