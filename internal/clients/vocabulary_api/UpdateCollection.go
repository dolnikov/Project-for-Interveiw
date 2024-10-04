package vocabulary_api

import (
	"context"
	VocabularyApiProto "gitlab.com/wbwapis/go-genproto/wbw/vocabulary/vocabulary_api/v1"
	"gitlab.com/wordbyword.io/microservices/pkg/constants"
	"gitlab.com/wordbyword.io/microservices/pkg/utils"
)

func (c *VocabularyApiClient) UpdateCollection(ctx context.Context, request *VocabularyApiProto.UpdateCollectionRequest) (response *VocabularyApiProto.UpdateCollectionResponse, err error) {
	requestId := utils.AnyToString(ctx.Value(constants.RequestIdKey))
	lgr := c.lgr.With().
		Str("api", "UpdateCollection").
		Str(constants.RequestIdKey, requestId).
		Interface("request", request).Logger()

	conn, err := c.pool.CreateConn(ctx)
	if err != nil {
		lgr.Error().Err(err).Msg("failed to connect to grpc connection worker pool")
		return nil, err
	}
	defer conn.Close()

	response, err = c.createVocabularyServiceClient(conn).
		UpdateCollection(ctx, request)
	if err != nil {
		lgr.Error().Err(err).Msg("response error")
		return nil, err
	}

	lgr.Debug().Interface("response", response).Msg("executed")

	return response, nil
}