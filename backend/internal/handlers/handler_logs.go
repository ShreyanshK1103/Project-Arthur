package handlers

import (
	"fmt"
	"net/http"

	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/models"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

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

	logs, err := cfg.DB.GetDeploymentLogs(
		r.Context(),
		id,
	)
	if err != nil {
		respondWithError(
			w,
			404,
			fmt.Sprintf(
				"could not fetch logs: %v",
				err,
			),
		)
		return
	}

	respondWithJSON(
		w,
		200,
		models.DeploymentLogsResponse{
			Logs:models.DeploymentLogsToStrings(logs),
		},
	)
}