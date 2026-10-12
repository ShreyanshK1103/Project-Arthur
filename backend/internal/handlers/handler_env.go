package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/auth"
	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/crypto"
	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

var envKeyRegex = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

const (
	maxEnvKeyLen         = 128
	maxEnvValueLen       = 4096
	maxEnvVarsPerProject = 50
)

// ownedProjectID parses {id}, checks the logged-in user owns that project,
// and writes the error response itself when something is wrong.
func (cfg *Config) ownedProjectID(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid project id")
		return uuid.Nil, false
	}

	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return uuid.Nil, false
	}

	project, err := cfg.DB.GetProjectByIDAndUser(
		r.Context(),
		database.GetProjectByIDAndUserParams{
			ID:     projectID,
			UserID: userID,
		},
	)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "project not found")
		return uuid.Nil, false
	}

	return project.ID, true
}

type envVarResponse struct {
	Key       string    `json:"key"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GET /v1/projects/{id}/env  -> keys only, never values
func (cfg *Config) HandlerListEnvVars(w http.ResponseWriter, r *http.Request) {
	projectID, ok := cfg.ownedProjectID(w, r)
	if !ok {
		return
	}

	rows, err := cfg.DB.ListEnvVarsByProject(r.Context(), projectID)
	if err != nil {
		log.Printf("list env vars: %v", err)
		respondWithError(w, http.StatusInternalServerError, "could not list variables")
		return
	}

	resp := make([]envVarResponse, 0, len(rows)) // [] not null in JSON
	for _, row := range rows {
		resp = append(resp, envVarResponse{Key: row.Key, UpdatedAt: row.UpdatedAt})
	}

	respondWithJSON(w, http.StatusOK, resp)
}

// PUT /v1/projects/{id}/env  body: {"key": "...", "value": "..."}
func (cfg *Config) HandlerSetEnvVar(w http.ResponseWriter, r *http.Request) {
	projectID, ok := cfg.ownedProjectID(w, r)
	if !ok {
		return
	}

	var params struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	r.Body = http.MaxBytesReader(w, r.Body, 16<<10)
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if len(params.Key) == 0 || len(params.Key) > maxEnvKeyLen || !envKeyRegex.MatchString(params.Key) {
		respondWithError(w, http.StatusBadRequest,
			"key must start with a letter or underscore and contain only letters, digits and underscores")
		return
	}

	if len(params.Value) > maxEnvValueLen {
		respondWithError(w, http.StatusBadRequest, "value is too long (max 4096 bytes)")
		return
	}

	if strings.ContainsAny(params.Value, "\r\n\x00") {
		respondWithError(w, http.StatusBadRequest, "value cannot contain newlines")
		return
	}

	// Enforce the per-project cap, but still allow updating an existing key.
	existing, err := cfg.DB.ListEnvVarsByProject(r.Context(), projectID)
	if err != nil {
		log.Printf("set env var: list: %v", err)
		respondWithError(w, http.StatusInternalServerError, "could not save variable")
		return
	}

	exists := false
	for _, row := range existing {
		if row.Key == params.Key {
			exists = true
			break
		}
	}

	if !exists && len(existing) >= maxEnvVarsPerProject {
		respondWithError(w, http.StatusBadRequest, "too many variables (max 50 per project)")
		return
	}

	encrypted, err := crypto.Encrypt(params.Value)
	if err != nil {
		log.Printf("set env var: encrypt: %v", err)
		respondWithError(w, http.StatusInternalServerError, "could not save variable")
		return
	}

	err = cfg.DB.UpsertEnvVar(r.Context(), database.UpsertEnvVarParams{
		ProjectID:      projectID,
		Key:            params.Key,
		ValueEncrypted: encrypted,
	})
	if err != nil {
		log.Printf("set env var: upsert: %v", err)
		respondWithError(w, http.StatusInternalServerError, "could not save variable")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"key": params.Key})
}

// DELETE /v1/projects/{id}/env/{key}
func (cfg *Config) HandlerDeleteEnvVar(w http.ResponseWriter, r *http.Request) {
	projectID, ok := cfg.ownedProjectID(w, r)
	if !ok {
		return
	}

	err := cfg.DB.DeleteEnvVar(r.Context(), database.DeleteEnvVarParams{
		ProjectID: projectID,
		Key:       chi.URLParam(r, "key"),
	})
	if err != nil {
		log.Printf("delete env var: %v", err)
		respondWithError(w, http.StatusInternalServerError, "could not delete variable")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "deleted"})
}
