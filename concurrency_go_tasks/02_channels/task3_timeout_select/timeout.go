package timeout

import (
	"context"
	"time"
)

func Work(ctx context.Context) error {
	select {
	case <-time.After(100 * time.Millisecond):
		return context.DeadlineExceeded
	case <-ctx.Done():
		return ctx.Err()
	}
}
