package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/config"
	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/database"
	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/storage"
	"github.com/google/uuid"
)

func processDeployment(job database.Deployment, db *database.Queries) error {
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

	addLog(
		db,
		job.ID,
		"Cloning repository.....",
	)

	log.Printf("Cloning repo: %s", job.RepoUrl)

	cmd := exec.Command(
		"git",
		"clone",
		job.RepoUrl,
		projectPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		addLog(
			db,
			job.ID,
			fmt.Sprintf(
				"Cloning Failed:\n%s",
				string(output),
			),
		)
		return fmt.Errorf(
			"git clone failed: %v\n%s",
			err, 
			string(output),
		)
	}

	log.Println("Clone Successful")
	addLog(
		db,
		job.ID,
		"Cloning Successful.....",
	)

	// ------------------------- INSTALLING THE PROJECT DEPENDENCIES --------------------------------

	log.Println("Installing dependencies.....")
	addLog(
		db,
		job.ID,
		"Installing Dependencies.......",
	)

	cmd = exec.Command(
		"docker",
		"run",
		"--rm",
		"-v", projectPath+":/app",
		"-w", "/app",
		"node:24",
		"npm",
		"install",
	)

	output, err = cmd.CombinedOutput()

	if err != nil {
		addLog(
			db,
			job.ID,
			fmt.Sprintf(
				"Npm Install Failed:\n%s",
				string(output),
			),
		)
		return fmt.Errorf(
			"npm install failed: %v\n%s",
			err, 
			string(output),
		)
	}
	log.Println("Dependencies installed")
	addLog(
		db,
		job.ID,
		"Dependencies Installed.....",
	)

	// -------------------------BUILDING THE PROJECT -----------------------------

	log.Println("Building project.....")
	addLog(
		db,
		job.ID,
		"Building the Project.......",
	)

	cmd = exec.Command(
		"docker",
		"run",
		"--rm",
		"-v", projectPath+":/app",
		"-w", "/app",
		"node:24",
		"npm",
		"run",
		"build",
	)
	output, err = cmd.CombinedOutput()

	if err != nil {
		addLog(
			db,
			job.ID,
			fmt.Sprintf(
				"Build Failed:\n%s",
				string(output),
			),
		)
		return fmt.Errorf(
			"build failed: %v\n%s",
			err,
			string(output),
		)
	}

	log.Println("Build Successful")
	addLog(
		db,
		job.ID,
		"Build Successful.......",
	)

	//----------------------------- VERIFYING THE BUILD -------------------------------

	distPath := filepath.Join(projectPath, "dist")

	_, err = os.Stat(distPath)
	if err != nil {
		addLog(
			db,
			job.ID,
			"No Dist Found",
		)
		return fmt.Errorf("Dist folder not found")
	}
	// --------------------UPLOADING TO S3---------------------

	addLog(
		db,
		job.ID,
		"Uploading Artifacts To S3",
	)

	err = storage.UploadDirectory(
		job.ID.String(),
		distPath,
	)

	if err != nil {
		addLog(
			db,
			job.ID,
			fmt.Sprintf(
				"S3 Upload Failed : \n%v",
				err,
			),
		)
		return err
	}

	addLog(
		db,
		job.ID,
		fmt.Sprintf(
			"Artifacts Uploaded To S3 Successfully (%s)",
			job.ID.String(),
		),
	)

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
				Valid: true,
			},
		},
	)
	if err != nil {
		addLog(
			db,
			job.ID,
			"Failed to Deploy",
		)
		return err
	}

	addLog(
		db,
		job.ID,
		"Deployment Successful",
	)


	return nil
}

func addLog(db *database.Queries, deploymentID uuid.UUID, message string) {
	err := db.CreateDeploymentLog(
		context.Background(),
		database.CreateDeploymentLogParams{
			DeploymentID: deploymentID,
			Log: message,
		},
	)

	if err != nil {
		log.Printf("failed to write log: %v", err)
	}
}

func main () {
	conn, _, err := config.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	db := database.New(conn)

	// This sets the building failed job to queued again so it can be retried
	err = db.ResetBuildingDeployments(
		context.Background(),
	)
	if err != nil {
		log.Fatal(err)
	}

	for {
		job, err := db.GetNextDeployment(context.Background())
		if err != nil {
			time.Sleep(3 * time.Second)
			continue
		}

		err = processDeployment(job, db)
		if err != nil {
			log.Printf("Processing failed: %v", err)

			db.MarkDeploymentFailed(
				context.Background(),
				job.ID,
			)

			continue
		}

		log.Printf("Successfully processed deployment %s", job.ID)

	}
}
