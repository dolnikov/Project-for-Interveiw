package service

import (
	"context"
	GatewayApiProto "gitlab.com/wbwapis/go-genproto/wbw/gateway/gateway_api/v1"
	SpeakerApiProto "gitlab.com/wbwapis/go-genproto/wbw/speaker/speaker_api/v1"
	_errors "gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/errors"
	"gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/services/gateway_api/mappers"
	"gitlab.com/wordbyword.io/microservices/pkg/constants"
	"gitlab.com/wordbyword.io/microservices/pkg/utils"
)

func (s *GatewayApiService) GetVoiceover(ctx context.Context, request *GatewayApiProto.GetVoiceoverRequest) (*GatewayApiProto.GetVoiceoverResponse, error) {
	requestId := utils.AnyToString(ctx.Value(constants.RequestIdKey))
	lgr := s.lgr.With().
		Str(constants.RequestIdKey, requestId).
		Str("api", "GetVoiceover").
		Interface("request", request).
		Logger()

	voiceover, err := s.speakerApiClient.GetVoiceover(ctx, &SpeakerApiProto.GetVoiceoverRequest{
		Text:       request.Text,
		LanguageId: request.LanguageId,
		Gender:     mappers.GatewayGenderTypeToSpeakerApiGenderType[request.Gender],
	})
	if err != nil {
		outerErr := _errors.FailedToGetVoiceover
		lgr.Error().Err(err).Msg(outerErr.ErrorMessage)
		return nil, outerErr
	}

	resp := &GatewayApiProto.GetVoiceoverResponse{
		Url: voiceover.Url,
	}
	lgr.Debug().Interface("response", resp).Msg("executed")
	return resp, nil
}
