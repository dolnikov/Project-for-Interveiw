package vocabulary_api

import (
	"context"
	VocabularyApiProto "gitlab.com/wbwapis/go-genproto/wbw/vocabulary/vocabulary_api/v1"
	"gitlab.com/wordbyword.io/microservices/pkg/constants"
	"gitlab.com/wordbyword.io/microservices/pkg/utils"
)

func (c *VocabularyApiClient) GetCollections(ctx context.Context, request *VocabularyApiProto.GetCollectionsRequest) (response *VocabularyApiProto.GetCollectionsResponse, err error) {
	requestId := utils.AnyToString(ctx.Value(constants.RequestIdKey))
	lgr := c.lgr.With().
		Str("api", "GetCollections").
		Str(constants.RequestIdKey, requestId).
		Interface("request", request).Logger()

	conn, err := c.pool.CreateConn(ctx)
	if err != nil {
		lgr.Error().Err(err).Msg("failed to connect to grpc connection worker pool")
		return nil, err
	}
	defer conn.Close()

	response, err = c.createVocabularyServiceClient(conn).
		GetCollections(ctx, request)
	if err != nil {
		lgr.Error().Err(err).Msg("response error")
		return nil, err
	}

	lgr.Debug().Msg("executed")

	return response, nil
}
