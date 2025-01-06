package message

import (
	"context"
	"log/slog"
	"medbloggers_cleaner_bot/internal/domain/entity"
)

type MessageService struct {
	readRepo readMessageRepo
	logger   *slog.Logger
}

func NewMessageService(readRepo readMessageRepo, logger *slog.Logger) *MessageService {
	return &MessageService{
		readRepo: readRepo,
		logger:   logger,
	}
}

func (ms *MessageService) GetAvailableHashtags(ctx context.Context, chatId int64) (*entity.Hashtags, error) {
	return ms.readRepo.GetAvailableHashtags(ctx, chatId)
}
