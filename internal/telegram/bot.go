package telegram

import (
	"context"
	"gopkg.in/telebot.v4"
	"log"
	"main.go/internal/db"
	"main.go/pkg/gokick"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Bot struct {
	bot telebot.Bot
}

func Create(api gokick.ApiKick, db *db.DB) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

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
	go func() {
		b.Start()
	}()
	<-ctx.Done()

	log.Println("Shutting down...")
	b.Stop()

	log.Println("Bot successfully stopped")
}
