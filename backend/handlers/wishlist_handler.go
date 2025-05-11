package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"cz.Finance/backend/models"
	"cz.Finance/backend/services"
	"cz.Finance/backend/utils"

	"github.com/gorilla/mux"
)

// WishlistHandlerImpl представляет реализацию обработчика списка желаний
type WishlistHandlerImpl struct {
	wishlistService services.WishlistService
}

// NewWishlistHandler создает новый экземпляр обработчика списка желаний
func NewWishlistHandler(wishlistService services.WishlistService) WishlistHandler {
	return &WishlistHandlerImpl{
		wishlistService: wishlistService,
	}
}

// CreateWishlistItem обрабатывает запрос на создание нового элемента списка желаний
func (h *WishlistHandlerImpl) CreateWishlistItem(w http.ResponseWriter, r *http.Request) {
	// Получаем ID пользователя из контекста
	userID, err := utils.GetUserIDFromContext(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Требуется авторизация", err.Error())
		return
	}

	// Декодируем тело запроса
	var request models.CreateWishlistItemRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Некорректный запрос", err.Error())
		return
	}

	// Валидируем запрос
	if err := utils.ValidateStruct(request); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Ошибка валидации", err.Error())
		return
	}

	// Создаем элемент списка желаний
	item, err := h.wishlistService.CreateWishlistItem(r.Context(), userID, &request)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Не удалось создать элемент списка желаний", err.Error())
		return
	}

	// Отправляем ответ
	utils.RespondWithJSON(w, http.StatusCreated, item)
}

// GetWishlistItem обрабатывает запрос на получение элемента списка желаний по ID
func (h *WishlistHandlerImpl) GetWishlistItem(w http.ResponseWriter, r *http.Request) {
	// Получаем ID пользователя из контекста
	userID, err := utils.GetUserIDFromContext(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Требуется авторизация", err.Error())
		return
	}

	// Получаем ID элемента из URL
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Некорректный ID", err.Error())
		return
	}

	// Получаем элемент списка желаний
	item, err := h.wishlistService.GetWishlistItem(r.Context(), id, userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Элемент списка желаний не найден", err.Error())
		return
	}

	// Отправляем ответ
	utils.RespondWithJSON(w, http.StatusOK, item)
}

// GetUserWishlist обрабатывает запрос на получение списка желаний пользователя
func (h *WishlistHandlerImpl) GetUserWishlist(w http.ResponseWriter, r *http.Request) {
	// Получаем ID пользователя из контекста
	userID, err := utils.GetUserIDFromContext(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Требуется авторизация", err.Error())
		return
	}

	// Получаем список желаний пользователя
	items, err := h.wishlistService.GetUserWishlist(r.Context(), userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Не удалось получить список желаний", err.Error())
		return
	}

	// Отправляем ответ
	utils.RespondWithJSON(w, http.StatusOK, items)
}

// UpdateWishlistItem обрабатывает запрос на обновление элемента списка желаний
func (h *WishlistHandlerImpl) UpdateWishlistItem(w http.ResponseWriter, r *http.Request) {
	// Получаем ID пользователя из контекста
	userID, err := utils.GetUserIDFromContext(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Требуется авторизация", err.Error())
		return
	}

	// Получаем ID элемента из URL
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Некорректный ID", err.Error())
		return
	}

	// Декодируем тело запроса
	var request models.UpdateWishlistItemRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Некорректный запрос", err.Error())
		return
	}

	// Валидируем запрос
	if err := utils.ValidateStruct(request); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Ошибка валидации", err.Error())
		return
	}

	// Обновляем элемент списка желаний
	updatedItem, err := h.wishlistService.UpdateWishlistItem(r.Context(), id, userID, &request)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Не удалось обновить элемент списка желаний", err.Error())
		return
	}

	// Отправляем ответ
	utils.RespondWithJSON(w, http.StatusOK, updatedItem)
}

// DeleteWishlistItem обрабатывает запрос на удаление элемента списка желаний
func (h *WishlistHandlerImpl) DeleteWishlistItem(w http.ResponseWriter, r *http.Request) {
	// Получаем ID пользователя из контекста
	userID, err := utils.GetUserIDFromContext(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Требуется авторизация", err.Error())
		return
	}

	// Получаем ID элемента из URL
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Некорректный ID", err.Error())
		return
	}

	// Удаляем элемент списка желаний
	err = h.wishlistService.DeleteWishlistItem(r.Context(), id, userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Не удалось удалить элемент списка желаний", err.Error())
		return
	}

	// Отправляем ответ
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Элемент успешно удален"})
}
