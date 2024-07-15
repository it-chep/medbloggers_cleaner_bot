package entity

type Hashtags struct {
	hashtags []string
	chatID   int64
}

func NewHashtags(hashtags []string, chatId int64) *Hashtags {
	return &Hashtags{
		hashtags,
		chatId,
	}
}

func (h *Hashtags) GetHashtags() []string {
	return h.hashtags
}

func (h *Hashtags) GetChatId() int64 {
	return h.chatID
}
