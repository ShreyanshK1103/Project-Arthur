package models

import (
	"time"

	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name string `json:"name"`
	Email string `json:"email"`
}

type Projects struct {
	ID uuid.UUID `json:"project_id"`
	Name string `json:"name"`
	UserID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

type Deployments struct {
	ID uuid.UUID `json:"id"`
    ProjectID uuid.UUID `json:"project_id"`
    Status string `json:"status"`
    RepoUrl string `json:"repo_url"`
    Url *string `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func DeploymentToResponse(d database.Deployment) Deployments {
	var url *string
	if d.Url.Valid {
		url = &d.Url.String
	}

	return Deployments{
		ID: d.ID,
		ProjectID: d.ProjectID,
		Status: d.Status,
		RepoUrl: d.RepoUrl,
		Url: url,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}