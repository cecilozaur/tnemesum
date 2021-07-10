package middleware

import (
	"github.com/cecilozaur/tnemesum/conf"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const defaultMaxConcurrent = 100

func RateControl(cfg conf.Config) gin.HandlerFunc {
	maxAllowed := defaultMaxConcurrent
	if cfg.MaxConcurrent > 0 {
		maxAllowed = cfg.MaxConcurrent
	}
	sem := make(chan struct{}, maxAllowed)
	return func(c *gin.Context) {
		requestStart := time.Now()
		sem <- struct{}{}
		defer func() {
			<-sem
		}()

		// if the request wait for more than 1 sec, just return 429?
		if time.Since(requestStart) > time.Second {
			c.AbortWithStatus(http.StatusTooManyRequests)
		} else {
			c.Next()
		}
	}
}
