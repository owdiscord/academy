package main

import (
	"log"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/owdiscord/academy/internal/config"
	"github.com/owdiscord/academy/internal/database"
	"github.com/owdiscord/academy/internal/handlers"
	"github.com/owdiscord/academy/internal/migrations"
)

func main() {
	config, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Run database migrations
	if err := migrations.Migrate(config.DatabaseURI); err != nil {
		log.Fatalf("cannot run migrations: %+v\n", err)
	}

	db, err := database.New(config.DatabaseURI)
	if err != nil {
		log.Fatalf("cannot connect to database: %+v\n", err)
	}

	h := handlers.New(db, config)

	// Start the web app
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.CORS("*"))

	e.GET("/api/auth/redirect", h.AuthRedirect)
	e.GET("/api/auth/callback", h.AuthCallback)

	g := e.Group("/api")
	g.Use(h.RequireAuth)
	g.Any("/auth/logout", h.AuthLogout)
	g.GET("/auth/me", h.Me)
	g.GET("/wave", h.Wave)
	g.GET("/threads", h.Threads)
	g.GET("/threads/:id", h.Thread)
	g.GET("/cases", h.Cases)
	g.GET("/cases/:id", h.Case)
	g.GET("/issues", h.GetIssues)
	g.GET("/issues/:id", h.GetIssue)
	g.PUT("/issues/id", h.CreateIssue)
	g.GET("/questions", h.Questions)
	g.GET("/avatar/:userID", h.Avatar)

	e.Start(":1323")
}
