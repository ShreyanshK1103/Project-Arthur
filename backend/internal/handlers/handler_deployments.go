package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/database"
	"github.com/google/uuid"
)

func (cfg *Config) HandlerCreateDeployment(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		ProjectID string `json:"project_id"`
		ProjectURL string `json:"repo_url"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error Parsing JSON : %v", err));
		return
	}

	projectID, err := uuid.Parse(params.ProjectID)
	if err != nil {
		respondWithError(w, 400, "Invalid project_id")
		return
	}

	jobs, err := cfg.DB.CreateDeployment(r.Context(), database.CreateDeploymentParams{
		ProjectID: projectID,
		RepoUrl:   params.ProjectURL,
		Status:    "queued",
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't Get Jobs: %v", err))
		return
	}
	respondWithJSON(w, 201, jobs)
}