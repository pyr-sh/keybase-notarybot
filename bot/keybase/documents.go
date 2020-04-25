package keybase

import (
	"context"
	"os"
	"path/filepath"

	"github.com/pyr-sh/keybase-notarybot/bot/alice"
)

func (b *Bot) ListUserDocs(ctx context.Context, username string) ([]os.FileInfo, error) {
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
