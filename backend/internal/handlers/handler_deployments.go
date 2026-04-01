package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/database"
	models "github.com/ShreyanshK1103/Project-Arthur/backend/internal/models"
	"github.com/go-chi/chi"
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
	respondWithJSON(w, 201, models.DeploymentToResponse(jobs))
}

func (cfg *Config) HandlerGetDeployment(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		respondWithError(w, 400, "Missing Deployment ID")
		return
	}

	id, err := uuid.Parse(idParam)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Invalid Deployment ID: %v", err))
	}

	deployment, err := cfg.DB.GetDeploymentByID(r.Context(), id)
	if err != nil {
		respondWithError(w, 404, fmt.Sprintf("Deployment Not Found: %v", err))
		return
	}

	respondWithJSON(w, 200, models.DeploymentToResponse(deployment))

}