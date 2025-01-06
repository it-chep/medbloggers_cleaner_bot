package callback

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"medbloggers_cleaner_bot/internal/controller/dto/tg"
	"medbloggers_cleaner_bot/pkg/client/telegram"
	"strconv"
	"strings"
)

type CallbackBotMessage struct {
	logger *slog.Logger
	bot    telegram.Bot
	tgUser tg.TgUserDTO
}

func NewCallbackBot(
	logger *slog.Logger,
	bot telegram.Bot,
	tgUser tg.TgUserDTO,

) CallbackBotMessage {
	return CallbackBotMessage{
		logger: logger,
		bot:    bot,
		tgUser: tgUser,
	}
}

// Execute место связи telegram и бизнес логи
func (c *CallbackBotMessage) Execute(ctx context.Context, messageDTO tg.MessageDTO, callbackData string) {
	var msg tgbotapi.MessageConfig

	if strings.Contains(callbackData, "approve") {
		stringBeforeParams := strings.Split(callbackData, ":")
		params := strings.Split(stringBeforeParams[1], ",")

		chatId, _ := strconv.ParseInt(strings.Split(params[0], "=")[1], 10, 64)
		userId, _ := strconv.ParseInt(strings.Split(params[1], "=")[1], 10, 64)

		escapedText := escapeMarkdownV2(messageDTO.Text)
		quotedText := fmt.Sprintf("❇️Объявление:\n\n%s\n\n \nОПУБЛИКОВАНО", quoteText(escapedText))
		msg = tgbotapi.NewMessage(messageDTO.Chat.ID, quotedText)
		msg.ParseMode = "MarkdownV2"
		c.bot.SendMessage(msg)

		c.bot.RemoveMessage(messageDTO.Chat.ID, messageDTO.MessageID)

		msg = tgbotapi.NewMessage(chatId, messageDTO.Text)
		c.bot.SendMessage(msg)

		msg = tgbotapi.NewMessage(userId, "Сообщение опубликовано")
		c.bot.SendMessage(msg)
	}

	if strings.Contains(callbackData, "reject") {
		stringBeforeParams := strings.Split(callbackData, ":")
		params := strings.Split(stringBeforeParams[1], ",")

		userId, _ := strconv.ParseInt(strings.Split(params[1], "=")[1], 10, 64)

		escapedText := escapeMarkdownV2(messageDTO.Text)
		quotedText := fmt.Sprintf("🚫Объявление:\n\n%s\n\n \nОТКЛОНЕНО", quoteText(escapedText))
		msg = tgbotapi.NewMessage(messageDTO.Chat.ID, quotedText)
		msg.ParseMode = "MarkdownV2"
		c.bot.SendMessage(msg)

		c.bot.RemoveMessage(messageDTO.Chat.ID, messageDTO.MessageID)

		msg = tgbotapi.NewMessage(userId, "Ваше сообщение не прошло модерацию, пожалуйста исправьте его")
		c.bot.SendMessage(msg)
	}
}

func escapeMarkdownV2(text string) string {
	replacer := strings.NewReplacer(
		"_", "\\_",
		"*", "\\*",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
		"~", "\\~",
		"`", "\\`",
		">", "\\>",
		"#", "\\#",
		"+", "\\+",
		"-", "\\-",
		"=", "\\=",
		"|", "\\|",
		"{", "\\{",
		"}", "\\}",
		".", "\\.",
		"!", "\\!",
	)
	return replacer.Replace(text)
}

func quoteText(text string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = "> " + line
	}
	return strings.Join(lines, "\n")
}
