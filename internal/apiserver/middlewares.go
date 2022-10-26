package apiserver

import (
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func corsMiddleware(a *APIServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		switch a.config.Server.CORS.SameSite {
		case "default":
			c.SetSameSite(http.SameSiteDefaultMode)
		case "lax":
			c.SetSameSite(http.SameSiteLaxMode)
		case "strict":
			c.SetSameSite(http.SameSiteStrictMode)
		case "none":
			c.SetSameSite(http.SameSiteNoneMode)
		}
		c.Header("Access-Control-Allow-Origin", a.config.Server.CORS.AccessControlAllowOrigin)
		if a.config.Server.CORS.AccessControlAllowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}
	}
}

func loggerMiddleware(a *APIServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		ip := net.ParseIP(c.ClientIP())
		if ip == nil {
			status := http.StatusForbidden
			a.logger.Warn().
				Dur("latency", time.Since(start)).
				Str("clientIP", c.ClientIP()).
				Int("status", status).
				Msg("could not parse client IP")
			c.AbortWithStatus(status)
			return
		}

		requestID := uuid.New()

		c.Set("requestID", requestID)

		c.Next()

		if raw != "" {
			path = path + "?" + raw
		}

		a.logger.Info().
			Dur("latency", time.Since(start)).
			Str("clientIP", c.ClientIP()).
			Int("status", c.Writer.Status()).
			Int("size", c.Writer.Size()).
			Str("requestID", requestID.String()).
			Str("endpoint", path).Send()
	}
}
