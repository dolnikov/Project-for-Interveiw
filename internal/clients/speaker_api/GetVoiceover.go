package speaker_api

import (
	"context"
	SpeakerApiProto "gitlab.com/wbwapis/go-genproto/wbw/speaker/speaker_api/v1"
	"gitlab.com/wordbyword.io/microservices/pkg/constants"
	"gitlab.com/wordbyword.io/microservices/pkg/utils"
)

func (c *SpeakerApiClient) GetVoiceover(ctx context.Context, request *SpeakerApiProto.GetVoiceoverRequest,
) (response *SpeakerApiProto.GetVoiceoverResponse, err error) {
	requestId := utils.AnyToString(ctx.Value(constants.RequestIdKey))
	lgr := c.lgr.With().
		Str("api", "GetVoiceover").
		Str(constants.RequestIdKey, requestId).
		Interface("request", request).Logger()

	conn, err := c.pool.CreateConn(ctx)
	if err != nil {
		lgr.Error().Err(err).Msg("failed to connect to grpc connection worker pool")
		return nil, err
	}
	defer conn.Close()

	response, err = c.createSpeakerServiceClient(conn).
		GetVoiceover(ctx, request)
	if err != nil {
		lgr.Error().Err(err).Msg("response error")
		return nil, err
	}

	lgr.Debug().Interface("response", response).Msg("executed")

	return response, nil
}
