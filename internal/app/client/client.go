package client

import (
	"fmt"
	"net"
)

// Opts - connection options.
type Opts struct {
	Config  Config
	Logger  Logger
	Service Service
}

// Connect - connect to server.
func Connect(opts Opts) error {
	const operationName = "client.Connect"

	conn, err := net.Dial("tcp", opts.Config.ServerAddress())
	if err != nil {
		opts.Logger.Error(err.Error(), "operationName", operationName)

		return fmt.Errorf("TCP dial: %w", err)
	}

	defer conn.Close()

	_, err = opts.Service.RequestResource(conn.LocalAddr().String(), conn)
	if err != nil {
		return fmt.Errorf("RequestResource: %w", err)
	}

	return nil
}
