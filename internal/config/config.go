package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/sethvargo/go-envconfig"
)

func NewConfig() (*Config, error) {
	var envFiles []string
	if _, err := os.Stat(".env"); err == nil {
		log.Println("found .env file, adding it to env config files list")
		envFiles = append(envFiles, ".env")
	}
	if os.Getenv("APP_ENV") != "" {
		appEnvName := fmt.Sprintf(".env.%s", os.Getenv("APP_ENV"))
		if _, err := os.Stat(appEnvName); err == nil {
			log.Println("found", appEnvName, "file, adding it to env config files list")
			envFiles = append(envFiles, appEnvName)
		}
	}
	if len(envFiles) > 0 {
		err := godotenv.Overload(envFiles...)
		if err != nil {
			return nil, errors.Wrapf(err, "error while opening env config: %s", err)
		}
	}
	cfg := &Config{}
	ctx := context.Background()

	err := envconfig.Process(ctx, cfg)
	if err != nil {
		return nil, errors.Wrapf(err, "error while parsing env config: %s", err)
	}
	return cfg, nil
}

type Config struct {
	Version     string            `env:"VERSION,default=local"`
	Environment string            `env:"ENVIRONMENT,default=local"`
	Log         LogConfig         `env:",prefix=LOG_"`
	Runtime     RuntimeConfig     `env:",prefix=RUNTIME_"`
	GRPC        GRPCConfig        `env:",prefix=GRPC_"`
	HTTP        HTTPConfig        `env:",prefix=HTTP_"`
	RateLimiter RateLimiterConfig `env:",prefix=RATE_LIMITER_"`
	HTTPCors    HTTPCorsConfig    `env:",prefix=CORS_"`
	HealthCheck HealthCheckConfig `env:",prefix=HEALTHCHECK_"`
	Metrics     MetricsConfig     `env:",prefix=METRICS_"`
	Profiling   ProfilingConfig   `env:",prefix=PROFILING_"`
	JWT         JwtConfig         `env:",prefix=JWT_"`
	Rabbit      RabbitConfig      `env:",prefix=RABBIT_"`

	// list of clients:
	AuthApi        AuthApiConfig        `env:",prefix=AUTH_API_"`
	UserApi        UserApiConfig        `env:",prefix=USER_API_"`
	ActionApi      ActionApiConfig      `env:",prefix=ACTION_API_"`
	VocabularyApi  VocabularyApiConfig  `env:",prefix=VOCABULARY_API_"`
	SpeakerApi     SpeakerApiConfig     `env:",prefix=SPEAKER_API_"`
	LanguageApi    LanguageApiConfig    `env:",prefix=LANGUAGE_API_"`
	TranslationApi TranslationApiConfig `env:",prefix=TRANSLATION_API_"`
	GoogleApi      GoogleApiConfig      `env:",prefix=GOOGLE_API_"`
}

type LogConfig struct {
	Level  string `env:"LEVEL,default=info"`
	Secure bool   `env:"SECURE"`
}

type RuntimeConfig struct {
	UseCPUs    int `env:"USE_CPUS,default=0"`
	MaxThreads int `env:"MAX_THREADS,default=0"`
}

type GRPCConfig struct {
	Network            string `env:"NETWORK,default=tcp"`
	Address            string `env:"ADDRESS,default=:18080"`
	MaxRequestBodySize int    `env:"MAX_REQUEST_BODY_SIZE,default=26214400"` // 25Мb
}

type HTTPConfig struct {
	ReadTimeout        time.Duration `env:"READ_TIMEOUT,default=30s"`
	WriteTimeout       time.Duration `env:"WRITE_TIMEOUT,default=30s"`
	IdleTimeout        time.Duration `env:"IDLE_TIMEOUT,default=30s"`
	MaxRequestBodySize int64         `env:"MAX_REQUEST_BODY_SIZE,default=26214400"` // 25Мb
	Network            string        `env:"NETWORK,default=tcp"`
	Address            string        `env:"ADDRESS,default=:8080"`
}

type HTTPCorsConfig struct {
	Enabled          bool   `env:"ENABLED,default=true"`
	AllowedOrigins   string `env:"ALLOWED_ORIGINS,default=http://localhost:3000"`
	AllowedMethods   string `env:"ALLOWED_METHODS,default=GET, POST, OPTIONS"`
	AllowedHeaders   string `env:"ALLOWED_HEADERS,default=Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, Connection, Accept-Language, User-Agent"`
	ExposedHeaders   string `env:"EXPOSED_HEADERS,default=Link"`
	AllowCredentials string `env:"ALLOW_CREDENTIALS,default=true"`
	MaxAge           string `env:"MAX_AGE,default=3600"`
}

type HealthCheckConfig struct {
	DiskspaceThreshold uint64 `env:"DISKSPACE_THRESHOLD,default=80"`
	GoroutineThreshold int    `env:"GOROUTINE_THRESHOLD,default=20"`
	GoroutineReadiness int    `env:"GOROUTINE_READINESS,default=10"`
}

type MetricsConfig struct {
	Enabled   bool   `env:"ENABLED,default=false"`
	Path      string `env:"PATH,default=/metrics"`
	Namespace string `env:"NAMESPACE,default=wbw_gateway_api"`
}

type ProfilingConfig struct {
	Enabled bool `env:"ENABLED,default=false"`
}

type RateLimiterConfig struct {
	Enabled          bool `env:"ENABLED,default=true"`
	SignUp           int  `env:"SIGN_UP,default=30"`
	SignIn           int  `env:"SIGN_IN,default=60"`
	RefreshTokens    int  `env:"REFRESH_TOKENS,default=30"`
	ConfirmEmail     int  `env:"CONFIRM_EMAIL,default=30"`
	AskResetPassword int  `env:"ASK_RESET_PASSWORD,default=15"`
	ResetPassword    int  `env:"RESET_PASSWORD,default=30"`
	GetLanguages     int  `env:"GET_LANGUAGES,default=120"`
	Logout           int  `env:"LOGOUT,default=30"`
	GetUser          int  `env:"GET_USER,default=120"`
	UpdateUser       int  `env:"UPDATE_USER,default=60"`
	CreateCollection int  `env:"CREATE_COLLECTION,default=60"`
	UpdateCollection int  `env:"UPDATE_COLLECTION,default=60"`
	GetCollections   int  `env:"GET_COLLECTIONS,default=120"`
	GetCollection    int  `env:"GET_COLLECTION,default=120"`
	DeleteCollection int  `env:"DELETE_COLLECTION,default=60"`
	CreateTerms      int  `env:"CREATE_TERMS,default=120"`
	UpdateTerm       int  `env:"UPDATE_TERM,default=120"`
	GetTerms         int  `env:"GET_TERMS,default=120"`
	ChangeTermStatus int  `env:"CHANGE_TERM_STATUS,default=120"`
	DeleteTerms      int  `env:"DELETE_TERMS,default=120"`
	GetVoiceover     int  `env:"GET_VOICEOVER,default=120"`
	GetTranslation   int  `env:"GET_VOICEOVER,default=60"`
}

type RabbitConfig struct {
	URI                      string                         `env:"URI,required"`
	NotificationApiSendEmail NotificationApiSendEmailConfig `env:",prefix=NOTIFICATION_API_SEND_EMAIL_"`
}

type NotificationApiSendEmailConfig struct {
	Queue    string `env:"QUEUE,default=wbw-notification-send-email"`
	Exchange string `env:"EXCHANGE,default=wbw-notification"`
}

type JwtConfig struct {
	AccessSecret  string `env:"ACCESS_SECRET,default=35ea96e76e3f41b98cd66e70aecd5505"`
	RefreshSecret string `env:"REFRESH_SECRET,default=f78e9d4fbd944e4785706a9b97bfad5a"`
}

type AuthApiConfig struct {
	URI     string `env:"URI,required"`
	WithTLS bool   `env:"WITH_TLS,default=false"`
}

type UserApiConfig struct {
	URI     string `env:"URI,required"`
	WithTLS bool   `env:"WITH_TLS,default=false"`
}

type ActionApiConfig struct {
	URI     string `env:"URI,required"`
	WithTLS bool   `env:"WITH_TLS,default=false"`
}

type VocabularyApiConfig struct {
	URI     string `env:"URI,required"`
	WithTLS bool   `env:"WITH_TLS,default=false"`
}

type SpeakerApiConfig struct {
	URI     string `env:"URI,required"`
	WithTLS bool   `env:"WITH_TLS,default=false"`
}

type LanguageApiConfig struct {
	URI     string `env:"URI,required"`
	WithTLS bool   `env:"WITH_TLS,default=false"`
}

type TranslationApiConfig struct {
	URI     string `env:"URI,required"`
	WithTLS bool   `env:"WITH_TLS,default=false"`
}

type GoogleApiConfig struct {
	URI     string `env:"URI,required"`
	Timeout int64  `env:"TIMEOUT,default=10"`
}
