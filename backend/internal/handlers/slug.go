package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"regexp"
	"strings"

	"github.com/lib/pq"
)

var nonAlnum = regexp.MustCompile(`[^a-z0-9]+`)

// slugBase turns "MY PORTFOLIO!!" into "my-portfolio",
// Capped at 30 chars so slug + "-xxx" stays well under the 63 char DNS label limit.
func slugBase(name string) string {
	s := nonAlnum.ReplaceAllString(strings.ToLower(name), "-")
	s = strings.Trim(s, "-")
	if len(s) > 30 {
		s = strings.Trim(s[:30], "-")
	}
	if s == "" {
		s = "project"
	}
	return s
}

// newSlug returns e.g. "my-portfolio-a1b2"
func newSlug(name string) (string, error) {
	b := make([]byte, 2)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return slugBase(name) + "-" + hex.EncodeToString(b), nil
}

func isUniqueViolation(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && pqErr.Code == "23505"
}
