package main

import (
	"github.com/kamilkn/pow-tcp-server-client/internal/pkg/lib/config"
)

func newConfigClient(c *config.Config) *configClient {
	return &configClient{
		c: c,
	}
}

type configClient struct {
	c *config.Config
}

func (cc *configClient) ServerAddress() string {
	return cc.c.Client.ServerAddress
}

func newConfigService(c *config.Config) *configService {
	return &configService{
		c: c,
	}
}

type configService struct {
	c *config.Config
}

func (cs *configService) PuzzleComputeMaxAttempts() int {
	return cs.c.Hashcash.ComputeMaxAttempts
}
