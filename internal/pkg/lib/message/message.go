package message

import (
	"fmt"
	"strings"
)

// Command - command type.
type Command int8

const (
	// CommandError - using when something went wrong by client or server.
	CommandError Command = iota

	// CommandRequestPuzzle - using when client requests puzzle from server.
	CommandRequestPuzzle

	// CommandResponsePuzzle - using when server sends puzzle to client.
	CommandResponsePuzzle

	// CommandRequestResource - using when client requests resource from server.
	CommandRequestResource

	// CommandResponseResource - using when server sends resource to client.
	CommandResponseResource
)

const (
	// DelimiterMessage - sign to divide messages from each other.
	DelimiterMessage = '\n'

	// DelimiterCommand - sign to divide command and payload in message.
	DelimiterCommand = ':'
)

// string has "command:payload" format where command could be 0-4.
func ParseMessage(msg string) (Message, error) {
	var (
		message          Message
		maxMessageLength = 2
	)

	msg = strings.TrimSpace(msg)

	if len(msg) < maxMessageLength {
		return message, ErrIncorrectMessageFormat
	}

	switch msg[:maxMessageLength] {
	case fmt.Sprintf("0%c", DelimiterCommand):
		message.Command = CommandError
	case fmt.Sprintf("1%c", DelimiterCommand):
		message.Command = CommandRequestPuzzle
	case fmt.Sprintf("2%c", DelimiterCommand):
		message.Command = CommandResponsePuzzle
	case fmt.Sprintf("3%c", DelimiterCommand):
		message.Command = CommandRequestResource
	case fmt.Sprintf("4%c", DelimiterCommand):
		message.Command = CommandResponseResource
	default:
		return message, ErrIncorrectMessageFormat
	}

	if len(msg) > maxMessageLength {
		message.Payload = strings.TrimSpace(msg[maxMessageLength:])
	}

	return message, nil
}

// Message - message with command and payload.
type Message struct {
	Command Command
	Payload string
}

// String - format message as string.
func (m Message) String() string {
	return fmt.Sprintf("%d%c%s%c", m.Command, DelimiterCommand, m.Payload, DelimiterMessage)
}

// Bytes - format message as bytes.
func (m Message) Bytes() []byte {
	return []byte(m.String())
}
