package dao

import "medbloggers_cleaner_bot/internal/domain/entity"

type HashtagsDAO struct {
	Hashtags []string `db:"hashtags"`
	ChatId   int64    `db:"chat_id"`
}

func (h HashtagsDAO) ToDomain() *entity.Hashtags {
	return entity.NewHashtags(
		h.Hashtags,
		h.ChatId,
	)
}
