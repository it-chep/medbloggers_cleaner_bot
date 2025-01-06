package controller

import (
	"github.com/gin-gonic/gin"
	"medbloggers_cleaner_bot/internal/config"
	"medbloggers_cleaner_bot/pkg/client/telegram"

	"log/slog"
	"net/http"

	botapi "medbloggers_cleaner_bot/internal/controller/bot"
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
	r.router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
	r.router.POST("/"+r.cfg.Bot.Token+"/", r.botApiController.BotWebhookHandler)
}

func (r *RestController) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}
