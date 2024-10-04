package service

import (
	"context"
	ActionApiProto "gitlab.com/wbwapis/go-genproto/wbw/action/action_api/v1"
	GatewayApiProto "gitlab.com/wbwapis/go-genproto/wbw/gateway/gateway_api/v1"
	_errors "gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/errors"
	"gitlab.com/wordbyword.io/microservices/pkg/constants"
	"gitlab.com/wordbyword.io/microservices/pkg/utils"
)

func (s *GatewayApiService) ResetPassword(ctx context.Context, request *GatewayApiProto.ResetPasswordRequest) (*GatewayApiProto.ResetPasswordResponse, error) {
	requestId := utils.AnyToString(ctx.Value(constants.RequestIdKey))
	lgr := s.lgr.With().
		Str(constants.RequestIdKey, requestId).
		Str("api", "ResetPassword").
		Interface("request", request).
		Logger()

	_, err := s.actionApiClient.ExecuteAction(ctx, &ActionApiProto.ExecuteActionRequest{
		ActionUuid: request.ActionUuid,
		ExecuteParams: &ActionApiProto.ExecuteActionRequest_ResetPasswordExecuteParams{
			ResetPasswordExecuteParams: &ActionApiProto.ResetPasswordExecuteParams{
				Password: request.Password,
			},
		},
	})
	if err != nil {
		outerErr := _errors.FailedToResetPassword
		lgr.Error().Err(err).Msg(outerErr.ErrorMessage)
		return nil, outerErr
	}

	resp := &GatewayApiProto.ResetPasswordResponse{}
	lgr.Debug().Interface("response", resp).Msg("executed")
	return resp, nil
}
