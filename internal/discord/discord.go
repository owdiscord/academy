// Package discord contains our basic Discord API functions, for getting
// tokens, getting the current user, etc.
package discord

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/owdiscord/academy/internal/config"
)

type DiscordUser struct {
	ID         string `json:"id"`
	Username   string `json:"username"`
	GlobalName string `json:"global_name"`
	Avatar     string `json:"avatar"`
}

// GetUser will fetch a Discord user
func GetUser(token string, userID string) (*DiscordUser, error) {
	req, _ := http.NewRequest("GET", "https://discord.com/api/v10/users/"+userID, nil)
	req.Header.Set("Authorization", token)
	meRes, err := http.DefaultClient.Do(req)
	if err != nil || meRes.StatusCode >= 400 {
		body, _ := io.ReadAll(meRes.Body)
		return nil, fmt.Errorf("could not get Discord user: %s", body)
	}
	defer meRes.Body.Close()

	var discordUser DiscordUser
	if err := json.NewDecoder(meRes.Body).Decode(&discordUser); err != nil {
		return nil, errors.New("could not parse discord user response")
	}

	return &discordUser, nil
}

// GetAccessToken will exchange an OAuth2 token code for an access token
func GetAccessToken(config config.Config, code string) (string, error) {
	params := url.Values{}
	params.Set("grant_type", "authorization_code")
	params.Set("client_id", config.ClientID)
	params.Set("client_secret", config.ClientSecret)
	params.Set("code", code)
	params.Set("scope", "identify")
	params.Set("redirect_uri", config.RedirectURI)

	authRes, err := http.Post(
		"https://discord.com/api/v10/oauth2/token",
		"application/x-www-form-urlencoded",
		strings.NewReader(params.Encode()),
	)
	if err != nil || authRes.StatusCode >= 400 {
		body, _ := io.ReadAll(authRes.Body)
		return "", fmt.Errorf("could not retrieve discord token: %s", body)
	}
	defer authRes.Body.Close()

	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(authRes.Body).Decode(&tokenResp); err != nil {
		return "", errors.New("could not parse discord token response")
	}

	return tokenResp.AccessToken, nil
}

func DownloadAvatar(user DiscordUser, outPath string) error {
	resp, err := http.Get("https://cdn.discordapp.com/avatars/" + user.ID + "/" + user.Avatar + ".png?size=256")
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to download image: status " + resp.Status)
	}

	out, err := os.Create(outPath + user.Avatar + ".png")
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
