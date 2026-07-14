package auth

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func MiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(
				w,
				"missing authorization header",
				http.StatusUnauthorized,
			)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(
				w,
				"invalid authorization header",
				http.StatusUnauthorized,
			)
			return
		}

		tokenString := strings.TrimPrefix(
			authHeader,
			"Bearer ",
		)

		claims, err := ValidAccessToken(tokenString)

		if err != nil {
			http.Error(
				w,
				"invalid token",
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