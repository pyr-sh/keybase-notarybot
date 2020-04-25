package keybase

import (
	"strings"

	"github.com/pyr-sh/keybase-notarybot/bot/alice"
	"samhofi.us/x/keybase/v2/types/chat1"
)

func (b *Bot) isChannelPrivate(sender chat1.MsgSender, channel chat1.ChatChannel) bool {
	if channel.Name == "" {
		return false
	}
	users := strings.Split(channel.Name, ",")
	if len(users) != 2 {
		return false
	}
	for _, user := range users {
		if user != b.Username && user != sender.Username {
			return false
		}
	}
	return true
}

func (b *Bot) privateChannel(sender chat1.MsgSender) alice.Channel {
	return &alice.ChatChannel{
		Name: b.Username + "," + sender.Username,
	}
}
