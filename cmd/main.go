package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"goledger-challenge-besu/configs/app"
	"goledger-challenge-besu/configs/besu"
	"goledger-challenge-besu/configs/db"
)

func main() {
	slog.Info("Starting the application...")
	app, err := appConfig.New()
	if err != nil {
		slog.Error("Error loading env variables", "error", err)
		os.Exit(1)
	}
	slog.Info("Application started", "app", app.Name, "env", app.Env)

	ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelCtx()

	slog.Info("Connecting to database...")
	db, err := dbConfig.New(&ctx)
	if err != nil {
		slog.Error("Error initializing database connection", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	slog.Info("Database connected succesfully")

	slog.Info("Migrating database...")
	err = db.Migrate()
	if err != nil {
		slog.Error("Error migrating database", "error", err)
		os.Exit(1)
	}
	slog.Info("Database migrated succesfully")

	slog.Info("Connecting to Besu node...")
	client, err := besuConfig.New(&ctx)
	if err != nil {
		slog.Error("Erro connecting to besu node", "error", err)
	}
	defer client.Close()
	slog.Info("Besu node connected succesfully")
}
