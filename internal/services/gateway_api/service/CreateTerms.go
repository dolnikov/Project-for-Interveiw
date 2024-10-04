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

func (s *GatewayApiService) CreateTerms(ctx context.Context, request *GatewayApiProto.CreateTermsRequest) (*GatewayApiProto.CreateTermsResponse, error) {
	requestId := utils.AnyToString(ctx.Value(constants.RequestIdKey))
	lgr := s.lgr.With().
		Str(constants.RequestIdKey, requestId).
		Str("api", "CreateTerms").
		Interface("request", request).
		Logger()

	tokenClaims, ok := ctx.Value(constants.TokenClaimsKey).(*_jwt.TokenClaims)
	if !ok {
		outerErr := _errors.TokenClaimsDoesNotSet
		lgr.Error().Msg(outerErr.ErrorMessage)
		return nil, outerErr
	}

	for i, term := range request.Terms {
		if i == 0 {
			continue
		}

		if request.Terms[i-1].CollectionId != term.CollectionId {
			outerErr := _errors.FailedToCreateTermsCollection
			lgr.Error().Msg(outerErr.ErrorMessage)
			return nil, outerErr
		}
	}

	getCollection, err := s.vocabularyApiClient.GetCollection(ctx, &VocabularyApiProto.GetCollectionRequest{
		CollectionId: request.Terms[0].CollectionId,
	})
	if err != nil {
		outerErr := _errors.FailedToGetCollection
		lgr.Error().Err(err).Msg(outerErr.ErrorMessage)
		return nil, outerErr
	}

	if getCollection.Collection.UserId != tokenClaims.UserId {
		outerErr := _errors.FailedToGetCollectionPrivateCollection
		lgr.Error().Err(err).Msg(outerErr.ErrorMessage)
		return nil, outerErr
	}

	terms := make([]*VocabularyApiProto.CreateTerms, len(request.Terms))
	for i, term := range request.Terms {
		terms[i] = &VocabularyApiProto.CreateTerms{
			CollectionId:      term.CollectionId,
			TermLanguageId:    term.TermLanguageId,
			MeaningLanguageId: term.MeaningLanguageId,
			Term:              term.Term,
			Meaning:           term.Meaning,
			Example:           term.Example,
			ImageUrl:          term.ImageUrl,
		}
	}

	_, err = s.vocabularyApiClient.CreateTerms(ctx, &VocabularyApiProto.CreateTermsRequest{
		Terms: terms,
	})
	if err != nil {
		outerErr := _errors.FailedToCreateTerms
		lgr.Error().Err(err).Msg(outerErr.ErrorMessage)
		return nil, outerErr
	}

	resp := &GatewayApiProto.CreateTermsResponse{}

	lgr.Debug().Msg("executed")
	return resp, nil
}
