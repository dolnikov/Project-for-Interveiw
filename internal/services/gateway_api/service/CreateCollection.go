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

func (s *GatewayApiService) CreateCollection(ctx context.Context, request *GatewayApiProto.CreateCollectionRequest) (*GatewayApiProto.CreateCollectionResponse, error) {
	requestId := utils.AnyToString(ctx.Value(constants.RequestIdKey))
	lgr := s.lgr.With().
		Str(constants.RequestIdKey, requestId).
		Str("api", "CreateCollection").
		Interface("request", request).
		Logger()

	tokenClaims, ok := ctx.Value(constants.TokenClaimsKey).(*_jwt.TokenClaims)
	if !ok {
		outerErr := _errors.TokenClaimsDoesNotSet
		lgr.Error().Msg(outerErr.ErrorMessage)
		return nil, outerErr
	}

	createCollection, err := s.vocabularyApiClient.CreateCollection(ctx, &VocabularyApiProto.CreateCollectionRequest{
		UserId:      tokenClaims.UserId,
		LanguageId:  request.LanguageId,
		Name:        request.Name,
		Description: request.Description,
		IsPublic:    request.IsPublic,
	})
	if err != nil {
		outerErr := _errors.FailedToCreateCollection
		lgr.Error().Err(err).Msg(outerErr.ErrorMessage)
		return nil, outerErr
	}

	resp := &GatewayApiProto.CreateCollectionResponse{
		CollectionId: createCollection.CollectionId,
	}

	lgr.Debug().Msg("executed")
	return resp, nil
}
