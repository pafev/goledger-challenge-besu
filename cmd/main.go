package main

import (
	"goledger-challenge-besu/configs/app"
	"log/slog"
	"os"
)

func main() {
	slog.Info("Starting the application...")
	app, err := appConfig.New()
	if err != nil {
		slog.Error("Error loading env variables", "error", err)
		os.Exit(1)
	}
	slog.Info("Application started", "app", app.Name, "env", app.Env)

}
