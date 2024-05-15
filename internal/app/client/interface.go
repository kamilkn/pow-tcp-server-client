package client

import (
	"io"
)

type Config interface {
	ServerAddress() string
}

type Logger interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
}

type Service interface {
	RequestResource(clientID string, rw io.ReadWriter) (resource string, err error)
}
