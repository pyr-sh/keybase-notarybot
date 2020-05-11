package keybase

import (
	"context"
	"fmt"

	"github.com/pyr-sh/keybase-notarybot/bot/alice"
	"samhofi.us/x/keybase/v2/types/chat1"
)

func errchan(err error) <-chan string {
	return strchan(err.Error())
}

func strchan(str string) <-chan string {
	ch := make(chan string)
	go func() {
		ch <- str
		close(ch)
	}()
	return ch
}

func (b *Bot) prompt(
	ctx context.Context,
	channel alice.Channel,
	username string,
	options []string,
	msgText string,
	args ...interface{},
) <-chan string {
	msg, err := b.Sendf(ctx, channel, msgText, args...)
	if err != nil {
		return errchan(err)
	}
	optionsMap := map[string]struct{}{}
	for _, option := range options {
		if _, err := b.Alice.Chat.React(ctx, channel, *msg.MessageID, option); err != nil {
			return errchan(err)
		}
		optionsMap[option] = struct{}{}
	}

	res := make(chan string)
	key := fmt.Sprintf("%s:reaction", username)
	var fn func(reaction chat1.MsgNotification)
	fn = func(reaction chat1.MsgNotification) {
		if reaction.Msg.ConvID != reaction.Msg.ConvID || reaction.Msg.Content.Reaction.MessageID != *msg.MessageID {
			return
		}
		choice := reaction.Msg.Content.Reaction.Body
		if _, ok := optionsMap[choice]; !ok {
			return
		}

		b.Bus.Unsubscribe(key, fn)

		res <- choice
		close(res)

		if _, err := b.Alice.Chat.React(
			ctx, channel, reaction.Msg.Content.Reaction.MessageID,
			":white_check_mark:",
		); err != nil {
			return
		}

	}
	if err := b.Bus.Subscribe(key, fn); err != nil {
		return errchan(err)
	}
	return res
}
