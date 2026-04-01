package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/config"
	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/database"
)

func main () {
	conn, _, err := config.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	db := database.New(conn)

	for {
		job, err := db.GetNextDeployment(context.Background())
		if err != nil {
			time.Sleep(3 * time.Second)
			continue
		}

		log.Println("Picked Job: ", job.ID)
		time.Sleep(5 * time.Second)

		err = db.MarkDeploymentSuccess(context.Background(), database.MarkDeploymentSuccessParams{
			ID : job.ID,
			Url : sql.NullString{
				String : "https://example.com",
				Valid: true,
			},
		})
		if err != nil {
			log.Printf("Failed to mark Success: %v", err)
		}

	}
}
