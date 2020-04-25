package api

import "github.com/gin-gonic/gin"

func (a *API) errorMiddleware(c *gin.Context) {
	c.Next()

	if len(c.Errors) > 0 {
		c.JSON(-1, c.Errors)
	}
}
