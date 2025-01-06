package chat_data

import (
	"medbloggers_cleaner_bot/internal/config"
)

func GetChatMapper(cfg config.Config) map[string]int64 {
	return map[string]int64{
		"#ищувмедблогерс":       cfg.Chats.JobChat,
		"#помогувмедблогерс":    cfg.Chats.JobChat,
		"#ищунарекламу":         cfg.Chats.AdvertisingChat,
		"#ищувп":                cfg.Chats.AdvertisingChat,
		"#рекламноепредложение": cfg.Chats.AdvertisingChat,
		"#отзывмедблогерс":      cfg.Chats.AdvertisingChat,
	}
}
