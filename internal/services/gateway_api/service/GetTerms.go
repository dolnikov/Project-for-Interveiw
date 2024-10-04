package service

import (
	"context"
	VocabularyApiProto "gitlab.com/wbwapis/go-genproto/wbw/vocabulary/vocabulary_api/v1"
	_errors "gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/errors"
	"gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/services/gateway_api/mappers"
	"gitlab.com/wordbyword.io/microservices/pkg/constants"
	_jwt "gitlab.com/wordbyword.io/microservices/pkg/jwt"
	"gitlab.com/wordbyword.io/microservices/pkg/utils"

	GatewayApiProto "gitlab.com/wbwapis/go-genproto/wbw/gateway/gateway_api/v1"
)

func (s *GatewayApiService) GetTerms(ctx context.Context, request *GatewayApiProto.GetTermsRequest) (*GatewayApiProto.GetTermsResponse, error) {
	requestId := utils.AnyToString(ctx.Value(constants.RequestIdKey))
	lgr := s.lgr.With().
		Str(constants.RequestIdKey, requestId).
		Str("api", "GetTerms").
		Interface("request", request).
		Logger()

	tokenClaims, ok := ctx.Value(constants.TokenClaimsKey).(*_jwt.TokenClaims)
	if !ok {
		outerErr := _errors.TokenClaimsDoesNotSet
		lgr.Error().Msg(outerErr.ErrorMessage)
		return nil, outerErr
	}

	getCollection, err := s.vocabularyApiClient.GetCollection(ctx, &VocabularyApiProto.GetCollectionRequest{
		CollectionId: request.CollectionId,
	})
	if err != nil {
		outerErr := _errors.FailedToGetCollection
		lgr.Error().Err(err).Msg(outerErr.ErrorMessage)
		return nil, outerErr
	}

	if getCollection.Collection.UserId != tokenClaims.UserId && getCollection.Collection.IsPublic != true {
		outerErr := _errors.FailedToGetCollectionPrivateCollection
		lgr.Error().Err(err).Msg(outerErr.ErrorMessage)
		return nil, outerErr
	}

	getCollections, err := s.vocabularyApiClient.GetTerms(ctx, &VocabularyApiProto.GetTermsRequest{
		CollectionId: request.CollectionId,
	})
	if err != nil {
		outerErr := _errors.FailedToGetTerms
		lgr.Error().Err(err).Msg(outerErr.ErrorMessage)
		return nil, outerErr
	}

	var terms = make([]*GatewayApiProto.Term, len(getCollections.Terms))
	for i, term := range getCollections.Terms {
		terms[i] = &GatewayApiProto.Term{
			TermId:            term.TermId,
			CollectionId:      term.CollectionId,
			TermLanguageId:    term.TermLanguageId,
			MeaningLanguageId: term.MeaningLanguageId,
			Term:              term.Term,
			Meaning:           term.Meaning,
			Example:           term.Example,
			ImageUrl:          term.ImageUrl,
			Status:            mappers.VocabularyStatusToGatewayStatus[term.Status],
			IsPhrase:          term.IsPhrase,
			RepeatedAt:        term.RepeatedAt,
			CreatedAt:         term.CreatedAt,
			UpdatedAt:         term.UpdatedAt,
		}
	}

	resp := &GatewayApiProto.GetTermsResponse{
		Terms: terms,
	}

	lgr.Debug().Msg("executed")
	return resp, nil
}
