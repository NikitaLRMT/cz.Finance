package handlers

import (
	"net/http"

	"cz.Finance/backend/services"
	"cz.Finance/backend/utils"
)

// CalculatorHandlerImpl представляет реализацию обработчика калькуляторов
type CalculatorHandlerImpl struct {
	calculatorService services.CalculatorService
}

// NewCalculatorHandler создает новый экземпляр обработчика калькуляторов
func NewCalculatorHandler() CalculatorHandler {
	// Создаем сервис калькуляторов
	calculatorService := services.NewCalculatorService()

	return &CalculatorHandlerImpl{
		calculatorService: calculatorService,
	}
}

// CompoundInterestCalculator обрабатывает запрос на расчет сложного процента
func (h *CalculatorHandlerImpl) CompoundInterestCalculator(w http.ResponseWriter, r *http.Request) {
	// Декодируем запрос
	var request struct {
		Principal float64 `json:"principal"`
		Rate      float64 `json:"rate"`
		Time      float64 `json:"time"`
		Frequency int     `json:"frequency"`
	}

	if err := utils.ParseJSON(r, &request); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Ошибка при разборе запроса", err.Error())
		return
	}

	// Проверка входных данных
	if request.Principal <= 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "Неверные входные данные", "Начальная сумма должна быть положительным числом")
		return
	}

	if request.Rate < 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "Неверные входные данные", "Процентная ставка не может быть отрицательной")
		return
	}

	if request.Time <= 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "Неверные входные данные", "Срок вклада должен быть положительным числом")
		return
	}

	if request.Frequency <= 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "Неверные входные данные", "Частота начисления должна быть положительным числом")
		return
	}

	// Вычисляем сложный процент
	result := h.calculatorService.CalculateCompoundInterest(
		request.Principal,
		request.Rate,
		request.Time,
		request.Frequency,
	)

	// Отправляем ответ
	utils.RespondWithJSON(w, http.StatusOK, result)
}

// MortgageCalculator обрабатывает запрос на расчет ипотеки
func (h *CalculatorHandlerImpl) MortgageCalculator(w http.ResponseWriter, r *http.Request) {
	// Декодируем запрос
	var request struct {
		Principal float64 `json:"principal"`
		Rate      float64 `json:"rate"`
		Years     int     `json:"years"`
	}

	if err := utils.ParseJSON(r, &request); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Ошибка при разборе запроса", err.Error())
		return
	}

	// Проверка входных данных
	if request.Principal <= 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "Неверные входные данные", "Сумма кредита должна быть положительным числом")
		return
	}

	if request.Rate < 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "Неверные входные данные", "Процентная ставка не может быть отрицательной")
		return
	}

	if request.Years <= 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "Неверные входные данные", "Срок кредита должен быть положительным числом")
		return
	}

	// Вычисляем параметры ипотеки
	result := h.calculatorService.CalculateMortgage(
		request.Principal,
		request.Rate,
		request.Years,
	)

	// Отправляем ответ
	utils.RespondWithJSON(w, http.StatusOK, result)
}
