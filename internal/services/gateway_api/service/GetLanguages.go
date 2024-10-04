package service

import (
	"context"
	LanguageApiProto "gitlab.com/wbwapis/go-genproto/wbw/language/language_api/v1"
	_errors "gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/errors"
	"gitlab.com/wordbyword.io/microservices/pkg/constants"
	"gitlab.com/wordbyword.io/microservices/pkg/utils"

	GatewayApiProto "gitlab.com/wbwapis/go-genproto/wbw/gateway/gateway_api/v1"
)

func (s *GatewayApiService) GetLanguages(ctx context.Context, request *GatewayApiProto.GetLanguagesRequest) (*GatewayApiProto.GetLanguagesResponse, error) {
	requestId := utils.AnyToString(ctx.Value(constants.RequestIdKey))
	lgr := s.lgr.With().
		Str(constants.RequestIdKey, requestId).
		Str("api", "GetLanguages").
		Interface("request", request).
		Logger()

	getLanguages, err := s.languageApiClient.GetLanguages(ctx, &LanguageApiProto.GetLanguagesRequest{})
	if err != nil {
		outerErr := _errors.FailedToGetLanguages
		lgr.Error().Err(err).Msg(outerErr.ErrorMessage)
		return nil, outerErr
	}

	languages := make([]*GatewayApiProto.Language, len(getLanguages.Languages))
	for i, language := range getLanguages.Languages {
		languages[i] = &GatewayApiProto.Language{
			LanguageId:   language.LanguageId,
			Code:         language.Code,
			ShortCode:    language.ShortCode,
			Name:         language.Name,
			I18NSlug:     language.I18NSlug,
			SiteLanguage: language.SiteLanguage,
			Order:        language.Order,
		}
	}

	resp := &GatewayApiProto.GetLanguagesResponse{
		Languages: languages,
	}

	lgr.Debug().Msg("executed")
	return resp, nil
}
