package services

import (
	"context"
	"errors"

	"cz.Finance/backend/models"
	"cz.Finance/backend/repositories"
	"cz.Finance/backend/utils"
)

// TelegramServiceImpl представляет реализацию сервиса телеграм
type TelegramServiceImpl struct {
	telegramRepo repositories.TelegramUserRepository
	userRepo     repositories.UserRepository
}

// NewTelegramService создает новый экземпляр сервиса телеграм
func NewTelegramService(telegramRepo repositories.TelegramUserRepository, userRepo repositories.UserRepository) TelegramService {
	return &TelegramServiceImpl{
		telegramRepo: telegramRepo,
		userRepo:     userRepo,
	}
}

// LinkAccount связывает аккаунт Telegram с аккаунтом пользователя
func (s *TelegramServiceImpl) LinkAccount(ctx context.Context, telegramID int64, username, firstName, lastName string, request *models.TelegramLinkRequest) (*models.User, error) {
	// Проверяем, существует ли уже связь для этого Telegram ID
	existingLink, err := s.telegramRepo.GetByTelegramID(ctx, telegramID)
	if err == nil && existingLink != nil {
		return nil, errors.New("этот Telegram аккаунт уже связан с пользователем")
	}

	// Получаем пользователя по email
	user, err := s.userRepo.GetByEmail(ctx, request.Email)
	if err != nil {
		return nil, errors.New("пользователь с таким email не найден")
	}

	// Проверяем пароль
	if !utils.CheckPasswordHash(request.Password, user.PasswordHash) {
		return nil, errors.New("неверный пароль")
	}

	// Создаем новую связь
	telegramUser := &models.TelegramUser{
		UserID:     user.ID,
		TelegramID: telegramID,
		Username:   username,
		FirstName:  firstName,
		LastName:   lastName,
	}

	_, err = s.telegramRepo.Create(ctx, telegramUser)
	if err != nil {
		return nil, errors.New("ошибка при связывании аккаунтов")
	}

	return user, nil
}

// GetUserByTelegramID получает пользователя по Telegram ID
func (s *TelegramServiceImpl) GetUserByTelegramID(ctx context.Context, telegramID int64) (*models.User, error) {
	// Получаем связь
	telegramUser, err := s.telegramRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return nil, errors.New("аккаунт Telegram не связан с пользователем")
	}

	// Получаем пользователя
	user, err := s.userRepo.GetByID(ctx, telegramUser.UserID)
	if err != nil {
		return nil, errors.New("пользователь не найден")
	}

	return user, nil
}

// UnlinkAccount отвязывает аккаунт Telegram от аккаунта пользователя
func (s *TelegramServiceImpl) UnlinkAccount(ctx context.Context, telegramID int64) error {
	// Получаем связь
	telegramUser, err := s.telegramRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return errors.New("аккаунт Telegram не связан с пользователем")
	}

	// Удаляем связь
	return s.telegramRepo.Delete(ctx, telegramUser.ID)
}
