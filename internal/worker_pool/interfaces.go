package worker_pool

import (
	"context"
)

type Task interface {
	Process(ctx context.Context) error
}
