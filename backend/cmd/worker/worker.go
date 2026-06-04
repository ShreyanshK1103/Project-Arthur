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
)

func processDeployment(job database.Deployment, db *database.Queries) error {
	projectPath := "/tmp/project-arthur/" + job.ID.String()
	defer os.RemoveAll(projectPath)

	err := os.MkdirAll(projectPath, 0755)
	if err != nil {
		return err
	}

	log.Printf("Cloning repo: %s", job.RepoUrl)

	cmd := exec.Command(
		"git",
		"clone",
		job.RepoUrl,
		projectPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"git clone failed: %v\n%s",
			err, 
			string(output),
		)
	}

	log.Println("Clone Successfull")

	log.Println("Installing dependencies.....")

	cmd = exec.Command("npm", "install")
	cmd.Dir = projectPath

	output, err = cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf(
			"npm install failed: %v\n%s",
			err, 
			string(output),
		)
	}
	log.Println("Dependencies installed")

	log.Println("Building project.....")

	cmd = exec.Command("npm", "run", "build")
	cmd.Dir = projectPath
	output, err = cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf(
			"build failed: %v\n%s",
			err,
			string(output),
		)
	}

	log.Println("Build Successful")

	distPath := filepath.Join(projectPath, "dist")

	_, err = os.Stat(distPath)
	if err != nil {
		return fmt.Errorf("Dist folder not found")
	}

	err = db.MarkDeploymentSuccess(
		context.Background(),
		database.MarkDeploymentSuccessParams{
			ID: job.ID,
			Url: sql.NullString{
				String: "https://success.local",
				Valid: true,
			},
		},
	)
	if err != nil {
		return err
	}


	return nil
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
