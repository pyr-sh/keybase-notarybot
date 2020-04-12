package alice

import (
	"bytes"
	"context"
	"encoding/json"

	"samhofi.us/x/keybase/v2/types/chat1"
)

type Chat struct {
	c *Client
}

type chatCall struct {
	Method string
	Params interface{}
}

func (c Chat) Send(ctx context.Context, channel Channel, msg string, opts *SendOpts) (*chat1.SendRes, error) {
	body, err := json.Marshal(chatCall{
		Method: "send",
		Params: J{
			"options": &sendArgs{
				channelScope: channel.scope(),
				SendOpts:     opts.Value(),
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
	reply := &struct{ Result *chat1.SendRes }{}
	if err := res.DecodeOnce(reply); err != nil {
		return nil, err
	}
	return reply.Result, nil
}

func (c Chat) React(ctx context.Context, channel Channel, msgID chat1.MessageID, reaction string) (*chat1.SendRes, error) {
	body, err := json.Marshal(chatCall{
		Method: "reaction",
		Params: J{
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
	reply := &struct{ Result *chat1.SendRes }{}
	if err := res.DecodeOnce(reply); err != nil {
		return nil, err
	}
	return reply.Result, nil
}
