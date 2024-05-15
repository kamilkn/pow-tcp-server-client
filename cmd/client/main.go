package main

import (
	"fmt"
	"os"

	"github.com/kamilkn/pow-tcp-server-client/internal/app/client"
	"github.com/kamilkn/pow-tcp-server-client/internal/pkg/lib/config"
	"github.com/kamilkn/pow-tcp-server-client/internal/pkg/lib/log"
	"github.com/kamilkn/pow-tcp-server-client/internal/pkg/service"
)

func main() {
	configuration, err := config.Parse("config")
	if err != nil {
		fmt.Println(err.Error()) //nolint:forbidigo // print error.
		os.Exit(1)
	}

	configService := newConfigService(configuration)
	configClient := newConfigClient(configuration)

	logger := log.New(log.Opts{
		Level: log.Level(configuration.Client.LogLevel),
		JSON:  configuration.Client.LogJSON,
	})

	mainService := service.NewClient(service.ClientOpts{
		Config: configService,
		Logger: logger,
	})

	logger.Debug("client configured",
		"server_address", configClient.ServerAddress(),
		"puzzle_compute_max_attempts", configService.PuzzleComputeMaxAttempts(),
	)

	err = client.Connect(client.Opts{
		Config:  configClient,
		Logger:  logger,
		Service: mainService,
	})
	if err != nil {
		fmt.Println(err.Error()) //nolint:forbidigo // print error.
		os.Exit(1)
	}
}
