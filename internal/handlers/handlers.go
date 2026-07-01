// Package handlers contains our HTTP handlers and central context struct
package handlers

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/owdiscord/academy/internal/cache"
	"github.com/owdiscord/academy/internal/config"
	"github.com/owdiscord/academy/internal/database"
	"github.com/owdiscord/academy/internal/discord"
	"github.com/owdiscord/academy/internal/periodic"
)

type Handlers struct {
	db           *database.DB
	config       *config.Config
	sessionCache *cache.Cache[string, database.Session]
	jobManager   *periodic.Manager
}

func New(db *database.DB, config *config.Config, jobManager *periodic.Manager) Handlers {
	return Handlers{
		db:           db,
		config:       config,
		sessionCache: cache.New[string, database.Session](time.Minute * 3),
		jobManager:   jobManager,
	}
}

// -- Endpoints ----

func (h *Handlers) AuthRedirect(c *echo.Context) error {
	url := "https://discord.com/oauth2/authorize?client_id=" + h.config.ClientID + "&response_type=code&redirect_uri=" + url.QueryEscape(h.config.RedirectURI) + "&scope=identify"
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *Handlers) AuthCallback(c *echo.Context) error {
	code := c.QueryParam("code")
	if len(code) < 10 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "no code provided"})
	}

	accessToken, err := discord.GetAccessToken(*h.config, code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not parse discord token response"})
	}

	discordUser, err := discord.GetUser("Bearer "+accessToken, "@me")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Look up user in DB
	user, err := h.db.LatestUserForDiscordID(c.Request().Context(), discordUser.ID)
	if err != nil || user == nil {
		c.Logger().Error("unknown user attempted to access", "discord_id", discordUser.ID, "name", discordUser.GlobalName)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "You are not authorized to access this page."})
	}

	// Create session
	token, expiry, err := h.db.CreateSession(c.Request().Context(), user.ID, user.WaveID)
	if err != nil {
		c.Logger().Error("failed to create session", "discord_id", discordUser.ID, "user_id", user.ID, "wave_id", user.WaveID, "err", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Failed to create this session."})
	}

	// Set cookie
	cookie := &http.Cookie{
		Name:     "academy_session",
		Value:    token,
		HttpOnly: true,
		Secure:   os.Getenv("DEV_SERVER") == "",
		SameSite: http.SameSiteLaxMode,
		Expires:  expiry,
		Path:     "/",
	}
	c.SetCookie(cookie)

	go func() {
		if err := h.db.UpdateUserDetails(context.Background(), *discordUser); err != nil {
			c.Logger().Error("could not update user details", "userid", discordUser.ID, "db_err", err)
		}

		if err := discord.DownloadAvatar(*discordUser, "./avatars/"); err != nil {
			c.Logger().Error("could not download avatar", "userid", discordUser.ID, "avatar_hash", discordUser.Avatar, "io_err", err)
		}
	}()

	return c.Redirect(http.StatusFound, "/")
}

func (h *Handlers) AuthLogout(c *echo.Context) error {
	cookie := &http.Cookie{
		Name:     "academy_session",
		Value:    "",
		HttpOnly: true,
		Secure:   os.Getenv("DEV_SERVER") == "",
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now(),
		Path:     "/",
	}
	c.SetCookie(cookie)

	return c.Redirect(http.StatusFound, "/")
}

func (h *Handlers) Me(c *echo.Context) error {
	session, exists := c.Get("session_value").(*database.Session)
	if !exists {
		c.Logger().Error("could not get session on authenticated route", "session", session)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "cannot retrieve your session"})
	}

	staff, err := h.db.GetStaffDetails(c.Request().Context(), session.UserID, session.WaveID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "cannot retrieve user details"})
	}

	return c.JSON(http.StatusOK, staff)
}

func (h *Handlers) Wave(c *echo.Context) error {
	waveID := c.Get("session_value").(*database.Session).WaveID
	wave, err := h.db.GetWaveByID(c.Request().Context(), waveID)
	if err != nil {
		c.Logger().Error("could not get wave", "err", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not get wave from database"})
	}

	trainees, err := h.db.GetWaveTrainees(c.Request().Context(), waveID)
	if err != nil {
		c.Logger().Error("could not get trainees", "err", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not get wave from database"})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"id":         wave.ID,
		"state":      wave.State,
		"begin_at":   wave.BeginAt.Unix(),
		"close_at":   wave.CloseAt.Unix(),
		"created_at": wave.CreatedAt.Unix(),
		"trainees":   trainees,
	})
}

func (h *Handlers) Threads(c *echo.Context) error {
	threads, err := h.db.GetAllThreads(c.Request().Context())
	if err != nil {
		c.Logger().Error("could not get threads", "db", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not retrieve threads"})
	}

	return c.JSON(http.StatusOK, threads)
}

func (h *Handlers) Thread(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "given ID was not a valid UUID"})
	}

	thread, err := h.db.GetThreadByID(c.Request().Context(), database.BinaryUUID(id))
	if err != nil {
		c.Logger().Error("could not get thread", "id", id, "db", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not retrieve thread"})
	}

	return c.JSON(http.StatusOK, thread)
}

func (h *Handlers) Cases(c *echo.Context) error {
	cases, err := h.db.GetAllCases(c.Request().Context())
	if err != nil {
		c.Logger().Error("could not get cases", "db", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not retrieve cases"})
	}

	return c.JSON(http.StatusOK, cases)
}

func (h *Handlers) Case(c *echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "given ID was not a valid integer"})
	}

	modCase, err := h.db.GetCaseByID(c.Request().Context(), id)
	if err != nil {
		c.Logger().Error("could not get case", "id", id, "db", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not retrieve case"})
	}

	return c.JSON(http.StatusOK, modCase)
}

func (h *Handlers) GetIssues(c *echo.Context) error {
	sess := c.Get("session_value").(*database.Session)
	role := sess.Role
	waveID := sess.WaveID

	if role == "admin" {
		issues, err := h.db.GetFullIssues(c.Request().Context(), waveID)
		if err != nil {
			c.Logger().Error("could not get issues", "db_err", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not get issues"})
		}

		return c.JSON(http.StatusOK, issues)
	}

	return c.JSON(http.StatusOK, []string{})
}

func (h *Handlers) GetIssue(c *echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{})
}

func (h *Handlers) CreateIssue(c *echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{})
}

func (h *Handlers) UpdateIssue(c *echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{})
}

func (h *Handlers) Questions(c *echo.Context) error {
	questions, err := h.db.GetQuestions(c.Request().Context())
	if err != nil {
		c.Logger().Error("could not get questions", "err", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not retrieve questions from database"})
	}

	return c.JSON(http.StatusOK, questions)
}

func (h *Handlers) Stats(c *echo.Context) error {
	sess := c.Get("session_value").(*database.Session)
	waveID := sess.WaveID

	stats, err := h.db.GetStatsOverview(c.Request().Context(), waveID)
	if err != nil {
		c.Logger().Error("could not get stats", "err", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not retrieve or calculate stats"})
	}

	return c.JSON(http.StatusOK, stats)
}

func (h *Handlers) BackImport(c *echo.Context) error {
	sess := c.Get("session_value").(*database.Session)
	if sess.Role != "admin" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "only admin users can trigger a back import!"})
	}

	waveID, err := strconv.Atoi(c.Param("waveID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "provided wave ID was not an integer"})
	}

	jobID, err := h.jobManager.TriggerImport(context.Background(), time.Date(2026, time.January, 1, 1, 1, 1, 1, time.UTC), nil, &waveID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not run job: " + err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{"message": "job triggered to run", "wave": waveID, "job_id": jobID.String()})
}

func (h *Handlers) Avatar(c *echo.Context) error {
	userID := c.Param("userID")[:len(c.Param("userID"))-len(filepath.Ext(c.Param("userID")))]
	if _, err := os.Stat("./avatars/" + userID + ".png"); err != nil {
		return c.File("./avatars/default.png")
	}

	return c.File("./avatars/" + userID + ".png")
}

// -- Middleware ----

func (h *Handlers) RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		cookie, err := c.Cookie("academy_session")
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "authentication expired or not provided"})
		}

		cached, exists := h.sessionCache.Get(cookie.Value)
		if exists {
			c.Set("session_id", cookie.Value)
			c.Set("session_value", &cached)
			return next(c)
		}

		inDB, err := h.db.GetSessionByToken(c.Request().Context(), cookie.Value)
		if err == nil {
			h.sessionCache.Set(cookie.Value, *inDB)

			c.Set("session_id", cookie.Value)
			c.Set("session_value", inDB)
			return next(c)
		}

		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "authentication expired or not provided"})
	}
}
