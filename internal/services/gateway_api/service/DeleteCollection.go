package service

import (
	"context"
	VocabularyApiProto "gitlab.com/wbwapis/go-genproto/wbw/vocabulary/vocabulary_api/v1"
	_errors "gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/errors"
	"gitlab.com/wordbyword.io/microservices/pkg/constants"
	_jwt "gitlab.com/wordbyword.io/microservices/pkg/jwt"
	"gitlab.com/wordbyword.io/microservices/pkg/utils"

	GatewayApiProto "gitlab.com/wbwapis/go-genproto/wbw/gateway/gateway_api/v1"
)

func (s *GatewayApiService) DeleteCollection(ctx context.Context, request *GatewayApiProto.DeleteCollectionRequest) (*GatewayApiProto.DeleteCollectionResponse, error) {
	requestId := utils.AnyToString(ctx.Value(constants.RequestIdKey))
	lgr := s.lgr.With().
		Str(constants.RequestIdKey, requestId).
		Str("api", "DeleteCollection").
		Interface("request", request).
		Logger()

	tokenClaims, ok := ctx.Value(constants.TokenClaimsKey).(*_jwt.TokenClaims)
	if !ok {
		outerErr := _errors.TokenClaimsDoesNotSet
		lgr.Error().Msg(outerErr.ErrorMessage)
		return nil, outerErr
	}

	_, err := s.vocabularyApiClient.DeleteCollection(ctx, &VocabularyApiProto.DeleteCollectionRequest{
		UserId:       tokenClaims.UserId,
		CollectionId: request.CollectionId,
	})
	if err != nil {
		outerErr := _errors.FailedToDeleteCollection
		lgr.Error().Err(err).Msg(outerErr.ErrorMessage)
		return nil, outerErr
	}

	lgr.Debug().Msg("executed")
	return &GatewayApiProto.DeleteCollectionResponse{}, nil
}