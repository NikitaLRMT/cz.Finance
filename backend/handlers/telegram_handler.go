package handlers

import (
	"net/http"
	"strconv"

	"cz.Finance/backend/models"
	"cz.Finance/backend/services"
	"cz.Finance/backend/utils"
	"github.com/gorilla/mux"
)

// TelegramHandlerImpl представляет реализацию обработчика телеграм
type TelegramHandlerImpl struct {
	telegramService services.TelegramService
	userService     services.UserService
}

// NewTelegramHandler создает новый экземпляр обработчика телеграм
func NewTelegramHandler(telegramService services.TelegramService, userService services.UserService) TelegramHandler {
	return &TelegramHandlerImpl{
		telegramService: telegramService,
		userService:     userService,
	}
}

// LinkTelegramAccount обрабатывает запрос на связывание аккаунта Telegram с аккаунтом пользователя
func (h *TelegramHandlerImpl) LinkTelegramAccount(w http.ResponseWriter, r *http.Request) {
	// Декодируем запрос
	var request struct {
		TelegramID int64  `json:"telegram_id"`
		Username   string `json:"username"`
		FirstName  string `json:"first_name"`
		LastName   string `json:"last_name"`
		Email      string `json:"email"`
		Password   string `json:"password"`
	}

	if err := utils.ParseJSON(r, &request); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Ошибка при разборе запроса", err.Error())
		return
	}

	// Создаем запрос на связывание
	linkRequest := &models.TelegramLinkRequest{
		Email:    request.Email,
		Password: request.Password,
	}

	// Связываем аккаунты
	_, err := h.telegramService.LinkAccount(r.Context(), request.TelegramID, request.Username, request.FirstName, request.LastName, linkRequest)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Ошибка при связывании аккаунтов", err.Error())
		return
	}

	// Создаем запрос на логин для получения токена
	loginRequest := &models.UserLogin{
		Email:    request.Email,
		Password: request.Password,
	}

	// Выполняем вход для получения токена
	tokenResponse, err := h.userService.Login(r.Context(), loginRequest)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Аккаунты связаны, но ошибка при генерации токена", err.Error())
		return
	}

	// Отправляем ответ с токеном
	utils.RespondWithJSON(w, http.StatusOK, tokenResponse)
}

// GetUserByTelegramID обрабатывает запрос на получение пользователя по Telegram ID
func (h *TelegramHandlerImpl) GetUserByTelegramID(w http.ResponseWriter, r *http.Request) {
	// Получаем Telegram ID из URL
	vars := mux.Vars(r)
	telegramIDStr := vars["telegram_id"]

	// Преобразуем строку в int64
	telegramID, err := strconv.ParseInt(telegramIDStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Неверный формат Telegram ID", err.Error())
		return
	}

	// Получаем пользователя
	user, err := h.telegramService.GetUserByTelegramID(r.Context(), telegramID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Пользователь не найден", err.Error())
		return
	}

	// Отправляем ответ
	utils.RespondWithJSON(w, http.StatusOK, user)
}

// UnlinkTelegramAccount обрабатывает запрос на отвязку аккаунта Telegram от аккаунта пользователя
func (h *TelegramHandlerImpl) UnlinkTelegramAccount(w http.ResponseWriter, r *http.Request) {
	// Получаем Telegram ID из URL
	vars := mux.Vars(r)
	telegramIDStr := vars["telegram_id"]

	// Преобразуем строку в int64
	telegramID, err := strconv.ParseInt(telegramIDStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Неверный формат Telegram ID", err.Error())
		return
	}

	// Отвязываем аккаунт
	err = h.telegramService.UnlinkAccount(r.Context(), telegramID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Ошибка при отвязке аккаунта", err.Error())
		return
	}

	// Отправляем ответ
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Аккаунт успешно отвязан"})
}
