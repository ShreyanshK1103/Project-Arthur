package builder

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/database"
	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/logs"
)

const (
	containerMemory = "2g"
	containerCPUs = "2"
	containerPIDs = "512"
	stepTimeout = 10 * time.Minute
)

// runInContainer runs a shell command in a node container with resource limits,
// a non-root user and a hard timeout, name must be unique per run 
func runInContainer(ctx context.Context, name, projectPath, envFile, command string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, stepTimeout)
	defer cancel()

	//Clear any leftover container with this name
	exec.Command("docker", "rm", "-f", name).Run()

	args := []string{
		"run", "--rm",
		"--name", name,
		"--memory", containerMemory,
		"--memory-swap", containerMemory, // no swap beyond the limit
		"--cpus", containerCPUs,
		"--pids-limit", containerPIDs, // stops fork bombs
		"--cap-drop", "ALL",
		"--security-opt", "no-new-privileges",
		"--user", fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid()),
		"-e", "HOME=/tmp", // npm needs a writable home
	}
	args = append(args, EnvFileArgs(envFile)...)
	args = append(args,
		"-v", projectPath+":/app",
		"-w", "/app",
		"node:24",
		"sh", "-c", command,
	)

	output, err := exec.CommandContext(ctx, "docker", args...).CombinedOutput()

	if ctx.Err() == context.DeadlineExceeded{
		//Killing the docker CLI does NOT stop the container, so remove it manually
		exec.Command("docker", "rm", "-f", name).Run()
		return output, fmt.Errorf("timed out after %s", stepTimeout)
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && exitErr.ExitCode() == 137 {
		err = fmt.Errorf("killed (exit 137), likely exceeded the %s memory limit: %w", containerMemory, err)
	}

	return output, err
}
func InstallDependencies(command string, env *EnvVars, envFile string, job database.Deployment, db *database.Queries, projectPath string) error {
	log.Println("Installing dependencies.....")
	logs.AddLog(db, job.ID, "Installing Dependencies.......")

	output, err := runInContainer(context.Background(), fmt.Sprintf("arthur-%s-install", job.ID), projectPath, envFile, command)
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

	output, err := runInContainer(context.Background(), fmt.Sprintf("arthur-%s-build", job.ID), projectPath, envFile, command)
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