package config

import (
	"os"
	"strings"
)

// FrontendURL is the dashboard's origin, e.g. http://localhost:5173 locally
// or https://app.yourdomain.com in production (no trailing slash).
func FrontendURL() string {
	v := strings.TrimRight(os.Getenv("FRONTEND_URL"), "/")
	if v == "" {
		return "http://localhost:5173"
	}
	return v
}
