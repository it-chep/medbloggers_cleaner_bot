package text_message

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"medbloggers_cleaner_bot/internal/config"
	"medbloggers_cleaner_bot/internal/controller/bot/chat_data"
	"medbloggers_cleaner_bot/internal/controller/dto/tg"
	"medbloggers_cleaner_bot/pkg/client/telegram"
	"regexp"
)

type TextBotMessage struct {
	logger *slog.Logger
	bot    telegram.Bot
	tgUser tg.TgUserDTO
	cfg    config.Config
}

func NewTextBotMessage(logger *slog.Logger, bot telegram.Bot, tgUser tg.TgUserDTO, cfg config.Config) TextBotMessage {
	return TextBotMessage{
		logger: logger,
		bot:    bot,
		tgUser: tgUser,
		cfg:    cfg,
	}
}

func (c TextBotMessage) Execute(ctx context.Context, messageDTO tg.MessageDTO) {
	adminChatId := c.cfg.Chats.AdminChat
	if messageDTO.Chat.ID == adminChatId || len(messageDTO.Text) == 0 {
		return
	}

	chatMapper := chat_data.GetChatMapper(c.cfg)

	re := regexp.MustCompile(`#\p{L}[\p{L}0-9_]*`)
	hashtags := re.FindAllString(messageDTO.Text, -1)

	var (
		toChatId   int64
		hasHashtag bool
	)

	for _, hashtag := range hashtags {
		if chatId, ok := chatMapper[hashtag]; ok {
			toChatId = chatId
			hasHashtag = true
			break
		}
	}

	if hasHashtag {
		// пересылаем сообщение в чат админов
		msgToAdmin := tgbotapi.NewMessage(adminChatId, "Новое объявление, прошу промодерировать его")
		c.bot.SendMessage(msgToAdmin)

		msgWithButtons := tgbotapi.NewMessage(adminChatId, messageDTO.Text)
		buttons := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					"Опубликовать",
					fmt.Sprintf("approve:chat_id=%d,user_id=%d", toChatId, c.tgUser.TgID),
				),
				tgbotapi.NewInlineKeyboardButtonData(
					"Отклонить",
					fmt.Sprintf("reject:chat_id=%d,user_id=%d", toChatId, c.tgUser.TgID),
				),
			),
		)
		msgWithButtons.ReplyMarkup = buttons
		c.bot.SendMessage(msgWithButtons)

		msgToUser := tgbotapi.NewMessage(messageDTO.Chat.ID, "Сообщение отправлено на модерацию")
		c.bot.SendMessage(msgToUser)

	} else {
		msg := tgbotapi.NewMessage(c.tgUser.TgID, "Пожалуйста ознакомьтесь с правилами постинга сообщений в чаты")
		c.bot.SendMessage(msg)
	}

	return
}
