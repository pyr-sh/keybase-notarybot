package alice

import (
	"context"
	"os"
	"os/exec"

	"golang.org/x/sync/errgroup"
)

type Client struct {
	Chat Chat
	FS   FS

	cfg *clientConfig

	service *exec.Cmd
	kbfs    *exec.Cmd
}

// Creates a new Keybase client using the given settings
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

	if cfg.kbfsEnabled {
		// kbfs executable path defaults to `kbfsfuse`
		if cfg.kbfsExecutablePath == "" {
			cfg.kbfsExecutablePath = "kbfsfuse"
		}
		var err error
		cfg.kbfsExecutablePath, err = exec.LookPath(cfg.kbfsExecutablePath)
		if err != nil {
			return nil, err
		}
	}

	// the home directory is used as-is
	client := &Client{
		cfg: cfg,
	}
	client.Chat = Chat{c: client}
	client.FS = FS{c: client}
	return client, nil
}

// Starts a library-managed Keybase service. Make sure to Stop() it to gracefully
// shut it down. The passed context gets directly passed by.
func (c *Client) Start(ctx context.Context) error {
	eg, _ := errgroup.WithContext(ctx)
	eg.Go(func() error {
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
	})
	eg.Go(func() error {
		if !c.cfg.kbfsEnabled {
			return nil
		}

		res, err := c.Exec(ctx, "ctl", "wait")
		if err != nil {
			return err
		}
		if err := res.RunOnce(); err != nil {
			return err
		}

		args := []string{"-debug"}
		if c.cfg.mountType == "" {
			c.cfg.mountType = "none"
		}
		args = append(args, "-mount-type="+c.cfg.mountType)
		if c.cfg.kbfsLogFilePath != "" {
			args = append(args, "-log-file="+c.cfg.kbfsLogFilePath)
		}
		c.kbfs = exec.CommandContext(ctx, c.cfg.kbfsExecutablePath, args...)
		c.kbfs.Env = append(os.Environ(), "KEYBASE_DEBUG=1")
		if c.cfg.kbfsLogFilePath == "" {
			c.kbfs.Stdout = os.Stdout
			c.kbfs.Stderr = os.Stderr
		}
		if err := c.kbfs.Start(); err != nil {
			return err
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

// Stop kills the managed Keybase service process.
func (c *Client) Stop(ctx context.Context) error {
	eg, _ := errgroup.WithContext(ctx)
	eg.Go(func() error {
		if c.service.Process == nil {
			return nil
		}
		if err := c.service.Process.Kill(); err != nil {
			return err
		}
		if _, err := c.service.Process.Wait(); err != nil {
			return err
		}
		return nil
	})
	eg.Go(func() error {
		if c.kbfs.Process == nil {
			return nil
		}
		if err := c.kbfs.Process.Kill(); err != nil {
			return err
		}
		if _, err := c.kbfs.Process.Wait(); err != nil {
			return err
		}
		return nil
	})
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

// Wait waits for the managed service to fully start up
func (c *Client) Wait(ctx context.Context) error {
	args := []interface{}{"--debug", "ctl", "wait"}
	if c.cfg.kbfsEnabled {
		args = append(args, "--include-kbfs")
	}
	res, err := c.Exec(ctx, args...)
	if err != nil {
		return err
	}
	return res.RunOnce()
}
