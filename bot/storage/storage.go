package storage

import (
	"io"

	"github.com/graymeta/stow"
	_ "github.com/graymeta/stow/google"
	_ "github.com/graymeta/stow/local"
	_ "github.com/graymeta/stow/s3"
	_ "github.com/graymeta/stow/swift"
	"go.uber.org/zap"
)

type Config struct {
	Log       *zap.Logger
	Kind      string
	Params    map[string]string
	Container string
}

type Storage struct {
	Config
	location  stow.Location
	container stow.Container
}

func New(cfg Config) (*Storage, error) {
	location, err := stow.Dial(cfg.Kind, stow.ConfigMap(cfg.Params))
	if err != nil {
		return nil, err
	}

	container, err := location.Container(cfg.Container)
	if err == stow.ErrNotFound {
		container, err = location.CreateContainer(cfg.Container)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return &Storage{
		Config:    cfg,
		location:  location,
		container: container,
	}, nil
}

func (s *Storage) Get(name string) (io.ReadCloser, error) {
	item, err := s.container.Item(name)
	if err != nil {
		return nil, err
	}
	return item.Open()
}

func (s *Storage) Put(name string, r io.Reader, size int64) (string, error) {
	item, err := s.container.Put(name, r, size, nil)
	if err != nil {
		return "", err
	}
	return item.URL().String(), nil
}

func (s *Storage) Close() error {
	if err := s.location.Close(); err != nil {
		return err
	}
	return nil
}
