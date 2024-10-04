package service

import (
	"github.com/rs/zerolog"
	"gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/config"
	"gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/prometheus"
	"gitlab.com/wordbyword.io/microservices/pkg/sequence_gen"
)

const serviceName = "gateway-api"

type GatewayApiService struct {
	cfg       *config.Config
	lgr       zerolog.Logger
	snowflake *sequence_gen.Snowflake
	prom      *prometheus.Exporter

	//list of clients:
	authApiClient         IAuthApiClient
	userApiClient         IUserApiClient
	actionApiClient       IActionApiClient
	notificationApiClient INotificationApiClient
	vocabularyApiClient   IVocabularyApiClient
	speakerApiClient      ISpeakerApiClient
	languageApiClient     ILanguageApiClient
	translationApiClient  ITranslationApiClient
	googleAuthApiClient   IGoogleAuthClient
}

var _ IGatewayApiService = (*GatewayApiService)(nil)

// NewGatewayApiService initializes a new Service struct.
func NewGatewayApiService(
	cfg *config.Config,
	lgr zerolog.Logger,
	prom *prometheus.Exporter,
	notificationApiClient INotificationApiClient,
	authApiClient IAuthApiClient,
	userApiClient IUserApiClient,
	actionApiClient IActionApiClient,
	vocabularyApiClient IVocabularyApiClient,
	speakerApiClient ISpeakerApiClient,
	languageApiClient ILanguageApiClient,
	translationApiClient ITranslationApiClient,
	googleAuthApiClient IGoogleAuthClient,
) *GatewayApiService {
	lgr = lgr.With().Str("service", serviceName).Logger()

	return &GatewayApiService{
		cfg:       cfg,
		lgr:       lgr,
		snowflake: sequence_gen.NewSnowflake(nil),
		prom:      prom,

		//list of clients:
		notificationApiClient: notificationApiClient,
		authApiClient:         authApiClient,
		userApiClient:         userApiClient,
		actionApiClient:       actionApiClient,
		vocabularyApiClient:   vocabularyApiClient,
		speakerApiClient:      speakerApiClient,
		languageApiClient:     languageApiClient,
		translationApiClient:  translationApiClient,
		googleAuthApiClient:   googleAuthApiClient,
	}
}

func (s *GatewayApiService) Shutdown() {

}
