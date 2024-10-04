package service

import (
	"context"
	"gitlab.com/wordbyword.io/microservices/pkg/constants"
	"gitlab.com/wordbyword.io/microservices/pkg/utils"

	AuthApiProto "gitlab.com/wbwapis/go-genproto/wbw/auth/auth_api/v1"
	GatewayApiProto "gitlab.com/wbwapis/go-genproto/wbw/gateway/gateway_api/v1"
	_errors "gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/errors"
	_jwt "gitlab.com/wordbyword.io/microservices/pkg/jwt"
)

func (s *GatewayApiService) Logout(ctx context.Context, request *GatewayApiProto.LogoutRequest) (*GatewayApiProto.LogoutResponse, error) {
	requestId := utils.AnyToString(ctx.Value(constants.RequestIdKey))
	lgr := s.lgr.With().
		Str(constants.RequestIdKey, requestId).
		Str("api", "Logout").
		Interface("request", request).
		Logger()

	tokenClaims, ok := ctx.Value(constants.TokenClaimsKey).(*_jwt.TokenClaims)
	if !ok {
		outerErr := _errors.TokenClaimsDoesNotSet
		lgr.Error().Msg(outerErr.ErrorMessage)
		return nil, outerErr
	}

	_, err := s.authApiClient.DeleteTokens(ctx, &AuthApiProto.DeleteTokensRequest{
		UserId:  tokenClaims.UserId,
		TokenId: tokenClaims.TokenId,
	})
	if err != nil {
		outerErr := _errors.FailedToDeleteTokens
		lgr.Error().Err(err).Msg(outerErr.ErrorMessage)
		return nil, outerErr
	}

	resp := &GatewayApiProto.LogoutResponse{}
	lgr.Debug().Interface("response", resp).Msg("executed")
	return resp, nil
}
