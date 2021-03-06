package alice

import (
	"fmt"
	"strings"
	"time"

	"samhofi.us/x/keybase/v2/types/chat1"
)

// Either ChatChannel or ConversationID.
type Channel interface {
	scope() channelScope
}

// A specified, non-restricted chat channel.
type ChatChannel chat1.ChatChannel

func (c ChatChannel) scope() channelScope { return channelScope{Channel: chat1.ChatChannel(c)} }

// Random conversation ID, required to use as a restricted bot.
type ConversationID chat1.ConvIDStr

func (c ConversationID) scope() channelScope { return channelScope{ConversationID: chat1.ConvIDStr(c)} }

type channelScope struct {
	Channel        chat1.ChatChannel
	ConversationID chat1.ConvIDStr `json:"conversation_id"`
}

type ephemeralLifetime struct {
	time.Duration
}

func (l *ephemeralLifetime) UnmarshalJSON(b []byte) (err error) {
	l.Duration, err = time.ParseDuration(strings.Trim(string(b), `"`))
	return err
}

func (l ephemeralLifetime) MarshalJSON() (b []byte, err error) {
	return []byte(fmt.Sprintf(`"%s"`, l.String())), nil
}
