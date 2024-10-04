package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/config"
	"gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/constants"
	"gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/prometheus"
	"gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/transports/http"
	"gitlab.com/wordbyword.io/microservices/pkg/logger"

	ActionApiClient "gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/clients/action_api"
	AuthApiClient "gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/clients/auth_api"
	GoogleAuthApiClient "gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/clients/google_auth_api"
	LanguageApiClient "gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/clients/language_api"
	NotificationApiClient "gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/clients/notification_api"
	SpeakerApiClient "gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/clients/speaker_api"
	TranslationApiClient "gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/clients/translation_api"
	UserApiClient "gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/clients/user_api"
	VocabularyApiClient "gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/clients/vocabulary_api"
	GatewayApiService "gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/services/gateway_api/service"
	HttpEndpoints "gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/transports/http/endpoints"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalln(err)
	}

	prom := prometheus.NewExporter(cfg.Metrics.Namespace)
	lgr, err := logger.NewLogger(os.Stdout, cfg.Log.Level)
	if err != nil {
		log.Fatalln(err)
	}

	lgr = lgr.With().
		Str("app", constants.AppName).
		Str("version", cfg.Version).
		CallerWithSkipFrameCount(2).
		Logger()

	initRuntime(cfg, lgr)
	runApp(cfg, lgr, prom)
}

func runApp(cfg *config.Config, lgr zerolog.Logger, prom *prometheus.Exporter) {
	var err error
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM)

	//---------------------------
	// 1) Clients Initialization
	//---------------------------
	notificationApiClient := NotificationApiClient.NewNotificationApiClient(cfg, lgr)
	defer notificationApiClient.Shutdown()

	authApiClient := AuthApiClient.NewAuthApiClient(cfg, lgr)
	defer authApiClient.Shutdown()

	userApiClient := UserApiClient.NewUserApiClient(cfg, lgr)
	defer userApiClient.Shutdown()

	actionApiClient := ActionApiClient.NewActionApiClient(cfg, lgr)
	defer actionApiClient.Shutdown()

	vocabularyApiClient := VocabularyApiClient.NewVocabularyApiClient(cfg, lgr)
	defer vocabularyApiClient.Shutdown()

	speakerApiClient := SpeakerApiClient.NewSpeakerApiClient(cfg, lgr)
	defer speakerApiClient.Shutdown()

	languageApiClient := LanguageApiClient.NewLanguageApiClient(cfg, lgr)
	defer languageApiClient.Shutdown()

	translationApiClient := TranslationApiClient.NewTranslationApiClient(cfg, lgr)
	defer translationApiClient.Shutdown()

	googleAuthApiClient := GoogleAuthApiClient.NewGoogleAuthApiClient(cfg, lgr)
	defer googleAuthApiClient.Shutdown()

	//---------------------------
	// 2) Services Initialization
	//---------------------------
	gatewayApiService := GatewayApiService.NewGatewayApiService(cfg, lgr, prom,
		notificationApiClient,
		authApiClient,
		userApiClient,
		actionApiClient,
		vocabularyApiClient,
		speakerApiClient,
		languageApiClient,
		translationApiClient,
		googleAuthApiClient,
	)
	defer gatewayApiService.Shutdown()

	//---------------------------
	// 3) Transports Initialization
	//---------------------------
	httpEndpoints, err := HttpEndpoints.NewHttpEndpoints(cfg, lgr, prom, gatewayApiService, authApiClient)
	if err != nil {
		lgr.Fatal().Err(err).Msg("failed to initialize http api endpoints")
	}

	//---------------------------
	// 4) Starting the HTTP server
	//---------------------------
	httpListener, err := net.Listen(cfg.HTTP.Network, cfg.HTTP.Address)
	if err != nil {
		lgr.Fatal().Err(err).Msg("failed to init net.Listen for http")
	}

	httpServer, err := http.NewServer(cfg, lgr, prom, httpListener, httpEndpoints)
	if err != nil {
		lgr.Fatal().Err(err).Stack().Msg("failed to init http server")
	}

	httpErrCh := httpServer.Start()

	//---------------------------
	// 5) Errors handling
	//---------------------------
	runningApp := true
	for runningApp {
		select {
		case err = <-httpErrCh:
			if err != nil {
				lgr.Error().Stack().Err(err).Msg("received http server error")
				shutdownCh <- os.Kill
			}
		case sig := <-shutdownCh:
			lgr.Info().Str("signal", sig.String()).Msg("shutdown signal received")

			ctxTimeout, timeoutCancelFunc := context.WithTimeout(ctx, 10*time.Second)
			defer timeoutCancelFunc()

			err = httpServer.Shutdown(ctxTimeout)
			if err != nil {
				lgr.Error().Stack().Err(err).Msg("received http shutdown error")
			}

			lgr.Info().Msg("server loop stopped")
			runningApp = false
			break
		}
	}
}
