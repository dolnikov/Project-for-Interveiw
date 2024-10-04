package vocabulary_api

import (
	"context"
	VocabularyApiProto "gitlab.com/wbwapis/go-genproto/wbw/vocabulary/vocabulary_api/v1"
	"gitlab.com/wordbyword.io/microservices/pkg/constants"
	"gitlab.com/wordbyword.io/microservices/pkg/utils"
)

func (c *VocabularyApiClient) GetCollection(ctx context.Context, request *VocabularyApiProto.GetCollectionRequest) (response *VocabularyApiProto.GetCollectionResponse, err error) {
	requestId := utils.AnyToString(ctx.Value(constants.RequestIdKey))
	lgr := c.lgr.With().
		Str("api", "GetCollection").
		Str(constants.RequestIdKey, requestId).
		Interface("request", request).Logger()

	conn, err := c.pool.CreateConn(ctx)
	if err != nil {
		lgr.Error().Err(err).Msg("failed to connect to grpc connection worker pool")
		return nil, err
	}
	defer conn.Close()

	response, err = c.createVocabularyServiceClient(conn).
		GetCollection(ctx, request)
	if err != nil {
		lgr.Error().Err(err).Msg("response error")
		return nil, err
	}

	lgr.Debug().Interface("response", response).Msg("executed")

	return response, nil
}
