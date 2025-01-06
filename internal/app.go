package internal

import (
	"context"
	"log"
	"log/slog"
	"medbloggers_cleaner_bot/internal/controller"
	"medbloggers_cleaner_bot/internal/domain/services/message"
	"medbloggers_cleaner_bot/internal/domain/services/user"
	"medbloggers_cleaner_bot/internal/storage/read_repo"
	"medbloggers_cleaner_bot/internal/storage/write_repo"
	"medbloggers_cleaner_bot/internal/worker_pool"
	"medbloggers_cleaner_bot/internal/worker_pool/tasks"
	"medbloggers_cleaner_bot/pkg/client/postgres"
	"medbloggers_cleaner_bot/pkg/client/telegram"
	"net/http"
	"os"
	"time"
)

func (app *App) InitLogger(ctx context.Context) *App {
	app.logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return app
}

func (app *App) InitPgxConn(ctx context.Context) *App {
	client, err := postgres.NewClient(ctx, app.config.StorageConfig)
	if err != nil {
		log.Fatal(err)
	}
	app.pgxClient = client
	app.logger.Info("init pgxclient", app.pgxClient)
	return app
}

func (app *App) InitStorage(ctx context.Context) *App {
	app.storages.readUserStorage = read_repo.NewUserStorage(app.pgxClient, app.logger)
	app.storages.readMessageStorage = read_repo.NewMessageStorage(app.pgxClient, app.logger)
	app.storages.writeUserStorage = write_repo.NewUserStorage(app.pgxClient, app.logger)
	return app
}

func (app *App) InitServices(ctx context.Context) *App {
	app.services.messageService = message.NewMessageService(
		app.storages.readMessageStorage,
		app.logger,
	)
	app.services.userService = user.NewUserService(
		app.storages.writeUserStorage,
		app.storages.readUserStorage,
		app.services.messageService,
		app.logger,
	)
	return app
}

func (app *App) InitTelegram(ctx context.Context) *App {
	app.bot = *telegram.NewTelegramBot(*app.config, app.logger)
	return app
}

func (app *App) InitControllers(ctx context.Context) *App {
	app.controller.telegramWebhookController = controller.NewRestController(*app.config, app.logger, app.bot, app.services.userService)
	app.controller.telegramWebhookController.InitController()

	app.server = &http.Server{
		Addr:         app.config.HTTPServer.Address,
		Handler:      app.controller.telegramWebhookController,
		ReadTimeout:  2 * time.Minute,
		WriteTimeout: 10 * time.Second,
	}
	return app
}

func (app *App) InitTasks(ctx context.Context) *App {
	app.periodicalTasks.cleanMessageTask = tasks.NewCleanLastMessageTask(app.logger, app.services.userService)
	return app
}

func (app *App) InitWorkers(ctx context.Context) *App {
	workers := []worker_pool.Worker{
		worker_pool.NewWorker(app.periodicalTasks.cleanMessageTask, 5*time.Minute),
	}
	app.workerPool = worker_pool.NewWorkerPool(workers)
	return app
}
