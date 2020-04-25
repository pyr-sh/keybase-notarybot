package models

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/pkg/errors"
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

type SigHash []byte

func (s SigHash) String() string {
	return hex.EncodeToString(s)
}

var ErrInvalidSigHash = errors.New("invalid sig hash")

func CreateSigHash(key []byte, username string, id string) (SigHash, error) {
	hashFn := hmac.New(sha256.New, key)
	if _, err := fmt.Fprintf(hashFn, "%s/%s", username, id); err != nil {
		return nil, err
	}
	return SigHash(hashFn.Sum(nil)), nil
}

func VerifySigHash(key []byte, username string, id string, hash SigHash) error {
	generated, err := CreateSigHash(key, username, id)
	if err != nil {
		return err
	}
	if subtle.ConstantTimeCompare(hash, generated) != 1 {
		return ErrInvalidSigHash
	}
	return nil
}
