package alice

import (
	"context"
	"strings"

	"samhofi.us/x/keybase/v2/types/keybase1"
)

// Provisions the device using oneshot mode.
func (c *Client) Oneshot(ctx context.Context, username string, paperkey string) error {
	res, err := c.ExecWithInput(ctx, strings.NewReader(paperkey), "oneshot", "-u", username)
	if err != nil {
		return err
	}
	return res.RunOnce()
}

type WhoamiUserResult struct {
	UID      keybase1.UID `json:"uid"`
	Username string       `json:"username"`
}

type WhoamiResult struct {
	Configured     bool              `json:"configured"`
	Registered     bool              `json:"registered"`
	LoggedIn       bool              `json:"loggedIn"`
	SessionIsValid bool              `json:"sessionIsValid"`
	User           *WhoamiUserResult `json:"user"`
	DeviceName     string            `json:"deviceName"`
}

// Returns details about the currently logged in user / device.
func (c *Client) Whoami(ctx context.Context) (*WhoamiResult, error) {
	res, err := c.Exec(ctx, "whoami", "-j")
	if err != nil {
		return nil, err
	}
	reply := &WhoamiResult{}
	if err := res.DecodeOnce(reply); err != nil {
		return nil, err
	}
	return reply, nil
}
