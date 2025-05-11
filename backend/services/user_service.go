package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cz.Finance/backend/models"
	"cz.Finance/backend/repositories"
	"cz.Finance/backend/utils"
)

// UserServiceImpl представляет реализацию сервиса пользователя
type UserServiceImpl struct {
	userRepo    repositories.UserRepository
	authService AuthService
}

// NewUserService создает новый экземпляр сервиса пользователя
func NewUserService(userRepo repositories.UserRepository, authService AuthService) UserService {
	return &UserServiceImpl{
		userRepo:    userRepo,
		authService: authService,
	}
}

// SignUp регистрирует нового пользователя
func (s *UserServiceImpl) SignUp(ctx context.Context, signup *models.UserSignup) (*models.TokenResponse, error) {
	// Проверка, что пользователя с таким email не существует
	existingUser, err := s.userRepo.GetByEmail(ctx, signup.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("пользователь с таким email уже существует")
	}

	// Проверка, что пользователя с таким username не существует
	existingUser, err = s.userRepo.GetByUsername(ctx, signup.Username)
	if err == nil && existingUser != nil {
		return nil, errors.New("пользователь с таким именем пользователя уже существует")
	}

	// Хешируем пароль
	passwordHash, err := utils.HashPassword(signup.Password)
	if err != nil {
		return nil, errors.New("ошибка при хешировании пароля")
	}

	// Создаем нового пользователя
	user := &models.User{
		Email:        signup.Email,
		Username:     signup.Username,
		PasswordHash: passwordHash,
		FirstName:    signup.FirstName,
		LastName:     signup.LastName,
		MonthlyLimit: 0,
		SavingsGoal:  0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Сохраняем пользователя в базе данных
	userID, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, errors.New("ошибка при создании пользователя")
	}

	// Устанавливаем ID пользователя
	user.ID = userID

	// Генерируем JWT токен
	token, expiresAt, err := s.authService.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, errors.New("ошибка при создании токена авторизации")
	}

	// Формируем ответ
	userResponse := user.ToUserResponse()
	return &models.TokenResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      userResponse,
	}, nil
}

// Login аутентифицирует пользователя
func (s *UserServiceImpl) Login(ctx context.Context, login *models.UserLogin) (*models.TokenResponse, error) {
	// Добавляем логирование для отладки
	fmt.Printf("Попытка входа с email: %s\n", login.Email)

	// Получаем пользователя по email
	user, err := s.userRepo.GetByEmail(ctx, login.Email)
	if err != nil {
		fmt.Printf("Ошибка при поиске пользователя по email %s: %v\n", login.Email, err)
		return nil, errors.New("неверный email или пароль")
	}

	fmt.Printf("Пользователь найден: ID=%d, Email=%s, Username=%s, PasswordHash=%s\n", user.ID, user.Email, user.Username, user.PasswordHash)

	// Проверяем пароль
	passwordMatch := utils.CheckPasswordHash(login.Password, user.PasswordHash)
	fmt.Printf("Проверка пароля: %v (введено: %s, хеш в БД: %s)\n", passwordMatch, login.Password, user.PasswordHash)

	if !passwordMatch {
		return nil, errors.New("неверный email или пароль")
	}

	// Генерируем JWT токен
	token, expiresAt, err := s.authService.GenerateToken(user.ID, user.Email)
	if err != nil {
		fmt.Printf("Ошибка генерации токена: %v\n", err)
		return nil, errors.New("ошибка при создании токена авторизации")
	}

	fmt.Printf("Токен успешно создан для пользователя ID=%d, Email=%s\n", user.ID, user.Email)

	// Формируем ответ
	userResponse := user.ToUserResponse()
	return &models.TokenResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      userResponse,
	}, nil
}

// GetUser получает информацию о пользователе
func (s *UserServiceImpl) GetUser(ctx context.Context, id int64) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("пользователь не найден")
	}

	// Выводим данные пользователя до преобразования
	fmt.Printf("Получен пользователь из базы: ID=%d, Username=%s, AvatarPath=%s\n", user.ID, user.Username, user.AvatarPath)

	userResponse := user.ToUserResponse()

	// Выводим данные после преобразования
	fmt.Printf("Подготовлен ответ: ID=%d, Username=%s, AvatarURL=%s\n", userResponse.ID, userResponse.Username, userResponse.AvatarURL)

	return &userResponse, nil
}

// UpdateUser обновляет информацию о пользователе
func (s *UserServiceImpl) UpdateUser(ctx context.Context, id int64, updateRequest *models.UpdateUserRequest) (*models.UserResponse, error) {
	// Валидируем запрос
	if err := utils.ValidateStruct(updateRequest); err != nil {
		return nil, fmt.Errorf("ошибка валидации: %w", err)
	}

	// Получаем текущего пользователя
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("пользователь не найден")
	}

	// Проверяем, не занят ли email другим пользователем, если он был изменен
	if updateRequest.Email != nil && *updateRequest.Email != user.Email {
		existingUser, err := s.userRepo.GetByEmail(ctx, *updateRequest.Email)
		if err == nil && existingUser != nil && existingUser.ID != user.ID {
			return nil, errors.New("email уже используется другим пользователем")
		}
	}

	// Проверяем, не занято ли имя пользователя другим пользователем, если оно было изменено
	if updateRequest.Username != nil && *updateRequest.Username != user.Username {
		existingUser, err := s.userRepo.GetByUsername(ctx, *updateRequest.Username)
		if err == nil && existingUser != nil && existingUser.ID != user.ID {
			return nil, errors.New("имя пользователя уже используется")
		}
	}

	// Обновляем поля, если они указаны в запросе
	if updateRequest.FirstName != nil {
		user.FirstName = *updateRequest.FirstName
	}
	if updateRequest.LastName != nil {
		user.LastName = *updateRequest.LastName
	}
	if updateRequest.Username != nil {
		user.Username = *updateRequest.Username
	}
	if updateRequest.Email != nil {
		user.Email = *updateRequest.Email
	}
	if updateRequest.MonthlyLimit != nil {
		user.MonthlyLimit = *updateRequest.MonthlyLimit
	}
	if updateRequest.SavingsGoal != nil {
		user.SavingsGoal = *updateRequest.SavingsGoal
	}

	// Обновляем время изменения
	user.UpdatedAt = time.Now()

	// Сохраняем изменения в базе данных
	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("ошибка при обновлении пользователя: %w", err)
	}

	userResponse := user.ToUserResponse()
	return &userResponse, nil
}

// DeleteUser удаляет пользователя
func (s *UserServiceImpl) DeleteUser(ctx context.Context, id int64) error {
	return s.userRepo.Delete(ctx, id)
}

// UploadAvatar загружает и устанавливает аватар пользователя
func (s *UserServiceImpl) UploadAvatar(ctx context.Context, userID int64, file *multipart.FileHeader) (*models.UserResponse, error) {
	// Проверяем существование пользователя
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Создаем директорию для хранения аватаров, если ее не существует
	avatarsDir := "uploads/avatars"
	if err := os.MkdirAll(avatarsDir, 0755); err != nil {
		return nil, fmt.Errorf("не удалось создать директорию для аватаров: %v", err)
	}

	// Генерируем уникальное имя файла с временной меткой для предотвращения кэширования браузером
	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("%d_%d_%s", userID, timestamp, filepath.Base(file.Filename))
	filePath := filepath.Join(avatarsDir, filename)

	// Открываем файл для записи
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать файл: %v", err)
	}
	defer dst.Close()

	// Открываем загруженный файл
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть загруженный файл: %v", err)
	}
	defer src.Close()

	// Копируем данные из загруженного файла в созданный файл
	if _, err = io.Copy(dst, src); err != nil {
		return nil, fmt.Errorf("не удалось скопировать файл: %v", err)
	}

	// Если у пользователя уже был аватар, удаляем его, если он не является аватаром по умолчанию
	if user.AvatarPath != "" && user.AvatarPath != "/default-avatar.png" {
		// Удаляем начальный слеш для получения относительного пути
		relPath := strings.TrimPrefix(user.AvatarPath, "/")

		// Проверяем существование файла и удаляем его
		if _, err := os.Stat(relPath); err == nil {
			if err := os.Remove(relPath); err != nil {
				// Логируем ошибку, но не прерываем процесс
				fmt.Printf("Не удалось удалить старый аватар: %v\n", err)
			}
		}
	}

	// Обновляем путь к аватару в базе данных
	// Для доступа через веб используем путь относительно URL
	avatarWebPath := fmt.Sprintf("/uploads/avatars/%s", filename)

	// Логируем путь аватара для отладки
	fmt.Printf("Новый путь аватара: %s\n", avatarWebPath)

	err = s.userRepo.UpdateAvatar(ctx, userID, avatarWebPath)
	if err != nil {
		return nil, err
	}

	// Получаем обновленные данные пользователя
	updatedUser, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	fmt.Printf("После обновления в базе: ID=%d, Username=%s, AvatarPath=%s\n",
		updatedUser.ID, updatedUser.Username, updatedUser.AvatarPath)

	// Преобразуем в UserResponse
	response := updatedUser.ToUserResponse()

	// Принудительно устанавливаем AvatarPath, если он пустой
	if response.AvatarURL == "" && updatedUser.AvatarPath != "" {
		response.AvatarURL = updatedUser.AvatarPath
		fmt.Printf("Принудительно установлен AvatarURL: %s\n", response.AvatarURL)
	}

	// Логируем значение avatar_url в ответе
	fmt.Printf("Значение avatar_url в ответе: %s\n", response.AvatarURL)

	return &response, nil
}

// RemoveAvatar удаляет аватар пользователя
func (s *UserServiceImpl) RemoveAvatar(ctx context.Context, userID int64) (*models.UserResponse, error) {
	// Проверяем существование пользователя
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Если у пользователя есть аватар, удаляем его
	if user.AvatarPath != "" && user.AvatarPath != "/default-avatar.png" {
		// Получаем полный путь к файлу
		// Удаляем начальный слеш для получения относительного пути
		relPath := strings.TrimPrefix(user.AvatarPath, "/")

		// Проверяем существование файла
		if _, err := os.Stat(relPath); err == nil {
			if err := os.Remove(relPath); err != nil {
				return nil, fmt.Errorf("не удалось удалить аватар: %v", err)
			}
		}
	}

	// Устанавливаем аватар по умолчанию
	err = s.userRepo.UpdateAvatar(ctx, userID, "/default-avatar.png")
	if err != nil {
		return nil, err
	}

	// Получаем обновленные данные пользователя
	updatedUser, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Преобразуем в UserResponse
	response := updatedUser.ToUserResponse()
	return &response, nil
}
