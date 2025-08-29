package clients

import (
	"context"
	"time"
)

type Client interface {
	MakeRequest(ctx context.Context, payloadSize int) (time.Duration, error)
	Close() error
}