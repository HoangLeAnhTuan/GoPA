package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"gopa/pkg/config"
	"gopa/pkg/database"
	"gopa/pkg/logger"
)

func main() {
	direction := flag.String("direction", "up", "migration direction: up or down")
	steps := flag.Int("steps", 1, "number of migrations to apply")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	log := logger.New(cfg.AppEnv, cfg.LogLevel)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	db, err := database.Open(ctx, cfg.PostgresDSN)
	if err != nil {
		log.Error("open database for migrations", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := database.Migrate(db, database.MigrationDirection(*direction), *steps); err != nil {
		log.Error("migration failed", "direction", *direction, "error", err)
		os.Exit(1)
	}
	log.Info("migration completed", "direction", *direction, "steps", *steps)
}
