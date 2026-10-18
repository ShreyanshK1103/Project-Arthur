package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// resolveGithubRepoID asks GitHub's public API for a repo's numeric ID, which
// is what push webhooks identify a repository by. If the lookup fails the
// result is invalid (NULL), and project creation carries on without it.
func resolveGithubRepoID(ctx context.Context, repoURL string) sql.NullInt64 {
	id, err := fetchGithubRepoID(ctx, repoURL)
	if err != nil {
		log.Printf("github repo id lookup failed for %q: %v", repoURL, err)
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: id, Valid: true}
}

func fetchGithubRepoID(ctx context.Context, repoURL string) (int64, error) {
	u, err := url.Parse(repoURL)
	if err != nil || !strings.EqualFold(u.Hostname(), "github.com") {
		return 0, fmt.Errorf("not a github.com URL")
	}

	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
		return 0, fmt.Errorf("cannot parse owner/repo from URL")
	}
	owner := parts[0]
	repo := strings.TrimSuffix(parts[1], ".git")

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	apiURL := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s",
		url.PathEscape(owner), url.PathEscape(repo),
	)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "arthur")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("github api returned status %d", resp.StatusCode)
	}

	var body struct {
		ID int64 `json:"id"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&body); err != nil {
		return 0, err
	}
	if body.ID == 0 {
		return 0, fmt.Errorf("no id in github response")
	}
	return body.ID, nil
}