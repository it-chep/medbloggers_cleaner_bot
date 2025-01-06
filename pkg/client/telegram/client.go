package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"medbloggers_cleaner_bot/internal/config"
)

type Bot struct {
	Bot    *tgbotapi.BotAPI
	logger *slog.Logger
}

func NewTelegramBot(cfg config.Config, logger *slog.Logger) *Bot {
	bot, err := tgbotapi.NewBotAPI(cfg.Bot.Token)
	bot.Debug = true
	if err != nil {
		panic("can't create bot instance")
	}

	wh, _ := tgbotapi.NewWebhook(cfg.Bot.WebhookURL + bot.Token + "/")
	_, err = bot.Request(wh)
	if err != nil {
		panic("can't while request set webhook")
	}

	_, err = bot.GetWebhookInfo()

	if err != nil {
		panic("error while getting webhook")
	}
	return &Bot{
		Bot:    bot,
		logger: logger,
	}
}

func (bot *Bot) SendMessage(msg tgbotapi.MessageConfig) {
	_, err := bot.Bot.Send(msg)
	if err != nil {
		bot.logger.Error(fmt.Sprintf("%s: Bot SendMessage", err))
	}
}

func (bot *Bot) RemoveMessage(chatId int64, messageId int64) {
	messageToDelete := tgbotapi.NewDeleteMessage(chatId, int(messageId))
	_, err := bot.Bot.Send(messageToDelete)
	if err != nil {
		return
	}
}

func (bot *Bot) RemoveUserFromChat(chatId int64, userId int64) {
	kickConfig := tgbotapi.KickChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: chatId,
			UserID: userId,
		},
	}

	_, err := bot.Bot.Request(kickConfig)
	if err != nil {
		bot.logger.Error(fmt.Sprintf("%s: Bot RemoveUserFromChat", err))
	}
}
