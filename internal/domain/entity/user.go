package entity

import "time"

type User struct {
	firstName            string
	tgID                 int64
	username             *string
	messageCount         int
	lastMessageTimestamp *time.Time
	chatId               int64
	admin                bool
}

func NewUser(tgId int64, firstName string, chatId int64, opts ...UserOpt) *User {
	u := &User{
		tgID:      tgId,
		firstName: firstName,
		chatId:    chatId,
	}

	for _, opt := range opts {
		opt(u)
	}

	return u
}

func (usr *User) GetFirstName() string {
	return usr.firstName
}

func (usr *User) GetTgId() int64 {
	return usr.tgID
}

func (usr *User) GetChatId() int64 {
	return usr.chatId
}

func (usr *User) GetUsername() *string {
	return usr.username
}

func (usr *User) GetMessageCount() int {
	return usr.messageCount
}

func (usr *User) GetLastMessageTimestamp() *time.Time {
	return usr.lastMessageTimestamp
}

func (usr *User) IsAdmin() bool {
	return usr.admin
}
