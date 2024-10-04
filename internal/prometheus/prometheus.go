package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Exporter struct {
	HttpReqDuration *prometheus.HistogramVec
	HttpRespCount   *prometheus.CounterVec

	AuthApiReqDuration *prometheus.HistogramVec
	AuthApiReqErrCount *prometheus.CounterVec

	UserApiReqDuration *prometheus.HistogramVec
	UserApiReqErrCount *prometheus.CounterVec

	ActionApiReqDuration *prometheus.HistogramVec
	ActionApiReqErrCount *prometheus.CounterVec
}

func NewExporter(namespace string) *Exporter {
	prom := new(Exporter)

	reqDurBuckets := []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}

	prom.HttpReqDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Name:      "http_request_duration",
		Help:      "http request duration",
	},
		[]string{"method"},
	)

	prom.HttpRespCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "http_response_count",
		Help:      "http response count",
	},
		[]string{"code"},
	)

	prom.AuthApiReqDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Name:      "auth_api_request_duration",
		Help:      "auth_api request duration",
		Buckets:   reqDurBuckets,
	},
		[]string{"method"},
	)
	prom.AuthApiReqErrCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "auth_api_request_err_count",
		Help:      "auth_api request err count",
	},
		[]string{"method"},
	)

	prom.UserApiReqDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Name:      "user_api_request_duration",
		Help:      "user_api request duration",
		Buckets:   reqDurBuckets,
	},
		[]string{"method"},
	)
	prom.UserApiReqErrCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "user_api_request_err_count",
		Help:      "user_api request err count",
	},
		[]string{"method"},
	)

	prom.ActionApiReqDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Name:      "action_api_request_duration",
		Help:      "action_api request duration",
		Buckets:   reqDurBuckets,
	},
		[]string{"method"},
	)
	prom.ActionApiReqErrCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "action_api_request_err_count",
		Help:      "action_api request err count",
	},
		[]string{"method"},
	)

	prometheus.MustRegister(
		prom.HttpReqDuration, prom.HttpRespCount,
		prom.AuthApiReqDuration, prom.AuthApiReqErrCount,
		prom.UserApiReqDuration, prom.UserApiReqErrCount,
		prom.ActionApiReqDuration, prom.ActionApiReqErrCount,
	)

	return prom
}
