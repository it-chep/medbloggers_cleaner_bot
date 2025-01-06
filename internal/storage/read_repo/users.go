package read_repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"medbloggers_cleaner_bot/internal/domain/entity"
	"medbloggers_cleaner_bot/internal/storage/dao"
	"medbloggers_cleaner_bot/pkg/client/postgres"
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

func (rs UserStorage) GetUserByTgID(ctx context.Context, userID, chatId int64) (user entity.User, err error) {
	op := "internal/storage/read_repo/users/GetUserByTgID"
	q := `select tg_id, name, username, last_message_timestamp, message_counter, is_blocked, chat_id, admin from tg_users where tg_id = $1 and chat_id = $2`

	var userDAO dao.UserDAO
	err = pgxscan.Get(ctx, rs.client, &userDAO, q, userID, chatId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, nil
		}
		rs.logger.Error(fmt.Sprintf("Error while scanning row: %s, op: %s", err, op))
		return entity.User{}, err
	}

	user = *userDAO.ToDomain()
	return user, nil
}
