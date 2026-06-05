package handlers

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi"
)

func (cfg *Config) HandlerServeDeployment(w http.ResponseWriter, r *http.Request) {
	deploymentID := chi.URLParam(r, "id")

	if deploymentID == "" {
		http.Error(w, "Missing Deployment id", http.StatusBadRequest)
		return
	}

	prefix := "/v1/deployments/" + deploymentID

	if r.URL.Path == prefix || r.URL.Path == prefix+"/" {
		r.URL.Path = "/"
	}

	deploymentPath := filepath.Join(
		"./storage/deployments",
		deploymentID,
	)

	log.Printf("Serving path: %s", r.URL.Path)

	fs := http.FileServer(http.Dir(deploymentPath))

	http.StripPrefix(
		"/v1/deployments/"+deploymentID,
		fs,
	).ServeHTTP(w, r)
}