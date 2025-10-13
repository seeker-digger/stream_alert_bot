package main

import (
	"main.go/config"
	db "main.go/db"
	"main.go/gokick"
	"main.go/telegram"
)

func main() {
	config.InitData()

	api, err := gokick.GetAuthToken()
	if err != nil {
		panic(err)
	}

	b := db.Init()
	telegram.Create(api, &b)
	
}
