package keybase

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/pyr-sh/keybase-notarybot/bot/alice"
	"github.com/pyr-sh/keybase-notarybot/bot/models"
)

func (b *Bot) ListUsersSigs(ctx context.Context, username string) ([]os.FileInfo, error) {
	sigPath := filepath.Join(b.PrivateDir(username), "signatures")
	files, err := b.Alice.FS.List(ctx, sigPath, &alice.ListOpts{
		All: true,
	})
	if err == alice.ErrNotExist {
		return nil, b.Alice.FS.Mkdir(ctx, sigPath)
	}
	if err != nil {
		return nil, err
	}
	result := []os.FileInfo{}
	for _, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}
		result = append(result, file)
	}
	return result, nil
}

func (b *Bot) ReadSig(ctx context.Context, path string) (*models.Signature, error) {
	rc, err := b.Alice.FS.Read(ctx, path, nil)
	defer rc.Close()
	if err != nil {
		return nil, err
	}
	var manifest models.Signature
	if err := json.NewDecoder(rc).Decode(&manifest); err != nil {
		return nil, err
	}
	return &manifest, nil
}
