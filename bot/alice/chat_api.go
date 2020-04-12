package alice

import (
	"fmt"
	"strings"
	"time"

	"samhofi.us/x/keybase/v2/types/chat1"
)

type Channel interface {
	apply(o *channelScope)
	scope() channelScope
}

type ChatChannel chat1.ChatChannel

func (c ChatChannel) apply(o *channelScope) { o.Channel = chat1.ChatChannel(c) }
func (c ChatChannel) scope() channelScope   { return channelScope{Channel: chat1.ChatChannel(c)} }

type ConversationID chat1.ConvIDStr

func (c ConversationID) apply(o *channelScope) { o.ConversationID = chat1.ConvIDStr(c) }
func (c ConversationID) scope() *channelScope {
	return &channelScope{ConversationID: chat1.ConvIDStr(c)}
}

type channelScope struct {
	Channel        chat1.ChatChannel
	ConversationID chat1.ConvIDStr `json:"conversation_id"`
}

type SendOpts struct {
	Nonblock          bool              `json:"nonblock"`
	MembersType       string            `json:"members_type"`
	EphemeralLifetime ephemeralLifetime `json:"exploding_lifetime"`
	ConfirmLumenSend  bool              `json:"confirm_lumen_send"`
	ReplyTo           *chat1.MessageID  `json:"reply_to"`
}

func (o *SendOpts) Value() SendOpts {
	if o == nil {
		return SendOpts{}
	}
	return *o
}

type sendArgs struct {
	channelScope
	SendOpts
	Message chat1.ChatMessage
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

type reactArgs struct {
	channelScope
	MessageID chat1.MessageID `json:"message_id"`
	Message   chat1.ChatMessage
}
