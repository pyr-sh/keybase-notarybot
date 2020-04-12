package database

import (
	"context"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type Config struct {
	Context context.Context
	Driver  string
	DSN     string

	Log *zap.Logger
}

type Database struct {
	Config

	db *sqlx.DB
}

func New(cfg Config) (*Database, error) {
	conn, err := sqlx.ConnectContext(cfg.Context, cfg.Driver, cfg.DSN)
	if err != nil {
		return nil, err
	}
	db := &Database{
		Config: cfg,
		db:     conn,
	}
	return db, nil
}

func (d *Database) Stop() error {
	return d.db.Close()
}
