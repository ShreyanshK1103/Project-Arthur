package builder

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/database"
	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/logs"
)

func CloneRepo(repoUrl string, job database.Deployment, db *database.Queries, projectPath string) error {

	logs.AddLog(
		db,
		job.ID,
		"Cloning repository.....",
	)

	log.Printf("Cloning repo: %s", repoUrl)

	cmd := exec.Command(
		"git",
		"clone",
		repoUrl,
		projectPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logs.AddLog(
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
	logs.AddLog(
		db,
		job.ID,
		"Cloning Successful.....",
	)

	return nil

}
