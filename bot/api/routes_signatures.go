package api

import "github.com/gin-gonic/gin"

func (a *API) signaturesCreate(c *gin.Context) {
	c.JSON(200, gin.H{
		"hello": "world",
	})
}
