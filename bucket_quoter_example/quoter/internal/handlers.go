package internal

import (
	"github.com/gin-gonic/gin"
)

// ping for testing
func (a *Api) getPing(c *gin.Context) {
	if a.core == nil {
		a.apiSendError(c, 502, "Internal error")
		return
	}

	a.apiSendOK(c, 200, "pong")
}

// api limiter
func (a *Api) isAPIAvailableWithLimiter(c *gin.Context) {
	if a.core == nil {
		a.apiSendError(c, 502, "Internal error")
		return
	}

	if limiter, ok := a.limiterMap[c.Request.Header.Get("X-Limiter-Subscription-ID")]; ok {
		if !limiter.IsAvailable() {
			a.apiSendError(c, 429, "Too Many Requests")
			return
		}
		limiter.UseWithSleep(1)
	} else {
		a.apiSendError(c, 503, "Service Unavailable")
		return
	}

	a.apiSendOK(c, 200, "")
}
