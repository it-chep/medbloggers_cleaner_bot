package bot

import (
	"context"
	"docstar_cleaner_bot/internal/config"
	"docstar_cleaner_bot/internal/controller/bot/commands/health_check"
	"docstar_cleaner_bot/internal/controller/bot/commands/start"
	"docstar_cleaner_bot/internal/controller/bot/commands/tech_support"
	textmessage "docstar_cleaner_bot/internal/controller/bot/text"
	"docstar_cleaner_bot/internal/controller/dto/tg"
	"docstar_cleaner_bot/pkg/client/telegram"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"log/slog"
	"net/http"
)

type TelegramWebhookController struct {
	cfg         config.Config
	logger      *slog.Logger
	bot         telegram.Bot
	userService userService
}

func NewTelegramWebhookController(
	cfg config.Config,
	logger *slog.Logger,
	bot telegram.Bot,
	userService userService,
) TelegramWebhookController {

	return TelegramWebhookController{
		cfg:         cfg,
		logger:      logger,
		bot:         bot,
		userService: userService,
	}
}

func (t TelegramWebhookController) BotWebhookHandler(c *gin.Context) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.logger.Error(fmt.Sprintf("%s", err))
		}
	}(c.Request.Body)

	var update tgbotapi.Update
	if err := c.ShouldBindJSON(&update); err != nil {
		t.logger.Error("Error binding JSON", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	tgUser := t.getUserFromWebhook(update)
	tgMessage := t.getMessageFromWebhook(update)

	// Сначала проверяем на команду, потом на текстовое сообщение, потом callback
	if update.Message != nil {
		ctx := context.WithValue(context.Background(), "userID", update.Message.From.ID)
		if update.Message.IsCommand() {
			t.ForkCommands(ctx, update, tgUser, tgMessage)
		} else {
			t.ForkMessages(ctx, tgUser, tgMessage)
		}
	} else if update.CallbackQuery != nil {
		ctx := context.WithValue(context.Background(), "userID", update.CallbackQuery.From.ID)
		t.ForkCallbacks(ctx, update, tgUser, tgMessage)
	} else if update.MyChatMember != nil {
		// событие о добавлении бота в чат
		c.JSON(http.StatusOK, gin.H{"status": "received"})
		return
	} else if update.ChannelPost != nil {
		// событие о посте в канале
		c.JSON(http.StatusOK, gin.H{"status": "received"})
		return
	} else {
		t.logger.Warn(fmt.Sprintf("Unhandled update type: %+v", update))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "received"})
	return
}

func (t TelegramWebhookController) ForkCommands(ctx context.Context, update tgbotapi.Update, tgUser tg.TgUserDTO, tgMessage tg.MessageDTO) {

	switch update.Message.Command() {
	case "start":
		command := start.NewStartBotCommand(t.logger, tgUser, t.bot, t.userService)
		command.Execute(ctx, tgMessage)
	case "health_check":
		command := health_check.NewHealthCheckCommand(t.logger, tgUser, t.bot)
		command.Execute(ctx, tgMessage)
	case "tech_support":
		command := tech_support.NewTechSupportCommand(t.logger, tgUser, t.bot)
		command.Execute(ctx, tgMessage)
	}
}

// todo в эти форки будут сыпаться все текстовые сообщения и колбэки
// todo и в зависимости от состояния пользователя ему будет выдаваться контент

func (t TelegramWebhookController) ForkMessages(ctx context.Context, tgUser tg.TgUserDTO, tgMessage tg.MessageDTO) {
	messageBot := textmessage.NewTextBotMessage(t.logger, t.bot, tgUser, t.userService)
	messageBot.Execute(ctx, tgMessage)
}

func (t TelegramWebhookController) ForkCallbacks(ctx context.Context, update tgbotapi.Update, tgUser tg.TgUserDTO, tgMessage tg.MessageDTO) {
	callbackData := update.CallbackData()
	t.logger.Info(callbackData)
}

func (t TelegramWebhookController) getUserFromWebhook(update tgbotapi.Update) tg.TgUserDTO {
	var tgUser tg.TgUserDTO
	var userJSON []byte
	var err error

	if update.Message == nil {
		return tg.TgUserDTO{}
	}
	t.logger.Info("update", update.CallbackQuery, update.Message.Chat)
	// Todo возможно улучшить
	if update.CallbackQuery != nil {
		userJSON, err = json.Marshal(update.CallbackQuery.From)
	} else if update.Message.From != nil {
		userJSON, err = json.Marshal(update.Message.From)
	} else {
		userJSON, err = json.Marshal(update.Message.Chat)
	}

	if err != nil {
		t.logger.Error(fmt.Sprintf("Error marshaling user to JSON: %s", err))
		return tg.TgUserDTO{}
	}

	if err = json.Unmarshal(userJSON, &tgUser); err != nil {
		t.logger.Error(fmt.Sprintf("Error decoding JSON: %s", err))
		return tg.TgUserDTO{}
	}

	tgUser.ChatId = &update.Message.Chat.ID

	return tgUser
}

func (t TelegramWebhookController) getMessageFromWebhook(update tgbotapi.Update) tg.MessageDTO {
	var tgMessage tg.MessageDTO
	var userJSON []byte
	var err error

	// Todo возможно улучшить
	if update.CallbackQuery != nil {
		userJSON, err = json.Marshal(update.CallbackQuery.Message)
	} else {
		userJSON, err = json.Marshal(update.Message)
	}

	if err != nil {
		t.logger.Error(fmt.Sprintf("Error marshaling user to JSON: %s", err))
		return tg.MessageDTO{}
	}

	if err = json.Unmarshal(userJSON, &tgMessage); err != nil {
		t.logger.Error(fmt.Sprintf("Error decoding JSON: %s", err))
		return tg.MessageDTO{}
	}

	return tgMessage
}
