package api

import (
	"testing"

	"main.go/internal/config"
	l "main.go/internal/logger"
)

func TestGetTokens(t *testing.T) {
	config.Init()
	tok, err := GetTokens()
	if err != nil {
		t.Error(err)
	}
	l.Log.Warn(tok)
}
