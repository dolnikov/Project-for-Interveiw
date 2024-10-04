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

func (s *GatewayApiService) UpdateTerm(ctx context.Context, request *GatewayApiProto.UpdateTermRequest) (*GatewayApiProto.UpdateTermResponse, error) {
	requestId := utils.AnyToString(ctx.Value(constants.RequestIdKey))
	lgr := s.lgr.With().
		Str(constants.RequestIdKey, requestId).
		Str("api", "UpdateTerm").
		Interface("request", request).
		Logger()

	tokenClaims, ok := ctx.Value(constants.TokenClaimsKey).(*_jwt.TokenClaims)
	if !ok {
		outerErr := _errors.TokenClaimsDoesNotSet
		lgr.Error().Msg(outerErr.ErrorMessage)
		return nil, outerErr
	}

	req := VocabularyApiProto.UpdateTermRequest{
		UserId:            tokenClaims.UserId,
		TermId:            request.TermId,
		TermLanguageId:    request.TermLanguageId,
		MeaningLanguageId: request.MeaningLanguageId,
		Term:              request.Term,
		Meaning:           request.Meaning,
		Example:           request.Example,
		ImageUrl:          request.ImageUrl,
	}

	if request.Term != nil {
		term := strings.TrimSpace(*request.Term)
		req.Term = &term
	}

	if request.Meaning != nil {
		meaning := strings.TrimSpace(*request.Meaning)
		req.Meaning = &meaning
	}

	if request.Example != nil {
		example := strings.TrimSpace(*request.Example)
		req.Example = &example
	}

	_, err := s.vocabularyApiClient.UpdateTerm(ctx, &req)
	if err != nil {
		outerErr := _errors.FailedToUpdateTerm
		lgr.Error().Err(err).Msg(outerErr.ErrorMessage)
		return nil, outerErr
	}

	resp := &GatewayApiProto.UpdateTermResponse{}
	lgr.Debug().Msg("executed")
	return resp, nil
}
