package config

import (
	"fmt"
	"os"
	"testing"
)

func TestGetDataPath(t *testing.T) {
	println(GetDataPath(".env"))

}

func TestInitData(t *testing.T) {
	InitData()
	a := os.Getenv("TEST")
	fmt.Println(a)
}
