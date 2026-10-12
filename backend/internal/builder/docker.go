package builder

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/database"
	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/logs"
)

func runInContainer(projectPath, envFile, command string) ([]byte, error) {
	args := []string{"run", "--rm"}
	args = append(args, EnvFileArgs(envFile)...)
	args = append(args,
		"-v", projectPath+":/app",
		"-w", "/app",
		"node:24",
		"sh", "-c", command,
	)

	return exec.Command("docker", args...).CombinedOutput()
}
func InstallDependencies(command string, env *EnvVars, envFile string, job database.Deployment, db *database.Queries, projectPath string) error {
	log.Println("Installing dependencies.....")
	logs.AddLog(db, job.ID, "Installing Dependencies.......")

	output, err := runInContainer(projectPath, envFile, command)
	if err != nil {
		safe := env.Redact(string(output))
		logs.AddLog(db, job.ID, fmt.Sprintf("Install Failed:\n%s", safe))
		return fmt.Errorf("install failed: %v\n%s", err, safe)
	}

	log.Println("Dependencies installed")
	logs.AddLog(db, job.ID, "Dependencies Installed.....")
	return nil
}

func BuildProject(command string, env *EnvVars, envFile string, job database.Deployment, db *database.Queries, projectPath string) error {
	log.Println("Building project.....")
	logs.AddLog(db, job.ID, "Building the Project.......")

	output, err := runInContainer(projectPath, envFile, command)
	if err != nil {
		safe := env.Redact(string(output))
		logs.AddLog(db, job.ID, fmt.Sprintf("Build Failed:\n%s", safe))
		return fmt.Errorf("build failed: %v\n%s", err, safe)
	}

	// Successful builds can print values too (e.g. a build script that echoes them).
	if safe := env.Redact(string(output)); safe != "" {
		logs.AddLog(db, job.ID, fmt.Sprintf("Build Output:\n%s", safe))
	}

	log.Println("Build Successful")
	logs.AddLog(db, job.ID, "Build Successful.......")
	return nil
}