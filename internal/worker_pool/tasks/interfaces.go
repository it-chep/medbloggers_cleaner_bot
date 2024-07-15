package tasks

import (
	"context"
)

type userService interface {
	CleanLastMessage(ctx context.Context) (err error)
}
