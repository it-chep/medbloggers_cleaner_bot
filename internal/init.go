package internal

import (
	"context"
	"docstar_cleaner_bot/internal/config"
	"docstar_cleaner_bot/internal/controller"
	"docstar_cleaner_bot/internal/domain/services/message"
	"docstar_cleaner_bot/internal/domain/services/user"
	"docstar_cleaner_bot/internal/storage/read_repo"
	"docstar_cleaner_bot/internal/storage/write_repo"
	"docstar_cleaner_bot/internal/worker_pool"
	"docstar_cleaner_bot/internal/worker_pool/tasks"
	"docstar_cleaner_bot/pkg/client/postgres"
	"docstar_cleaner_bot/pkg/client/telegram"
	"log/slog"
	"net/http"
)

type controllers struct {
	telegramWebhookController *controller.RestController
}

type services struct {
	userService    user.UserService
	messageService *message.MessageService
}

type storages struct {
	readUserStorage    read_repo.UserStorage
	readMessageStorage *read_repo.MessageStorage
	writeUserStorage   write_repo.UserStorage
}

type periodicalTasks struct {
	cleanMessageTask tasks.CleanLastMessageTask
}

type App struct {
	logger          *slog.Logger
	config          *config.Config
	controller      controllers
	services        services
	storages        storages
	periodicalTasks periodicalTasks
	workerPool      worker_pool.WorkerPool
	bot             telegram.Bot
	pgxClient       postgres.Client
	server          *http.Server
}

func NewApp(ctx context.Context) *App {
	cfg := config.NewConfig()

	app := &App{
		config: cfg,
	}

	app.InitLogger(ctx).
		InitPgxConn(ctx).
		InitStorage(ctx).
		InitServices(ctx).
		InitTasks(ctx).
		InitWorkers(ctx).
		InitTelegram(ctx).
		InitControllers(ctx)

	return app
}

func (app *App) Run(ctx context.Context) error {
	app.logger.Info("start server")
	app.workerPool.Run(ctx)
	return app.server.ListenAndServe()
}
