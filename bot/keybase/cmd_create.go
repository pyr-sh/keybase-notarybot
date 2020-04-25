package keybase

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/dchest/uniuri"
	"samhofi.us/x/keybase/v2/types/chat1"

	"github.com/pyr-sh/keybase-notarybot/bot/alice"
)

const createUsageMsg = "Usage: `!notary [create|new] [signature|sig]`"

func (b *Bot) handleCreate(ctx context.Context, msg chat1.MsgNotification, channel alice.Channel, args []string) error {
	if len(args) < 2 {
		if _, err := b.Alice.Chat.Send(ctx, channel, createUsageMsg, nil); err != nil {
			return err
		}
		return nil
	}
	switch args[1] {
	case "signature", "sig":
		// Actual signature uploads are performed through the HTTP interface, so we simply
		// need to provide the user with a MAC'd ID.
		sigID := uniuri.NewLen(uniuri.UUIDLen)
		mac := hmac.New(sha256.New, b.HMACKey)
		if _, err := mac.Write([]byte(sigID)); err != nil {
			return err
		}
		completeURL := b.Config.HTTPURL + "/signature/" + msg.Msg.Sender.Username + "/" + sigID + "/" + hex.EncodeToString(mac.Sum(nil))

		// We always want to send a signature in a private message.
		privateChannel := b.privateChannel(msg.Msg.Sender)
		if b.isChannelPrivate(msg.Msg.Sender, msg.Msg.Channel) {
			privateChannel = channel
		} else {
			if _, err := b.Alice.Chat.Send(
				ctx, channel,
				fmt.Sprintf("Sent a new signature upload link to @%s.", msg.Msg.Sender.Username),
				nil,
			); err != nil {
				return err
			}
		}
		if _, err := b.Alice.Chat.Send(
			ctx, privateChannel,
			fmt.Sprintf("Click here to upload your signature:\n%s", completeURL),
			nil,
		); err != nil {
			return err
		}
	default:
		if _, err := b.Alice.Chat.Send(ctx, channel, createUsageMsg, nil); err != nil {
			return err
		}
	}
	return nil
}
