package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/database"
)

func (cfg *Config) HandlerGithubWebhook(w http.ResponseWriter, r *http.Request) {

	type GithubPayload struct {
		Repository struct {
			ID int64 `json:"id"`
		} `json:"repository"`

		Ref string `json:"ref"`
	}

	payload := GithubPayload{}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		respondWithError(
			w,
			400,
			"invalid payload",
		)

		return
	}

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
				"message":"project ignored",
			},
		)
		return
	}

	deployment, err := cfg.DB.CreateDeployment(
		r.Context(),
		database.CreateDeploymentParams{
			ProjectID: project.ID,
			Status: "queued",
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

	respondWithJSON(
		w,
		201,
		deployment,
	)
}