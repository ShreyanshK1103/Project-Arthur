package deployment

import (
	"context"
	"log"
	"time"

	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/database"
	"github.com/ShreyanshK1103/Project-Arthur/backend/internal/logs"
)

const (
	reaperInterval = time.Minute
	// Must stay longer than the deploymentTimeout, otherwise we would fail
	// builds that are still legitimately running.
	reaperMargin = 5 * time.Minute
)

// This blocks until ctx is cancelled,so run it in a goroutine
func StartReaper(ctx context.Context, db *database.Queries) {
	ticker := time.NewTicker(reaperInterval)
	defer ticker.Stop()

	reap(ctx, db) // once at startup, then every tick

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			reap(ctx, db)
		}
	}
}

func reap(ctx context.Context, db *database.Queries) {
	minutes := int32((deploymentTimeout + reaperMargin) / time.Minute)

	ids, err := db.FailStaleDeployments(ctx, minutes)
	if err != nil {
		log.Printf("reaper: failed to reap stale deployments: %v", err)
		return
	}

	for _, id := range ids {
		log.Printf("reaper: marked deployment %s as failed (stuck in building)", id)
		logs.AddLog(db, id, "Deployment failed: worker stopped responding while building")
	}
}