package api

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/pyr-sh/keybase-notarybot/bot/models"
)

var alphanumericRE = regexp.MustCompile("[^a-zA-Z0-9]+")

func (a *API) signaturesCreate(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		c.AbortWithError(http.StatusMethodNotAllowed, errors.New("only POST is supported"))
		return
	}

	username := c.Query("u")
	id := c.Query("id")
	sig := c.Query("sig")
	name := alphanumericRE.ReplaceAllString(c.Query("name"), "")
	percentage := c.Query("p")
	if username == "" || id == "" || sig == "" || name == "" {
		// validation error
		c.AbortWithError(http.StatusBadRequest, errors.New("invalid args"))
		return
	}
	sigRaw, err := hex.DecodeString(sig)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("invalid sig format"))
		return
	}
	if err := models.VerifySigHash(a.HMACKey, username, id, models.SigHash(sigRaw)); err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.Wrap(err, "invalid signature"))
		return
	}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, "failed to read the body"))
		return
	}

	fmt.Printf("Got the following details:\nUsername: %s\nID: %s\nSig: %s (valid)\nPercentage: %s\nBody: %s\n", username, id, sig, percentage, string(body))
	c.JSON(http.StatusCreated, gin.H{
		"hello": "world",
	})
}
