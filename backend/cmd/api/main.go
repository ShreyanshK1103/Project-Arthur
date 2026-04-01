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



func main () {

	conn, portString, err := config.ConnectDB()
	db := database.New(conn)

	apiCfg := handlers.Config{
		DB : db,
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
	v1Router.Post("/deployments", apiCfg.HandlerCreateDeployment)

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler : router,
		Addr: ":" + portString,
	}

	log.Printf("Server Starting on port %v", portString)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("PORT: ", portString)
}