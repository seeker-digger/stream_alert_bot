package telegram

import (
	"errors"
	"fmt"
	"gopkg.in/telebot.v4"
	db "main.go/db"
	"main.go/gokick"
	"slices"
	"strconv"
	"strings"
	"time"
)

type notification struct {
	name     string
	title    string
	startsIn time.Time
}

func onStart(b *db.DB) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		id := c.Sender().ID
		u := db.User{}

		if _, err := b.GetUser(id); errors.Is(err, db.ErrKeyNotExist) {
			err = b.SetUser(id, u)
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}

		return c.Send("*Hello*, "+strconv.FormatInt(id, 10)+"\\!", &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2})
	}
}

func onAdd(b *db.DB, kick gokick.ApiKick) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		tags := c.Args()
		if len(tags) == 0 {
			err := c.Send("*⚠️ Добавьте аргумент:* `\\/add kick.com\\/username`", &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2})
			if err != nil {
				return err
			}
			return nil
		} else if len(tags) > 1 {
			err := c.Send("*⚠️ Слишком много аргументов:* `\\/add kick.com\\/username`", &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2})
			if err != nil {
				return err
			}
			return nil
		}

		slug, err := kick.GetSlugByURL(tags[0])
		if errors.Is(err, gokick.ErrInvalidURL) {
			err = c.Send("*⚠️ Неверная ссылка:* `\\/add kick.com\\/username`", &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2})
			if err != nil {
				return err
			}
			return nil
		} else if errors.Is(err, gokick.ErrUserDoesNotExist) {
			err = c.Send("*⚠️ Такого пользователя не существует\\!*", &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2})
			if err != nil {
				return err
			}
			return nil
		}
		id := c.Sender().ID
		u, err := b.GetUser(id)
		if err != nil {
			return err
		}
		if slices.Contains(u.Kick, slug) {
			err = c.Send("*⚠️ Стример уже отслеживается\\!*", &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2})
			if err != nil {
				return err
			}
			return nil
		}
		u.Kick = append(u.Kick, slug)

		err = b.SetUser(id, u)
		if err != nil {
			return err
		}
		scheduleUpd()

		text := fmt.Sprintf("✅ *Стример: [%s](kick.com/%s) добавлен успешно\\!*", escapeMarkdownV2Text(slug), slug)
		err = c.Send(text, &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2})
		if err != nil {
			return err
		}
		return nil
	}
}

func onRemove(b *db.DB, kick gokick.ApiKick) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		tags := c.Args()
		if len(tags) == 0 {
			err := c.Send("*⚠️ Добавьте аргумент:* `\\/remove arg1`", &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2})
			if err != nil {
				return err
			}
			return nil
		} else if len(tags) > 1 {
			err := c.Send("*⚠️ Слишком много аргументов:* `\\/remove arg1`", &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2})
			if err != nil {
				return err
			}
			return nil
		}
		var slug string
		var err error
		if strings.Contains(tags[0], "/") {
			slug, err = kick.GetSlugByURL(tags[0])
			if errors.Is(err, gokick.ErrInvalidURL) {
				err = c.Send("*⚠️ Некорректная ссылка*", &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2})
				if err != nil {
					return err
				}
				return nil
			} else if errors.Is(err, gokick.ErrUserDoesNotExist) {
				err = c.Send("*⚠️ Такого пользователя не существует\\!*", &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2})
				if err != nil {
					return err
				}
				return nil
			}
		} else {
			slug = tags[0]
		}
		id := c.Sender().ID
		u, err := b.GetUser(id)
		if err != nil {
			return err
		}
		if !slices.Contains(u.Kick, slug) {
			text := fmt.Sprintf("*⚠️ Стример: [%s](kick.com/%s) не отслеживается\\!*", escapeMarkdownV2Text(slug), slug)
			err = c.Send(text, &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2})
			if err != nil {
				return err
			}
			return nil
		}
		u.Kick = removeAll(u.Kick, slug)

		err = b.SetUser(id, u)
		if err != nil {
			return err
		}
		scheduleUpd()

		text := fmt.Sprintf("✅ *Стример: [%s](kick.com/%s) удалён успешно\\!*", escapeMarkdownV2Text(slug), slug)
		err = c.Send(text, &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2})
		if err != nil {
			return err
		}
		return nil
	}
}

func onList(b *db.DB) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		id := c.Sender().ID
		u, err := b.GetUser(id)
		if err != nil {
			return err
		}
		text := fmt.Sprintf("*🎯Отслеживаются:*")
		for j, i := range u.Kick {
			if j == len(u.Kick)-1 {
				text += fmt.Sprintf(" [%s](kick.com/%s)", escapeMarkdownV2Text(i), i)
			} else {
				text += fmt.Sprintf(" [%s](kick.com/%s),", escapeMarkdownV2Text(i), i)
			}
		}
		return c.Send(text, &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2})
	}
}
