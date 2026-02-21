package errors

import (
	"fmt"
	"time"
)

var (
	ErrNoEndpoints      = fmt.Errorf("no endpoints available")
	ErrNoTransport      = fmt.Errorf("no transport available")
	ErrNotConnected     = fmt.Errorf("not connected")
	ErrTransportClosed  = fmt.Errorf("transport is closed")
	ErrNilMessage       = fmt.Errorf("message is nil")
	ErrIdRequired       = fmt.Errorf("id is required")
	ErrDestinationEmpty = fmt.Errorf("destination is empty")
	ErrTimeout          = fmt.Errorf("timeout")
	ErrReplyChannel     = fmt.Errorf("reply channel is not set")
	ErrNoRouteFound     = fmt.Errorf("no route found for topic")
)

// RetriableError is an error that can be retried.
// It is used by the client to retry sending a message.
// It contains the original error and the time to wait before retrying.
type RetriableError struct {
	Err  error
	Wait time.Duration
}

func (e *RetriableError) Error() string {
	if e.Err == nil {
		return ""
	}
	return e.Err.Error()
}

func (e *RetriableError) Unwrap() error {
	return e.Err
}
