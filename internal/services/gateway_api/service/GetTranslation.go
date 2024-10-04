package service

import (
	"context"
	TranslationApiProto "gitlab.com/wbwapis/go-genproto/wbw/translation/translation_api/v1"
	_errors "gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/errors"
	"gitlab.com/wordbyword.io/microservices/pkg/constants"
	"gitlab.com/wordbyword.io/microservices/pkg/utils"

	GatewayApiProto "gitlab.com/wbwapis/go-genproto/wbw/gateway/gateway_api/v1"
)

func (s *GatewayApiService) GetTranslation(ctx context.Context, request *GatewayApiProto.GetTranslationRequest) (*GatewayApiProto.GetTranslationResponse, error) {
	requestId := utils.AnyToString(ctx.Value(constants.RequestIdKey))
	lgr := s.lgr.With().
		Str(constants.RequestIdKey, requestId).
		Str("api", "GetTranslation").
		Interface("request", request).
		Logger()

	getTranslation, err := s.translationApiClient.GetTranslation(ctx, &TranslationApiProto.GetTranslationRequest{
		Text:           request.Text,
		SourceLanguage: request.SourceLanguage,
		TargetLanguage: request.TargetLanguage,
	})
	if err != nil {
		outerErr := _errors.FailedToGetTranslation
		lgr.Error().Err(err).Msg(outerErr.ErrorMessage)
		return nil, outerErr
	}

	translations := make([]*GatewayApiProto.Translations, len(getTranslation.Translations))
	for i, tr := range getTranslation.Translations {
		translations[i] = &GatewayApiProto.Translations{
			Text: tr.Text,
		}
	}

	resp := &GatewayApiProto.GetTranslationResponse{
		Translations: translations,
	}

	lgr.Debug().Interface("response", resp).Msg("executed")
	return resp, nil
}
