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

func (s *GatewayApiService) GetCollection(ctx context.Context, request *GatewayApiProto.GetCollectionRequest) (*GatewayApiProto.GetCollectionResponse, error) {
	requestId := utils.AnyToString(ctx.Value(constants.RequestIdKey))
	lgr := s.lgr.With().
		Str(constants.RequestIdKey, requestId).
		Str("api", "GetCollection").
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

	resp := &GatewayApiProto.GetCollectionResponse{
		Collection: &GatewayApiProto.Collection{
			CollectionId: getCollection.Collection.CollectionId,
			UserId:       getCollection.Collection.UserId,
			LanguageId:   getCollection.Collection.LanguageId,
			Name:         getCollection.Collection.Name,
			Description:  getCollection.Collection.Description,
			IsPinned:     getCollection.Collection.IsPinned,
			IsPublic:     getCollection.Collection.IsPublic,
			TotalTerms:   getCollection.Collection.TotalTerms,
			LearnedTerms: getCollection.Collection.LearnedTerms,
			OpenedAt:     getCollection.Collection.OpenedAt,
			CreatedAt:    getCollection.Collection.CreatedAt,
			UpdatedAt:    getCollection.Collection.UpdatedAt,
		},
	}

	lgr.Debug().Msg("executed")
	return resp, nil
}
