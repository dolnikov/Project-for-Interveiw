package notification_api

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/wagslane/go-rabbitmq"
	"gitlab.com/wordbyword.io/microservices/pkg/constants"
	"gitlab.com/wordbyword.io/microservices/pkg/utils"
	"google.golang.org/protobuf/proto"

	NotificationApiProto "gitlab.com/wbwapis/go-genproto/wbw/notification/notification_api/v1"
)

// SendEmail send message for notification worker
func (c *NotificationApiClient) SendEmail(ctx context.Context, lgr zerolog.Logger, request *NotificationApiProto.SendEmailRequest) (err error) {
	requestId := utils.AnyToString(ctx.Value(constants.RequestIdKey))
	acceptLanguage := utils.AnyToString(ctx.Value(constants.AcceptLanguageKey))
	lgr = lgr.With().
		Str("consumer", NotificationApiConsumer).
		Str("api", "SendEmail").
		Interface("request", request).
		Str(constants.RequestIdKey, requestId).
		Str(constants.AcceptLanguageKey, acceptLanguage).
		Logger()

	requestBytes, _ := proto.Marshal(request)
	err = c.publisher.Publish(
		requestBytes,
		[]string{c.cfg.Rabbit.NotificationApiSendEmail.Queue},
		rabbitmq.WithPublishOptionsContentType("application/x-protobuf"),
		rabbitmq.WithPublishOptionsExchange(c.cfg.Rabbit.NotificationApiSendEmail.Exchange),
		rabbitmq.WithPublishOptionsHeaders(map[string]interface{}{
			constants.RequestIdKey:      requestId,
			constants.AcceptLanguageKey: acceptLanguage,
		}),
		rabbitmq.WithPublishOptionsMandatory,
		rabbitmq.WithPublishOptionsPersistentDelivery,
	)
	if err != nil {
		lgr.Error().Err(err).Msg("AMQP publish error")
		return err
	}

	lgr.Debug().Msg("executed")

	return nil
}
