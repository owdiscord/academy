package main

import (
	"log"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/owdiscord/academy/internal/config"
	"github.com/owdiscord/academy/internal/database"
	"github.com/owdiscord/academy/internal/handlers"
	"github.com/owdiscord/academy/internal/logger"
	"github.com/owdiscord/academy/internal/migrations"
	"github.com/owdiscord/academy/internal/periodic"
	"github.com/vinovest/sqlx"
)

func main() {
	logger.Config()

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

	// Cron-timed periodic jobs
	modmailDB, err := sqlx.Connect("mysql", config.ModmailDBURI)
	if err != nil {
		log.Fatalf("could not connect to modmail database: %v", err)
	}

	athenaDB, err := sqlx.Connect("mysql", config.AthenaDBURI)
	if err != nil {
		log.Fatalf("could not connect to athena database: %v", err)
	}

	jobs, err := periodic.NewManager(*config, athenaDB, modmailDB, db)
	if err != nil {
		log.Fatalf("cannot create job manager: %+v\n", err)
	}

	jobs.AddImportJob()
	jobs.Start()

	h := handlers.New(db, config, jobs)

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
	g.GET("/stats", h.Stats)
	g.GET("/import/:waveID", h.BackImport)
	g.GET("/avatar/:userID/:avatarHash", h.Avatar)

	// Serve static frontend
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:  "./frontend",
		Index: "index.html",
		HTML5: true, // fall back to index.html for unmatched paths
		Skipper: func(c *echo.Context) bool {
			return strings.HasPrefix(c.Path(), "/api")
		},
	}))

	e.Start(":1323")
}
