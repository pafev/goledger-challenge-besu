package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"goledger-challenge-besu/configs/app"
	"goledger-challenge-besu/configs/besu"
	"goledger-challenge-besu/configs/db"
	"goledger-challenge-besu/configs/http"
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
	ethClient, err := besuConfig.New(&ctx)
	if err != nil {
		slog.Error("Erro connecting to besu node", "error", err)
	}
	defer ethClient.Close()
	slog.Info("Besu node connected succesfully")

	slog.Info("Starting the HTTP server...")
	http, err := httpConfig.New()
	if err != nil {
		slog.Error("Error trying setup the HTTP server", "error", err)
		os.Exit(1)
	}

	err = http.Route(&ctx, db, ethClient)
	if err != nil {
		slog.Error("Error building the HTTP routes", "error", err)
		os.Exit(1)
	}

	fmt.Println()
	slog.Info("\nHTTP Server running\n", "address", fmt.Sprintf("http://%s", http.Address), "port", http.Port)
	fmt.Println()
	err = http.Serve()
	if err != nil {
		slog.Error("Error starting the HTTP server", "error", err)
		os.Exit(1)
	}
}
