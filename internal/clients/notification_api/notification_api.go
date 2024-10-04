package notification_api

import (
	"github.com/rs/zerolog"
	"github.com/wagslane/go-rabbitmq"
	"gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/config"
	"time"
)

const clientName = "rabbit_mq"

const (
	NotificationApiConsumer = "notification-api"
)

type NotificationApiClient struct {
	cfg        *config.Config
	lgr        zerolog.Logger
	connection *rabbitmq.Conn
	publisher  *rabbitmq.Publisher
}

func NewNotificationApiClient(cfg *config.Config, lgr zerolog.Logger) *NotificationApiClient {
	lgr = lgr.With().Str("client", clientName).Logger()
	connection, err := rabbitmq.NewConn(
		cfg.Rabbit.URI,
		rabbitmq.WithConnectionOptionsLogging,
		rabbitmq.WithConnectionOptionsReconnectInterval(5*time.Second),
	)
	if err != nil {
		lgr.Fatal().Err(err).Msg("AMQP connection error")
	}

	publisher, err := rabbitmq.NewPublisher(
		connection,
		rabbitmq.WithPublisherOptionsLogging,
	)
	if err != nil {
		lgr.Fatal().Err(err).Msg("AMQP connection error")
	}

	return &NotificationApiClient{
		cfg:        cfg,
		lgr:        lgr,
		connection: connection,
		publisher:  publisher,
	}
}

func (c *NotificationApiClient) Shutdown() {
	_ = c.connection.Close()
	c.publisher.Close()
}
