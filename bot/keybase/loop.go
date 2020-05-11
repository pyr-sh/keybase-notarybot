package keybase

import (
	"context"
	"fmt"
	"log"
	"strings"

	"go.uber.org/zap"
	"samhofi.us/x/keybase/v2/types/chat1"

	"github.com/pyr-sh/keybase-notarybot/bot/alice"
)

const usageMsg = "Usage: `!notary [help|create|list|update]`"

func (b *Bot) handlersLoop() {
	for {
		select {
		case <-b.Context.Done():
			return
		default:
		}

		if err := b.startHandler(b.Context); err != nil {
			b.Log.With(zap.Error(err)).Warn("Restarting bot chat handler...")
		}
	}
}

func (b *Bot) startHandler(ctx context.Context) error {
	if err := b.Alice.Chat.AdvertiseCommands(ctx, &alice.Advertisement{
		Alias: "Notary Bot",
		Advertisements: []*chat1.AdvertiseCommandAPIParam{
			{
				Typ: "public",
				Commands: []chat1.UserBotCommandInput{
					{
						Name:        "notary",
						Description: "Allows you to interact with the document signing functionality.",
						Usage:       "[create|list|delete] [signatures] ...",
					},
				},
			},
		},
	}); err != nil {
		return err
	}

	ch, err := b.Alice.Chat.Listen(ctx, nil, nil)
	if err != nil {
		return err
	}
	defer ch.Close()
	b.Log.Info("Listening to new Keybase messages...")
	for msg := range ch.Messages() {
		if msg.Error != nil {
			b.Log.With(zap.String("error", *msg.Error)).Warn("Bot handler received an error")
			continue
		}
		if msg.Msg != nil {
			b.Bus.Publish(fmt.Sprintf("%s:%s", msg.Msg.Sender.Username, msg.Msg.Content.TypeName), msg)
		}
		if msg.Msg == nil || msg.Msg.Content.TypeName != "text" {
			continue
		}
		channel := alice.ConversationID(msg.Msg.ConvID)
		if !strings.HasPrefix(msg.Msg.Content.Text.Body, "!notary ") {
			if _, err := b.Alice.Chat.Send(ctx, channel, usageMsg, nil); err != nil {
				return err
			}
			continue
		}
		msgParts := strings.Split(msg.Msg.Content.Text.Body, " ")
		if len(msgParts) == 1 {
			if _, err := b.Alice.Chat.Send(ctx, channel, usageMsg, nil); err != nil {
				return err
			}
			continue
		}
		if _, err := b.Alice.Chat.React(ctx, channel, msg.Msg.Id, ":eyes:"); err != nil {
			return err
		}

		go func(msg chat1.MsgNotification) {
			args := msgParts[1:]
			switch args[0] {
			case "help":
				if err := b.handleHelp(ctx, msg, channel, args); err != nil {
					log.Println(err)
				}
			case "create", "new":
				if err := b.handleCreate(ctx, msg, channel, args); err != nil {
					log.Println(err)
				}
			case "list":
				if err := b.handleList(ctx, msg, channel, args); err != nil {
					log.Println(err)
				}
			case "delete", "del", "rm", "remove":
				if err := b.handleDelete(ctx, msg, channel, args); err != nil {
					log.Println(err)
				}
			case "sign":
				if err := b.handleSign(ctx, msg, channel, args); err != nil {
					log.Println(err)
				}
			default:
				if _, err := b.Alice.Chat.Send(ctx, channel, usageMsg, nil); err != nil {
					log.Println(err)
				}
			}
		}(msg)
	}
	if err := ch.Err(); err != nil {
		return err
	}
	return nil
}
