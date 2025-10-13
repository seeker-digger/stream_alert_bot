package telegram

import (
	"gopkg.in/telebot.v4"
	"log"
	"main.go/db"
	"main.go/gokick"
	"os"
	"time"
)

type Bot struct {
	bot telebot.Bot
}

func Create(api gokick.ApiKick, db *db.DB) {
	preference := telebot.Settings{
		Token:  os.Getenv("TELEGRAM_BOT_API"),
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := telebot.NewBot(preference)
	if err != nil {
		log.Fatal(err)
	}
	l := initLoop()
	l.startAPILoop(db, api)
	l.startMailingLoop(db, b)
	b.Handle("/start", onStart(db))
	b.Handle("/add", onAdd(db, api))
	b.Handle("/remove", onRemove(db, api))
	b.Handle("/list", onList(db))
	b.Start()
}
