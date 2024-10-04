package service

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	AuthApiProto "gitlab.com/wbwapis/go-genproto/wbw/auth/auth_api/v1"
	GatewayApiProto "gitlab.com/wbwapis/go-genproto/wbw/gateway/gateway_api/v1"
	UserApiProto "gitlab.com/wbwapis/go-genproto/wbw/user/user_api/v1"
	_errors "gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/errors"
	"gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/services/gateway_api/mappers"
	"gitlab.com/wordbyword.io/microservices/pkg/constants"
	"gitlab.com/wordbyword.io/microservices/pkg/utils"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strconv"
)

func (s *GatewayApiService) SignIn(ctx context.Context, request *GatewayApiProto.SignInRequest) (resp *GatewayApiProto.SignInResponse, err error) {
	requestId := utils.AnyToString(ctx.Value(constants.RequestIdKey))
	clientIP := utils.AnyToString(ctx.Value(constants.ClientIPKey))
	device := utils.AnyToString(ctx.Value(constants.DeviceKey))
	lgr := s.lgr.With().
		Str(constants.RequestIdKey, requestId).
		Str("api", "SignIn").
		Interface("request", request).
		Logger()

	var user *GatewayApiProto.User
	if request.GetGoogleToken() != "" {
		user, err = s.getOrCreateUserByGoogleToken(ctx, lgr, request)
	} else {
		user, err = s.getUserByEmailOrUsername(ctx, lgr, request)
	}

	if err != nil {
		outerErr := _errors.FailedToSignIn
		lgr.Error().Err(err).Msg(outerErr.ErrorMessage)
		return nil, outerErr
	}

	generateTokensResp, err := s.authApiClient.GenerateTokens(ctx, &AuthApiProto.GenerateTokensRequest{
		UserId: user.UserId,
		Ip:     clientIP,
		Device: device,
	})
	if err != nil {
		outerErr := _errors.FailedToGenerateTokens
		lgr.Error().Err(err).Msg(outerErr.ErrorMessage)
		return nil, outerErr
	}

	resp = &GatewayApiProto.SignInResponse{
		AccessToken:  generateTokensResp.AccessToken,
		RefreshToken: generateTokensResp.RefreshToken,
		User:         user,
	}

	lgr.Debug().Interface("response", resp).Msg("executed")
	return resp, nil
}

func (s *GatewayApiService) getOrCreateUserByGoogleToken(ctx context.Context, lgr zerolog.Logger, request *GatewayApiProto.SignInRequest) (*GatewayApiProto.User, error) {
	// 1) Инвалидираем токен
	googleUser, err := s.googleAuthApiClient.GetGoogleUser(ctx, request.GetGoogleToken())
	if err != nil {
		outerErr := _errors.FailedToGetGoogleUser
		lgr.Error().Err(err).Msg(outerErr.ErrorMessage)
		return nil, outerErr
	}

	if !googleUser.VerifiedEmail {
		outerErr := _errors.FailedToSignInEmailNotConfirmed
		lgr.Error().Err(outerErr).Msg(outerErr.ErrorMessage)
		return nil, outerErr
	}

	// 2) Проверем есть ли такой пользователь
	getUser, err := s.userApiClient.GetUser(ctx, &UserApiProto.GetUserRequest{
		FindBy: &UserApiProto.GetUserRequest_Email{
			Email: googleUser.Email,
		},
	})
	if err == nil {
		return &GatewayApiProto.User{
			UserId:          getUser.User.UserId,
			Email:           getUser.User.Email,
			Username:        getUser.User.Username,
			EmailVerifiedAt: getUser.User.EmailVerifiedAt,
			CreatedAt:       getUser.User.CreatedAt,
			Settings: &GatewayApiProto.Settings{
				SpeakerGender:       mappers.UserApiGenderTypeToGatewayGenderType[getUser.User.Settings.SpeakerGender],
				InterfaceLanguageId: getUser.User.Settings.InterfaceLanguageId,
			},
		}, nil
	} else {
		if err.Error() != _errors.FailedToGetUserFormDB {
			outerErr := _errors.FailedToGetUser
			lgr.Error().Err(err).Msg(outerErr.ErrorMessage)
			return nil, outerErr
		}
	}

	// 3) Если нет то регистрируем методом SignUp
	username, _ := s.snowflake.GenerateId()
	password, _ := s.snowflake.GenerateId()
	numberString := strconv.FormatUint(username, 10)

	createUserResp, err := s.userApiClient.CreateUser(ctx, &UserApiProto.CreateUserRequest{
		Email:           googleUser.Email,
		Username:        fmt.Sprintf("user_%s", numberString[len(numberString)-10:]),
		Password:        fmt.Sprintf("pass_%d", password),
		EmailVerifiedAt: timestamppb.Now(),
	})
	if err != nil {
		outerErr := _errors.FailedToCreateUser
		lgr.Error().Err(err).Msg(outerErr.ErrorMessage)
		return nil, outerErr
	}

	return &GatewayApiProto.User{
		UserId:          createUserResp.User.UserId,
		Email:           createUserResp.User.Email,
		Username:        createUserResp.User.Username,
		EmailVerifiedAt: createUserResp.User.EmailVerifiedAt,
		CreatedAt:       createUserResp.User.CreatedAt,
		Settings: &GatewayApiProto.Settings{
			SpeakerGender:       mappers.UserApiGenderTypeToGatewayGenderType[createUserResp.User.Settings.SpeakerGender],
			InterfaceLanguageId: createUserResp.User.Settings.InterfaceLanguageId,
		},
	}, nil
}

func (s *GatewayApiService) getUserByEmailOrUsername(ctx context.Context, lgr zerolog.Logger, request *GatewayApiProto.SignInRequest) (*GatewayApiProto.User, error) {
	if request.Password == nil {
		outerErr := _errors.FailedToSignInWrongPassword
		lgr.Error().Err(outerErr).Msg(outerErr.ErrorMessage)
		return nil, outerErr
	}

	req := &UserApiProto.GetUserByCredentialsRequest{
		Password: *request.Password,
	}

	if request.GetEmail() != "" {
		req.LoginBy = &UserApiProto.GetUserByCredentialsRequest_Email{Email: request.GetEmail()}
	} else {
		req.LoginBy = &UserApiProto.GetUserByCredentialsRequest_Username{Username: request.GetUsername()}
	}

	getUserByCredentialsResp, err := s.userApiClient.GetUserByCredentials(ctx, req)
	if err != nil {
		outerErr := _errors.FailedToSignIn
		lgr.Error().Err(err).Msg(outerErr.ErrorMessage)
		return nil, outerErr
	}

	return &GatewayApiProto.User{
		UserId:          getUserByCredentialsResp.User.UserId,
		Email:           getUserByCredentialsResp.User.Email,
		Username:        getUserByCredentialsResp.User.Username,
		EmailVerifiedAt: getUserByCredentialsResp.User.EmailVerifiedAt,
		CreatedAt:       getUserByCredentialsResp.User.CreatedAt,
		Settings: &GatewayApiProto.Settings{
			SpeakerGender:       mappers.UserApiGenderTypeToGatewayGenderType[getUserByCredentialsResp.User.Settings.SpeakerGender],
			InterfaceLanguageId: getUserByCredentialsResp.User.Settings.InterfaceLanguageId,
		},
	}, nil
}
