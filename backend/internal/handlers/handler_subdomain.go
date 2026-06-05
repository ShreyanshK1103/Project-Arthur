package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func (cfg *Config) HandlerServeByDomain (w http.ResponseWriter, r *http.Request) {
	// THIS FUNCTION HANDLES THE SERVING LOGIC INSTEAD OF STATIC HANDLING 
	// {deploymentID}.localhost:8080 is the root in this logic cause of the file rendering issues prior to the 
	// localhost:8080/v1/deployments/{deploymentID} logic was flawed cause the assests was being searched in the dir localhost:8080
	// But in reality the assests were being served in the given path, due to this subdomain approach was considered as an viable 
	// architectural decision 
	host := r.Host 

	host = strings.Split(host, ":")[0]

	deploymentID := strings.Split(host, ".")[0]

	deploymentPath := filepath.Join(
		"./storage/deployments",
		deploymentID,
	)

	_, err := os.Stat(deploymentPath)
	if err != nil {
		http.Error(
			w,
			"deployment not found",
			http.StatusNotFound,
		)
		return 
	}

	fs := http.FileServer(
		http.Dir(deploymentPath),
	)

	fs.ServeHTTP(w, r)
}