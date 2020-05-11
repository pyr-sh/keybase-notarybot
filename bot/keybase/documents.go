package keybase

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/pyr-sh/keybase-notarybot/bot/alice"
	"github.com/pyr-sh/keybase-notarybot/bot/models"
)

func (b *Bot) ListUsersDocs(ctx context.Context, username string) ([]os.FileInfo, error) {
	docPath := filepath.Join(b.PrivateDir(username), "documents")
	files, err := b.Alice.FS.List(ctx, docPath, &alice.ListOpts{
		All: true,
	})
	if err == alice.ErrNotExist {
		return nil, b.Alice.FS.Mkdir(ctx, docPath)
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

func (b *Bot) ReadDoc(ctx context.Context, path string) (*models.Document, error) {
	rc, err := b.Alice.FS.Read(ctx, path, nil)
	defer rc.Close()
	if err != nil {
		return nil, err
	}
	var manifest models.Document
	if err := json.NewDecoder(rc).Decode(&manifest); err != nil {
		return nil, err
	}
	return &manifest, nil
}
