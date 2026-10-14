package siteurl

import (
	"fmt"
	"os"
)

// ForSlug builds the public URL for a site
// In production set SITE_SCHEME=https and SITE_BASE_DOMAIN=yourdomain.com
func ForSlug(slug string) string {
	scheme := getenv("SITE_SCHEME", "http")
	base := getenv("SITE_BASE_DOMAIN", "localhost:8080")
	return fmt.Sprintf("%s://%s.%s", scheme, slug, base)
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
