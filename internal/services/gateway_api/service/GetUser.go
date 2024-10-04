package service

import (
	"context"
	UserApiProto "gitlab.com/wbwapis/go-genproto/wbw/user/user_api/v1"
	_errors "gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/errors"
	"gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/services/gateway_api/mappers"
	"gitlab.com/wordbyword.io/microservices/pkg/constants"
	_jwt "gitlab.com/wordbyword.io/microservices/pkg/jwt"
	"gitlab.com/wordbyword.io/microservices/pkg/utils"

	GatewayApiProto "gitlab.com/wbwapis/go-genproto/wbw/gateway/gateway_api/v1"
)

func (s *GatewayApiService) GetUser(ctx context.Context, request *GatewayApiProto.GetUserRequest) (*GatewayApiProto.GetUserResponse, error) {
	requestId := utils.AnyToString(ctx.Value(constants.RequestIdKey))
	lgr := s.lgr.With().
		Str(constants.RequestIdKey, requestId).
		Str("api", "GetUser").
		Interface("request", request).
		Logger()

	tokenClaims, ok := ctx.Value(constants.TokenClaimsKey).(*_jwt.TokenClaims)
	if !ok {
		outerErr := _errors.TokenClaimsDoesNotSet
		lgr.Error().Msg(outerErr.ErrorMessage)
		return nil, outerErr
	}

	getUser, err := s.userApiClient.GetUser(ctx, &UserApiProto.GetUserRequest{
		FindBy: &UserApiProto.GetUserRequest_UserId{
			UserId: tokenClaims.UserId,
		},
	})
	if err != nil {
		outerErr := _errors.FailedToGetUser
		lgr.Error().Err(err).Msg(outerErr.ErrorMessage)
		return nil, outerErr
	}

	resp := &GatewayApiProto.GetUserResponse{
		User: &GatewayApiProto.User{
			UserId:          getUser.User.UserId,
			Email:           getUser.User.Email,
			Username:        getUser.User.Username,
			EmailVerifiedAt: getUser.User.EmailVerifiedAt,
			CreatedAt:       getUser.User.CreatedAt,
			Settings: &GatewayApiProto.Settings{
				SpeakerGender:       mappers.UserApiGenderTypeToGatewayGenderType[getUser.User.Settings.SpeakerGender],
				InterfaceLanguageId: getUser.User.Settings.InterfaceLanguageId,
			},
		},
	}

	lgr.Debug().Interface("response", resp).Msg("executed")
	return resp, nil
}
