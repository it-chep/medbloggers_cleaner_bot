package health_check

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"medbloggers_cleaner_bot/internal/controller/dto/tg"
	"medbloggers_cleaner_bot/pkg/client/telegram"
)

type HealthCheckCommand struct {
	logger *slog.Logger
	bot    telegram.Bot
	tgUser tg.TgUserDTO
}

func NewHealthCheckCommand(logger *slog.Logger, tgUser tg.TgUserDTO, bot telegram.Bot) HealthCheckCommand {
	return HealthCheckCommand{
		logger: logger,
		tgUser: tgUser,
		bot:    bot,
	}
}

// Execute место связи telegram и бизнес логи
func (c *HealthCheckCommand) Execute(ctx context.Context, message tg.MessageDTO) {
	msg := tgbotapi.NewMessage(
		c.tgUser.TgID,
		fmt.Sprintf("Жив здоров"),
	)
	c.bot.SendMessage(msg)
	return
}
