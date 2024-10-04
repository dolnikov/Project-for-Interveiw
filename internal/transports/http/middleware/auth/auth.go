package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/config"
	_errors "gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/errors"
	"gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/prometheus"
	"gitlab.com/wordbyword.io/microservices/pkg/constants"
	"gitlab.com/wordbyword.io/microservices/pkg/errors/outer"
	_jwt "gitlab.com/wordbyword.io/microservices/pkg/jwt"
	"strconv"
	"strings"
)

type authHeader struct {
	Token string `header:"Authorization"`
}

func Auth(cfg *config.Config, lgr zerolog.Logger, prom *prometheus.Exporter) gin.HandlerFunc {
	var err error
	return func(c *gin.Context) {
		h := new(authHeader)
		if err = c.ShouldBindHeader(h); err != nil {
			outerErr := _errors.BadAuthorizationTokenError(_errors.FailedToGetAuthorizationToken)
			code, obj := outer.GetHTTPError(outerErr)
			lgr.Warn().Err(err).Msg(outerErr.ErrorMessage)
			prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
			c.PureJSON(code, obj)
			c.Abort()
			return
		}

		if h.Token == "" {
			outerErr := _errors.BadAuthorizationTokenError(_errors.AuthorizationTokenIsEmpty)
			code, obj := outer.GetHTTPError(outerErr)
			lgr.Warn().Err(err).Msg(outerErr.ErrorMessage)
			prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
			c.PureJSON(code, obj)
			c.Abort()
			return
		}

		parts := strings.Split(h.Token, " ")
		if len(parts) != 2 {
			outerErr := _errors.BadAuthorizationTokenError(_errors.AuthorizationTokenIsInvalid)
			code, obj := outer.GetHTTPError(outerErr)
			lgr.Warn().Err(err).Msg(outerErr.ErrorMessage)
			prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
			c.PureJSON(code, obj)
			c.Abort()
		}
		h.Token = parts[1]

		var tokenClaims *_jwt.TokenClaims
		tokenClaims, err = _jwt.VerifyToken(h.Token, cfg.JWT.AccessSecret)
		if err != nil {
			outerErr := _errors.BadAuthorizationTokenError(_errors.AuthorizationTokenIsInvalid)
			code, obj := outer.GetHTTPError(outerErr)
			lgr.Warn().Err(err).Msg(outerErr.ErrorMessage)
			prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
			c.PureJSON(code, obj)
			c.Abort()
			return
		}

		c.Set(constants.TokenClaimsKey, tokenClaims)
		c.Next()
	}
}
