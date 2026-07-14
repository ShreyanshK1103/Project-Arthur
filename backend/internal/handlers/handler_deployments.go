package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/auth"
	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/database"
	models "github.com/ShreyanshK1103/Project-Arthur/backend/internal/models"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (cfg *Config) HandlerCreateDeployment(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		ProjectID string `json:"project_id"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error Parsing JSON : %v", err))
		return
	}

	projectID, err := uuid.Parse(params.ProjectID)
	if err != nil {
		respondWithError(w, 400, "Invalid project_id")
		return
	}

	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	project, err := cfg.DB.GetProjectByIDAndUser(
		r.Context(),
		database.GetProjectByIDAndUserParams{
			ID:     projectID,
			UserID: userID,
		},
	)

	if err != nil {
		respondWithError(w, 403, "Project not found")
		return
	}

	jobs, err := cfg.DB.CreateDeployment(r.Context(), database.CreateDeploymentParams{
		ProjectID: project.ID,
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

	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	deployment, err := cfg.DB.GetDeploymentByIDAndUser(
		r.Context(),
		database.GetDeploymentByIDAndUserParams{
			ID:     id,
			UserID: userID,
		},
	)

	respondWithJSON(w, 200, models.DeploymentToResponse(deployment))

}
