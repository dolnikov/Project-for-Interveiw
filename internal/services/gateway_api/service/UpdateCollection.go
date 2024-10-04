package service

import (
	"context"
	VocabularyApiProto "gitlab.com/wbwapis/go-genproto/wbw/vocabulary/vocabulary_api/v1"
	_errors "gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/errors"
	"gitlab.com/wordbyword.io/microservices/pkg/constants"
	_jwt "gitlab.com/wordbyword.io/microservices/pkg/jwt"
	"gitlab.com/wordbyword.io/microservices/pkg/utils"
	"strings"

	GatewayApiProto "gitlab.com/wbwapis/go-genproto/wbw/gateway/gateway_api/v1"
)

func (s *GatewayApiService) UpdateCollection(ctx context.Context, request *GatewayApiProto.UpdateCollectionRequest) (*GatewayApiProto.UpdateCollectionResponse, error) {
	requestId := utils.AnyToString(ctx.Value(constants.RequestIdKey))
	lgr := s.lgr.With().
		Str(constants.RequestIdKey, requestId).
		Str("api", "UpdateCollection").
		Interface("request", request).
		Logger()

	tokenClaims, ok := ctx.Value(constants.TokenClaimsKey).(*_jwt.TokenClaims)
	if !ok {
		outerErr := _errors.TokenClaimsDoesNotSet
		lgr.Error().Msg(outerErr.ErrorMessage)
		return nil, outerErr
	}

	req := VocabularyApiProto.UpdateCollectionRequest{
		UserId:       tokenClaims.UserId,
		CollectionId: request.CollectionId,
		LanguageId:   request.LanguageId,
		IsPublic:     request.IsPublic,
		IsPinned:     request.IsPinned,
	}

	if request.Name != nil {
		name := strings.TrimSpace(*request.Name)
		req.Name = &name
	}

	if request.Description != nil {
		description := strings.TrimSpace(*request.Description)
		req.Description = &description
	}

	_, err := s.vocabularyApiClient.UpdateCollection(ctx, &req)
	if err != nil {
		outerErr := _errors.FailedToUpdateCollection
		lgr.Error().Err(err).Msg(outerErr.ErrorMessage)
		return nil, outerErr
	}

	resp := &GatewayApiProto.UpdateCollectionResponse{}
	lgr.Debug().Msg("executed")
	return resp, nil
}
