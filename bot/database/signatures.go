package database

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/pyr-sh/keybase-notarybot/bot/models"
)

func (d *Database) PutSignature(ctx context.Context, model *models.Signature) (*models.Signature, error) {
	query, args, err := sq.Insert("signatures").
		SetMap(models.Export(model)).
		Suffix(models.Returning(model)).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}
	if err := d.db.GetContext(ctx, model, query, args...); err != nil {
		return nil, err
	}
	return model, nil
}
