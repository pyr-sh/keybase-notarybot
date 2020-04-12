package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Signature struct {
	ID        string          `db:"id" json:"id"`
	CreatedAt time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt time.Time       `db:"updated_at" json:"updated_at"`
	Username  string          `db:"username" json:"username"`
	Name      string          `db:"name" json:"name"`
	FileURL   string          `db:"file_url" json:"file_url"`
	LinePos   decimal.Decimal `db:"line_pos" json:"line_pos"`
}
