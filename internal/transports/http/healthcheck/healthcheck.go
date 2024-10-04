package healthcheck

import (
	"context"
	"fmt"
	"github.com/etherlabsio/healthcheck/checkers"
	"github.com/etherlabsio/healthcheck/v2"
	"github.com/gin-gonic/gin"
	"github.com/rabbitmq/amqp091-go"
	"gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/config"
	"net"
	"runtime"
	"time"
)

const (
	DefaultPrefix = "/healthcheck"
)

func getPrefix(prefixOptions ...string) string {
	prefix := DefaultPrefix
	if len(prefixOptions) > 0 {
		prefix = prefixOptions[0]
	}
	return prefix
}

func Register(cfg *config.Config, r *gin.Engine, prefixOptions ...string) {
	RouteRegister(cfg, &(r.RouterGroup), prefixOptions...)
}

func RouteRegister(cfg *config.Config, rg *gin.RouterGroup, prefixOptions ...string) {
	prefixRouter := rg.Group(getPrefix(prefixOptions...))
	prefixRouter.GET("/_live", gin.WrapF(healthcheck.HandlerFunc(
		// Checking the application address
		healthcheck.WithChecker(
			"tcp", healthcheck.CheckerFunc(
				func(ctx context.Context) error {
					return TCPDialCheck(cfg.HTTP.Address, 50*time.Millisecond)
				},
			),
		),

		// Checking for the number of Goroutine
		healthcheck.WithChecker(
			"goroutine", healthcheck.CheckerFunc(
				func(ctx context.Context) error {
					count := runtime.NumGoroutine()
					if count > cfg.HealthCheck.GoroutineThreshold {
						return fmt.Errorf("too many goroutines (%d > %d)", count, cfg.HealthCheck.GoroutineThreshold)
					}
					return nil
				},
			),
		),
	)))

	prefixRouter.GET("/_ready", gin.WrapF(
		healthcheck.HandlerFunc(
			//Total timeout
			healthcheck.WithTimeout(5*time.Second),

			//Checking the AuthApi connection
			healthcheck.WithObserver(
				"auth_api_connection", healthcheck.CheckerFunc(
					func(ctx context.Context) error {
						return TCPDialCheck(cfg.AuthApi.URI, 200*time.Millisecond)
					},
				),
			),

			//Checking the AuthApi connection
			healthcheck.WithObserver(
				"user_api_connection", healthcheck.CheckerFunc(
					func(ctx context.Context) error {
						return TCPDialCheck(cfg.UserApi.URI, 200*time.Millisecond)
					},
				),
			),

			//Checking the AuthApi connection
			healthcheck.WithObserver(
				"action_api_connection", healthcheck.CheckerFunc(
					func(ctx context.Context) error {
						return TCPDialCheck(cfg.ActionApi.URI, 200*time.Millisecond)
					},
				),
			),

			//Checking the VocabularyApi connection
			healthcheck.WithObserver(
				"vocabulary_api_connection", healthcheck.CheckerFunc(
					func(ctx context.Context) error {
						return TCPDialCheck(cfg.VocabularyApi.URI, 200*time.Millisecond)
					},
				),
			),

			//Checking the SpeakerApi connection
			healthcheck.WithObserver(
				"speaker_api_connection", healthcheck.CheckerFunc(
					func(ctx context.Context) error {
						return TCPDialCheck(cfg.SpeakerApi.URI, 200*time.Millisecond)
					},
				),
			),

			//Checking the LanguageApi connection
			healthcheck.WithObserver(
				"language_api_connection", healthcheck.CheckerFunc(
					func(ctx context.Context) error {
						return TCPDialCheck(cfg.LanguageApi.URI, 200*time.Millisecond)
					},
				),
			),

			//Checking the LanguageApi connection
			healthcheck.WithObserver(
				"translation_api_connection", healthcheck.CheckerFunc(
					func(ctx context.Context) error {
						return TCPDialCheck(cfg.TranslationApi.URI, 200*time.Millisecond)
					},
				),
			),

			//-------------------
			// Checking the Rabbit connection
			//-------------------
			healthcheck.WithObserver(
				"rabbit_connection", healthcheck.CheckerFunc(
					func(ctx context.Context) error {
						connection, err := amqp091.ParseURI(cfg.Rabbit.URI)
						if err != nil {
							return err
						}
						return TCPDialCheck(fmt.Sprintf("%s:%d", connection.Host, connection.Port), 50*time.Millisecond)
					},
				),
			),

			//Number Goroutine
			healthcheck.WithObserver(
				"goroutine", healthcheck.CheckerFunc(
					func(ctx context.Context) error {
						count := runtime.NumGoroutine()
						if count > cfg.HealthCheck.GoroutineReadiness {
							return fmt.Errorf("too many goroutines (%d > %d)", count, cfg.HealthCheck.GoroutineReadiness)
						}
						return nil
					},
				),
			),

			//We output an error if the disk is occupied by more than a threshold (Take out in a hundred)
			healthcheck.WithObserver(
				"diskspace", checkers.DiskSpace("/", cfg.HealthCheck.DiskspaceThreshold),
			),
		),
	))
}

func TCPDialCheck(addr string, timeout time.Duration) error {
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return err
	}
	return conn.Close()
}
