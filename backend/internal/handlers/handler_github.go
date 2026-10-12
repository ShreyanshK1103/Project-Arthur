package handlers

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/database"
	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/models"
)

func (cfg *Config) HandlerGithubWebhook(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1<<20))
	if err != nil {
		respondWithError(w, http.StatusRequestEntityTooLarge, "payload too large")
		return
	}

	if !validGithubSignature(body, r.Header.Get("X-Hub-Signature-256")) {
		log.Println("github webhook: invalid or missing signature")
		respondWithError(w, http.StatusUnauthorized, "invalid signature")
		return
	}

	event := r.Header.Get("X-GitHub-Event")

	switch event {
	case "ping":
		respondWithJSON(w, 200, map[string]string{
			"message": "pong",
		})
		return

	case "push":
		// continue

	default:
		respondWithJSON(w, 200, map[string]string{
			"message": "ignored",
		})
		return
	}

	log.Println("========== Github Webhook =========")

	type GithubPayload struct {
		Repository struct {
			ID int64 `json:"id"`
		} `json:"repository"`

		Ref     string `json:"ref"`
		Deleted bool   `json:"deleted"`
	}

	payload := GithubPayload{}

	err = json.Unmarshal(body, &payload)
	if err != nil {
		respondWithError(
			w,
			400,
			"invalid payload",
		)

		return
	}
	if payload.Deleted {
		respondWithJSON(w, 200, map[string]string{"message": "branch deletion ignored"})
		return
	}

	log.Printf("Repository ID: %d", payload.Repository.ID)
	log.Printf("Ref: %s", payload.Ref)

	project, err := cfg.DB.GetProjectByGithubRepoID(
		r.Context(),
		sql.NullInt64{
			Int64: payload.Repository.ID,
			Valid: true,
		},
	)

	if err != nil {
		respondWithJSON(
			w,
			200,
			map[string]string{
				"message": "project ignored",
			},
		)
		return
	}

	expectedRef := "refs/heads/" + project.Branch

	if payload.Ref != expectedRef {
		log.Printf(
			"Ignoring push on %s (project tracks %s)",
			payload.Ref,
			expectedRef,
		)

		respondWithJSON(w, 200, map[string]string{
			"message": "branch ignored",
		})
		return
	}

	log.Printf(
		"Matched project: %s (%s)",
		project.Name,
		project.ID,
	)

	deployment, err := cfg.DB.CreateDeployment(
		r.Context(),
		database.CreateDeploymentParams{
			ProjectID: project.ID,
			Status:    "queued",
		},
	)

	if err != nil {
		respondWithError(
			w,
			500,
			err.Error(),
		)
		return
	}

	log.Printf(
		"Queued deployment: %s",
		deployment.ID,
	)

	respondWithJSON(
		w,
		202,
		models.DeploymentToResponse(deployment),
	)
}
