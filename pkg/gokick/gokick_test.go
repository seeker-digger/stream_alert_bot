package gokick

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"main.go/internal/config"
)

func TestGetAuthToken(t *testing.T) {
	config.Init()
	a, err := GetAuthToken()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(a)
}

func TestAuthToken_GetChannel(t *testing.T) {
	a, err := GetAuthToken()
	if err != nil {
		t.Error(err)
	}
	var lst = []string{"ppfdgbdfgdf"}
	c, err := a.GetChannel(lst)
	if err != nil {
		t.Error(err)
	}

	jsonD, err := json.MarshalIndent(c, " ", "")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(string(jsonD))
}

func TestAuthToken_GetSlugByURL(t *testing.T) {
	a, err := GetAuthToken()
	if err != nil {
		t.Error(err)
	}
	url := "https://kick.com/ppfdgbdfgdf"
	s, err := a.GetSlugByURL(url)
	if errors.Is(err, ErrInvalidURL) {
		t.Error("Invalid URL")
		return
	} else if errors.Is(err, ErrUserDoesNotExist) {
		t.Log("User Does Not Exist")
		return
	} else if err != nil {
		t.Error(err)
		return
	}
	t.Log(s)

}
