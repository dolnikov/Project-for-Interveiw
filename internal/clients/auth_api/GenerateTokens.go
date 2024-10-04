package auth_api

import (
	"context"
	AuthApiProto "gitlab.com/wbwapis/go-genproto/wbw/auth/auth_api/v1"
	"gitlab.com/wordbyword.io/microservices/pkg/constants"
	"gitlab.com/wordbyword.io/microservices/pkg/utils"
)

// GenerateTokens generate new pair of tokens
func (c *AuthApiClient) GenerateTokens(ctx context.Context, request *AuthApiProto.GenerateTokensRequest,
) (response *AuthApiProto.GenerateTokensResponse, err error) {
	requestId := utils.AnyToString(ctx.Value(constants.RequestIdKey))
	lgr := c.lgr.With().
		Str("api", "GenerateTokens").
		Str(constants.RequestIdKey, requestId).
		Interface("request", request).Logger()

	conn, err := c.pool.CreateConn(ctx)
	if err != nil {
		lgr.Error().Err(err).Msg("failed to connect to grpc connection worker pool")
		return nil, err
	}
	defer conn.Close()

	response, err = c.createAuthServiceClient(conn).
		GenerateTokens(ctx, request)
	if err != nil {
		lgr.Warn().Err(err).Msg("response error")
		return nil, err
	}

	lgr.Debug().Interface("response", response).Msg("executed")

	return response, nil
}
