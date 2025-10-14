package telegram

import (
	"errors"
	"fmt"
	"gopkg.in/telebot.v4"
	"log"
	"main.go/internal/db"
	"main.go/pkg/gokick"
	"time"
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

func (l *Loop) startAPILoop(b *db.DB, kick gokick.ApiKick) {
	l.toggle = true
	go func() {
		ticker := time.NewTicker(delay)
		slugs, err := b.GetAllUniqueValues()
		if err != nil {
			log.Println("error getting slugs: ", err)
		}
		for {
			select {
			case <-needUpdate:
				slugs, err = b.GetAllUniqueValues()
				if err != nil {
					log.Println("error getting slugs: ", err)
				}
			case <-ticker.C:
				chunks := chunkSlice(slugs, 50)
				for _, s := range chunks {
					r, err := kick.GetChannel(s)
					if err != nil {
						log.Println("error getting channels: ", err.Error())
					}
					for _, i := range r.Data {
						if i.Stream.IsLive == true {
							if time.Since(time.Unix(l.targets[i.Slug], 0)) > 10*time.Minute {
								l.onlSlug <- i
							}
							l.targets[i.Slug] = time.Now().Unix()
						}
					}
				}
			}
		}
	}()
}

func (l *Loop) pauseAPILoop() {
	l.toggle = false
}

func (l *Loop) startMailingLoop(b *db.DB, bot *telebot.Bot) {
	go func() {
		for {
			select {
			case s := <-l.onlSlug:
				ids, err := b.GetAllIdsByValueKick(s.Slug)
				if err != nil {
					log.Println("error getting slugs: ", err)
				}
				for _, i := range ids {
					text := fmt.Sprintf("*Канал*: [%s](kick.com\\/%s) запустил стрим\\!\n>%s", escapeMarkdownV2Text(s.Slug), s.Slug, escapeMarkdownV2Text(s.StreamTitle))
					_, err := bot.Send(&telebot.Chat{ID: i}, text, &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2})
					if errors.Is(err, telebot.ErrUserIsDeactivated) || errors.Is(err, telebot.ErrBlockedByUser) || errors.Is(err, telebot.ErrNotStartedByUser) {
						err = b.RemoveUser(i)
						if err != nil {
							log.Println("error removing user: ", err)
						}
						continue
					} else if err != nil {
						log.Println("error sending message: ", err)
					}
				}
			}
		}
	}()
}
