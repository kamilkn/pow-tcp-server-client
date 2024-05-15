package tcp

import (
	"errors"
	"net"
)

// NewConnErrorChecker - create new connection error checker.
func NewConnErrorChecker() *ConnErrorChecker {
	return &ConnErrorChecker{}
}

// ConnErrorChecker - tcp connection error checker.
type ConnErrorChecker struct{}

// IsTimeout - define that tcp connection was timed out.
func (ec *ConnErrorChecker) IsTimeout(err error) bool {
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}

	return false
}
