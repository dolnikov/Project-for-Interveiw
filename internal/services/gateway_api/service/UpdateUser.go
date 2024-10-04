package service

import (
	"context"
	UserApiProto "gitlab.com/wbwapis/go-genproto/wbw/user/user_api/v1"
	_errors "gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/errors"
	"gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/services/gateway_api/mappers"
	"gitlab.com/wordbyword.io/microservices/pkg/constants"
	_jwt "gitlab.com/wordbyword.io/microservices/pkg/jwt"
	"gitlab.com/wordbyword.io/microservices/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	GatewayApiProto "gitlab.com/wbwapis/go-genproto/wbw/gateway/gateway_api/v1"
)

func (s *GatewayApiService) UpdateUser(ctx context.Context, request *GatewayApiProto.UpdateUserRequest) (*GatewayApiProto.UpdateUserResponse, error) {
	requestId := utils.AnyToString(ctx.Value(constants.RequestIdKey))
	lgr := s.lgr.With().
		Str(constants.RequestIdKey, requestId).
		Str("api", "UpdateUser").
		Interface("request", request).
		Logger()

	tokenClaims, ok := ctx.Value(constants.TokenClaimsKey).(*_jwt.TokenClaims)
	if !ok {
		outerErr := _errors.TokenClaimsDoesNotSet
		lgr.Error().Msg(outerErr.ErrorMessage)
		return nil, outerErr
	}

	newUserData := UserApiProto.User{
		UserId:   tokenClaims.UserId,
		Settings: &UserApiProto.Settings{},
	}

	if request.Username != nil {
		newUserData.Username = *request.Username
	}

	if request.Settings != nil {
		if request.Settings.SpeakerGender != 0 {
			newUserData.Settings.SpeakerGender = mappers.GatewayGenderTypeToUserApiGenderType[request.Settings.SpeakerGender]
		}

		if request.Settings.InterfaceLanguageId != 0 {
			newUserData.Settings.InterfaceLanguageId = request.Settings.InterfaceLanguageId
		}
	}

	_, err := s.userApiClient.UpdateUser(ctx, &UserApiProto.UpdateUserRequest{
		User: &newUserData,
	})
	if err != nil {
		if errGrpc, ok := status.FromError(err); ok {
			if errGrpc.Code() == codes.AlreadyExists {
				outerErr := _errors.AlreadyExistError(errGrpc.Message())
				lgr.Error().Err(err).Msg(outerErr.ErrorMessage)
				return nil, outerErr
			}
		}

		outerErr := _errors.FailedToUpdateUser
		lgr.Error().Err(err).Msg(outerErr.ErrorMessage)
		return nil, outerErr
	}

	resp := &GatewayApiProto.UpdateUserResponse{}
	lgr.Debug().Interface("response", resp).Msg("executed")
	return resp, nil
}
