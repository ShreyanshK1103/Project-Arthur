package auth

import (
	"context"
	"net/http"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func MiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("arthur_access_token")

		if err != nil {
			http.Error(
				w,
				"unauthorized",
				http.StatusUnauthorized,
			)
			return
		}

		token := cookie.Value

		claims, err := ValidAccessToken(token)

		if err != nil {
			http.Error(
				w,
				"unauthorized",
				http.StatusUnauthorized,
			)
			return
		}

		ctx := context.WithValue(
			r.Context(),
			UserIDKey,
			claims.UserID,
		)

		next.ServeHTTP(
			w,
			r.WithContext(ctx),
		)
	})
}
