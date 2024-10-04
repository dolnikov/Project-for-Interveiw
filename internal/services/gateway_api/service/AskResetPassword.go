package service

import (
	"context"
	ActionApiProto "gitlab.com/wbwapis/go-genproto/wbw/action/action_api/v1"
	GatewayApiProto "gitlab.com/wbwapis/go-genproto/wbw/gateway/gateway_api/v1"
	NotificationApiProto "gitlab.com/wbwapis/go-genproto/wbw/notification/notification_api/v1"
	UserApiProto "gitlab.com/wbwapis/go-genproto/wbw/user/user_api/v1"
	_errors "gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/errors"
	"gitlab.com/wordbyword.io/microservices/pkg/constants"
	"gitlab.com/wordbyword.io/microservices/pkg/utils"
)

func (s *GatewayApiService) AskResetPassword(ctx context.Context, request *GatewayApiProto.AskResetPasswordRequest) (*GatewayApiProto.AskResetPasswordResponse, error) {
	requestId := utils.AnyToString(ctx.Value(constants.RequestIdKey))
	lgr := s.lgr.With().
		Str(constants.RequestIdKey, requestId).
		Str("api", "AskResetPassword").
		Interface("request", request).
		Logger()

	user, err := s.userApiClient.GetUser(ctx, &UserApiProto.GetUserRequest{
		FindBy: &UserApiProto.GetUserRequest_Email{
			Email: request.Email,
		},
	})
	if err != nil {
		outerErr := _errors.FailedToGetUser
		lgr.Error().Err(err).Msg(outerErr.ErrorMessage)
		return nil, outerErr
	}

	createActionResp, err := s.actionApiClient.CreateAction(ctx, &ActionApiProto.CreateActionRequest{
		Action: ActionApiProto.ActionType_ACTION_TYPE_RESET_PASSWORD,
		Params: &ActionApiProto.CreateActionRequest_ResetPasswordParams{
			ResetPasswordParams: &ActionApiProto.ResetPasswordParams{
				Email: user.User.Email,
			},
		},
	})
	if err != nil {
		outerErr := _errors.FailedToCreateAction
		lgr.Error().Err(err).Msg(outerErr.ErrorMessage)
		return nil, outerErr
	}

	err = s.notificationApiClient.SendEmail(ctx, lgr, &NotificationApiProto.SendEmailRequest{
		Email: request.Email,
		Type:  NotificationApiProto.EmailType_EMAIL_TYPE_RESET_PASSWORD,
		Params: &NotificationApiProto.SendEmailRequest_ResetPasswordParams{
			ResetPasswordParams: &NotificationApiProto.ResetPasswordParams{
				ActionUuid: createActionResp.ActionUuid,
			},
		},
	})
	if err != nil {
		outerErr := _errors.FailedToSendResetPasswordEmail
		lgr.Error().Err(err).Msg(outerErr.ErrorMessage)
		return nil, outerErr
	}

	resp := &GatewayApiProto.AskResetPasswordResponse{}
	lgr.Debug().Interface("response", resp).Msg("executed")
	return resp, nil
}
