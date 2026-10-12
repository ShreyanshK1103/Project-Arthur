package handlers

import (
	"net/http"
	"os"
)

const (
	AccessTokenCookie  = "arthur_access_token"
	RefreshTokenCookie = "arthur_refresh_token"
)

var secureCookies = os.Getenv("COOKIE_SECURE") == "true"

func SetAuthCookies(w http.ResponseWriter, accessToken, refreshToken string) {
	http.SetCookie(w, &http.Cookie{
		Name: AccessTokenCookie, Value: accessToken, Path: "/",
		MaxAge: 15 * 60, HttpOnly: true, Secure: secureCookies, SameSite: http.SameSiteLaxMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name: RefreshTokenCookie, Value: refreshToken, Path: "/v1/auth",
		MaxAge: 7 * 24 * 60 * 60, HttpOnly: true, Secure: secureCookies, SameSite: http.SameSiteLaxMode,
	})
}

func ClearAuthCookies(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name: AccessTokenCookie, Value: "", Path: "/",
		MaxAge: -1, HttpOnly: true, Secure: secureCookies, SameSite: http.SameSiteLaxMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name: RefreshTokenCookie, Value: "", Path: "/v1/auth",
		MaxAge: -1, HttpOnly: true, Secure: secureCookies, SameSite: http.SameSiteLaxMode,
	})
}
