package app

import (
	"errors"
	"main.go/internal/config"
	"main.go/internal/db"
	l "main.go/internal/logger"
	"main.go/internal/telegram"
	"main.go/pkg/gokick"
	"syscall"
)

var State bool = false

func Run() error {

	if !(syscall.Geteuid() == 0) {
		l.Log.Error("Root's rights required!")
		return errors.New("root's rights required")
	}
	l.InitLogger()

	config.Init()

	l.Log.Println("Initializing bot...")

	api, err := gokick.GetAuthToken()
	if err != nil {
		l.Log.Panic(err)
	}
	l.Log.Println("Kick auth token successfully received")

	b := db.Init()
	State = true
	l.Log.Println(State)
	telegram.Create(api, &b)
	return nil
}
