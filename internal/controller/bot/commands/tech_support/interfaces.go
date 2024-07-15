package tech_support

import (
	"context"
	"docstar_cleaner_bot/internal/controller/dto/tg"
	"docstar_cleaner_bot/pkg/client/telegram"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"os"
	"strconv"
)

type TechSupportCommand struct {
	logger *slog.Logger
	bot    telegram.Bot
	tgUser tg.TgUserDTO
}

func NewTechSupportCommand(logger *slog.Logger, tgUser tg.TgUserDTO, bot telegram.Bot) TechSupportCommand {
	return TechSupportCommand{
		logger: logger,
		tgUser: tgUser,
		bot:    bot,
	}
}

// Execute место связи telegram и бизнес логи
func (c *TechSupportCommand) Execute(ctx context.Context, message tg.MessageDTO) {
	msg := tgbotapi.NewMessage(
		c.tgUser.TgID,
		fmt.Sprintf("Зову техподдержку"),
	)
	c.bot.SendMessage(msg)

	adminId, err := strconv.Atoi(os.Getenv("ADMIN_ID"))

	if err != nil {
		return
	}

	msg = tgbotapi.NewMessage(
		int64(adminId),
		fmt.Sprintf("Кто-то жмет ТП %s", c.tgUser.UserName),
	)
	c.bot.SendMessage(msg)

	return
}
