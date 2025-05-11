package handlers

import (
	"net/http"

	"cz.Finance/backend/models"
	"cz.Finance/backend/services"
	"cz.Finance/backend/utils"
)

// UserHandlerImpl представляет реализацию обработчика пользователя
type UserHandlerImpl struct {
	userService services.UserService
}

// NewUserHandler создает новый экземпляр обработчика пользователя
func NewUserHandler(userService services.UserService) UserHandler {
	return &UserHandlerImpl{
		userService: userService,
	}
}

// SignUp обрабатывает запрос на регистрацию нового пользователя
func (h *UserHandlerImpl) SignUp(w http.ResponseWriter, r *http.Request) {
	// Декодируем запрос
	var signup models.UserSignup
	if err := utils.ParseJSON(r, &signup); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Ошибка при разборе запроса", err.Error())
		return
	}

	// Регистрируем пользователя
	tokenResponse, err := h.userService.SignUp(r.Context(), &signup)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Ошибка при регистрации", err.Error())
		return
	}

	// Отправляем ответ
	utils.RespondWithJSON(w, http.StatusCreated, tokenResponse)
}

// Login обрабатывает запрос на вход пользователя
func (h *UserHandlerImpl) Login(w http.ResponseWriter, r *http.Request) {
	// Декодируем запрос
	var login models.UserLogin
	if err := utils.ParseJSON(r, &login); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Ошибка при разборе запроса", err.Error())
		return
	}

	// Аутентифицируем пользователя
	tokenResponse, err := h.userService.Login(r.Context(), &login)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Ошибка аутентификации", err.Error())
		return
	}

	// Отправляем ответ
	utils.RespondWithJSON(w, http.StatusOK, tokenResponse)
}

// GetUser обрабатывает запрос на получение информации о пользователе
func (h *UserHandlerImpl) GetUser(w http.ResponseWriter, r *http.Request) {
	// Получаем ID пользователя из контекста
	userID, err := utils.GetUserIDFromContext(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Требуется авторизация", err.Error())
		return
	}

	// Получаем информацию о пользователе
	user, err := h.userService.GetUser(r.Context(), userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Пользователь не найден", err.Error())
		return
	}

	// Отправляем ответ
	utils.RespondWithJSON(w, http.StatusOK, user)
}

// UpdateUser обрабатывает запрос на обновление информации о пользователе
func (h *UserHandlerImpl) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Получаем ID пользователя из контекста
	userID, err := utils.GetUserIDFromContext(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Требуется авторизация", err.Error())
		return
	}

	// Декодируем запрос
	var updateRequest models.UpdateUserRequest
	if err := utils.ParseJSON(r, &updateRequest); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Ошибка при разборе запроса", err.Error())
		return
	}

	// Обновляем информацию о пользователе
	user, err := h.userService.UpdateUser(r.Context(), userID, &updateRequest)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Ошибка при обновлении пользователя", err.Error())
		return
	}

	// Отправляем ответ
	utils.RespondWithJSON(w, http.StatusOK, user)
}

// DeleteUser обрабатывает запрос на удаление пользователя
func (h *UserHandlerImpl) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Получаем ID пользователя из контекста
	userID, err := utils.GetUserIDFromContext(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Требуется авторизация", err.Error())
		return
	}

	// Удаляем пользователя
	err = h.userService.DeleteUser(r.Context(), userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Ошибка при удалении пользователя", err.Error())
		return
	}

	// Отправляем ответ
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Пользователь успешно удален"})
}

// UploadAvatar обрабатывает запрос на загрузку аватара
func (h *UserHandlerImpl) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	// Получаем ID пользователя из контекста
	userID, err := utils.GetUserIDFromContext(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Требуется авторизация", err.Error())
		return
	}

	// Проверяем размер файла
	r.ParseMultipartForm(10 << 20) // Ограничение 10 МБ
	file, handler, err := r.FormFile("avatar")
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Ошибка загрузки файла", err.Error())
		return
	}
	defer file.Close()

	// Проверяем тип файла
	contentType := handler.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/gif" {
		utils.RespondWithError(w, http.StatusBadRequest, "Неподдерживаемый тип файла", "Поддерживаются только JPEG, PNG и GIF")
		return
	}

	// Загружаем аватар
	userResponse, err := h.userService.UploadAvatar(r.Context(), userID, handler)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Не удалось загрузить аватар", err.Error())
		return
	}

	// Отправляем ответ
	utils.RespondWithJSON(w, http.StatusOK, userResponse)
}

// RemoveAvatar обрабатывает запрос на удаление аватара
func (h *UserHandlerImpl) RemoveAvatar(w http.ResponseWriter, r *http.Request) {
	// Получаем ID пользователя из контекста
	userID, err := utils.GetUserIDFromContext(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Требуется авторизация", err.Error())
		return
	}

	// Удаляем аватар
	userResponse, err := h.userService.RemoveAvatar(r.Context(), userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Не удалось удалить аватар", err.Error())
		return
	}

	// Отправляем ответ
	utils.RespondWithJSON(w, http.StatusOK, userResponse)
}
