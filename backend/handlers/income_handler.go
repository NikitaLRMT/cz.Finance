package handlers

import (
	"net/http"
	"time"

	"cz.Finance/backend/models"
	"cz.Finance/backend/services"
	"cz.Finance/backend/utils"
)

// IncomeHandlerImpl представляет реализацию обработчика накоплений
type IncomeHandlerImpl struct {
	incomeService services.IncomeService
}

// NewIncomeHandler создает новый экземпляр обработчика накоплений
func NewIncomeHandler(incomeService services.IncomeService) IncomeHandler {
	return &IncomeHandlerImpl{
		incomeService: incomeService,
	}
}

// CreateIncome обрабатывает запрос на создание нового накопления
func (h *IncomeHandlerImpl) CreateIncome(w http.ResponseWriter, r *http.Request) {
	// Получаем ID пользователя из контекста
	userID, err := utils.GetUserIDFromContext(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Требуется авторизация", err.Error())
		return
	}

	// Декодируем запрос
	var request models.CreateIncomeRequest
	if err := utils.ParseJSON(r, &request); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Ошибка при разборе запроса", err.Error())
		return
	}

	// Создаем накопление
	income, err := h.incomeService.CreateIncome(r.Context(), userID, &request)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Ошибка при создании накопления", err.Error())
		return
	}

	// Отправляем ответ
	utils.RespondWithJSON(w, http.StatusCreated, income)
}

// GetIncome обрабатывает запрос на получение накопления по ID
func (h *IncomeHandlerImpl) GetIncome(w http.ResponseWriter, r *http.Request) {
	// Получаем ID пользователя из контекста
	userID, err := utils.GetUserIDFromContext(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Требуется авторизация", err.Error())
		return
	}

	// Получаем ID накопления из URL
	incomeID, err := utils.GetIDParam(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Неверный ID накопления", err.Error())
		return
	}

	// Получаем накопление
	income, err := h.incomeService.GetIncome(r.Context(), incomeID, userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Накопление не найдено", err.Error())
		return
	}

	// Отправляем ответ
	utils.RespondWithJSON(w, http.StatusOK, income)
}

// GetUserIncomes обрабатывает запрос на получение списка накоплений пользователя
func (h *IncomeHandlerImpl) GetUserIncomes(w http.ResponseWriter, r *http.Request) {
	// Получаем ID пользователя из контекста
	userID, err := utils.GetUserIDFromContext(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Требуется авторизация", err.Error())
		return
	}

	// Получаем параметры пагинации из запроса
	limit := utils.GetIntQueryParam(r, "limit", 10)
	offset := utils.GetIntQueryParam(r, "offset", 0)

	// Проверяем наличие параметров периода
	startDateStr := utils.GetQueryParam(r, "start_date")
	endDateStr := utils.GetQueryParam(r, "end_date")

	// Если указаны даты начала и конца периода, получаем накопления за период
	if startDateStr != "" && endDateStr != "" {
		startDate, err := time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Неверный формат даты начала периода", err.Error())
			return
		}

		endDate, err := time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Неверный формат даты конца периода", err.Error())
			return
		}

		incomes, err := h.incomeService.GetUserIncomesByPeriod(r.Context(), userID, startDate, endDate)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Ошибка при получении накоплений", err.Error())
			return
		}

		// Отправляем ответ
		utils.RespondWithJSON(w, http.StatusOK, incomes)
		return
	}

	// Иначе получаем накопления с пагинацией
	incomes, err := h.incomeService.GetUserIncomes(r.Context(), userID, limit, offset)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Ошибка при получении накоплений", err.Error())
		return
	}

	// Отправляем ответ
	utils.RespondWithJSON(w, http.StatusOK, incomes)
}

// GetIncomeSummary обрабатывает запрос на получение сводки по накоплениям
func (h *IncomeHandlerImpl) GetIncomeSummary(w http.ResponseWriter, r *http.Request) {
	// Получаем ID пользователя из контекста
	userID, err := utils.GetUserIDFromContext(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Требуется авторизация", err.Error())
		return
	}

	// Получаем даты периода из запроса
	startDateStr := utils.GetQueryParam(r, "start_date")
	endDateStr := utils.GetQueryParam(r, "end_date")

	// Если даты не указаны, используем текущий месяц
	var startDate, endDate time.Time
	if startDateStr == "" || endDateStr == "" {
		now := time.Now()
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		endDate = startDate.AddDate(0, 1, -1).Add(time.Hour * 23).Add(time.Minute * 59).Add(time.Second * 59)
	} else {
		var err error
		startDate, err = time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Неверный формат даты начала периода", err.Error())
			return
		}

		endDate, err = time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Неверный формат даты конца периода", err.Error())
			return
		}
	}

	// Получаем сводку по накоплениям
	summary, err := h.incomeService.GetIncomeSummary(r.Context(), userID, startDate, endDate)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Ошибка при получении сводки по накоплениям", err.Error())
		return
	}

	// Отправляем ответ
	utils.RespondWithJSON(w, http.StatusOK, summary)
}

// UpdateIncome обрабатывает запрос на обновление накопления
func (h *IncomeHandlerImpl) UpdateIncome(w http.ResponseWriter, r *http.Request) {
	// Получаем ID пользователя из контекста
	userID, err := utils.GetUserIDFromContext(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Требуется авторизация", err.Error())
		return
	}

	// Получаем ID накопления из URL
	incomeID, err := utils.GetIDParam(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Неверный ID накопления", err.Error())
		return
	}

	// Декодируем запрос
	var request models.UpdateIncomeRequest
	if err := utils.ParseJSON(r, &request); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Ошибка при разборе запроса", err.Error())
		return
	}

	// Обновляем накопление
	income, err := h.incomeService.UpdateIncome(r.Context(), incomeID, userID, &request)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Ошибка при обновлении накопления", err.Error())
		return
	}

	// Отправляем ответ
	utils.RespondWithJSON(w, http.StatusOK, income)
}

// DeleteIncome обрабатывает запрос на удаление накопления
func (h *IncomeHandlerImpl) DeleteIncome(w http.ResponseWriter, r *http.Request) {
	// Получаем ID пользователя из контекста
	userID, err := utils.GetUserIDFromContext(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Требуется авторизация", err.Error())
		return
	}

	// Получаем ID накопления из URL
	incomeID, err := utils.GetIDParam(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Неверный ID накопления", err.Error())
		return
	}

	// Удаляем накопление
	err = h.incomeService.DeleteIncome(r.Context(), incomeID, userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Ошибка при удалении накопления", err.Error())
		return
	}

	// Отправляем ответ
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Накопление успешно удалено"})
}
