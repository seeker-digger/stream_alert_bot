package config

import (
	"fmt"
	"os"
	"testing"
)

func TestGetDataPath(t *testing.T) {
	println(GetDataPath("qwerty"))

}

func TestInitData(t *testing.T) {
	Init()
	a := os.Getenv("KICK_CLIENT_ID")
	fmt.Println(a)
}
