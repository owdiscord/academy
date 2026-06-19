package main

import (
	"log"
	"os"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
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

	clientID := os.Getenv("DISCORD_CLIENT_ID")
	if clientID == "" {
		log.Fatalf("Missing env variable: DISCORD_CLIENT_ID")
	}

	secretID := os.Getenv("DISCORD_CLIENT_SECRET")
	if secretID == "" {
		log.Fatalf("Missing env variable: DISCORD_CLIENT_SECRET")
	}

	redirectURI := os.Getenv("DISCORD_REDIRECT_URI")
	if redirectURI == "" {
		log.Fatalf("Missing env variable: DISCORD_REDIRECT_URI")
	}

	// Run database migrations
	if err := migrations.Migrate(databaseURI); err != nil {
		log.Fatalf("cannot run migrations: %+v\n", err)
	}

	db, err := database.New(databaseURI)
	if err != nil {
		log.Fatalf("cannot connect to database: %+v\n", err)
	}

	h := handlers.New(db, clientID, secretID, redirectURI)

	// Start the web app
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.CORS("*"))

	e.GET("/api/auth/redirect", h.AuthRedirect)
	e.GET("/api/auth/callback", h.AuthCallback)

	g := e.Group("/api")
	g.Use(h.RequireAuth)
	g.GET("/auth/me", h.Me)
	g.GET("/wave", h.Wave)
	g.GET("/threads", h.Threads)
	g.GET("/threads/:id", h.Thread)
	g.GET("/cases", h.Cases)
	g.GET("/cases/:id", h.Case)
	g.GET("/questions", h.Questions)

	e.Start(":1323")
}
