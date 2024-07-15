package text_message

import (
	"context"
	"docstar_cleaner_bot/internal/controller/dto/tg"
)

type userService interface {
	ValidateMessage(ctx context.Context, userDTO tg.TgUserDTO, messageDTO tg.MessageDTO) error
}
