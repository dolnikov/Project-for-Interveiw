package gateway_api

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/config"
	"gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/prometheus"
	GatewayApiService "gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/services/gateway_api/service"
)

type IGatewayApiHttpEndpoint interface {
	SignUp(c *gin.Context)
	SignIn(c *gin.Context)
	Logout(c *gin.Context)
	RefreshTokens(c *gin.Context)
	GetUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	ConfirmEmail(c *gin.Context)
	CreateCollection(c *gin.Context)
	UpdateCollection(c *gin.Context)
	GetCollections(c *gin.Context)
	GetCollection(c *gin.Context)
	DeleteCollection(c *gin.Context)
	CreateTerms(c *gin.Context)
	UpdateTerm(c *gin.Context)
	GetTerms(c *gin.Context)
	DeleteTerms(c *gin.Context)
	ChangeTermStatus(c *gin.Context)
	GetLanguages(c *gin.Context)
	GetVoiceover(c *gin.Context)
	GetTranslation(c *gin.Context)
}

type GatewayApiHttpEndpoint struct {
	IGatewayApiHttpEndpoint

	cfg  *config.Config
	lgr  zerolog.Logger
	prom *prometheus.Exporter

	//list of services:
	gatewayApiService GatewayApiService.IGatewayApiService
}

var _ IGatewayApiHttpEndpoint = (*GatewayApiHttpEndpoint)(nil)

func NewGatewayApiHttpEndpoint(cfg *config.Config, lgr zerolog.Logger, prom *prometheus.Exporter, gatewayApiService GatewayApiService.IGatewayApiService) *GatewayApiHttpEndpoint {
	return &GatewayApiHttpEndpoint{
		cfg:  cfg,
		lgr:  lgr,
		prom: prom,

		//list of services:
		gatewayApiService: gatewayApiService,
	}
}
