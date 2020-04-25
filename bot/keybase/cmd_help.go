package keybase

import (
	"context"

	"github.com/pyr-sh/keybase-notarybot/bot/alice"
	"samhofi.us/x/keybase/v2/types/chat1"
)

const helpMsg = `Hello! I'm @notarybot, your contract signing companion.

` + "`" + `!notary [create|new] [signature|sig]` + "`" + `
Messages you with a link to upload a new signature.`

func (b *Bot) handleHelp(ctx context.Context, msg chat1.MsgNotification, channel alice.Channel, args []string) error {
	_, err := b.Alice.Chat.Send(ctx, channel, helpMsg, nil)
	return err
}
