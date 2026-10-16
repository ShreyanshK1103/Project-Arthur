package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/auth"
	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/database"
	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/models"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

// Must match the limit in GetDeploymentLogsAfter
const logsPageSize = 500

func (cfg *Config) HandlerGetDeploymentLogs(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		respondWithError(
			w,
			400,
			"Invalid Deployment ID",
		)
		return
	}

	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		respondWithError(w, 401, "unauthorized")
		return
	}

	var after int64
	if v := r.URL.Query().Get("after"); v != "" {
		after, err = strconv.ParseInt(v, 10, 64)
		if err != nil || after < 0 {
			respondWithError(w, 400, "Invalid after cursor")
			return
		}
	}
	// OwnerShip Check
	deployment, err := cfg.DB.GetDeploymentByIDAndUser(r.Context(), database.GetDeploymentByIDAndUserParams{ID: id, UserID: userID})
	if errors.Is(err, sql.ErrNoRows) {
		respondWithError(w, 404, "deployment not found")
		return
	}
	if err != nil {
		respondWithError(w, 500, "Failed to load deployment")
		return
	}

	logs, err := cfg.DB.GetDeploymentLogsAfter(
		r.Context(),
		database.GetDeploymentLogsAfterParams{DeploymentID: id, Seq: after},
	)
	if err != nil {
		respondWithError(w, 500, "Failed to fetch logs")
		return
	}

	next := after
	if len(logs) > 0 {
		next = logs[len(logs)-1].Seq
	}

	finished := deployment.Status == "success" || deployment.Status == "failed"

	respondWithJSON(w, 200, models.DeploymentLogsResponse{
		Logs:      models.DeploymentLogsToLines(logs),
		Status:    deployment.Status,
		NextAfter: next,
		Done:      finished && len(logs) < logsPageSize,
	})
}
