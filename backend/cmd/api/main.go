package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/config"
	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/database"
	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func main() {

	conn, portString, err := config.ConnectDB()
	db := database.New(conn)

	apiCfg := handlers.Config{
		DB: db,
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	v1Router.Get("/healthz", handlers.HandlerReadiness)
	v1Router.Get("/deployments/{id}", apiCfg.HandlerGetDeployment)
	v1Router.Get("/deployments/{id}/logs", apiCfg.HandlerGetDeploymentLogs)
	v1Router.Get("/projects/{id}/deployments", apiCfg.HandlerGetProjectDeployment)

	v1Router.Post("/deployments", apiCfg.HandlerCreateDeployment)
	v1Router.Post("/projects", apiCfg.HandlerCreateProject)
	v1Router.Post("/projects/{id}/redeploy", apiCfg.HandlerRedeployProject)
	v1Router.Post("/github/webhook", apiCfg.HandlerGithubWebhook)
	v1Router.Post("/auth/register", apiCfg.HandlerRegister)
	v1Router.Post("/auth/login", apiCfg.HandlerLogin)

	router.Mount("/v1", v1Router)

	// THIS ROUTE SERVES THE WEBSITE USING THE SUBDOMAIN APPROACH
	router.Handle("/*", http.HandlerFunc(apiCfg.HandlerServeByDomain))

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server Starting on port %v", portString)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("PORT: ", portString)
}
