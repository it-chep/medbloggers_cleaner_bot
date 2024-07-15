package start

import (
	"context"
	"docstar_cleaner_bot/internal/controller/dto/tg"
	"docstar_cleaner_bot/pkg/client/telegram"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
)

type StartBotCommand struct {
	logger      *slog.Logger
	bot         telegram.Bot
	tgUser      tg.TgUserDTO
	userService userService
}

func NewStartBotCommand(logger *slog.Logger, tgUser tg.TgUserDTO, bot telegram.Bot, userService userService) StartBotCommand {
	return StartBotCommand{
		logger:      logger,
		tgUser:      tgUser,
		bot:         bot,
		userService: userService,
	}
}

// Execute место связи telegram и бизнес логи
func (c *StartBotCommand) Execute(ctx context.Context, message tg.MessageDTO) {
	_, err := c.userService.RegisterNewUser(ctx, c.tgUser)
	if err != nil {
		return
	}

	msg := tgbotapi.NewMessage(
		c.tgUser.TgID,
		fmt.Sprintf("Привет от @medbloggers_cleaner_bot, меня создал @maxim_jordan на языке golang для очистки чатов от сообщений без нужных хэштегов"),
	)
	c.bot.SendMessage(msg)
	return
}
