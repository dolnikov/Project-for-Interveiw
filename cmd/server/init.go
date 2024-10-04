package main

import (
	"github.com/rs/zerolog"
	"runtime"
	"runtime/debug"

	"gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/config"
)

func initRuntime(cfg *config.Config, lgr zerolog.Logger) {
	cpu := cfg.Runtime.UseCPUs
	threads := cfg.Runtime.MaxThreads

	if cpu == 0 {
		cpu = runtime.NumCPU()
		runtime.GOMAXPROCS(runtime.NumCPU())
	} else {
		runtime.GOMAXPROCS(cpu)
	}
	lgr.Info().Msgf("set to use %d CPUs", cpu)
	if threads == 0 {
		threads = 10000
	} else {
		debug.SetMaxThreads(threads)
	}
	lgr.Info().Msgf("set to use maximum %d threads", threads)
}
