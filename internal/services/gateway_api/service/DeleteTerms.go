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

func (s *GatewayApiService) DeleteTerms(ctx context.Context, request *GatewayApiProto.DeleteTermsRequest) (*GatewayApiProto.DeleteTermsResponse, error) {
	requestId := utils.AnyToString(ctx.Value(constants.RequestIdKey))
	lgr := s.lgr.With().
		Str(constants.RequestIdKey, requestId).
		Str("api", "DeleteTerms").
		Interface("request", request).
		Logger()

	tokenClaims, ok := ctx.Value(constants.TokenClaimsKey).(*_jwt.TokenClaims)
	if !ok {
		outerErr := _errors.TokenClaimsDoesNotSet
		lgr.Error().Msg(outerErr.ErrorMessage)
		return nil, outerErr
	}

	termIds := make([]uint64, len(request.TermIds))
	for i, termId := range request.TermIds {
		termIds[i] = termId
	}

	_, err := s.vocabularyApiClient.DeleteTerms(ctx, &VocabularyApiProto.DeleteTermsRequest{
		UserId:       tokenClaims.UserId,
		CollectionId: request.CollectionId,
		TermIds:      termIds,
	})
	if err != nil {
		outerErr := _errors.FailedToDeleteTerms
		lgr.Error().Err(err).Msg(outerErr.ErrorMessage)
		return nil, outerErr
	}

	lgr.Debug().Msg("executed")
	return &GatewayApiProto.DeleteTermsResponse{}, nil
}
