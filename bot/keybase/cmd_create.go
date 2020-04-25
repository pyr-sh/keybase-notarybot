package keybase

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"samhofi.us/x/keybase/v2/types/chat1"

	"github.com/pyr-sh/keybase-notarybot/bot/alice"
	"github.com/pyr-sh/keybase-notarybot/bot/models"
)

const createUsageMsg = "Usage: `!notary [create|new] [signature|sig|document|doc] (name){3-64}`"

func (b *Bot) handleCreate(ctx context.Context, msg chat1.MsgNotification, channel alice.Channel, args []string) error {
	if len(args) < 2 {
		if _, err := b.Alice.Chat.Send(ctx, channel, createUsageMsg, nil); err != nil {
			return err
		}
		return nil
	}
	switch args[1] {
	case "document", "doc":
		// Figure out the name arg
		if len(args) != 3 {
			if _, err := b.Alice.Chat.Send(ctx, channel, createUsageMsg, nil); err != nil {
				return err
			}
			return nil
		}
		name := models.NonAlphanumericRE.ReplaceAllString(args[2], "")
		if len(name) < 3 || len(name) > 64 {
			if _, err := b.Alice.Chat.Send(ctx, channel, createUsageMsg, nil); err != nil {
				return err
			}
			return nil
		}
		log.Println(name)

		// Make sure that it doesn't already exist
		sigs, err := b.ListUsersSigs(ctx, msg.Msg.Sender.Username)
		if err != nil {
			return err
		}
		for _, sig := range sigs {
			if strings.TrimSuffix(filepath.Base(sig.Name()), filepath.Ext(sig.Name())) == name {
				if _, err := b.Alice.Chat.Send(ctx, channel, "Document with that name already exists.", nil); err != nil {
					return err
				}
				return nil
			}
		}

		// Actual document uploads are performed through the HTTP interface, so we simply
		// need to provide the user with a MAC'd ID.
		sigHash, err := models.CreateSigHash(b.HMACKey, msg.Msg.Sender.Username, name)
		if err != nil {
			return err
		}
		completeURL := fmt.Sprintf(
			"%s/document/%s/%s/%s",
			b.Config.HTTPURL,
			msg.Msg.Sender.Username,
			name,
			sigHash,
		)

		// We always want to send a document in a private message.
		privateChannel := b.privateChannel(msg.Msg.Sender)
		if b.isChannelPrivate(msg.Msg.Sender, msg.Msg.Channel) {
			privateChannel = channel
		} else {
			if _, err := b.Alice.Chat.Send(
				ctx, channel,
				fmt.Sprintf("Sent a new document upload link to @%s.", msg.Msg.Sender.Username),
				nil,
			); err != nil {
				return err
			}
		}
		if _, err := b.Alice.Chat.Send(
			ctx, privateChannel,
			fmt.Sprintf("Click here to upload your document:\n%s", completeURL),
			nil,
		); err != nil {
			return err
		}
		return nil
	case "signature", "sig":
		// Figure out the name arg
		if len(args) != 3 {
			if _, err := b.Alice.Chat.Send(ctx, channel, createUsageMsg, nil); err != nil {
				return err
			}
			return nil
		}
		name := models.NonAlphanumericRE.ReplaceAllString(args[2], "")
		if len(name) < 3 || len(name) > 64 {
			if _, err := b.Alice.Chat.Send(ctx, channel, createUsageMsg, nil); err != nil {
				return err
			}
			return nil
		}

		// Make sure that it doesn't already exist
		sigs, err := b.ListUsersSigs(ctx, msg.Msg.Sender.Username)
		if err != nil {
			return err
		}
		for _, sig := range sigs {
			if strings.TrimSuffix(filepath.Base(sig.Name()), filepath.Ext(sig.Name())) == name {
				if _, err := b.Alice.Chat.Send(ctx, channel, "Signature with that name already exists.", nil); err != nil {
					return err
				}
				return nil
			}
		}

		// Actual signature uploads are performed through the HTTP interface, so we simply
		// need to provide the user with a MAC'd ID.
		sigHash, err := models.CreateSigHash(b.HMACKey, msg.Msg.Sender.Username, name)
		if err != nil {
			return err
		}
		completeURL := fmt.Sprintf(
			"%s/signature/%s/%s/%s",
			b.Config.HTTPURL,
			msg.Msg.Sender.Username,
			name,
			sigHash,
		)

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
		return nil
	default:
		if _, err := b.Alice.Chat.Send(ctx, channel, createUsageMsg, nil); err != nil {
			return err
		}
	}
	return nil
}
