package main

import (
	"context"
	"log"
	"time"

	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/config"
	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/deployment"
	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/database"
)


func main () {
	conn, _, err := config.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	db := database.New(conn)

	// This sets the building failed job to queued again so it can be retried
	err = db.ResetBuildingDeployments(
		context.Background(),
	)
	if err != nil {
		log.Fatal(err)
	}

	for {
		job, err := db.GetNextDeployment(context.Background())
		if err != nil {
			time.Sleep(3 * time.Second)
			continue
		}

		err = deployment.ProcessDeployment(job, db)
		if err != nil {
			log.Printf("Processing failed: %v", err)

			db.MarkDeploymentFailed(
				context.Background(),
				job.ID,
			)

			continue
		}

		log.Printf("Successfully processed deployment %s", job.ID)

	}
}
