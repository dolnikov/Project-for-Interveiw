package user_api

import (
	"context"
	UserApiProto "gitlab.com/wbwapis/go-genproto/wbw/user/user_api/v1"
	"gitlab.com/wordbyword.io/microservices/pkg/constants"
	"gitlab.com/wordbyword.io/microservices/pkg/utils"
)

func (c *UserApiClient) CreateUser(ctx context.Context, request *UserApiProto.CreateUserRequest,
) (response *UserApiProto.CreateUserResponse, err error) {
	requestId := utils.AnyToString(ctx.Value(constants.RequestIdKey))
	lgr := c.lgr.With().
		Str("api", "CreateUser").
		Str(constants.RequestIdKey, requestId).
		Interface("request", request).Logger()

	conn, err := c.pool.CreateConn(ctx)
	if err != nil {
		lgr.Error().Err(err).Msg("failed to connect to grpc connection worker pool")
		return nil, err
	}
	defer conn.Close()

	response, err = c.createUserServiceClient(conn).
		CreateUser(ctx, request)
	if err != nil {
		lgr.Error().Err(err).Msg("response error")
		return nil, err
	}

	lgr.Debug().Interface("response", response).Msg("executed")

	return response, nil
}
