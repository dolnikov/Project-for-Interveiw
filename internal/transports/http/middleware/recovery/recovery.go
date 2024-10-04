package recovery

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/errors"
	"gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/prometheus"
	"gitlab.com/wordbyword.io/microservices/pkg/errors/outer"
	"strconv"
)

type Recoverer struct {
	prom *prometheus.Exporter
}

func NewRecoverer(prom *prometheus.Exporter) *Recoverer {
	return &Recoverer{prom: prom}
}

func (r *Recoverer) RecoveryFunc(c *gin.Context, recovered interface{}) {
	if err, ok := recovered.(error); ok {
		code, obj := outer.GetHTTPError(errors.InternalError(err))
		r.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
	}
	c.Abort()
}
