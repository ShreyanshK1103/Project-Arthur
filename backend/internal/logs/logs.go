package logs

import (
	"context"
	"log"

	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/database"
	"github.com/google/uuid"
)

func AddLog(db *database.Queries, deploymentID uuid.UUID, message string) {
	err := db.CreateDeploymentLog(
		context.Background(),
		database.CreateDeploymentLogParams{
			DeploymentID: deploymentID,
			Log:          message,
		},
	)

	if err != nil {
		log.Printf("failed to write log: %v", err)
	}
}
