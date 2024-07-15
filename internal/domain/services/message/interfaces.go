package message

import (
	"context"
	"docstar_cleaner_bot/internal/domain/entity"
)

type readMessageRepo interface {
	GetAvailableHashtags(ctx context.Context, chatId int64) (hashtags *entity.Hashtags, err error)
}
