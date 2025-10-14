package main

import (
	"log"
	"main.go/internal/config"
	"main.go/internal/db"
	"main.go/internal/telegram"
	"main.go/pkg/gokick"
	"syscall"
)

func main() {
	if !(syscall.Geteuid() == 0) {
		log.Fatal("Root's rights are required!")
	}

	config.InitData()

	api, err := gokick.GetAuthToken()
	if err != nil {
		panic(err)
	}

	b := db.Init()
	telegram.Create(api, &b)

}
