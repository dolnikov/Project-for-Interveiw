package service

import (
	"context"
	"gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/services/gateway_api/mappers"

	"gitlab.com/wordbyword.io/microservices/pkg/constants"
	"gitlab.com/wordbyword.io/microservices/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ActionApiProto "gitlab.com/wbwapis/go-genproto/wbw/action/action_api/v1"
	AuthApiProto "gitlab.com/wbwapis/go-genproto/wbw/auth/auth_api/v1"
	GatewayApiProto "gitlab.com/wbwapis/go-genproto/wbw/gateway/gateway_api/v1"
	NotificationApiProto "gitlab.com/wbwapis/go-genproto/wbw/notification/notification_api/v1"
	UserApiProto "gitlab.com/wbwapis/go-genproto/wbw/user/user_api/v1"
	_errors "gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/errors"
)

func (s *GatewayApiService) SignUp(ctx context.Context, request *GatewayApiProto.SignUpRequest) (*GatewayApiProto.SignUpResponse, error) {
	requestId := utils.AnyToString(ctx.Value(constants.RequestIdKey))
	clientIP := utils.AnyToString(ctx.Value(constants.ClientIPKey))
	device := utils.AnyToString(ctx.Value(constants.DeviceKey))
	lgr := s.lgr.With().
		Str(constants.RequestIdKey, requestId).
		Str("api", "SignUp").
		Interface("request", request).
		Logger()

	settings := UserApiProto.Settings{}
	if request.Settings != nil {
		settings.SpeakerGender = mappers.GatewayGenderTypeToUserApiGenderType[request.Settings.SpeakerGender]
		settings.InterfaceLanguageId = request.Settings.InterfaceLanguageId
	}

	createUserResp, err := s.userApiClient.CreateUser(ctx, &UserApiProto.CreateUserRequest{
		Password: request.Password,
		Email:    request.Email,
		Username: request.Username,
		Settings: &settings,
	})
	if err != nil {
		if errGrpc, ok := status.FromError(err); ok {
			if errGrpc.Code() == codes.AlreadyExists {
				outerErr := _errors.AlreadyExistError(errGrpc.Message())
				lgr.Error().Err(err).Msg(outerErr.ErrorMessage)
				return nil, outerErr
			}
		}

		outerErr := _errors.FailedToCreateUser
		lgr.Error().Err(err).Msg(outerErr.ErrorMessage)
		return nil, outerErr
	}

	generateTokensResp, err := s.authApiClient.GenerateTokens(ctx, &AuthApiProto.GenerateTokensRequest{
		UserId: createUserResp.User.UserId,
		Ip:     clientIP,
		Device: device,
	})
	if err != nil {
		outerErr := _errors.FailedToGenerateTokens
		lgr.Error().Err(err).Msg(outerErr.ErrorMessage)
		return nil, outerErr
	}

	createActionRes, err := s.actionApiClient.CreateAction(ctx, &ActionApiProto.CreateActionRequest{
		Action: ActionApiProto.ActionType_ACTION_TYPE_EMAIL_CONFIRMATION,
		Params: &ActionApiProto.CreateActionRequest_EmailConfirmationParams{
			EmailConfirmationParams: &ActionApiProto.EmailConfirmationParams{
				UserId: createUserResp.User.UserId,
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
		Type:  NotificationApiProto.EmailType_EMAIL_TYPE_EMAIL_CONFIRMATION,
		Params: &NotificationApiProto.SendEmailRequest_EmailConfirmationParams{
			EmailConfirmationParams: &NotificationApiProto.EmailConfirmationParams{
				ActionUuid: createActionRes.ActionUuid,
			},
		},
	})
	if err != nil {
		outerErr := _errors.FailedToSendConfirmationEmail
		lgr.Error().Err(err).Msg(outerErr.ErrorMessage)
		return nil, outerErr
	}

	resp := &GatewayApiProto.SignUpResponse{
		AccessToken:  generateTokensResp.AccessToken,
		RefreshToken: generateTokensResp.RefreshToken,
		User: &GatewayApiProto.User{
			UserId:          createUserResp.User.UserId,
			Email:           createUserResp.User.Email,
			Username:        createUserResp.User.Username,
			EmailVerifiedAt: createUserResp.User.EmailVerifiedAt,
			CreatedAt:       createUserResp.User.CreatedAt,
			Settings: &GatewayApiProto.Settings{
				SpeakerGender:       mappers.UserApiGenderTypeToGatewayGenderType[createUserResp.User.Settings.SpeakerGender],
				InterfaceLanguageId: createUserResp.User.Settings.InterfaceLanguageId,
			},
		},
	}

	lgr.Debug().Interface("response", resp).Msg("executed")
	return resp, nil
}
