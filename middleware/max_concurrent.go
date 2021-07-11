package middleware

import (
	"github.com/cecilozaur/tnemesum/conf"
	"github.com/gin-gonic/gin"
	"net/http"
)

const defaultMaxConcurrent = 100

func MaxConcurrent(cfg conf.Config) gin.HandlerFunc {
	maxConcurrent := defaultMaxConcurrent
	if cfg.MaxConcurrent > 0 {
		maxConcurrent = cfg.MaxConcurrent
	}
	sem := make(chan struct{}, maxConcurrent)
	return func(c *gin.Context) {
		defer func() {
			<-sem
		}()

		select {
		case sem <- struct{}{}:
			c.Next()
		default:
			c.AbortWithStatus(http.StatusTooManyRequests)
		}
	}
}
