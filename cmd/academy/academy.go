package main

import (
	"log"
	"os"

	"github.com/labstack/echo/v5"
	"github.com/owdiscord/academy/internal/database"
	"github.com/owdiscord/academy/internal/handlers"
	"github.com/owdiscord/academy/internal/migrations"
)

func main() {
	// Get necessary env variables
	databaseURI := os.Getenv("DATABASE_URI")
	if databaseURI == "" {
		log.Fatalf("Missing env variable: DATABASE_URI")
	}

	// Run database migrations
	if err := migrations.Migrate(databaseURI); err != nil {
		log.Fatalf("cannot run migrations: %+v\n", err)
	}

	db, err := database.New(databaseURI)
	if err != nil {
		log.Fatalf("cannot run migrations: %+v\n", err)
	}

	h := handlers.New(db)

	// Start the web app
	e := echo.New()

	e.Start(":1323")
}
