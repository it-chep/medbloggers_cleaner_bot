package user

import (
	"context"
	"medbloggers_cleaner_bot/internal/domain/entity"
)

type readUserStorage interface {
	GetUserByTgID(ctx context.Context, tgID, chatId int64) (user entity.User, err error)
}

type writeUserStorage interface {
	CreateUser(ctx context.Context, user entity.User) (userId int64, err error)
	UpdateLastMessageCounter(ctx context.Context, user entity.User, messageCount int) (err error)
	UpdateLastMessageTimestamp(ctx context.Context, user entity.User) (err error)
	BlockUsers(ctx context.Context) (err error)
	CleanLastMessage(ctx context.Context) (err error)
}

type messageService interface {
	GetAvailableHashtags(ctx context.Context, chatId int64) (*entity.Hashtags, error)
}
