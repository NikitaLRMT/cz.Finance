package handlers

import (
	"net/http"
	"time"

	"cz.Finance/backend/models"
	"cz.Finance/backend/services"
	"cz.Finance/backend/utils"
)

// ExpenseHandlerImpl представляет реализацию обработчика трат
type ExpenseHandlerImpl struct {
	expenseService services.ExpenseService
}

// NewExpenseHandler создает новый экземпляр обработчика трат
func NewExpenseHandler(expenseService services.ExpenseService) ExpenseHandler {
	return &ExpenseHandlerImpl{
		expenseService: expenseService,
	}
}

// CreateExpense обрабатывает запрос на создание новой траты
func (h *ExpenseHandlerImpl) CreateExpense(w http.ResponseWriter, r *http.Request) {
	// Получаем ID пользователя из контекста
	userID, err := utils.GetUserIDFromContext(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Требуется авторизация", err.Error())
		return
	}

	// Декодируем запрос
	var request models.CreateExpenseRequest
	if err := utils.ParseJSON(r, &request); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Ошибка при разборе запроса", err.Error())
		return
	}

	// Создаем трату
	expense, err := h.expenseService.CreateExpense(r.Context(), userID, &request)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Ошибка при создании траты", err.Error())
		return
	}

	// Отправляем ответ
	utils.RespondWithJSON(w, http.StatusCreated, expense)
}

// GetExpense обрабатывает запрос на получение траты по ID
func (h *ExpenseHandlerImpl) GetExpense(w http.ResponseWriter, r *http.Request) {
	// Получаем ID пользователя из контекста
	userID, err := utils.GetUserIDFromContext(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Требуется авторизация", err.Error())
		return
	}

	// Получаем ID траты из URL
	expenseID, err := utils.GetIDParam(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Неверный ID траты", err.Error())
		return
	}

	// Получаем трату
	expense, err := h.expenseService.GetExpense(r.Context(), expenseID, userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Трата не найдена", err.Error())
		return
	}

	// Отправляем ответ
	utils.RespondWithJSON(w, http.StatusOK, expense)
}

// GetUserExpenses обрабатывает запрос на получение списка трат пользователя
func (h *ExpenseHandlerImpl) GetUserExpenses(w http.ResponseWriter, r *http.Request) {
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

	// Если указаны даты начала и конца периода, получаем траты за период
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

		expenses, err := h.expenseService.GetUserExpensesByPeriod(r.Context(), userID, startDate, endDate)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Ошибка при получении трат", err.Error())
			return
		}

		// Отправляем ответ
		utils.RespondWithJSON(w, http.StatusOK, expenses)
		return
	}

	// Иначе получаем траты с пагинацией
	expenses, err := h.expenseService.GetUserExpenses(r.Context(), userID, limit, offset)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Ошибка при получении трат", err.Error())
		return
	}

	// Отправляем ответ
	utils.RespondWithJSON(w, http.StatusOK, expenses)
}

// GetExpenseSummary обрабатывает запрос на получение сводки по тратам
func (h *ExpenseHandlerImpl) GetExpenseSummary(w http.ResponseWriter, r *http.Request) {
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

	// Получаем сводку по тратам
	summary, err := h.expenseService.GetExpenseSummary(r.Context(), userID, startDate, endDate)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Ошибка при получении сводки по тратам", err.Error())
		return
	}

	// Отправляем ответ
	utils.RespondWithJSON(w, http.StatusOK, summary)
}

// UpdateExpense обрабатывает запрос на обновление траты
func (h *ExpenseHandlerImpl) UpdateExpense(w http.ResponseWriter, r *http.Request) {
	// Получаем ID пользователя из контекста
	userID, err := utils.GetUserIDFromContext(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Требуется авторизация", err.Error())
		return
	}

	// Получаем ID траты из URL
	expenseID, err := utils.GetIDParam(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Неверный ID траты", err.Error())
		return
	}

	// Декодируем запрос
	var request models.UpdateExpenseRequest
	if err := utils.ParseJSON(r, &request); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Ошибка при разборе запроса", err.Error())
		return
	}

	// Обновляем трату
	expense, err := h.expenseService.UpdateExpense(r.Context(), expenseID, userID, &request)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Ошибка при обновлении траты", err.Error())
		return
	}

	// Отправляем ответ
	utils.RespondWithJSON(w, http.StatusOK, expense)
}

// DeleteExpense обрабатывает запрос на удаление траты
func (h *ExpenseHandlerImpl) DeleteExpense(w http.ResponseWriter, r *http.Request) {
	// Получаем ID пользователя из контекста
	userID, err := utils.GetUserIDFromContext(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Требуется авторизация", err.Error())
		return
	}

	// Получаем ID траты из URL
	expenseID, err := utils.GetIDParam(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Неверный ID траты", err.Error())
		return
	}

	// Удаляем трату
	err = h.expenseService.DeleteExpense(r.Context(), expenseID, userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Ошибка при удалении траты", err.Error())
		return
	}

	// Отправляем ответ
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Трата успешно удалена"})
}
