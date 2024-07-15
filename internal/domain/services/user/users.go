package user

import (
	"context"
	"docstar_cleaner_bot/internal/controller/dto/tg"
	"docstar_cleaner_bot/internal/domain/entity"
	"docstar_cleaner_bot/internal/enums"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"strings"
	"time"
)

type UserService struct {
	readRepo       readUserStorage
	writeRepo      writeUserStorage
	messageService messageService
	logger         *slog.Logger
}

func NewUserService(
	writeRepo writeUserStorage,
	readRepo readUserStorage,
	messageService messageService,
	logger *slog.Logger,
) UserService {
	return UserService{
		writeRepo:      writeRepo,
		readRepo:       readRepo,
		messageService: messageService,
		logger:         logger,
	}
}

func (u UserService) GetUser(ctx context.Context, tgId, chatId int64) (user entity.User, err error) {
	op := "sorkin_bot.internal.domain.services.user.users.GetUser"

	//if not in cache
	user, err = u.readRepo.GetUserByTgID(ctx, tgId, chatId)

	if err != nil {
		u.logger.Error(fmt.Sprintf("error: %s, place %s", err, op))
		return entity.User{}, err
	}
	return user, nil
}

func (u UserService) RegisterNewUser(ctx context.Context, dto tg.TgUserDTO) (user entity.User, err error) {
	_ = "sorkin_bot.internal.domain.services.user.users.RegisterNewUser"
	u.logger.Info("fvfv", dto)
	user, err = u.GetUser(ctx, dto.TgID, *dto.ChatId)
	if err != nil {
		return entity.User{}, err
	}

	if !reflect.ValueOf(user.GetTgId()).IsZero() {
		return user, nil
	}

	newUser := dto.ToDomain(
		[]entity.UserOpt{
			entity.WithUsrUsername(&dto.UserName),
			entity.WithAdmin(dto.Admin),
		},
	)

	//// Save new user
	_, err = u.writeRepo.CreateUser(ctx, newUser)
	if err != nil {
		return entity.User{}, err
	}

	return newUser, nil
}

func (u UserService) ValidateMessage(ctx context.Context, userDTO tg.TgUserDTO, messageDTO tg.MessageDTO) error {
	valid := false
	user, err := u.RegisterNewUser(ctx, userDTO)
	if err != nil {
		return err
	}
	// Если это сообщение от чата
	if userDTO.TgID == messageDTO.Chat.ID || user.IsAdmin() {
		return nil
	}

	valid, err = u.validateMessageText(ctx, messageDTO)
	if err != nil {
		return err
	}

	messageCount := user.GetMessageCount() + 1

	// Если человек пишет больше 3 раз за минуту, то обновляем ему timestamp
	if messageCount >= 3 {
		err = u.writeRepo.UpdateLastMessageTimestamp(ctx, user)
		if err != nil {
			return err
		}
	}

	if !valid {
		err = u.writeRepo.UpdateLastMessageCounter(ctx, user, messageCount)
		if err != nil {
			return err
		}
		if messageCount == 5 {
			return errors.New(enums.NeedDeleteUser)
		}
		return errors.New(enums.InvalidHashtag)
	}

	// todo user pointer
	// todo minutes to config

	if user.GetLastMessageTimestamp() != nil && user.GetLastMessageTimestamp().After(time.Now().Add(-5*time.Minute)) {
		err = u.writeRepo.UpdateLastMessageCounter(ctx, user, messageCount)
		if err != nil {
			return err
		}
		if messageCount == 5 {
			return errors.New(enums.NeedDeleteUser)
		}
		return errors.New(enums.LittleInterval)
	}

	if user.GetLastMessageTimestamp() != nil && user.GetLastMessageTimestamp().After(time.Now().Add(-7*24*time.Hour)) {
		err = u.writeRepo.UpdateLastMessageCounter(ctx, user, messageCount)
		if err != nil {
			return err
		}
		return errors.New(enums.LittleWeekInterval)
	}

	return nil
}

func (u UserService) validateMessageText(ctx context.Context, messageDTO tg.MessageDTO) (valid bool, err error) {
	valid = false
	hashtags, err := u.messageService.GetAvailableHashtags(ctx, messageDTO.Chat.ID)
	if err != nil {
		return valid, err
	}
	for _, hashtag := range hashtags.GetHashtags() {
		if strings.Contains(messageDTO.Text, hashtag) {
			valid = true
		}
	}
	return valid, nil
}

func (u UserService) CleanLastMessage(ctx context.Context) (err error) {
	err = u.writeRepo.BlockUsers(ctx)
	if err != nil {
		return err
	}
	err = u.writeRepo.CleanLastMessage(ctx)
	if err != nil {
		return err
	}
	return nil
}
