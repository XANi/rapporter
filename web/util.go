package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func bodySizeMiddleware(c *gin.Context) {
	var maxBytes int64 = 1024 * 1024 * 1
	var w http.ResponseWriter = c.Writer
	c.Request.Body = http.MaxBytesReader(w, c.Request.Body, maxBytes)
	c.Next()

}
