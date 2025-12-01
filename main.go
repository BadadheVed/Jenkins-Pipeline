package main

import (
	"jenkins-prac/routes"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		slog.Info("No .env file found; continuing with environment variables")
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
		slog.Info("PORT not set in environment; using default", "port", port)
		return
	}
	r := routes.SetupRouter()

	err := r.Run(":" + port)
	if err != nil {
		slog.Error("Failed to run server", "error", err)
		return
	}

}
