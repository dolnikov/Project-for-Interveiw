package endpoints

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/config"
	"gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/constants"
	"gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/pkg/rate_limiter"
	"gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/prometheus"
	"gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/transports/http/middleware/auth"
	"gitlab.com/wordbyword.io/microservices/pkg/middleware/extractor"

	AuthApiClient "gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/clients/auth_api"
	GatewayApiService "gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/services/gateway_api/service"
	GatewayApiHttpEndpoint "gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/transports/http/endpoints/gateway_api"
)

const (
	// DefaultPrefix url prefix of HTTP API
	DefaultPrefix = "/"
)

func getPrefix(prefixOptions ...string) string {
	prefix := DefaultPrefix
	if len(prefixOptions) > 0 {
		if prefixOptions[0] != "" {
			prefix = prefixOptions[0]
		}
	}
	return prefix
}

type HttpEndpoints struct {
	cfg         *config.Config
	lgr         zerolog.Logger
	prom        *prometheus.Exporter
	rateLimiter *rate_limiter.RateLimiter

	//list of endpoints:
	gatewayApiHttpEndpoint *GatewayApiHttpEndpoint.GatewayApiHttpEndpoint

	//list of clients:
	authApiClient *AuthApiClient.AuthApiClient
}

func NewHttpEndpoints(
	cfg *config.Config,
	lgr zerolog.Logger,
	prom *prometheus.Exporter,
	gatewayApiService GatewayApiService.IGatewayApiService,
	authApiClient *AuthApiClient.AuthApiClient,
) (*HttpEndpoints, error) {

	return &HttpEndpoints{
		cfg:  cfg,
		lgr:  lgr,
		prom: prom,
		rateLimiter: rate_limiter.NewRateLimiter(map[string]int{
			constants.SignUp:           cfg.RateLimiter.SignUp,
			constants.SignIn:           cfg.RateLimiter.SignIn,
			constants.RefreshTokens:    cfg.RateLimiter.RefreshTokens,
			constants.ConfirmEmail:     cfg.RateLimiter.ConfirmEmail,
			constants.AskResetPassword: cfg.RateLimiter.AskResetPassword,
			constants.ResetPassword:    cfg.RateLimiter.ResetPassword,
			constants.GetLanguages:     cfg.RateLimiter.GetLanguages,
			constants.Logout:           cfg.RateLimiter.Logout,
			constants.GetUser:          cfg.RateLimiter.GetUser,
			constants.UpdateUser:       cfg.RateLimiter.UpdateUser,
			constants.CreateCollection: cfg.RateLimiter.CreateCollection,
			constants.UpdateCollection: cfg.RateLimiter.UpdateCollection,
			constants.GetCollections:   cfg.RateLimiter.GetCollections,
			constants.GetCollection:    cfg.RateLimiter.GetCollection,
			constants.DeleteCollection: cfg.RateLimiter.DeleteCollection,
			constants.CreateTerms:      cfg.RateLimiter.CreateTerms,
			constants.UpdateTerm:       cfg.RateLimiter.UpdateTerm,
			constants.GetTerms:         cfg.RateLimiter.GetTerms,
			constants.ChangeTermStatus: cfg.RateLimiter.ChangeTermStatus,
			constants.DeleteTerms:      cfg.RateLimiter.DeleteTerms,
			constants.GetVoiceover:     cfg.RateLimiter.GetVoiceover,
			constants.GetTranslation:   cfg.RateLimiter.GetTranslation,
		}),

		//list of endpoints:
		gatewayApiHttpEndpoint: GatewayApiHttpEndpoint.NewGatewayApiHttpEndpoint(cfg, lgr, prom, gatewayApiService),

		//list of clients:
		authApiClient: authApiClient,
	}, nil
}

// RegisterServer the standard HandlerFuncs
func (ep *HttpEndpoints) RegisterServer(r *gin.Engine, prefixOptions ...string) {
	ep.routeRegister(&(r.RouterGroup), prefixOptions...)
}

// routeRegister the standard HandlerFuncs from the net/http/pprof package with
// the provided gin.GrouterGroup. prefixOptions is a optional. If not prefixOptions,
// the default path prefix is used, otherwise first prefixOptions will be path prefix.
func (ep *HttpEndpoints) routeRegister(rg *gin.RouterGroup, prefixOptions ...string) {
	prefix := getPrefix(prefixOptions...)
	r := rg.Group(prefix)

	r.Use(extractor.ExtractRequestId())
	r.Use(extractor.ExtractClientIP())
	r.Use(extractor.ExtractAcceptLanguage())
	r.Use(rate_limiter.HttpMiddleware(ep.cfg, ep.lgr, ep.prom, ep.rateLimiter))

	r.POST(constants.SignUp, ep.gatewayApiHttpEndpoint.SignUp)
	r.POST(constants.SignIn, ep.gatewayApiHttpEndpoint.SignIn)
	r.POST(constants.RefreshTokens, ep.gatewayApiHttpEndpoint.RefreshTokens)
	r.POST(constants.ConfirmEmail, ep.gatewayApiHttpEndpoint.ConfirmEmail)
	r.POST(constants.AskResetPassword, ep.gatewayApiHttpEndpoint.AskResetPassword)
	r.POST(constants.ResetPassword, ep.gatewayApiHttpEndpoint.ResetPassword)
	r.POST(constants.GetLanguages, ep.gatewayApiHttpEndpoint.GetLanguages)
	r.POST(constants.GetVoiceover, ep.gatewayApiHttpEndpoint.GetVoiceover)

	// With auth:
	authorized := r.Group("", auth.Auth(ep.cfg, ep.lgr, ep.prom))
	authorized.POST(constants.Logout, ep.gatewayApiHttpEndpoint.Logout)
	authorized.POST(constants.GetUser, ep.gatewayApiHttpEndpoint.GetUser)
	authorized.POST(constants.UpdateUser, ep.gatewayApiHttpEndpoint.UpdateUser)
	authorized.POST(constants.CreateCollection, ep.gatewayApiHttpEndpoint.CreateCollection)
	authorized.POST(constants.GetCollections, ep.gatewayApiHttpEndpoint.GetCollections)
	authorized.POST(constants.GetCollection, ep.gatewayApiHttpEndpoint.GetCollection)
	authorized.POST(constants.GetTerms, ep.gatewayApiHttpEndpoint.GetTerms)
	authorized.POST(constants.GetTranslation, ep.gatewayApiHttpEndpoint.GetTranslation)

	// Only owner:
	authorized.POST(constants.UpdateCollection, ep.gatewayApiHttpEndpoint.UpdateCollection)
	authorized.POST(constants.DeleteCollection, ep.gatewayApiHttpEndpoint.DeleteCollection)
	authorized.POST(constants.CreateTerms, ep.gatewayApiHttpEndpoint.CreateTerms)
	authorized.POST(constants.DeleteTerms, ep.gatewayApiHttpEndpoint.DeleteTerms)
	authorized.POST(constants.UpdateTerm, ep.gatewayApiHttpEndpoint.UpdateTerm)
	authorized.POST(constants.ChangeTermStatus, ep.gatewayApiHttpEndpoint.ChangeTermStatus)
}
