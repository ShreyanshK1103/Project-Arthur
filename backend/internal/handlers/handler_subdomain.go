package handlers

import (
	"database/sql"
	"errors"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/storage"
	"github.com/google/uuid"
)

// Deployment permalinks are the first 12 hex chars of the deployment UUID.
// Validating the shape also stops a crafted Host like "_" or "%" from
// matching arbitrary rows in the LIKE-based prefix lookup
var deploymentPrefixRe = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{3}$`)

func (cfg *Config) HandlerServeByDomain(w http.ResponseWriter, r *http.Request) {
	// THIS FUNCTION HANDLES THE SERVING LOGIC INSTEAD OF STATIC HANDLING
	// {deploymentID}.localhost:8080 is the root in this logic cause of the file rendering issues prior to the
	// localhost:8080/v1/deployments/{deploymentID} logic was flawed cause the assests was being searched in the dir localhost:8080
	// But in reality the assests were being served in the given path, due to this subdomain approach was considered as an viable
	// architectural decision
	host := r.Host

	host = strings.Split(host, ":")[0]

	slug := strings.Split(host, ".")[0]

	var deploymentID uuid.UUID
	viaProjectSlug := false

	project, err := cfg.DB.GetProjectBySlug(r.Context(), slug)
	switch {
	case err == nil:
		if !project.ActiveDeploymentID.Valid {
			http.Error(w, "this project has no live deployment yet", http.StatusNotFound)
			return
		}
		deploymentID = project.ActiveDeploymentID.UUID
		viaProjectSlug = true
	case errors.Is(err, sql.ErrNoRows):
		if !deploymentPrefixRe.MatchString(slug) {
			http.Error(w, "site not found", http.StatusNotFound)
			return
		}
		d, err := cfg.DB.GetDeploymentByPrefix(r.Context(), sql.NullString{String: slug, Valid: true})
		if err != nil {
			http.Error(w, "deployment not found", http.StatusNotFound)
			return
		}
		deploymentID = d.ID
	default:
		log.Printf("serve: project lookup failed for %q: %v", slug, err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	obj, err := storage.GetDeploymentFile(
		deploymentID.String(),
		r.URL.Path,
	)
	if err != nil {
		http.Error(
			w,
			"file not found",
			http.StatusNotFound,
		)
		return
	}

	defer obj.Body.Close()

	if obj.ContentType != nil {
		w.Header().Set(
			"Content-Type",
			*obj.ContentType,
		)
	}

	if viaProjectSlug {
		w.Header().Set("Cache-Control", "no-cache")
	}

	io.Copy(
		w,
		obj.Body,
	)
}
