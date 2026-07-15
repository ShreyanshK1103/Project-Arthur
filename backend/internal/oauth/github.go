package oauth

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var GithubConfig *oauth2.Config

func InitGithubOAuth(
	clientID string,
	clientSecret string,
	callbackURL string,
) {
	if clientID == "" || clientSecret == "" || callbackURL == "" {
		panic("GitHub OAuth environment variables are missing")
	}
	GithubConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  callbackURL,
		Scopes: []string{
			"read:user",
			"user:email",
		},
		Endpoint: github.Endpoint,
	}
}
