package tasks

import (
	"context"
	"log/slog"
)

type CleanLastMessageTask struct {
	logger      *slog.Logger
	userService userService
}

func NewCleanLastMessageTask(logger *slog.Logger, userService userService) CleanLastMessageTask {
	return CleanLastMessageTask{
		logger:      logger,
		userService: userService,
	}
}

func (task CleanLastMessageTask) Process(ctx context.Context) error {
	err := task.userService.CleanLastMessage(ctx)
	if err != nil {
		return err
	}
	return nil
}
