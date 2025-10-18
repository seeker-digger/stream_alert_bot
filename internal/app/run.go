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
	telegram.Create(api, &b)
	return nil
}
