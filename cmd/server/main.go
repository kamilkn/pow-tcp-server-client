package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kamilkn/pow-tcp-server-client/internal/app/server"
	"github.com/kamilkn/pow-tcp-server-client/internal/pkg/lib/cache"
	"github.com/kamilkn/pow-tcp-server-client/internal/pkg/lib/config"
	"github.com/kamilkn/pow-tcp-server-client/internal/pkg/lib/log"
	"github.com/kamilkn/pow-tcp-server-client/internal/pkg/lib/tcp"
	"github.com/kamilkn/pow-tcp-server-client/internal/pkg/service"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	configuration, err := config.Parse("config")
	if err != nil {
		fmt.Println(err.Error()) //nolint:forbidigo // print error.
		os.Exit(1)
	}

	configService := newConfigService(configuration)
	configServer := newConfigServer(configuration)

	logger := log.New(log.Opts{
		Level: log.Level(configuration.Server.LogLevel),
		JSON:  configuration.Server.LogJSON,
	})

	puzzleCache := cache.New[string, struct{}](ctx, cache.Opts{
		CleanInterval: configService.PuzzleTTL(),
		Logger:        logger,
	})

	resourceCache := cache.New[int, string](ctx, cache.Opts{
		Logger: logger,
	})
	for i, r := range resources {
		resourceCache.Add(i, r)
	}

	mainService := service.NewServer(&service.ServerOpts{
		Config:        configService,
		Logger:        logger,
		PuzzleCache:   puzzleCache,
		ResourceCache: resourceCache,
		ErrorChecker:  tcp.NewConnErrorChecker(),
	})

	mainServer, err := server.Listen(ctx, server.Opts{
		Config:  configServer,
		Logger:  logger,
		Service: mainService,
	})
	if err != nil {
		fmt.Println(err.Error()) //nolint:forbidigo // print error.
		os.Exit(1)
	}

	logger.Debug("server started",
		"address", configServer.Address(),
		"shutdown_timeout", configServer.ShutdownTimeout(),
		"connection_timeout", configServer.ConnectionTimeout(),
		"puzzle_ttl", configService.PuzzleTTL(),
		"puzzle_zero_bits", configService.PuzzleZeroBits(),
	)

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	<-signalChannel
	cancel()
	mainServer.Shutdown()
}
