package gokick

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

const channelapiurl = "https://api.kick.com/public/v1/channels"

var ErrUserDoesNotExist = errors.New("user does not exist")
var ErrInvalidURL = errors.New("invalid URL")

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresAt   int64  `json:"expires_in"`
}

func GetAuthToken() (Token, error) {
	var clientid = os.Getenv("KICK_CLIENT_ID")
	var clientsecret = os.Getenv("KICK_CLIENT_SECRET")

	data := url.Values{}
	data.Set("client_id", clientid)
	data.Set("client_secret", clientsecret)
	data.Set("grant_type", "client_credentials")

	resp, err := http.Post(
		"https://id.kick.com/oauth/token",
		"application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return Token{}, fmt.Errorf("error making POST request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return Token{}, fmt.Errorf("kick API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Token{}, fmt.Errorf("failed to read Response body: %w", err)
	}

	var token Token
	err = json.Unmarshal(body, &token)
	if err != nil {
		return Token{}, fmt.Errorf("failed to unmarshal Response: %w", err)
	}
	token.ExpiresAt = token.ExpiresAt + time.Now().Unix()

	return token, nil
}

type Response struct {
	Data    []ChannelData `json:"data"`
	Message string        `json:"message"`
}

type ChannelData struct {
	BannerPicture      string   `json:"banner_picture"`
	BroadcasterUserID  int      `json:"broadcaster_user_id"`
	Category           Category `json:"category"`
	ChannelDescription string   `json:"channel_description"`
	Slug               string   `json:"slug"`
	Stream             Stream   `json:"stream"`
	StreamTitle        string   `json:"stream_title"`
}

type Category struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Thumbnail string `json:"thumbnail"`
}

type Stream struct {
	IsLive      bool   `json:"is_live"`
	IsMature    bool   `json:"is_mature"`
	Key         string `json:"key"`
	Language    string `json:"language"`
	StartTime   string `json:"start_time"`
	Thumbnail   string `json:"thumbnail"`
	URL         string `json:"url"`
	ViewerCount int    `json:"viewer_count"`
}

// GetChannel getting all channel's info
func (a Token) GetChannel(slug []string) (Response, error) {
	var r Response
	queryParams := url.Values{}
	for _, i := range slug {
		queryParams.Add("slug", i)
	}

	fullURL := fmt.Sprintf("%s?%s", channelapiurl, queryParams.Encode())

	client := &http.Client{}
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return r, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.AccessToken))

	resp, err := client.Do(req)

	if err != nil {
		return r, fmt.Errorf("failed to make request: %w", err)
	} else if resp.StatusCode != http.StatusOK {
		return r, fmt.Errorf("kick API returned status %d: %s", resp.StatusCode, resp.Status)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return r, fmt.Errorf("failed to read Response body: %w", err)
	}
	err = json.Unmarshal(body, &r)

	if err != nil {
		return r, fmt.Errorf("failed to unmarshal Response: %w", err)
	}
	if r.Message == "Unauthorized" {
		return r, fmt.Errorf("kick API returned 401 Unauthorized")
	}
	return r, nil
}

func (a Token) GetSlugByURL(rawurl string) (string, error) {
	if !strings.HasPrefix(rawurl, "http://") && !strings.HasPrefix(rawurl, "https://") {
		rawurl = "http://" + rawurl
	}
	ur, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}
	hostname := ur.Hostname()
	if hostname != "kick.com" && ur.Hostname() != "www.kick.com" {
		return "", ErrInvalidURL
	}

	clearPath := path.Clean(ur.Path)
	if clearPath == "/" || clearPath == "" {
		return "", ErrInvalidURL
	} else if strings.Count(clearPath, "/") > 1 {
		return "", ErrInvalidURL
	}

	slug := path.Base(clearPath)

	c, err := a.GetChannel([]string{slug})
	if err != nil {
		return "", err
	}
	if len(c.Data) == 0 {
		return "", ErrUserDoesNotExist
	}
	return slug, nil

}
