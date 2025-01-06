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

type MessageStorage struct {
	client postgres.Client
	logger *slog.Logger
}

func NewMessageStorage(client postgres.Client, logger *slog.Logger) *MessageStorage {
	return &MessageStorage{
		client: client,
		logger: logger,
	}
}

func (rs *MessageStorage) GetAvailableHashtags(ctx context.Context, chatId int64) (hashtags *entity.Hashtags, err error) {
	op := "storage/read_repo/messages/GetAvailableHashtags"
	q := `
		select array_agg(h.hashtag) as hashtags, ch.chat_id 
		from hashtags h 
		left join chats ch 
		on ch.id = h.chat_id 
		where ch.chat_id = $1
		group by ch.chat_id;
	`

	var hashtagsDAO dao.HashtagsDAO
	err = pgxscan.Get(ctx, rs.client, &hashtagsDAO, q, chatId)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		rs.logger.Error(fmt.Sprintf("Error while scanning row: %s, op: %s", err, op))
		return nil, err
	}

	hashtags = hashtagsDAO.ToDomain()
	return hashtags, nil
}
