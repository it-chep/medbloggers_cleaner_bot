package controller

import (
	"context"
	"docstar_cleaner_bot/internal/controller/dto/tg"
	"docstar_cleaner_bot/internal/domain/entity"
)

type userService interface {
	RegisterNewUser(ctx context.Context, dto tg.TgUserDTO) (user entity.User, err error)
	ValidateMessage(ctx context.Context, userDTO tg.TgUserDTO, messageDTO tg.MessageDTO) error
}
