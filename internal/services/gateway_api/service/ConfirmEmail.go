package service

import (
	"context"
	ActionApiProto "gitlab.com/wbwapis/go-genproto/wbw/action/action_api/v1"
	GatewayApiProto "gitlab.com/wbwapis/go-genproto/wbw/gateway/gateway_api/v1"
	_errors "gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/errors"
	"gitlab.com/wordbyword.io/microservices/pkg/constants"
	"gitlab.com/wordbyword.io/microservices/pkg/utils"
)

func (s *GatewayApiService) ConfirmEmail(ctx context.Context, request *GatewayApiProto.ConfirmEmailRequest) (*GatewayApiProto.ConfirmEmailResponse, error) {
	requestId := utils.AnyToString(ctx.Value(constants.RequestIdKey))
	lgr := s.lgr.With().
		Str(constants.RequestIdKey, requestId).
		Str("api", "ConfirmEmail").
		Interface("request", request).
		Logger()

	_, err := s.actionApiClient.ExecuteAction(ctx, &ActionApiProto.ExecuteActionRequest{
		ActionUuid: request.ActionUuid,
		ExecuteParams: &ActionApiProto.ExecuteActionRequest_EmailConfirmationExecuteParams{
			EmailConfirmationExecuteParams: &ActionApiProto.EmailConfirmationExecuteParams{},
		},
	})
	if err != nil {
		outerErr := _errors.FailedToConfirmEmail
		lgr.Error().Err(err).Msg(outerErr.ErrorMessage)
		return nil, outerErr
	}

	resp := &GatewayApiProto.ConfirmEmailResponse{}
	lgr.Debug().Interface("response", resp).Msg("executed")
	return resp, nil
}
