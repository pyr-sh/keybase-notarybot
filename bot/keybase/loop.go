package keybase

import (
	"context"
	"strings"

	"go.uber.org/zap"

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
		args := msgParts[1:]
		switch args[0] {
		case "help":
			if err := b.handleHelp(ctx, msg, channel, args); err != nil {
				return err
			}
		case "create", "new":
			if err := b.handleCreate(ctx, msg, channel, args); err != nil {
				return err
			}
		default:
			if _, err := b.Alice.Chat.Send(ctx, channel, usageMsg, nil); err != nil {
				return err
			}
		}
	}
	if err := ch.Err(); err != nil {
		return err
	}
	return nil
}
