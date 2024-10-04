package google_auth_api

import (
	"context"
	"fmt"
	"gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/entities"
	"gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/errors"
	"gitlab.com/wordbyword.io/microservices/pkg/constants"
	"gitlab.com/wordbyword.io/microservices/pkg/utils"
)

func (c *GoogleAuthApiClient) GetGoogleUser(ctx context.Context, accessToken string) (response *entities.GoogleAuthUser, err error) {
	requestId := utils.AnyToString(ctx.Value(constants.RequestIdKey))
	lgr := c.lgr.With().
		Str("api", "GetGoogleUser").
		Str(constants.RequestIdKey, requestId).
		Logger()

	response = &entities.GoogleAuthUser{}
	err = c.Call(ctx, fmt.Sprintf("oauth2/v1/userinfo?alt=json&access_token=%s", accessToken), nil, response)
	if err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToGetGoogleUserMsg)
		return nil, err
	}

	lgr.Debug().Interface("response", response).Msg("executed")

	return response, nil
}
