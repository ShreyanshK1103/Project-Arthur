package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/auth"
	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/config"
	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/database"
	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/handlers"
	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/oauth"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func main() {

	conn, portString, err := config.ConnectDB()
	oauth.InitGithubOAuth(
		os.Getenv("GITHUB_CLIENT_ID"),
		os.Getenv("GITHUB_CLIENT_SECRET"),
		os.Getenv("GITHUB_CALLBACK_URL"),
	)
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
	v1Router.Get("/auth/github/login", apiCfg.HandlerGithubLogin)
	v1Router.With(auth.MiddleWare).Get("/deployments/{id}", apiCfg.HandlerGetDeployment)
	v1Router.With(auth.MiddleWare).Get("/deployments/{id}/logs", apiCfg.HandlerGetDeploymentLogs)
	v1Router.With(auth.MiddleWare).Get("/projects/{id}/deployments", apiCfg.HandlerGetProjectDeployment)
	v1Router.With(auth.MiddleWare).Get("/auth/me", apiCfg.HandlerMe)

	v1Router.With(auth.MiddleWare).Post("/deployments", apiCfg.HandlerCreateDeployment)
	v1Router.With(auth.MiddleWare).Post("/projects", apiCfg.HandlerCreateProject)
	v1Router.With(auth.MiddleWare).Post("/projects/{id}/redeploy", apiCfg.HandlerRedeployProject)
	v1Router.Post("/github/webhook", apiCfg.HandlerGithubWebhook)
	v1Router.Post("/auth/register", apiCfg.HandlerRegister)
	v1Router.Post("/auth/login", apiCfg.HandlerLogin)
	v1Router.Post("/auth/refresh", apiCfg.HandlerRefresh)
	v1Router.Post("/auth/logout", apiCfg.HandlerLogout)

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
