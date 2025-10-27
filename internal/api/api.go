package api

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"main.go/internal/config"
	l "main.go/internal/logger"
	"main.go/pkg/gokick"
)

const DayInSeconds = int64(86400)

type Tokens struct {
	Kick gokick.Token
	// * Twitch  gotwitch.Token
	// * YouTube goyoutube.Token
}

func GetTokens() (Tokens, error) {
	var tokens Tokens

	if err := os.MkdirAll(config.GetDataPath(""), 0766); err != nil {
		return Tokens{}, fmt.Errorf("error creating data directory: %v", err)
	}
	fl, err := os.OpenFile(config.GetDataPath("token.json"), os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return Tokens{}, fmt.Errorf("error opening token.json: %v", err)
	}
	defer fl.Close()

	bytejson, err := io.ReadAll(fl)
	if err != nil {
		return Tokens{}, fmt.Errorf("error reading token.json: %v", err)
	}

	err = json.Unmarshal(bytejson, &tokens)
	if err != nil {
		l.Log.Warn("Error unmarshaling token.json: " + err.Error())
	} else {
		now := time.Now().Unix()
		exp := int64(tokens.Kick.ExpiresAt)

		if exp == 0 {
			l.Log.Warn("Kick token missing expiration, refreshing...")
		} else if exp <= now {
			l.Log.Warn("Kick token is expired, refreshing...")
		} else if exp <= now+3*DayInSeconds {
			l.Log.Warn("Kick token will expire in less than 3 days, refreshing...")
		} else if exp <= now+7*DayInSeconds {
			l.Log.Warn("Kick token will expire in less than 7 days.")
			return tokens, nil
		} else {
			return tokens, nil
		}
	}

	kickToken, err := gokick.GetAuthToken()
	if err != nil {
		return Tokens{}, fmt.Errorf("error getting Kick auth token: %v", err)
	}

	tokens.Kick = kickToken
	// tokens.Twitch = twitchToken
	// tokens.YouTube = youtubeToken

	fl.Truncate(0) // Clear file
	fl.Seek(0, 0)  // Move to beginning

	tokensBytes, err := json.MarshalIndent(tokens, "", "  ")
	if err != nil {
		return Tokens{}, fmt.Errorf("error marshaling tokens: %v", err)
	}

	_, err = fl.Write(tokensBytes)
	if err != nil {
		return Tokens{}, fmt.Errorf("error writing tokens to file: %v", err)
	}

	return tokens, nil
}
