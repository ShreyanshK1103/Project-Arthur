package deployment

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/builder"
	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/database"
	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/logs"
	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/storage"
)

func ProcessDeployment(job database.Deployment, db *database.Queries) error {
	// -----------------TEMP FOLDER CREATION FOR THE DEPLOYMENT ----------------------
	projectPath := "/tmp/project-arthur/" + job.ID.String()
	defer func() {
		if err := os.RemoveAll(projectPath); err != nil {
			log.Printf(
				"failed to cleanup %s: %v",
				projectPath,
				err,
			)
		}
	}()

	err := os.MkdirAll(projectPath, 0755)
	if err != nil {
		return err
	}

	// --------------------------CLONING THE REPO ---------------------------------------

	err = builder.CloneRepo(job, db, projectPath)
	if err != nil {
		return err
	}

	// ------------------------- INSTALLING THE PROJECT DEPENDENCIES --------------------------------

	err = builder.InstallDependencies(job, db, projectPath)
	if err != nil {
		return err
	}

	// -------------------------BUILDING THE PROJECT -----------------------------
	err = builder.BuildProject(job, db, projectPath)
	if err != nil {
		return err
	}
	//----------------------------- VERIFYING THE BUILD -------------------------------

	distPath := filepath.Join(projectPath, "dist")

	_, err = os.Stat(distPath)
	if err != nil {
		logs.AddLog(
			db,
			job.ID,
			"No Dist Found",
		)
		return fmt.Errorf("Dist folder not found")
	}
	// --------------------UPLOADING TO S3---------------------
	logs.AddLog(
		db,
		job.ID,
		"Uploading to S3",
	)

	err = storage.UploadDirectory(
		job.ID.String(),
		distPath,
	)

	if err != nil {
		return err
	}
	// -----------------------------MARKING SUCCESS OF THE PROJECT DEPLOYMENT --------------------------------

	slug := job.ID.String()[:12]

	deploymentURL := fmt.Sprintf(
		"http://%s.localhost:8080",
		slug,
	)
	err = db.MarkDeploymentSuccess(
		context.Background(),
		database.MarkDeploymentSuccessParams{
			ID: job.ID,
			Url: sql.NullString{
				String: deploymentURL,
				Valid:  true,
			},
		},
	)
	if err != nil {
		logs.AddLog(
			db,
			job.ID,
			"Failed to Deploy",
		)
		return err
	}

	logs.AddLog(
		db,
		job.ID,
		"Deployment Successful",
	)

	return nil
}
