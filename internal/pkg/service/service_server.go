package service

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"time"

	"github.com/kamilkn/pow-tcp-server-client/internal/pkg/lib/hashcash"
	"github.com/kamilkn/pow-tcp-server-client/internal/pkg/lib/message"
)

// Opts - options to create new cache instance.
type ServerOpts struct {
	Logger        Logger
	Config        ServerConfig
	PuzzleCache   PuzzleCache
	ResourceCache ResourceCache
	ErrorChecker  ErrorChecker
}

// NewServer - create new server-side service.
func NewServer(opts *ServerOpts) *Server {
	return &Server{
		logger:        opts.Logger,
		config:        opts.Config,
		puzzleCache:   opts.PuzzleCache,
		resourceCache: opts.ResourceCache,
		errorChecker:  opts.ErrorChecker,
	}
}

// Server - server-side service.
type Server struct {
	logger        Logger
	config        ServerConfig
	puzzleCache   PuzzleCache
	resourceCache ResourceCache
	errorChecker  ErrorChecker
}

// HandleMessages - handle client messages.
func (s *Server) HandleMessages(clientID string, reader io.ReadWriter) {
	const operationName = "service.Server.HandleMessages"

	s.logger.Info("connected new client", "clientID", clientID)

	for {
		rawMsg, err := bufio.NewReader(reader).ReadString(message.DelimiterMessage)
		if err != nil {
			clientErr := ErrInternalError
			if s.errorChecker.IsTimeout(err) {
				clientErr = ErrTimeoutExceeded
				s.logger.Info(clientErr.Error(), "clientID", clientID)
			} else {
				s.logger.Error(err.Error(), "op", operationName, "clientID", clientID)
			}

			s.writeError(clientID, clientErr, reader)

			return
		}

		msg, err := message.ParseMessage(rawMsg)
		if err != nil {
			s.logger.Info(ErrIncorrectMessageFormat.Error(), "clientID", clientID, "message", rawMsg)
			s.writeError(clientID, ErrIncorrectMessageFormat, reader)

			return
		}

		switch msg.Command {
		case message.CommandRequestPuzzle:
			s.responsePuzzle(clientID, msg.Payload, reader)
		case message.CommandRequestResource:
			s.responseResource(clientID, msg.Payload, reader)

			return
		case message.CommandError, message.CommandResponsePuzzle, message.CommandResponseResource:
			s.responseResource(clientID, msg.Payload, reader)

			return
		default:
			s.writeError(clientID, ErrIncorrectMessageFormat, reader)

			return
		}
	}
}

func (s *Server) responsePuzzle(clientID, _ string, w io.Writer) {
	const operationName = "service.Server.responsePuzzle"

	s.logger.Info("requested new puzzle", "clientID", clientID)

	mainHashcash, err := hashcash.New(s.config.PuzzleZeroBits(), clientID)
	if err != nil {
		s.logger.Error(err.Error(), "op", operationName, "clientID", clientID)
		s.writeError(clientID, ErrInternalError, w)

		return
	}

	exp := time.Now().Add(s.config.PuzzleTTL())
	s.puzzleCache.AddWithExp(mainHashcash.Key(), struct{}{}, exp)

	msg := message.Message{
		Command: message.CommandResponsePuzzle,
		Payload: string(mainHashcash.Header()),
	}

	s.writeMsg(clientID, msg, w)
	s.logger.Info("puzzle sent", "clientID", clientID, "puzzle", msg.Payload)
}

func (s *Server) responseResource(clientID, payload string, w io.Writer) {
	const operationName = "service.Server.responseResource"

	s.logger.Info("requested resource", "clientID", clientID, "solution", payload)

	mainHashcash, err := hashcash.ParseHeader(payload)
	if err != nil {
		s.logger.Info(ErrHashcashHeaderNotCorrect.Error(), "clientID", clientID, "header", payload)
		s.writeError(clientID, ErrHashcashHeaderNotCorrect, w)

		return
	}

	if _, ok := s.puzzleCache.Get(mainHashcash.Key()); !ok {
		s.logger.Info(ErrHashcashHeaderNotFound.Error(), "clientID", clientID, "header", payload)
		s.writeError(clientID, ErrHashcashHeaderNotFound, w)

		return
	}

	if !mainHashcash.EqualResource(clientID) {
		s.logger.Info(ErrHashcashHeaderNotFound.Error(), "clientID", clientID, "header", payload)
		s.writeError(clientID, ErrHashcashHeaderNotFound, w)

		return
	}

	if !mainHashcash.IsActual(s.config.PuzzleTTL()) {
		s.logger.Info(ErrHashcashExpirationExceeded.Error(), "clientID", clientID, "header", payload)
		s.writeError(clientID, ErrHashcashExpirationExceeded, w)

		return
	}

	isHashCorrect, err := mainHashcash.Header().IsHashCorrect(mainHashcash.Bits())
	if err != nil {
		s.logger.Error(err.Error(), "op", operationName, "clientID", clientID)
		s.writeError(clientID, ErrInternalError, w)

		return
	}

	if !isHashCorrect {
		s.logger.Info(ErrHashcashHeaderNotCorrect.Error(), "clientID", clientID, "header", payload)
		s.writeError(clientID, ErrHashcashHeaderNotCorrect, w)

		return
	}

	resource, err := s.randomResource()
	if err != nil {
		s.logger.Error(err.Error(), "op", operationName, "clientID", clientID)
		s.writeError(clientID, ErrInternalError, w)

		return
	}

	msg := message.Message{
		Command: message.CommandResponseResource,
		Payload: resource,
	}

	s.writeMsg(clientID, msg, w)
	s.puzzleCache.Delete(mainHashcash.Key())
	s.logger.Info("resource sent", "clientID", clientID, "resource", msg.Payload)
}

func (s *Server) randomResource() (string, error) {
	keys := s.resourceCache.Keys()
	if len(keys) == 0 {
		return "", nil
	}

	randKeyIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(keys))))
	if err != nil {
		return "", fmt.Errorf("randKeyIndex get error: %w", err)
	}

	key := int(randKeyIndex.Int64())
	resource, _ := s.resourceCache.Get(key)

	return resource, nil
}

func (s *Server) writeMsg(clientID string, msg message.Message, w io.Writer) {
	const operationName = "service.Server.writeMsg"

	if _, err := w.Write(msg.Bytes()); err != nil {
		s.logger.Error(err.Error(), "op", operationName, "clientID", clientID)
	}
}

func (s *Server) writeError(clientID string, handleErr error, w io.Writer) {
	const operationName = "service.Server.writeError"

	if _, err := w.Write(errorMessage(handleErr).Bytes()); err != nil {
		s.logger.Error(err.Error(), "op", operationName, "clientID", clientID)
	}
}
