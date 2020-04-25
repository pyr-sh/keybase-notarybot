package keybase

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/dustin/go-humanize"
	"samhofi.us/x/keybase/v2/types/chat1"

	"github.com/pyr-sh/keybase-notarybot/bot/alice"
)

const listUsageMsg = "Usage: `!notary list [sig(nature)(s)]`"

func (b *Bot) handleList(ctx context.Context, msg chat1.MsgNotification, channel alice.Channel, args []string) error {
	if len(args) < 2 {
		if _, err := b.Alice.Chat.Send(ctx, channel, listUsageMsg, nil); err != nil {
			return err
		}
		return nil
	}
	switch args[1] {
	case "signature", "sig", "signatures", "sigs":
		// Make sure that it doesn't already exist
		sigs, err := b.ListUsersSigs(ctx, msg.Msg.Sender.Username)
		if err != nil {
			return err
		}
		lines := []string{
			fmt.Sprintf("@%s, you have uploaded the following signatures:", msg.Msg.Sender.Username),
		}
		for _, sig := range sigs {
			lines = append(
				lines,
				fmt.Sprintf(
					"%s (created %s) - %s",
					strings.TrimSuffix(filepath.Base(sig.Name()), filepath.Ext(sig.Name())),
					humanize.Time(sig.ModTime()),
					strings.Replace(sig.Name(), ".json", ".png", 1),
				),
			)
		}
		if len(sigs) == 0 {
			lines = append(lines, "_No signatures found_")
		}
		if _, err := b.Alice.Chat.Send(
			ctx, channel,
			strings.Join(lines, "\n"),
			nil,
		); err != nil {
			return err
		}
	default:
		if _, err := b.Alice.Chat.Send(ctx, channel, listUsageMsg, nil); err != nil {
			return err
		}
	}
	return nil
}
