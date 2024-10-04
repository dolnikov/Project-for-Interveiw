package google_auth_api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
	"gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/config"
	"gitlab.com/wordbyword.io/microservices/pkg/constants"
	"net/http"
	"time"
)

const clientName = "google_auth_api"

type GoogleAuthApiClient struct {
	cfg    *config.Config
	lgr    zerolog.Logger
	Client http.Client
}

func NewGoogleAuthApiClient(cfg *config.Config, lgr zerolog.Logger) *GoogleAuthApiClient {
	lgr = lgr.With().Str("client", clientName).Logger()
	client := http.Client{
		Timeout: time.Duration(cfg.GoogleApi.Timeout) * time.Second,
	}

	return &GoogleAuthApiClient{
		cfg:    cfg,
		lgr:    lgr,
		Client: client,
	}
}

func (c *GoogleAuthApiClient) Shutdown() {

}

func (c *GoogleAuthApiClient) Call(ctx context.Context, method string, body interface{}, response interface{}) (err error) {
	var message []byte
	if body != nil {
		message, err = json.Marshal(body)
		if err != nil {
			return err
		}
	}

	host := c.cfg.GoogleApi.URI
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/%s", host, method),
		bytes.NewBuffer(message),
	)
	if err != nil {
		return err
	}

	requestID := ctx.Value(constants.RequestIdKey).(string)
	req.Header.Add(constants.RequestIdKey, requestID)
	req.Header.Add("Accept", `application/json`)

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return err
	}

	return nil
}
