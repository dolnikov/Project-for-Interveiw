package translation_api

import (
	"context"
	TranslationApiProto "gitlab.com/wbwapis/go-genproto/wbw/translation/translation_api/v1"
	"gitlab.com/wordbyword.io/microservices/pkg/constants"
	"gitlab.com/wordbyword.io/microservices/pkg/utils"
)

func (c *TranslationApiClient) GetTranslation(ctx context.Context, request *TranslationApiProto.GetTranslationRequest) (response *TranslationApiProto.GetTranslationResponse, err error) {
	requestId := utils.AnyToString(ctx.Value(constants.RequestIdKey))
	lgr := c.lgr.With().
		Str("api", "GetTranslation").
		Str(constants.RequestIdKey, requestId).
		Interface("request", request).Logger()

	conn, err := c.pool.CreateConn(ctx)
	if err != nil {
		lgr.Error().Err(err).Msg("failed to connect to grpc connection worker pool")
		return nil, err
	}
	defer conn.Close()

	response, err = c.createTranslationServiceClient(conn).
		GetTranslation(ctx, request)
	if err != nil {
		lgr.Error().Err(err).Msg("response error")
		return nil, err
	}

	lgr.Debug().Msg("executed")

	return response, nil
}
