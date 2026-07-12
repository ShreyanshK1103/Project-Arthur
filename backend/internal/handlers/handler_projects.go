package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/database"
	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/models"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (cfg *Config) HandlerCreateProject(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name           string `json:"name"`
		UserID         string `json:"user_id"`
		RepoUrl        string `json:"repo_url"`
		InstallCommand string `json:"install_command"`
		BuildCommand   string `json:"build_command"`
		OutputDir      string `json:"output_dir"`
		GithubRepoID   *int64 `json:"github_repo_id"`
		Branch         string `json:"branch"`
		AutoDeploy     *bool  `json:"auto_deploy"`
	}

	params := parameters{}

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(
			w,
			400,
			fmt.Sprintf("Invalid JSON: %v", err),
		)
		return
	}

	userID, err := uuid.Parse(params.UserID)
	if err != nil {
		respondWithError(
			w,
			400,
			"Invalid User ID",
		)
		return
	}

	if params.InstallCommand == "" {
		params.InstallCommand = "npm install"
	}

	if params.BuildCommand == "" {
		params.BuildCommand = "npm run build"
	}

	if params.OutputDir == "" {
		params.OutputDir = "dist"
	}

	if params.Branch == "" {
		params.Branch = "main"
	}

	githubRepoID := sql.NullInt64{}

	if params.GithubRepoID != nil {
		githubRepoID = sql.NullInt64{
			Int64: *params.GithubRepoID,
			Valid: true,
		}
	}

	autoDeploy := true

	if params.AutoDeploy != nil {
		autoDeploy = *params.AutoDeploy
	}
	project, err := cfg.DB.CreateProject(
		r.Context(),
		database.CreateProjectParams{
			Name:           params.Name,
			UserID:         userID,
			RepoUrl:        params.RepoUrl,
			InstallCommand: params.InstallCommand,
			BuildCommand:   params.BuildCommand,
			OutputDir:      params.OutputDir,
			GithubRepoID:   githubRepoID,
			Branch:         params.Branch,
			AutoDeploy:     autoDeploy,
		},
	)

	if err != nil {
		respondWithError(
			w,
			500,
			fmt.Sprintf("Cannot create project: %v", err),
		)
		return
	}

	respondWithJSON(
		w,
		201,
		models.ProjectToResponse(project),
	)
}

func (cfg *Config) HandlerGetProjectDeployment(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(
		r,
		"id",
	)

	projectID, err := uuid.Parse(idParam)

	if err != nil {
		respondWithError(
			w,
			400,
			"Invalid project id",
		)
		return
	}

	deployments, err := cfg.DB.GetDeploymentsByProject(
		r.Context(),
		projectID,
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
		200,
		models.DeploymentsToResponse(
			deployments,
		),
	)
}

func (cfg *Config) HandlerRedeployProject(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(
		r,
		"id",
	)

	projectID, err := uuid.Parse(
		idParam,
	)

	if err != nil {
		respondWithError(
			w,
			400,
			"invalid project id",
		)
		return
	}

	deployment, err := cfg.DB.CreateDeployment(
		r.Context(),
		database.CreateDeploymentParams{
			ProjectID: projectID,
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

	respondWithJSON(
		w,
		201,
		models.DeploymentToResponse(
			deployment,
		),
	)
}
