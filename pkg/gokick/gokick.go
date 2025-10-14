package gokick

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"main.go/internal/config"
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

type ApiKick interface {
	GetChannel(slug []string) (Response, error)
	GetSlugByURL(slug string) (string, error)
}

type authToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	ChangeIn    int64
}

type appSecrets struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func GetAuthToken() (ApiKick, error) {
	var clientid = os.Getenv("KICK_CLIENT_ID")
	var clientsecret = os.Getenv("KICK_CLIENT_SECRET")

	var authToken authToken

	fl, err := os.OpenFile(config.GetDataPath("token.json"), os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return nil, fmt.Errorf("error opening token.json: %v", err)
	}
	defer fl.Close()

	bytejson, err := io.ReadAll(fl)
	if err != nil {
		return nil, fmt.Errorf("error reading token.json: %v", err)
	}

	err = json.Unmarshal(bytejson, &authToken)
	if authToken.ChangeIn <= time.Now().Unix() {
		err = errors.New("token is expired")
	} else if err == nil {
		log.Println("Kick auth token successfully read")
		return &authToken, nil
	} else {
		fmt.Println("Token is expired")
	}

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
		return nil, fmt.Errorf("error making POST request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("kick API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read Response body: %w", err)
	}

	err = json.Unmarshal(body, &authToken)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal Response: %w", err)
	}
	authToken.ChangeIn = time.Now().Unix() + int64(authToken.ExpiresIn)*95/100 // After 57 days

	fl, err = os.Create(config.GetDataPath("token.json"))
	if err != nil {
		return nil, fmt.Errorf("error creating token.json: %v", err)
	}
	defer fl.Close()

	aTF, err := json.MarshalIndent(authToken, "", "  ")

	_, err = fl.Write(aTF)
	if err != nil {
		return nil, fmt.Errorf("error writing token.json: %v", err)
	}

	log.Print("Kick api successfully initialized")
	return authToken, nil
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
func (a authToken) GetChannel(slug []string) (Response, error) {
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

func (a authToken) GetSlugByURL(rawurl string) (string, error) {
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
