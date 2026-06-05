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
	"github.com/google/uuid"
)

func processDeployment(job database.Deployment, db *database.Queries) error {
	// -----------------TEMP FOLDER CREATION FOR THE DEPLOYMENT ----------------------
	projectPath := "/tmp/project-arthur/" + job.ID.String()

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
			"Cloning Failed",
		)
		return fmt.Errorf(
			"git clone failed: %v\n%s",
			err, 
			string(output),
		)
	}

	log.Println("Clone Successfull")
	addLog(
		db,
		job.ID,
		"Cloning Completed.....",
	)

	// ------------------------- INSTALLING THE PROJECT DEPENDENCIES --------------------------------

	log.Println("Installing dependencies.....")
	addLog(
		db,
		job.ID,
		"Installing Dependencies.......",
	)

	cmd = exec.Command("npm", "install")
	cmd.Dir = projectPath

	output, err = cmd.CombinedOutput()

	if err != nil {
		addLog(
			db,
			job.ID,
			"npm failed",
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

	cmd = exec.Command("npm", "run", "build")
	cmd.Dir = projectPath
	output, err = cmd.CombinedOutput()

	if err != nil {
		addLog(
			db,
			job.ID,
			"Build Failed",
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
		"Build Successfull.......",
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

	// ------------------- COPYING THE DIST ARTIFACT -----------------------------------------

	storageBasePath := "./storage/deployments"

	err = os.MkdirAll(storageBasePath, 0755)
	if err != nil {
		return err
	}

	deploymentStoragePath := filepath.Join(
		storageBasePath,
		job.ID.String(),
	)

	err = os.MkdirAll(deploymentStoragePath, 0755)
	if err != nil {
		return err
	}

	log.Println("Copying build artifacts....")
	addLog(
		db,
		job.ID,
		"Copying Build Artifacts.....",
	)

	cmd = exec.Command(
		"cp",
		"-r",
		distPath+"/.",
		deploymentStoragePath,
	)

	output, err = cmd.CombinedOutput()
	if err != nil {
		addLog(
			db,
			job.ID,
			"Artifacts failed to copy",
		)
		return fmt.Errorf(
			"copy failed: %v\n%s",
			err, 
			string(output),
		)
	}

	// ------------------ VERIFYING THE DIST COPY --------------------------------------

	files, err := os.ReadDir(deploymentStoragePath)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return fmt.Errorf("no artifacts copied")
	}

	log.Println("Stored artifacts:")

	addLog(
		db,
		job.ID,
		"Artifacts Copied Successfully",
	)

	for _, file := range files {
		log.Println(file.Name())
	}
	// ------------------------------------ REMOVING THE TEMP FILES ----------------------------------------------
	err = os.RemoveAll(projectPath)
	if err != nil {
		return err
	}

	log.Printf("Storage path: %s", storageBasePath)

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
	db.CreateDeploymentLog(
		context.Background(),
		database.CreateDeploymentLogParams{
			DeploymentID: deploymentID,
			Log: message,
		},
	)
}

func main () {
	conn, _, err := config.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	db := database.New(conn)

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
