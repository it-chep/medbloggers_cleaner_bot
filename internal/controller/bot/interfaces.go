package bot

import (
	"context"
	"medbloggers_cleaner_bot/internal/controller/dto/tg"
	"medbloggers_cleaner_bot/internal/domain/entity"
)

type userService interface {
	RegisterNewUser(ctx context.Context, dto tg.TgUserDTO) (user entity.User, err error)
	ValidateMessage(ctx context.Context, userDTO tg.TgUserDTO, messageDTO tg.MessageDTO) error
}
