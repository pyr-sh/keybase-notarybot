package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

func (a *API) signaturesCreate(c *gin.Context) {
	id := c.Query("id")
	mac := c.Query("mac")
	if id == "" || mac == "" {
		// validation error
	}

	mac := hmac.New(sha256.New, b.HMACKey)
	if _, err := mac.Write([]byte(id)); err != nil {
		// unable to hash
	}
	if !hmac.Equal(hex.EncodeToString(mac.Sum(nil))) {
		// invalid hmac
	}

	c.JSON(200, gin.H{
		"hello": "world",
	})
}
