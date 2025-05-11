package handlers

import (
	"net/http"
	"strconv"

	"cz.Finance/backend/services"
	"cz.Finance/backend/utils"
	"github.com/gorilla/mux"
)

// DashboardHandlerImpl представляет реализацию обработчика панели статистики
type DashboardHandlerImpl struct {
	dashboardService services.DashboardService
}

// NewDashboardHandler создает новый экземпляр обработчика панели статистики
func NewDashboardHandler(dashboardService services.DashboardService) DashboardHandler {
	return &DashboardHandlerImpl{
		dashboardService: dashboardService,
	}
}

// GetDashboardSummary обрабатывает запрос на получение сводки панели статистики
func (h *DashboardHandlerImpl) GetDashboardSummary(w http.ResponseWriter, r *http.Request) {
	// Получаем ID пользователя из контекста
	userID, err := utils.GetUserIDFromContext(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Требуется авторизация", err.Error())
		return
	}

	// Получаем параметры запроса
	limit := utils.GetIntQueryParam(r, "limit", 5)

	// Получаем сводку
	summary, err := h.dashboardService.GetDashboardSummary(r.Context(), userID, limit)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Ошибка при получении сводки", err.Error())
		return
	}

	// Отправляем ответ
	utils.RespondWithJSON(w, http.StatusOK, summary)
}

// GetMonthlyStats обрабатывает запрос на получение статистики за месяц
func (h *DashboardHandlerImpl) GetMonthlyStats(w http.ResponseWriter, r *http.Request) {
	// Получаем ID пользователя из контекста
	userID, err := utils.GetUserIDFromContext(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Требуется авторизация", err.Error())
		return
	}

	// Получаем год и месяц из URL
	vars := mux.Vars(r)
	year, err := strconv.Atoi(vars["year"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Неверный формат года", err.Error())
		return
	}

	month, err := strconv.Atoi(vars["month"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Неверный формат месяца", err.Error())
		return
	}

	// Проверяем корректность месяца
	if month < 1 || month > 12 {
		utils.RespondWithError(w, http.StatusBadRequest, "Месяц должен быть в диапазоне от 1 до 12", "")
		return
	}

	// Получаем статистику за месяц
	stats, err := h.dashboardService.GetMonthlyStats(r.Context(), userID, year, month)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Ошибка при получении статистики за месяц", err.Error())
		return
	}

	// Отправляем ответ
	utils.RespondWithJSON(w, http.StatusOK, stats)
}

// GetYearlyStats обрабатывает запрос на получение статистики за год
func (h *DashboardHandlerImpl) GetYearlyStats(w http.ResponseWriter, r *http.Request) {
	// Получаем ID пользователя из контекста
	userID, err := utils.GetUserIDFromContext(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Требуется авторизация", err.Error())
		return
	}

	// Получаем год из URL
	vars := mux.Vars(r)
	year, err := strconv.Atoi(vars["year"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Неверный формат года", err.Error())
		return
	}

	// Получаем статистику за год
	stats, err := h.dashboardService.GetYearlyStats(r.Context(), userID, year)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Ошибка при получении статистики за год", err.Error())
		return
	}

	// Отправляем ответ
	utils.RespondWithJSON(w, http.StatusOK, stats)
}

// GetBudgetGoals обрабатывает запрос на получение бюджетных целей пользователя
func (h *DashboardHandlerImpl) GetBudgetGoals(w http.ResponseWriter, r *http.Request) {
	// Получаем ID пользователя из контекста
	userID, err := utils.GetUserIDFromContext(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Требуется авторизация", err.Error())
		return
	}

	// Получаем бюджетные цели
	goals, err := h.dashboardService.GetBudgetGoals(r.Context(), userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Ошибка при получении бюджетных целей", err.Error())
		return
	}

	// Отправляем ответ
	utils.RespondWithJSON(w, http.StatusOK, goals)
}

// SetBudgetGoal обрабатывает запрос на установку бюджетной цели
func (h *DashboardHandlerImpl) SetBudgetGoal(w http.ResponseWriter, r *http.Request) {
	// Получаем ID пользователя из контекста
	userID, err := utils.GetUserIDFromContext(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Требуется авторизация", err.Error())
		return
	}

	// Декодируем запрос
	var request struct {
		Category string  `json:"category"`
		Amount   float64 `json:"amount"`
	}

	if err := utils.ParseJSON(r, &request); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Ошибка при разборе запроса", err.Error())
		return
	}

	// Устанавливаем бюджетную цель
	err = h.dashboardService.SetBudgetGoal(r.Context(), userID, request.Category, request.Amount)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Ошибка при установке бюджетной цели", err.Error())
		return
	}

	// Отправляем ответ
	utils.RespondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"message":  "Бюджетная цель успешно установлена",
		"category": request.Category,
		"amount":   request.Amount,
	})
}
