package handlers

import (
	"database/sql"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/storage"
)

func (cfg *Config) HandlerServeByDomain(w http.ResponseWriter, r *http.Request) {
	// THIS FUNCTION HANDLES THE SERVING LOGIC INSTEAD OF STATIC HANDLING
	// {deploymentID}.localhost:8080 is the root in this logic cause of the file rendering issues prior to the
	// localhost:8080/v1/deployments/{deploymentID} logic was flawed cause the assests was being searched in the dir localhost:8080
	// But in reality the assests were being served in the given path, due to this subdomain approach was considered as an viable
	// architectural decision
	host := r.Host

	host = strings.Split(host, ":")[0]

	slug := strings.Split(host, ".")[0]

	deployment, err := cfg.DB.GetDeploymentByPrefix(
		r.Context(),
		sql.NullString{
			String: slug,
			Valid:  true,
		},
	)
	if err != nil {
		http.Error(
			w,
			"deployment not found",
			http.StatusNotFound,
		)
		return
	}

	// deploymentPath := filepath.Join(
	// 	"./storage/deployments",
	// 	deployment.ID.String(),
	// )

	// _, err = os.Stat(deploymentPath)
	// if err != nil {
	// 	http.Error(
	// 		w,
	// 		"deployment not found",
	// 		http.StatusNotFound,
	// 	)
	// 	return
	// }

	// fs := http.FileServer(
	// 	http.Dir(deploymentPath),
	// )

	// fs.ServeHTTP(w, r)

	obj, err := storage.GetDeploymentFile(
		deployment.ID.String(),
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

	if obj.ContentType != nil {
		w.Header().Set(
			"Content-Type",
			*obj.ContentType,
		)
	}

	log.Printf(
		"Sending Content-Type: %s",
		w.Header().Get("Content-Type"),
	)

	io.Copy(
		w,
		obj.Body,
	)
}
