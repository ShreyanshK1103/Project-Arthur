package models

import (
	"time"

	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/database"
	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/siteurl"
	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	AvatarURL     *string   `json:"avatar_url,omitempty"`
	Provider      string    `json:"provider"`
	EmailVerified bool      `json:"email_verified"`
}

type Projects struct {
	ID                 uuid.UUID  `json:"project_id"`
	Name               string     `json:"name"`
	UserID             uuid.UUID  `json:"user_id"`
	RepoUrl            string     `json:"repo_url"`
	InstallCommand     string     `json:"install_command"`
	BuildCommand       string     `json:"build_command"`
	OutputDir          string     `json:"output_dir"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	GithubRepoID       *int64     `json:"github_repo_id"`
	Branch             string     `json:"branch"`
	AutoDeploy         bool       `json:"auto_deploy"`
	Slug               string     `json:"slug"`
	URL                string     `json:"url"`
	ActiveDeploymentID *uuid.UUID `json:"active_deployment_id"`
}

type Deployments struct {
	ID         uuid.UUID `json:"id"`
	ProjectID  uuid.UUID `json:"project_id"`
	Status     string    `json:"status"`
	Url        *string   `json:"url"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	ProjectURL string    `json:"project_url,omitempty"`
	IsActive   *bool     `json:"is_active,omitempty"`
}

type DeploymentLogsResponse struct {
	Logs      []LogLine `json:"logs"`
	Status    string    `json:"status"`
	NextAfter int64     `json:"next_after"`
	Done      bool      `json:"done"`
}

type LogLine struct {
	Seq       int64     `json:"seq"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

func DeploymentToResponse(d database.Deployment) Deployments {
	var url *string
	if d.Url.Valid {
		url = &d.Url.String
	}

	return Deployments{
		ID:        d.ID,
		ProjectID: d.ProjectID,
		Status:    d.Status,
		Url:       url,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

func ProjectToResponse(p database.Project) Projects {
	var repoID *int64

	if p.GithubRepoID.Valid {
		repoID = &p.GithubRepoID.Int64
	}

	var activeID *uuid.UUID
	if p.ActiveDeploymentID.Valid {
		id := p.ActiveDeploymentID.UUID
		activeID = &id
	}

	return Projects{
		ID:                 p.ID,
		Name:               p.Name,
		UserID:             p.UserID,
		RepoUrl:            p.RepoUrl,
		InstallCommand:     p.InstallCommand,
		BuildCommand:       p.BuildCommand,
		OutputDir:          p.OutputDir,
		CreatedAt:          p.CreatedAt,
		UpdatedAt:          p.UpdatedAt,
		GithubRepoID:       repoID,
		Branch:             p.Branch,
		AutoDeploy:         p.AutoDeploy,
		Slug:               p.Slug,
		URL:                siteurl.ForSlug(p.Slug),
		ActiveDeploymentID: activeID,
	}
}

func DeploymentsToResponse(deployments []database.Deployment) []Deployments {
	result := []Deployments{}

	for _, deployment := range deployments {
		result = append(
			result,
			DeploymentToResponse(deployment),
		)
	}

	return result
}

func UserToResponse(u database.User) User {
	return User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func DeploymentForProjectToResponse(deployments []database.Deployment, p database.Project) []Deployments {
	result := []Deployments{}
	projectURL := siteurl.ForSlug(p.Slug)

	for _, d := range deployments {
		resp := DeploymentToResponse(d)
		active := p.ActiveDeploymentID.Valid && p.ActiveDeploymentID.UUID == d.ID
		resp.ProjectURL = projectURL
		resp.IsActive = &active
		result = append(result, resp)
	}

	return result
}
func DeploymentLogsToLines(logs []database.DeploymentLog) []LogLine {
	result := []LogLine{}

	for _, l := range logs {
		result = append(result, LogLine{
			Seq:       l.Seq,
			Message:   l.Log,
			CreatedAt: l.CreatedAt,
		})
	}

	return result
}
