// Package handlers contains our HTTP handlers and central context struct
package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/owdiscord/academy/internal/cache"
	"github.com/owdiscord/academy/internal/database"
)

type Handlers struct {
	db           *database.DB
	clientID     string
	clientSecret string
	redirectURI  string
	sessionCache *cache.Cache[string, database.Session]
}

func New(db *database.DB, clientID string, clientSecret string, redirectURI string) Handlers {
	return Handlers{
		db:           db,
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,
		sessionCache: cache.New[string, database.Session](time.Minute * 3),
	}
}

// -- Endpoints ----

func (h *Handlers) AuthRedirect(c *echo.Context) error {
	url := "https://discord.com/oauth2/authorize?client_id=" + h.clientID + "&response_type=code&redirect_uri=" + url.QueryEscape(h.redirectURI) + "&scope=identify"
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *Handlers) AuthCallback(c *echo.Context) error {
	code := c.QueryParam("code")
	if len(code) < 10 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "no code provided"})
	}

	// Exchange code for token
	params := url.Values{}
	params.Set("grant_type", "authorization_code")
	params.Set("client_id", h.clientID)
	params.Set("client_secret", h.clientSecret)
	params.Set("code", code)
	params.Set("scope", "identify")
	params.Set("redirect_uri", h.redirectURI)

	authRes, err := http.Post(
		"https://discord.com/api/v10/oauth2/token",
		"application/x-www-form-urlencoded",
		strings.NewReader(params.Encode()),
	)
	if err != nil || authRes.StatusCode >= 400 {
		body, _ := io.ReadAll(authRes.Body)
		c.Logger().Error("discord token error", "http_body", body)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "could not retrieve your discord token"})
	}
	defer authRes.Body.Close()

	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(authRes.Body).Decode(&tokenResp); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not parse discord token response"})
	}

	// Fetch Discord user
	req, _ := http.NewRequest("GET", "https://discord.com/api/v10/users/@me", nil)
	req.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)
	meRes, err := http.DefaultClient.Do(req)
	if err != nil || meRes.StatusCode >= 400 {
		body, _ := io.ReadAll(meRes.Body)
		c.Logger().Error("discord @me error", "http_body", body)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not retrieve your discord account details, despite getting a token successfully."})
	}
	defer meRes.Body.Close()

	var discordUser struct {
		ID         string `json:"id"`
		GlobalName string `json:"global_name"`
	}
	if err := json.NewDecoder(meRes.Body).Decode(&discordUser); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not parse discord user response"})
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

	return c.Redirect(http.StatusFound, "/academy")
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
	return c.JSON(http.StatusOK, []string{})
}

func (h *Handlers) Thread(c *echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{})
}

func (h *Handlers) Cases(c *echo.Context) error {
	return c.JSON(http.StatusOK, []string{})
}

func (h *Handlers) Case(c *echo.Context) error {
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
