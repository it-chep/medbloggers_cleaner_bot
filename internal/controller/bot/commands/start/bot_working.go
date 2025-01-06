package start

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"medbloggers_cleaner_bot/internal/controller/dto/tg"
	"medbloggers_cleaner_bot/pkg/client/telegram"
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
		fmt.Sprintf(
			"Отправьте сообщение, которое хотите опубликовать в \n🟣 "+
				"«[Чат вакансий врачей-блогеров MEDBLOGERS](https://t.me/docstar_job/198)» \n\nили в\n\n🟠 «[Чат рекламы врачей-блогеров MEDBLOGERS](https://t.me/docstar_ad/44)»"+
				"\n\nНе забудьте указать корректный хэштег для публикации в зависимости от запроса, например:"+
				" #ищувмедблогерс #помогувмедблогерс #ищунарекламу и тд\n\nПодробнее о хэштегах в закрепах обоих чатов.",
		),
	)
	msg.ParseMode = "Markdown"
	msg.DisableWebPagePreview = true
	c.bot.SendMessage(msg)
	return
}
