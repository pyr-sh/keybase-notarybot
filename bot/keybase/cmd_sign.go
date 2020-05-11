package keybase

import (
	"context"
	"encoding/json"
	"log"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/pyr-sh/keybase-notarybot/bot/alice"
	"github.com/pyr-sh/keybase-notarybot/bot/models"
	"samhofi.us/x/keybase/v2/types/chat1"
)

const signUsageMsg = "Usage: `!notary sign (doc name){3-64} [username1] [username2] [username3...]`"

func (b *Bot) handleSign(ctx context.Context, msg chat1.MsgNotification, channel alice.Channel, args []string) error {
	if len(args) < 3 {
		if _, err := b.Alice.Chat.Send(ctx, channel, signUsageMsg, nil); err != nil {
			return err
		}
		return nil
	}

	id := filepath.Base(args[1])
	if len(id) == 0 || strings.HasPrefix(id, "..") {
		if _, err := b.Alice.Chat.Send(ctx, channel, signUsageMsg, nil); err != nil {
			return err
		}
		return nil
	}

	base := filepath.Join(b.PrivateDir(msg.Msg.Sender.Username), "documents", id)
	pdfPath := base + ".pdf"
	jsonPath := base + ".json"

	doc, err := b.ReadDoc(ctx, jsonPath)
	if err != nil {
		if _, err := b.Sendf(ctx, channel, "Failed to read that document: %s", err); err != nil {
			return err
		}
		return nil
	}

	var (
		usernames           = []string{}
		signatureNames      = map[string][]string{}
		signatureNameToPath = map[string]map[string]string{}
	)
	for _, username := range args[2:] {
		short := models.NonAlphanumericRE.ReplaceAllString(username, "")
		if short == "" || len(short) < 3 || len(short) > 32 {
			if _, err := b.Sendf(
				ctx, channel,
				"@%s: %s is most likely an invalid username",
				msg.Msg.Sender.Username, username,
			); err != nil {
				return err
			}
			return nil
		}
		usernames = append(usernames, short)

		// Make sure that everyone has signatures with us, since the flow isn't that interactive yet.
		sigs, err := b.ListUsersSigs(ctx, short)
		if err != nil {
			if _, err := b.Sendf(
				ctx, channel,
				"@%s, I couldn't load the signatures of @%s: %s",
				msg.Msg.Sender.Username, short, err.Error(),
			); err != nil {
				return err
			}
			return nil
		}
		if len(sigs) == 0 {
			if _, err := b.Sendf(
				ctx, channel,
				"@%s, unfortunately @%s has no signatures",
				msg.Msg.Sender.Username, short,
			); err != nil {
				return err
			}
			return nil
		}
		names := []string{}
		nameToPath := map[string]string{}
		for _, sig := range sigs {
			name := strings.TrimSuffix(filepath.Base(sig.Name()), filepath.Ext(sig.Name()))
			names = append(names, name)
			nameToPath[name] = sig.Name()
		}
		sort.Strings(names)
		signatureNames[short] = names
		signatureNameToPath[short] = nameToPath
	}

	log.Println(pdfPath)

	// We're effectively trying to figure out who's what signatory
	if len(doc.Signatories) == 0 {
		if _, err := b.Alice.Chat.Send(ctx, channel, "Invalid signatories count in the manifest", nil); err != nil {
			return err
		}
		return nil
	}

	var (
		wg            sync.WaitGroup
		fieldToUserMu sync.Mutex
		fieldToUser   = map[string]string{}
	)
	for _, signatory := range doc.Signatories {
		ch := b.prompt(
			ctx, channel,
			msg.Msg.Sender.Username,
			usernames,
			"@%s, please select who is supposed to sign the field \"%s\" on the document \"%s\"",
			msg.Msg.Sender.Username,
			signatory.Name,
			id,
		)
		wg.Add(1)
		signatory := signatory
		go func() {
			defer wg.Done()
			choice := <-ch
			fieldToUserMu.Lock()
			fieldToUser[signatory.Name] = choice
			fieldToUserMu.Unlock()
		}()
	}
	wg.Wait()

	// At this point "fieldToUser" contains the mapping of field names to the
	// users' signatures. Now we need to map the signatures to the actual files.

	// First, figure out who we're going to ask for signatures
	uniqueUsersMap := map[string]struct{}{}
	userToFields := map[string][]string{}
	for field, user := range fieldToUser {
		if _, ok := uniqueUsersMap[user]; !ok {
			uniqueUsersMap[user] = struct{}{}
		}
		if _, ok := userToFields[user]; !ok {
			userToFields[user] = []string{}
		}
		userToFields[user] = append(userToFields[user], field)
	}
	uniqueUsers := []string{}
	for user := range uniqueUsersMap {
		uniqueUsers = append(uniqueUsers, user)
	}
	sort.Strings(uniqueUsers)

	// Then proceed to ask them for the signatures!
	userFieldToSignatureChoices := map[string]map[string]string{}
	for _, user := range uniqueUsers {
		fields := userToFields[user]
		sort.Strings(fields)
		userFieldToSignatureChoices[user] = map[string]string{}

		for _, field := range fields {
			choice := <-b.prompt(
				ctx,
				channel,
				user,
				signatureNames[user],
				"@%s, which one of your signatures would you like to use for the field \"%s\"?",
				user,
				field,
			)
			userFieldToSignatureChoices[user][field] = choice
		}
	}

	// At this point we can do the final transformation - field to path.
	fieldToPath := map[string]string{}
	for _, signatory := range doc.Signatories {
		user := fieldToUser[signatory.Name]
		if _, ok := userFieldToSignatureChoices[user]; !ok {
			if _, err := b.Sendf(ctx, channel, "Invalid resolved user: %s", user); err != nil {
				return err
			}
			return nil
		}
		choice, ok := userFieldToSignatureChoices[user][signatory.Name]
		if !ok {
			if _, err := b.Sendf(ctx, channel, "Field %s was not resolved for user %s", signatory.Name, user); err != nil {
				return err
			}
			return nil
		}
		path, ok := signatureNameToPath[user][choice]
		if !ok {
			if _, err := b.Sendf(ctx, channel, "Invalid choice %s for field %s selected by %s", choice, signatory.Name, user); err != nil {
				return err
			}
			return nil
		}
		fieldToPath[signatory.Name] = path
	}

	jsb, _ := json.Marshal(fieldToPath)
	b.Sendf(ctx, channel, "@%s, here's what I would do:\n%s", msg.Msg.Sender.Username, string(jsb))

	return nil
}
