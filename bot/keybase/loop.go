package keybase

import (
	"context"
	"log"

	"go.uber.org/zap"

	"github.com/pyr-sh/keybase-notarybot/bot/alice"
)

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
		if msg.Msg == nil {
			continue
		}

		log.Printf("got msg %#v", msg)
		channel := alice.ConversationID(msg.Msg.ConvID)
		x, err := b.Alice.Chat.Send(
			ctx,
			channel,
			"hello world",
			nil,
		)
		if err != nil {
			return err
		}
		log.Printf("%#v", x)
		y, err := b.Alice.Chat.React(ctx, channel, *x.MessageID, ":wave:")
		if err != nil {
			return err
		}
		log.Printf("%#v", y)
	}
	if err := ch.Err(); err != nil {
		return err
	}
	return nil
}
