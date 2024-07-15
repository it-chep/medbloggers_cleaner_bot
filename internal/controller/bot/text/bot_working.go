package text_message

import (
	"context"
	"docstar_cleaner_bot/internal/controller/dto/tg"
	"docstar_cleaner_bot/internal/enums"
	"docstar_cleaner_bot/pkg/client/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
)

type TextBotMessage struct {
	logger      *slog.Logger
	bot         telegram.Bot
	tgUser      tg.TgUserDTO
	userService userService
}

func NewTextBotMessage(logger *slog.Logger, bot telegram.Bot, tgUser tg.TgUserDTO, userService userService) TextBotMessage {
	return TextBotMessage{
		logger:      logger,
		bot:         bot,
		tgUser:      tgUser,
		userService: userService,
	}
}

func (c TextBotMessage) Execute(ctx context.Context, messageDTO tg.MessageDTO) {
	var _ tgbotapi.MessageConfig

	err := c.userService.ValidateMessage(ctx, c.tgUser, messageDTO)

	if err == nil {
		return
	}

	if err.Error() == enums.NeedDeleteUser {
		c.bot.RemoveMessage(messageDTO.Chat.ID, int(messageDTO.MessageID))
		c.bot.RemoveUserFromChat(messageDTO.Chat.ID, c.tgUser.TgID)
		return
	}

	if err.Error() == enums.InvalidHashtag || err.Error() == enums.LittleInterval || err.Error() == enums.LittleWeekInterval {
		c.bot.RemoveMessage(messageDTO.Chat.ID, int(messageDTO.MessageID))
		return
	}

	if err != nil {
		c.logger.Info("Bot catch error: ", err.Error(), "Chat ID:", *c.tgUser.ChatId)
	}
	return
}
