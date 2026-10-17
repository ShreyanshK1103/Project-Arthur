package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/auth"
	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/database"
	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/models"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

// GET /v1/projects — list projects owned by the current user
func (cfg *Config) HandlerListProjects(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	projects, err := cfg.DB.GetProjectsByUsers(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to list projects")
		return
	}

	resp := make([]models.Projects, 0, len(projects))
	for _, p := range projects {
		resp = append(resp, models.ProjectToResponse(p))
	}

	respondWithJSON(w, http.StatusOK, resp)
}

// GET /v1/projects/{id} — get a single project (ownership required)
func (cfg *Config) HandlerGetProject(w http.ResponseWriter, r *http.Request) {
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid project id")
		return
	}

	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
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
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "project not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to load project")
		return
	}

	respondWithJSON(w, http.StatusOK, models.ProjectToResponse(project))
}

// DELETE /v1/projects/{id} — delete a project (ownership required)
func (cfg *Config) HandlerDeleteProject(w http.ResponseWriter, r *http.Request) {
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid project id")
		return
	}

	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Confirm ownership first so we can return 404 vs silent no-op
	_, err = cfg.DB.GetProjectByIDAndUser(
		r.Context(),
		database.GetProjectByIDAndUserParams{
			ID:     projectID,
			UserID: userID,
		},
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "project not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to load project")
		return
	}

	err = cfg.DB.DeleteProjectByIDAndUser(
		r.Context(),
		database.DeleteProjectByIDAndUserParams{
			ID:     projectID,
			UserID: userID,
		},
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete project")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "project deleted"})
}

func (cfg *Config) HandlerCreateProject(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name           string `json:"name"`
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

	userID, ok := auth.GetUserID(r.Context())
	if !ok {
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
	var project database.Project
	var createErr error

	if !validRepoURL(params.RepoUrl) {
		respondWithError(w, 400, "repo_url must be a public https URL")
		return
	}

	for attempt := 0; attempt < 5; attempt++ {
		slug, err := newSlug(params.Name)
		if err != nil {
			respondWithError(w, 500, "Cannot generate project slug")
			return
		}

		project, createErr = cfg.DB.CreateProject(
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
				Slug:           slug,
			},
		)

		// retry only if the slug collided with an existing one
		if createErr == nil || !isUniqueViolation(createErr) {
			break
		}
	}

	if createErr != nil {
		respondWithError(
			w,
			500,
			fmt.Sprintf("Cannot create project: %v", createErr),
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

	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		respondWithError(w, 401, "unauthorized")
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
		respondWithError(w, 403, "project not found")
		return
	}

	deployments, err := cfg.DB.GetDeploymentsByProject(
		r.Context(),
		project.ID,
	)
	if err != nil {
		respondWithError(
			w,
			500,
			"Failed to fetch deployments",
		)
		return
	}

	respondWithJSON(
		w,
		200,
		models.DeploymentForProjectToResponse(
			deployments,
			project,
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
		respondWithError(w, 404, "Project not found")
		return
	}

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

	respondWithJSON(
		w,
		201,
		models.DeploymentToResponse(
			deployment,
		),
	)
}

// HandlerPromoteDeployment makes an older successful deployment the live one
// for the project. Nothing is rebuilt: the files are already in S3.
func (cfg *Config) HandlerPromoteDeployment(w http.ResponseWriter, r *http.Request) {
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, 400, "Invalid project id")
		return
	}

	deploymentID, err := uuid.Parse(chi.URLParam(r, "deploymentID"))
	if err != nil {
		respondWithError(w, 400, "Invalid deployment id")
		return
	}

	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		respondWithError(w, 401, "unauthorized")
		return
	}

	// One statement checks: project is yours, deployment belongs to that
	// project, and the deployment succeeded. Any miss - no row.
	_, err = cfg.DB.PromoteDeployment(
		r.Context(),
		database.PromoteDeploymentParams{
			ID:     projectID,
			UserID: userID,
			ID_2:   deploymentID,
		},
	)
	if errors.Is(err, sql.ErrNoRows) {
		respondWithError(w, 404, "deployment not found for this project, or it did not succeed")
		return
	}
	if err != nil {
		respondWithError(w, 500, "Failed to promote deployment")
		return
	}

	// Return the updated project so the client sees the new active_deployment_id
	project, err := cfg.DB.GetProjectByIDAndUser(
		r.Context(),
		database.GetProjectByIDAndUserParams{ID: projectID, UserID: userID},
	)
	if err != nil {
		respondWithError(w, 500, "Promoted, but failed to load project")
		return
	}

	respondWithJSON(w, 200, models.ProjectToResponse(project))
}
