package alice

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"samhofi.us/x/keybase/v2/types/chat1"
	"samhofi.us/x/keybase/v2/types/keybase1"
)

type Chat struct {
	c *Client
}

type chatCall struct {
	Method string
	Params interface{}
}

type SendOpts struct {
	Nonblock          bool              `json:"nonblock"`
	MembersType       string            `json:"members_type"`
	EphemeralLifetime ephemeralLifetime `json:"exploding_lifetime"`
	ConfirmLumenSend  bool              `json:"confirm_lumen_send"`
	ReplyTo           *chat1.MessageID  `json:"reply_to"`
}

func (o *SendOpts) value() SendOpts {
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

// Sends a message to the specified channel.
func (c Chat) Send(ctx context.Context, channel Channel, msg string, opts *SendOpts) (*chat1.SendRes, error) {
	body, err := json.Marshal(chatCall{
		Method: "send",
		Params: jm{
			"options": &sendArgs{
				channelScope: channel.scope(),
				SendOpts:     opts.value(),
				Message: chat1.ChatMessage{
					Body: msg,
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	res, err := c.c.ExecWithInput(ctx, bytes.NewReader(body), "chat", "api")
	if err != nil {
		return nil, err
	}
	reply := &struct{ Result *chat1.SendRes }{}
	if err := res.DecodeOnce(reply); err != nil {
		return nil, err
	}
	return reply.Result, nil
}

type reactArgs struct {
	channelScope
	MessageID chat1.MessageID `json:"message_id"`
	Message   chat1.ChatMessage
}

// Reacts to a specified message in a channel.
func (c Chat) React(ctx context.Context, channel Channel, msgID chat1.MessageID, reaction string) (*chat1.SendRes, error) {
	body, err := json.Marshal(chatCall{
		Method: "reaction",
		Params: jm{
			"options": &reactArgs{
				channelScope: channel.scope(),
				MessageID:    msgID,
				Message: chat1.ChatMessage{
					Body: reaction,
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	res, err := c.c.ExecWithInput(ctx, bytes.NewReader(body), "chat", "api")
	if err != nil {
		return nil, err
	}
	reply := &struct{ Result *chat1.SendRes }{}
	if err := res.DecodeOnce(reply); err != nil {
		return nil, err
	}
	return reply.Result, nil
}

type Advertisement struct {
	Alias          string `json:"alias,omitempty"`
	Advertisements []*chat1.AdvertiseCommandAPIParam
}

// Advertises specific commands
func (c Chat) AdvertiseCommands(ctx context.Context, advertisement *Advertisement) error {
	body, err := json.Marshal(chatCall{
		Method: "advertisecommands",
		Params: jm{
			"options": advertisement,
		},
	})
	if err != nil {
		return err
	}
	res, err := c.c.ExecWithInput(ctx, bytes.NewReader(body), "chat", "api")
	if err != nil {
		return err
	}
	return res.RunOnce()
}

// Clears the advertised commands
func (c Chat) ClearCommands(ctx context.Context) error {
	body, err := json.Marshal(chatCall{
		Method: "clearcommands",
		Params: jm{}, // explicitly pass {}
	})
	if err != nil {
		return err
	}
	res, err := c.c.ExecWithInput(ctx, bytes.NewReader(body), "chat", "api")
	if err != nil {
		return err
	}
	return res.RunOnce()
}

type ChatListenOpts struct {
	ShowLocal         bool
	ShowExploding     bool
	SubscribeToConvs  bool
	SubscribeToDev    bool
	SubscribeToWallet bool
}

// Subscribes to new chat API notifications. Pass channels to enable filtering, pass
// options to customize parameters such as enabling wallet / local notifications.
func (c Chat) Listen(ctx context.Context, channels []ChatChannel, opts *ChatListenOpts) (*ChatListenStream, error) {
	args := []interface{}{"chat", "api-listen"}
	if opts != nil {
		if opts.ShowLocal {
			args = append(args, "--local")
		}
		if !opts.ShowExploding {
			args = append(args, "--hide-exploding")
		}
		if opts.SubscribeToConvs {
			args = append(args, "--convs")
		}
		if opts.SubscribeToDev {
			args = append(args, "--dev")
		}
		if opts.SubscribeToDev {
			args = append(args, "--wallet")
		}
	}
	if len(channels) > 0 {
		args = append(args, "--filter-channels", channels)
	}

	reply, err := c.c.Exec(ctx, args...)
	if err != nil {
		return nil, err
	}
	stream, err := reply.Stream()
	if err != nil {
		return nil, err
	}

	result := &ChatListenStream{
		res: stream,
	}
	if opts == nil || opts.ShowLocal {
		// If we're not asking for local messages, we need to know our device's details
		// to skip our messages
		whoami, err := c.c.Whoami(ctx)
		if err != nil {
			if err2 := stream.Close(); err2 != nil {
				return nil, errors.Wrapf(err, "also failed to close: %s", err)
			}
			return nil, err
		}

		if whoami.LoggedIn {
			result.uid = whoami.User.UID
			result.deviceName = whoami.DeviceName
		}
	}
	return result, nil
}

type ChatListenStream struct {
	uid        keybase1.UID
	deviceName string

	ctx       context.Context
	cancel    context.CancelFunc
	res       *StreamedResult
	readError error
}

// Returns a channel populated with API notifications.
func (c *ChatListenStream) Messages() chan chat1.MsgNotification {
	c.ctx, c.cancel = context.WithCancel(c.res.Context)
	ch := make(chan chat1.MsgNotification)
	go func() {
		defer close(ch)

		for c.res.Next() {
			select {
			case <-c.ctx.Done():
				c.readError = context.Canceled
				return
			default:
			}
			frame := chat1.MsgNotification{}
			if err := c.res.Decode(&frame); err != nil {
				c.readError = err
				return
			}
			if c.uid != "" && c.deviceName != "" && frame.Msg != nil &&
				frame.Msg.Sender.Uid == c.uid && frame.Msg.Sender.DeviceName == c.deviceName {
				// Skip our messages
				continue
			}
			ch <- frame
		}
	}()
	return ch
}

// If any error occured during reading / decoding, it's returned here.
func (c *ChatListenStream) Err() error {
	if c.readError != nil {
		return c.readError
	}
	return c.res.Err()
}

// Closes the stream, killing the client process.
func (c *ChatListenStream) Close() error {
	if c.cancel != nil {
		c.cancel()
	}
	return c.res.Close()
}
