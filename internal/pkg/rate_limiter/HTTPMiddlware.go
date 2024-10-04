package rate_limiter

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/config"
	"gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/errors"
	"gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/prometheus"
	"gitlab.com/wordbyword.io/microservices/pkg/constants"
	"gitlab.com/wordbyword.io/microservices/pkg/errors/outer"
	"gitlab.com/wordbyword.io/microservices/pkg/utils"
	"strconv"
)

func HttpMiddleware(cfg *config.Config, lgr zerolog.Logger, prom *prometheus.Exporter, rl *RateLimiter) gin.HandlerFunc {
	var err error
	return func(c *gin.Context) {

		if cfg.RateLimiter.Enabled {
			clientIP := utils.AnyToString(c.Value(constants.ClientIPKey))
			if !rl.Allow(c.Request.RequestURI, clientIP) {
				err = errors.TooManyRequests
				lgr.Error().Err(err).Msg(errors.TooManyRequestsMsg)
				code, obj := outer.GetHTTPError(err)
				prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
				c.PureJSON(code, obj)
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
