package keybase

import (
	"context"
	"fmt"
	"path/filepath"

	"golang.org/x/sync/errgroup"
	"samhofi.us/x/keybase/v2/types/chat1"

	"github.com/pkg/errors"
	"github.com/pyr-sh/keybase-notarybot/bot/alice"
	"github.com/pyr-sh/keybase-notarybot/bot/models"
)

const deleteUsageMsg = "Usage: `!notary [del|delete|remove|rm] [sig(nature)(s)|doc(ument)(s)] (name)`"

func (b *Bot) handleDelete(ctx context.Context, msg chat1.MsgNotification, channel alice.Channel, args []string) error {
	if len(args) < 2 {
		if _, err := b.Alice.Chat.Send(ctx, channel, deleteUsageMsg, nil); err != nil {
			return err
		}
		return nil
	}
	switch args[1] {
	case "signature", "sig", "signatures", "sigs":
		// Require the user to pass the name
		if len(args) != 3 {
			if _, err := b.Alice.Chat.Send(ctx, channel, deleteUsageMsg, nil); err != nil {
				return err
			}
			return nil
		}
		name := models.NonAlphanumericRE.ReplaceAllString(args[2], "")
		if len(name) < 3 || len(name) > 64 {
			if _, err := b.Alice.Chat.Send(ctx, channel, deleteUsageMsg, nil); err != nil {
				return err
			}
			return nil
		}

		// Make sure that the files exists
		privateDir := filepath.Join(b.PrivateDir(msg.Msg.Sender.Username), "signatures")
		eg, _ := errgroup.WithContext(ctx)
		eg.Go(func() error {
			if _, err := b.Alice.FS.Stat(
				ctx,
				filepath.Join(privateDir, name+".json"),
				nil,
			); err != nil {
				if err == alice.ErrNotExist {
					return nil
				}
				return errors.Wrapf(err, "failed to stat %s", filepath.Join(privateDir, name+".json"))
			}
			if err := b.Alice.FS.Remove(
				ctx,
				filepath.Join(privateDir, name+".json"),
				nil,
			); err != nil {
				return errors.Wrapf(err, "failed to remove %s", filepath.Join(privateDir, name+".json"))
			}
			return nil
		})
		eg.Go(func() error {
			if _, err := b.Alice.FS.Stat(
				ctx,
				filepath.Join(privateDir, name+".png"),
				nil,
			); err != nil {
				if err == alice.ErrNotExist {
					return nil
				}
				return errors.Wrapf(err, "failed to stat %s", filepath.Join(privateDir, name+".png"))
			}
			if err := b.Alice.FS.Remove(
				ctx,
				filepath.Join(privateDir, name+".png"),
				nil,
			); err != nil {
				return errors.Wrapf(err, "failed to remove %s", filepath.Join(privateDir, name+".png"))
			}
			return nil
		})
		if err := eg.Wait(); err != nil {
			return errors.Wrap(err, "failed to delete the signature")
		}
		if _, err := b.Alice.Chat.Send(
			ctx, channel,
			fmt.Sprintf("@%s: Deleted signature %s.", msg.Msg.Sender.Username, name),
			nil,
		); err != nil {
			return err
		}
	case "document", "doc", "documents", "docs":
		// Require the user to pass the name
		if len(args) != 3 {
			if _, err := b.Alice.Chat.Send(ctx, channel, deleteUsageMsg, nil); err != nil {
				return err
			}
			return nil
		}
		name := models.NonAlphanumericRE.ReplaceAllString(args[2], "")
		if len(name) < 3 || len(name) > 64 {
			if _, err := b.Alice.Chat.Send(ctx, channel, deleteUsageMsg, nil); err != nil {
				return err
			}
			return nil
		}

		// Make sure that the files exists
		privateDir := filepath.Join(b.PrivateDir(msg.Msg.Sender.Username), "documents")
		eg, _ := errgroup.WithContext(ctx)
		eg.Go(func() error {
			if _, err := b.Alice.FS.Stat(
				ctx,
				filepath.Join(privateDir, name+".json"),
				nil,
			); err != nil {
				if err == alice.ErrNotExist {
					return nil
				}
				return errors.Wrapf(err, "failed to stat %s", filepath.Join(privateDir, name+".json"))
			}
			if err := b.Alice.FS.Remove(
				ctx,
				filepath.Join(privateDir, name+".json"),
				nil,
			); err != nil {
				return errors.Wrapf(err, "failed to remove %s", filepath.Join(privateDir, name+".json"))
			}
			return nil
		})
		eg.Go(func() error {
			if _, err := b.Alice.FS.Stat(
				ctx,
				filepath.Join(privateDir, name+".pdf"),
				nil,
			); err != nil {
				if err == alice.ErrNotExist {
					return nil
				}
				return errors.Wrapf(err, "failed to stat %s", filepath.Join(privateDir, name+".pdf"))
			}
			if err := b.Alice.FS.Remove(
				ctx,
				filepath.Join(privateDir, name+".pdf"),
				nil,
			); err != nil {
				return errors.Wrapf(err, "failed to remove %s", filepath.Join(privateDir, name+".pdf"))
			}
			return nil
		})
		if err := eg.Wait(); err != nil {
			return errors.Wrap(err, "failed to delete the signature")
		}
		if _, err := b.Alice.Chat.Send(
			ctx, channel,
			fmt.Sprintf("@%s: Deleted document %s.", msg.Msg.Sender.Username, name),
			nil,
		); err != nil {
			return err
		}
	default:
		if _, err := b.Alice.Chat.Send(ctx, channel, deleteUsageMsg, nil); err != nil {
			return err
		}
	}
	return nil
}
