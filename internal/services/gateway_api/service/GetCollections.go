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

func (s *GatewayApiService) GetCollections(ctx context.Context, request *GatewayApiProto.GetCollectionsRequest) (*GatewayApiProto.GetCollectionsResponse, error) {
	requestId := utils.AnyToString(ctx.Value(constants.RequestIdKey))
	lgr := s.lgr.With().
		Str(constants.RequestIdKey, requestId).
		Str("api", "GetCollections").
		Interface("request", request).
		Logger()

	tokenClaims, ok := ctx.Value(constants.TokenClaimsKey).(*_jwt.TokenClaims)
	if !ok {
		outerErr := _errors.TokenClaimsDoesNotSet
		lgr.Error().Msg(outerErr.ErrorMessage)
		return nil, outerErr
	}

	getCollections, err := s.vocabularyApiClient.GetCollections(ctx, &VocabularyApiProto.GetCollectionsRequest{
		UserId: tokenClaims.UserId,
	})
	if err != nil {
		outerErr := _errors.FailedToGetCollections
		lgr.Error().Err(err).Msg(outerErr.ErrorMessage)
		return nil, outerErr
	}

	collections := make([]*GatewayApiProto.Collection, len(getCollections.Collections))
	for i, collection := range getCollections.Collections {
		collections[i] = &GatewayApiProto.Collection{
			CollectionId: collection.CollectionId,
			UserId:       collection.UserId,
			LanguageId:   collection.LanguageId,
			Name:         collection.Name,
			Description:  collection.Description,
			IsPinned:     collection.IsPinned,
			IsPublic:     collection.IsPublic,
			TotalTerms:   collection.TotalTerms,
			LearnedTerms: collection.LearnedTerms,
			OpenedAt:     collection.OpenedAt,
			CreatedAt:    collection.CreatedAt,
			UpdatedAt:    collection.UpdatedAt,
		}
	}

	resp := &GatewayApiProto.GetCollectionsResponse{
		Collections: collections,
	}

	lgr.Debug().Msg("executed")
	return resp, nil
}
