package service

import (
	"errors"
	"fmt"

	"github.com/kamilkn/pow-tcp-server-client/internal/pkg/lib/message"
)

var (
	ErrIncorrectMessageFormat     = errors.New("incorrect message format")
	ErrTimeoutExceeded            = errors.New("timeout exceeded")
	ErrUnknownCommand             = errors.New("unknown command")
	ErrHashcashHeaderNotFound     = errors.New("hashcash header not found")
	ErrHashcashHeaderNotCorrect   = errors.New("hashcash header not correct")
	ErrHashcashExpirationExceeded = errors.New("hashcash expiration exceeded")
	ErrInternalError              = errors.New("internal error")
	ErrResponseCommandNotcorrect  = errors.New("response command is not correct")
	ErrCheckResMessage            = func(resMsg message.Message) error { //nolint:gochecknoglobals // pure functions.
		return fmt.Errorf("checkResMessage: %s", resMsg.Payload) //nolint:goerr113 // error message.
	}
)

func errorMessage(err error) message.Message {
	return message.Message{
		Command: message.CommandError,
		Payload: err.Error(),
	}
}
