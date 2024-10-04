package http

import (
	"context"
	"errors"
	"github.com/gin-contrib/pprof"
	limits "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/config"
	"gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/prometheus"
	"gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/transports/http/healthcheck"
	"gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/transports/http/metrics"
	"gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/transports/http/middleware/cors"
	"gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/transports/http/middleware/recovery"
	"net"
	"net/http"
	"os"
)

type IServer interface {
	Start() chan error
	Shutdown(ctx context.Context) error
}

type Server struct {
	cfg      *config.Config
	lgr      zerolog.Logger
	listener net.Listener
	srv      *http.Server
}

// Compile time assertion that Server implements IServer.
var _ IServer = (*Server)(nil)

type Endpointer interface {
	RegisterServer(r *gin.Engine, prefixOptions ...string)
}

func NewServer(
	cfg *config.Config,
	lgr zerolog.Logger,
	prom *prometheus.Exporter,
	listener net.Listener,
	ep Endpointer,
) (*Server, error) {
	if ep == nil {
		return nil, errors.New("nil Endpointer not allowed")
	}

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = os.Stdout

	router := gin.New()

	router.Use(gin.CustomRecovery(recovery.NewRecoverer(prom).RecoveryFunc))
	router.Use(cors.Cors(cfg))
	router.Use(limits.RequestSizeLimiter(cfg.HTTP.MaxRequestBodySize))

	if cfg.Profiling.Enabled {
		pprof.Register(router, "/debug/pprof")
	}
	if cfg.Metrics.Enabled {
		metrics.Register(router, "/metrics")
	}
	healthcheck.Register(cfg, router, "/")

	ep.RegisterServer(router, "/")

	httpSrv := &Server{
		cfg:      cfg,
		lgr:      lgr,
		listener: listener,
		srv: &http.Server{
			Handler:      router,
			TLSConfig:    nil,
			ReadTimeout:  cfg.HTTP.ReadTimeout,
			WriteTimeout: cfg.HTTP.WriteTimeout,
			IdleTimeout:  cfg.HTTP.IdleTimeout,
		},
	}

	return httpSrv, nil
}

func (g *Server) Start() chan error {
	listenErrCh := make(chan error, 1)
	g.lgr.Info().Msgf("starting http server on %s", g.listener.Addr())
	go func() {
		if err := g.srv.Serve(g.listener); err != nil {
			g.lgr.Error().Stack().Err(err).Msg("received listener error")
			listenErrCh <- err
		}
	}()
	return listenErrCh
}

func (g *Server) Shutdown(ctx context.Context) error {
	g.lgr.Info().Msg("stopping http server")
	defer g.lgr.Info().Msg("http server stopped")
	return g.srv.Shutdown(ctx)
}
