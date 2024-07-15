package write_repo

import (
	"context"
	"docstar_cleaner_bot/internal/domain/entity"
	"docstar_cleaner_bot/pkg/client/postgres"
	"fmt"
	"log/slog"
)

type UserStorage struct {
	client postgres.Client
	logger *slog.Logger
}

func NewUserStorage(client postgres.Client, logger *slog.Logger) UserStorage {
	return UserStorage{
		client: client,
		logger: logger,
	}
}

func (ws UserStorage) CreateUser(ctx context.Context, user entity.User) (userID int64, err error) {
	op := "internal/storage/write_repo/CreateUser"
	q := `
		insert into tg_users (tg_id, name, username, last_message_timestamp, chat_id, admin) 
		values ($1, $2, $3, now(), $4, $5) returning id;
	`

	err = ws.client.QueryRow(
		ctx, q, user.GetTgId(), user.GetFirstName(), user.GetUsername(), user.GetChatId(), user.IsAdmin(),
	).Scan(&userID)
	if err != nil {
		ws.logger.Error(fmt.Sprintf("%s op %s", err, op))
		return -1, err
	}

	return userID, nil
}

func (ws UserStorage) BlockUsers(ctx context.Context) (err error) {
	op := "internal/storage/write_repo/UpdateUserIsBlocked"
	q := `
		update tg_users 
		set is_blocked = true 
		where last_message_timestamp >= now() - interval '5 minutes' and message_counter >= 5;
	`

	_, err = ws.client.Query(ctx, q)
	if err != nil {
		ws.logger.Error(fmt.Sprintf("%s op %s", err, op))
		return err

	}

	return err
}

func (ws UserStorage) CleanLastMessage(ctx context.Context) (err error) {
	op := "internal.storage.write_repo.UpdateLastMessageTimestamp"
	q := `
		update tg_users 
		set message_counter = 0
		where last_message_timestamp <= now() - interval '5 minutes';
	`
	_, err = ws.client.Exec(ctx, q)
	if err != nil {
		ws.logger.Error(fmt.Sprintf("%s op %s", err, op))
		return err
	}
	return nil
}

func (ws UserStorage) UpdateLastMessageCounter(ctx context.Context, user entity.User, messageCount int) (err error) {
	op := "internal.storage.write_repo.UpdateLastMessageCounter"
	q := `
		update tg_users 
		set message_counter = $1
		where tg_id = $2 and chat_id = $3;
	`

	_, err = ws.client.Exec(ctx, q, messageCount, user.GetTgId(), user.GetChatId())
	if err != nil {
		ws.logger.Error(fmt.Sprintf("%s op %s", err, op))
		return err
	}

	return nil
}

func (ws UserStorage) UpdateLastMessageTimestamp(ctx context.Context, user entity.User) (err error) {
	op := "internal.storage.write_repo.UpdateLastMessageCounter"
	q := `
		update tg_users 
		set last_message_timestamp = now()
		where tg_id = $1 and chat_id = $2;
	`

	_, err = ws.client.Exec(ctx, q, user.GetTgId(), user.GetChatId())
	if err != nil {
		ws.logger.Error(fmt.Sprintf("%s op %s", err, op))
		return err
	}

	return nil
}
