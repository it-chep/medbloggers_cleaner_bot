package dao

import (
	"docstar_cleaner_bot/internal/domain/entity"
	"time"
)

type UserDAO struct {
	TgId                 int64      `db:"tg_id"`
	FirstName            string     `db:"name"`
	Username             *string    `db:"username"`
	LastMessageTimestamp *time.Time `db:"last_message_timestamp"`
	MessageCounter       int        `db:"message_counter"`
	IsBlocked            bool       `db:"is_blocked"`
	ChatId               int64      `db:"chat_id"`
	Admin                bool       `db:"admin"`
}

func NewUserDAO() *UserDAO {
	return &UserDAO{}
}

func (dao *UserDAO) ToDomain() *entity.User {
	return entity.NewUser(
		dao.TgId,
		dao.FirstName,
		dao.ChatId,
		entity.WithUsrUsername(dao.Username),
		entity.WithLastMessageTimestamp(dao.LastMessageTimestamp),
		entity.WithMessageCounter(dao.MessageCounter),
		entity.WithAdmin(dao.Admin),
	)
}
