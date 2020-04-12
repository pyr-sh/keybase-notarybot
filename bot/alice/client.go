package alice

import (
	"context"
	"os"
	"os/exec"
	"sync"
)

type Client struct {
	Chat Chat

	cfg *clientConfig

	serviceMu sync.Mutex
	service   *exec.Cmd
}

func New(opts ...ClientOption) (*Client, error) {
	cfg := &clientConfig{}
	for _, opt := range opts {
		if err := opt.apply(cfg); err != nil {
			return nil, err
		}
	}

	// executable path defaults to `keybase`
	if cfg.executablePath == "" {
		cfg.executablePath = "keybase"
	}
	var err error
	cfg.executablePath, err = exec.LookPath(cfg.executablePath)
	if err != nil {
		return nil, err
	}

	// the home directory is used as-is
	client := &Client{
		cfg: cfg,
	}
	client.Chat = Chat{c: client}
	return client, nil
}

func (c *Client) Start(ctx context.Context) error {
	args := c.commonArgs()
	if c.cfg.botLiteMode {
		args = append(args, "--enable-bot-lite-mode")
	}
	if c.cfg.logFilePath != "" {
		args = append(args, "--log-file", c.cfg.logFilePath)
	}
	args = append(args, "--debug", "service")
	c.service = exec.CommandContext(ctx, c.cfg.executablePath, args...)
	if c.cfg.logFilePath == "" {
		c.service.Stdout = os.Stdout
		c.service.Stderr = os.Stderr
	}

	if err := c.service.Start(); err != nil {
		return err
	}
	return nil
}

func (c *Client) Stop(ctx context.Context) error {
	return nil
}

func (c *Client) Wait(ctx context.Context) error {
	res, err := c.Exec(ctx, "ctl", "wait")
	if err != nil {
		return err
	}
	return res.RunOnce()
}
