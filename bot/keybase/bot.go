package keybase

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/pyr-sh/keybase-notarybot/bot/alice"
)

type Config struct {
	BinaryPath string
	HomeDir    string
	LogPath    string
	Username   string
	PaperKey   string

	HTTPURL string
	HMACKey []byte

	Context context.Context
	Log     *zap.Logger
}

type Bot struct {
	Config
	Alice *alice.Client
}

func New(cfg Config) (*Bot, error) {
	opts := []alice.ClientOption{}
	if cfg.BinaryPath != "" {
		opts = append(opts, alice.ExecutablePath(cfg.BinaryPath))
	}
	if cfg.HomeDir != "" {
		opts = append(opts, alice.HomeDir(cfg.HomeDir))
	}
	if cfg.LogPath != "" {
		opts = append(opts, alice.LogFilePath(cfg.LogPath))
	}
	client, err := alice.New(opts...)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create the client")
	}

	bot := &Bot{
		Config: cfg,
		Alice:  client,
	}
	return bot, nil
}

func (b *Bot) Start(ctx context.Context) error {
	if err := b.Alice.Start(ctx); err != nil {
		return err
	}
	if err := b.Alice.Wait(ctx); err != nil {
		return err
	}
	if b.Config.Username != "" && b.Config.PaperKey != "" {
		if err := b.Alice.Oneshot(b.Config.Context, b.Config.Username, b.Config.PaperKey); err != nil {
			return errors.Wrap(err, "unable to provision using oneshot")
		}
	}

	whoami, err := b.Alice.Whoami(ctx)
	if err != nil {
		return err
	}
	if !whoami.LoggedIn {
		return errors.New("keybase service is not logged in")
	}
	b.Config.Log.With(
		zap.Bool("logged_in", whoami.LoggedIn),
		zap.String("username", whoami.User.Username),
		zap.String("uid", string(whoami.User.UID)),
		zap.String("device_name", whoami.DeviceName),
	).Info("Connected to a Keybase service")

	go b.handlersLoop()

	return nil
}

func (b *Bot) Stop(ctx context.Context) error {
	return b.Alice.Stop(ctx)
}
