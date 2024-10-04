package cors

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/config"
	"net/http"
)

func Cors(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cfg.HTTPCors.Enabled {
			c.Writer.Header().Set("Access-Control-Allow-Origin", cfg.HTTPCors.AllowedOrigins)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", cfg.HTTPCors.AllowCredentials)
			c.Writer.Header().Set("Access-Control-Allow-Headers", cfg.HTTPCors.AllowedHeaders)
			c.Writer.Header().Set("Access-Control-Allow-Methods", cfg.HTTPCors.AllowedMethods)
			c.Writer.Header().Set("Access-Control-Max-Age", cfg.HTTPCors.MaxAge)

			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(http.StatusNoContent)
				return
			}
		}

		c.Next()
	}
}
