package language_api

import (
	"context"
	LanguageApiProto "gitlab.com/wbwapis/go-genproto/wbw/language/language_api/v1"
	"gitlab.com/wordbyword.io/microservices/pkg/constants"
	"gitlab.com/wordbyword.io/microservices/pkg/utils"
)

func (c *LanguageApiClient) GetLanguages(ctx context.Context, request *LanguageApiProto.GetLanguagesRequest) (response *LanguageApiProto.GetLanguagesResponse, err error) {
	requestId := utils.AnyToString(ctx.Value(constants.RequestIdKey))
	lgr := c.lgr.With().
		Str("api", "GetLanguages").
		Str(constants.RequestIdKey, requestId).
		Interface("request", request).Logger()

	conn, err := c.pool.CreateConn(ctx)
	if err != nil {
		lgr.Error().Err(err).Msg("failed to connect to grpc connection worker pool")
		return nil, err
	}
	defer conn.Close()

	response, err = c.createLanguageServiceClient(conn).
		GetLanguages(ctx, request)
	if err != nil {
		lgr.Error().Err(err).Msg("response error")
		return nil, err
	}

	lgr.Debug().Msg("executed")

	return response, nil
}
