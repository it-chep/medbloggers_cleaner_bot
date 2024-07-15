package controller

import (
	"docstar_cleaner_bot/internal/config"
	"docstar_cleaner_bot/pkg/client/telegram"
	"github.com/gin-gonic/gin"

	"log/slog"
	"net/http"

	botapi "docstar_cleaner_bot/internal/controller/bot"
)

type RestController struct {
	router           *gin.Engine
	cfg              config.Config
	logger           *slog.Logger
	botApiController botapi.TelegramWebhookController
	userService      userService
}

func NewRestController(
	cfg config.Config,
	logger *slog.Logger,
	bot telegram.Bot,
	userService userService,
) *RestController {
	router := gin.New()
	router.Use(gin.Recovery())

	botApiController := botapi.NewTelegramWebhookController(
		cfg, logger, bot, userService,
	)

	return &RestController{
		router:           router,
		cfg:              cfg,
		logger:           logger,
		botApiController: botApiController,
	}
}

func (r *RestController) InitController() {
	r.router.POST("/"+r.cfg.Bot.Token+"/", r.botApiController.BotWebhookHandler)
}

func (r *RestController) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}
