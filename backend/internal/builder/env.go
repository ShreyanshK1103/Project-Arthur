package builder

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/crypto"
	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/database"
	"github.com/google/uuid"
)

// Values shorter than this are not redacted from logs; replacing a value like "1"
// or "true" everywhere would mangle the build output

const minRedactLen = 4

type EnvVars struct {
	vars map[string]string
}

func LoadEnvVars(ctx context.Context, db *database.Queries, projectID uuid.UUID) (*EnvVars, error) {
	rows, err := db.ListEnvVarsByProject(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("load env vars: %w", err)
	}

	vars := make(map[string]string, len(rows))
	for _, row := range rows {
		value, err := crypto.Decrypt(row.ValueEncrypted)
		if err != nil {
			return nil, fmt.Errorf("decrypt env var %s: %w", row.Key, err)
		}
		vars[row.Key] = value
	}

	return &EnvVars{vars: vars}, nil
}

// WriteEnvFile writes KEY=value lines to a private temp file (0600) and returns
// its path plus a cleanup func. With no variables it returns an empty path.
func (e *EnvVars) WriteEnvFile() (string, func(), error) {
	if e == nil || len(e.vars) == 0 {
		return "", func() {}, nil
	}

	keys := make([]string, 0, len(e.vars))
	for k := range e.vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		sb.WriteString(k + "=" + e.vars[k] + "\n")
	}

	f, err := os.CreateTemp("", "arthur-env-*") // created with mode 0600
	if err != nil {
		return "", func() {}, err
	}

	cleanup := func() { os.Remove(f.Name()) }

	if _, err := f.WriteString(sb.String()); err != nil {
		f.Close()
		cleanup()
		return "", func() {}, err
	}
	if err := f.Close(); err != nil {
		cleanup()
		return "", func() {}, err
	}

	return f.Name(), cleanup, nil
}

// EnvFileArgs returns the docker run flags for an env-file path ("" -> no flags).
func EnvFileArgs(path string) []string {
	if path == "" {
		return nil
	}
	return []string{"--env-file", path}
}

// Redact replaces every secret value in text with ***, longest values first.
func (e *EnvVars) Redact(text string) string {
	if e == nil {
		return text
	}

	values := make([]string, 0, len(e.vars))
	for _, v := range e.vars {
		if len(v) >= minRedactLen {
			values = append(values, v)
		}
	}

	sort.Slice(values, func(i, j int) bool { return len(values[i]) > len(values[j]) })

	for _, v := range values {
		text = strings.ReplaceAll(text, v, "***")
	}

	return text
}
