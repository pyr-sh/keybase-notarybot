package api

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/pyr-sh/keybase-notarybot/bot/alice"
	"github.com/pyr-sh/keybase-notarybot/bot/models"
	"go.uber.org/zap"
)

const documentPrefix = "data:application/pdf;base64,"

func (a *API) documentsCreate(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		c.AbortWithError(http.StatusMethodNotAllowed, errors.New("only POST is supported"))
		return
	}

	username := c.Query("u")
	id := c.Query("id")
	sig := c.Query("sig")
	sigsText := c.Query("sigs")
	if username == "" || id == "" || sig == "" || sigsText == "" {
		// validation error
		c.AbortWithError(http.StatusBadRequest, errors.New("invalid args"))
		return
	}

	// We don't have validate username and ID since they're signed by us
	sigRaw, err := hex.DecodeString(sig)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("invalid sig format"))
		return
	}
	if err := models.VerifySigHash(a.HMACKey, username, id, models.SigHash(sigRaw)); err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.Wrap(err, "invalid signature"))
		return
	}

	// Sigs is a JSON array
	var sigs []*models.Signatory
	if err := json.Unmarshal([]byte(sigsText), &sigs); err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.Wrap(err, "invalid sigs arg"))
		return
	}

	// We expect a PDF data URI in the body
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, "failed to read the body"))
		return
	}
	if !bytes.HasPrefix(body, []byte(documentPrefix)) {
		c.AbortWithError(http.StatusBadRequest, errors.New("invalid data uri"))
		return
	}
	body = bytes.TrimSpace(body[len(documentPrefix):])
	decodedFile := make([]byte, base64.StdEncoding.DecodedLen(len(body)))
	n, err := base64.StdEncoding.Decode(decodedFile, body)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.Wrap(err, "unable to decode the data uri"))
		return
	}
	decodedFile = decodedFile[:n]

	// Either read the list of existing documents or prepare the directory
	sigPath := filepath.Join("/keybase/private", a.Username+","+username, "documents")
	existingSigs, err := a.Alice.FS.List(c, sigPath, &alice.ListOpts{
		All: true,
	})
	if err == alice.ErrNotExist {
		if err := a.Alice.FS.Mkdir(c, sigPath); err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err,
				"failed to create the documents dir",
			))
			return
		}
	} else if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err,
			"failed to read the documents dir",
		))
		return
	} else {
		// Make sure that the document doesn't already exist
		for _, file := range existingSigs {
			baseName := filepath.Base(file.Name())
			if strings.HasPrefix(baseName, id+".") {
				c.AbortWithError(http.StatusInternalServerError, errors.Errorf(
					"a document with the id %s already exists", id,
				))
				return
			}
		}
	}

	// At this point we're certain we can save the files
	manifestBytes, err := json.Marshal(models.Document{
		Signatories: sigs,
	})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, "failed to encode the manifest"))
		return
	}
	if err := a.Alice.FS.Write(c, filepath.Join(sigPath, id+".pdf"), bytes.NewReader(decodedFile), nil); err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, "failed to write the document image"))
		return
	}
	if err := a.Alice.FS.Write(c, filepath.Join(sigPath, id+".json"), bytes.NewReader(manifestBytes), nil); err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, "failed to write the document manifest"))
		return
	}
	if _, err := a.Alice.Chat.Send(
		c,
		alice.ChatChannel{
			Name: a.Username + "," + username,
		},
		fmt.Sprintf("A new document has been uploaded (name: %s)", id),
		nil,
	); err != nil {
		a.Log.With(
			zap.Error(err),
			zap.String("username", username),
			zap.String("id", id),
		).Warn("Failed to notify the user")
	}
	c.JSON(http.StatusCreated, gin.H{
		"ok": true,
	})
}
