package builder

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/database"
	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/logs"
)

func InstallDependencies(command string, job database.Deployment, db *database.Queries, projectPath string) error {

	log.Println("Installing dependencies.....")
	logs.AddLog(
		db,
		job.ID,
		"Installing Dependencies.......",
	)

	cmd := exec.Command(
		"docker",
		"run",
		"--rm",
		"-v", projectPath+":/app",
		"-w", "/app",
		"node:24",
		"sh",
		"-c",
		command,
	)

	output, err := cmd.CombinedOutput()

	if err != nil {
		logs.AddLog(
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
	logs.AddLog(
		db,
		job.ID,
		"Dependencies Installed.....",
	)
	return nil
}

func BuildProject(job database.Deployment, db *database.Queries, projectPath string) error {
	log.Println("Building project.....")
	logs.AddLog(
		db,
		job.ID,
		"Building the Project.......",
	)

	cmd := exec.Command(
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
	output, err := cmd.CombinedOutput()

	if err != nil {
		logs.AddLog(
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
	logs.AddLog(
		db,
		job.ID,
		"Build Successful.......",
	)
	return nil
}
