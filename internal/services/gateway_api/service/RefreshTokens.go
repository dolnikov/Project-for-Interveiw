package service

import (
	"context"
	"gitlab.com/wordbyword.io/microservices/pkg/constants"
	"gitlab.com/wordbyword.io/microservices/pkg/utils"

	AuthApiProto "gitlab.com/wbwapis/go-genproto/wbw/auth/auth_api/v1"
	GatewayApiProto "gitlab.com/wbwapis/go-genproto/wbw/gateway/gateway_api/v1"
	_errors "gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/errors"
)

func (s *GatewayApiService) RefreshTokens(ctx context.Context, request *GatewayApiProto.RefreshTokensRequest) (*GatewayApiProto.RefreshTokensResponse, error) {
	requestId := utils.AnyToString(ctx.Value(constants.RequestIdKey))
	lgr := s.lgr.With().
		Str(constants.RequestIdKey, requestId).
		Str("api", "RefreshTokens").
		Interface("request", request).
		Logger()

	generateTokensResp, err := s.authApiClient.RefreshTokens(ctx, &AuthApiProto.RefreshTokensRequest{
		RefreshToken: request.RefreshToken,
	})
	if err != nil {
		outerErr := _errors.FailedToRefreshTokens
		lgr.Error().Err(err).Msg(outerErr.ErrorMessage)
		return nil, outerErr
	}

	resp := &GatewayApiProto.RefreshTokensResponse{
		AccessToken:  generateTokensResp.AccessToken,
		RefreshToken: generateTokensResp.RefreshToken,
	}

	lgr.Debug().Interface("response", resp).Msg("executed")
	return resp, nil
}
