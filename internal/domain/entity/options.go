package entity

import "time"

type UserOpt func(usr *User) *User

func WithUsrUsername(username *string) UserOpt {
	return func(usr *User) *User {
		usr.username = username
		return usr
	}
}

func WithAdmin(admin bool) UserOpt {
	return func(usr *User) *User {
		usr.admin = admin
		return usr
	}
}
func WithLastMessageTimestamp(lastMessageTimestamp *time.Time) UserOpt {
	return func(usr *User) *User {
		usr.lastMessageTimestamp = lastMessageTimestamp
		return usr
	}
}

func WithMessageCounter(messageCount int) UserOpt {
	return func(usr *User) *User {
		usr.messageCount = messageCount
		return usr
	}
}
