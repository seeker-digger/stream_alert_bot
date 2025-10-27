package telegram

import (
	"errors"
	"fmt"
	"main.go/internal/api"
	"time"

	"gopkg.in/telebot.v4"
	"main.go/internal/db"
	l "main.go/internal/logger"
	"main.go/pkg/gokick"
)

var needUpdate chan bool

const delay = 10 * time.Second

type Loop struct {
	targets map[string]int64
	onlSlug chan gokick.ChannelData
	toggle  bool
}

func initLoop() Loop {
	targets := make(map[string]int64)
	needUpdate = make(chan bool, 2)
	onlSlug := make(chan gokick.ChannelData, 1)
	return Loop{targets: targets, onlSlug: onlSlug, toggle: false}
}

// UGLY!!
func (loop *Loop) startAPILoop(b *db.DB, tokens api.Tokens) {
	loop.toggle = true
	go func() {
		ticker := time.NewTicker(delay)
		slugs, err := b.GetAllUniqueValues()
		if err != nil {
			l.Log.Println("error getting slugs: ", err)
		}
		for {
			select {
			case <-needUpdate:
				slugs, err = b.GetAllUniqueValues()
				if err != nil {
					l.Log.Println("error getting slugs: ", err)
				}
			case <-ticker.C:
				chunks := chunkSlice(slugs, 50)
				for _, s := range chunks {
					r, err := tokens.Kick.GetChannel(s)
					if err != nil {
						l.Log.Println("error getting channels: ", err.Error())
					}
					for _, i := range r.Data {
						if i.Stream.IsLive == true {
							if time.Since(time.Unix(loop.targets[i.Slug], 0)) > 10*time.Minute {
								loop.onlSlug <- i
							}
							loop.targets[i.Slug] = time.Now().Unix()
						}
					}
				}
			}
		}
	}()
}

func (loop *Loop) pauseAPILoop() {
	loop.toggle = false // <-- THIS IS SHIT!
}

// VERY UGLY!!
func (loop *Loop) startMailingLoop(b *db.DB, bot *telebot.Bot) {
	go func() {
		for {
			select {
			case s := <-loop.onlSlug:
				ids, err := b.GetAllIdsByValueKick(s.Slug)
				if err != nil {
					l.Log.Println("error getting slugs: ", err)
				}
				for _, i := range ids {
					text := fmt.Sprintf("*Канал*: [%s](kick.com\\/%s) запустил стрим\\!\n>%s", escapeMarkdownV2Text(s.Slug), s.Slug, escapeMarkdownV2Text(s.StreamTitle))
					_, err := bot.Send(&telebot.Chat{ID: i}, text, &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2})
					if errors.Is(err, telebot.ErrUserIsDeactivated) || errors.Is(err, telebot.ErrBlockedByUser) || errors.Is(err, telebot.ErrNotStartedByUser) {
						err = b.RemoveUser(i)
						if err != nil {
							l.Log.Println("error removing user: ", err)
						}
						continue
					} else if err != nil {
						l.Log.Println("error sending message: ", err)
					}
				}
			}
		}
	}()
}
