package telegram

import (
	"context"
	"gopkg.in/telebot.v4"
	"main.go/internal/db"
	l "main.go/internal/logger"
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
		l.Log.Fatal(err)
	}
	loop := initLoop()
	loop.startAPILoop(db, api)
	loop.startMailingLoop(db, b)
	b.Handle("/start", onStart(db))
	b.Handle("/add", onAdd(db, api))
	b.Handle("/remove", onRemove(db, api))
	b.Handle("/list", onList(db))
	go func() {
		b.Start()
	}()
	<-ctx.Done()

	l.Log.Warn("Shutting down...")
	b.Stop()

	l.Log.Println("Bot successfully stopped")
}
