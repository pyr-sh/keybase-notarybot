package models

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"

	"github.com/pkg/errors"
)

type Signature struct {
	Name       string   `json:"name"`
	Tags       []string `json:"tags"`
	Percentage *float64 `json:"percentage"`
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
