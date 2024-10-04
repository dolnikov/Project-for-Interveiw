package action_api

import (
	"context"
	"crypto/tls"
	grpcpool "github.com/processout/grpc-go-pool"
	"github.com/rs/zerolog"
	ActionApiProto "gitlab.com/wbwapis/go-genproto/wbw/action/action_api/v1"
	"gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/config"
	_grpc "gitlab.com/wordbyword.io/microservices/pkg/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

const clientName = "action_api"

type ActionApiClient struct {
	cfg  *config.Config
	lgr  zerolog.Logger
	pool *_grpc.Pool
}

func NewActionApiClient(cfg *config.Config, lgr zerolog.Logger) *ActionApiClient {
	lgr = lgr.With().Str("client", clientName).Logger()

	var credentialsOption grpc.DialOption
	if cfg.ActionApi.WithTLS {
		credentialsOption = grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: true,
		}))
	} else {
		credentialsOption = grpc.WithTransportCredentials(insecure.NewCredentials())
	}

	factory := func() (*grpc.ClientConn, error) {
		ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
		conn, err := grpc.DialContext(ctx, cfg.ActionApi.URI,
			grpc.WithDefaultCallOptions(
				grpc.MaxCallRecvMsgSize(4194304),
				grpc.MaxCallSendMsgSize(4194304),
			),
			credentialsOption,
			grpc.WithChainUnaryInterceptor(
				_grpc.AddRequestIdToOutgoingContext,
				_grpc.AddAcceptLanguageToOutgoingContext,
			),
		)
		if err != nil {
			lgr.Fatal().Err(err).Msg(clientName + " connection failed")
		}
		return conn, err
	}

	pool, err := grpcpool.New(
		factory,
		10,
		10,
		30*time.Second,
		30*time.Second,
	)
	if err != nil {
		lgr.Fatal().Err(err).Msg(clientName + " create connection pool failed")
	}

	return &ActionApiClient{
		cfg:  cfg,
		lgr:  lgr,
		pool: &_grpc.Pool{Pool: pool},
	}
}

func (c *ActionApiClient) createActionServiceClient(conn *grpcpool.ClientConn) ActionApiProto.ActionApiClient {
	return ActionApiProto.NewActionApiClient(conn)
}

func (c *ActionApiClient) Shutdown() {
	c.pool.Close()
}
