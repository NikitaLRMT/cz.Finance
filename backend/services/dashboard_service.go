package services

import (
	"context"
	"errors"
	"time"

	"cz.Finance/backend/repositories"
)

// DashboardServiceImpl представляет реализацию сервиса информационной панели
type DashboardServiceImpl struct {
	expenseRepo repositories.ExpenseRepository
	incomeRepo  repositories.IncomeRepository
	userRepo    repositories.UserRepository
}

// NewDashboardService создает новый экземпляр сервиса информационной панели
func NewDashboardService(
	expenseRepo repositories.ExpenseRepository,
	incomeRepo repositories.IncomeRepository,
	userRepo repositories.UserRepository,
) DashboardService {
	return &DashboardServiceImpl{
		expenseRepo: expenseRepo,
		incomeRepo:  incomeRepo,
		userRepo:    userRepo,
	}
}

// GetDashboardSummary получает сводку для панели мониторинга
func (s *DashboardServiceImpl) GetDashboardSummary(ctx context.Context, userID int64, limit int) (map[string]interface{}, error) {
	// Получаем пользователя
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("пользователь не найден")
	}

	// Определяем начало и конец текущего месяца
	now := time.Now()
	currentMonthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	currentMonthEnd := currentMonthStart.AddDate(0, 1, -1).Add(time.Hour * 23).Add(time.Minute * 59).Add(time.Second * 59)

	// Получаем общую сумму трат за текущий месяц
	currentMonthExpenses, err := s.expenseRepo.GetTotalAmountByUserIDAndPeriod(ctx, userID, currentMonthStart, currentMonthEnd)
	if err != nil {
		return nil, errors.New("ошибка при получении трат за текущий месяц")
	}

	// Получаем общую сумму накоплений за текущий месяц
	currentMonthIncomes, err := s.incomeRepo.GetTotalAmountByUserIDAndPeriod(ctx, userID, currentMonthStart, currentMonthEnd)
	if err != nil {
		return nil, errors.New("ошибка при получении накоплений за текущий месяц")
	}

	// Определяем начало и конец всего периода (используем очень раннюю дату и текущую дату)
	allTimeStart := time.Date(2000, 1, 1, 0, 0, 0, 0, now.Location())
	allTimeEnd := now

	// Получаем общую сумму трат за все время
	allTimeExpenses, err := s.expenseRepo.GetTotalAmountByUserIDAndPeriod(ctx, userID, allTimeStart, allTimeEnd)
	if err != nil {
		return nil, errors.New("ошибка при получении трат за все время")
	}

	// Получаем общую сумму накоплений за все время
	allTimeIncomes, err := s.incomeRepo.GetTotalAmountByUserIDAndPeriod(ctx, userID, allTimeStart, allTimeEnd)
	if err != nil {
		return nil, errors.New("ошибка при получении накоплений за все время")
	}

	// Получаем сводку трат по категориям за текущий месяц
	expensesByCategory, err := s.expenseRepo.GetCategorySummaryByUserIDAndPeriod(ctx, userID, currentMonthStart, currentMonthEnd)
	if err != nil {
		return nil, errors.New("ошибка при получении сводки трат по категориям")
	}

	// Получаем сводку накоплений по источникам за текущий месяц
	incomesBySource, err := s.incomeRepo.GetSourceSummaryByUserIDAndPeriod(ctx, userID, currentMonthStart, currentMonthEnd)
	if err != nil {
		return nil, errors.New("ошибка при получении сводки накоплений по источникам")
	}

	// Получаем последние 5 трат
	recentExpenses, err := s.expenseRepo.GetByUserID(ctx, userID, 5, 0)
	if err != nil {
		return nil, errors.New("ошибка при получении последних трат")
	}

	// Получаем последние 5 накоплений
	recentIncomes, err := s.incomeRepo.GetByUserID(ctx, userID, 5, 0)
	if err != nil {
		return nil, errors.New("ошибка при получении последних накоплений")
	}

	// Формируем ответ
	result := map[string]interface{}{
		"user": map[string]interface{}{
			"monthly_limit": user.MonthlyLimit,
			"savings_goal":  user.SavingsGoal,
		},
		"current_month": map[string]interface{}{
			"expenses":         currentMonthExpenses,
			"incomes":          currentMonthIncomes,
			"expenses_percent": calculatePercentage(currentMonthExpenses, user.MonthlyLimit),
			"savings_percent":  calculatePercentage(currentMonthIncomes, user.SavingsGoal),
			"balance":          currentMonthIncomes - currentMonthExpenses,
		},
		"all_time": map[string]interface{}{
			"expenses": allTimeExpenses,
			"incomes":  allTimeIncomes,
			"balance":  allTimeIncomes - allTimeExpenses,
		},
		"expenses_by_category": expensesByCategory,
		"incomes_by_source":    incomesBySource,
		"recent_expenses":      recentExpenses,
		"recent_incomes":       recentIncomes,
	}

	return result, nil
}

// GetMonthlyStats получает статистику за указанный месяц
func (s *DashboardServiceImpl) GetMonthlyStats(ctx context.Context, userID int64, year int, month int) (map[string]interface{}, error) {
	// Определяем начало и конец указанного месяца
	monthStart := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Now().Location())
	monthEnd := monthStart.AddDate(0, 1, -1).Add(time.Hour * 23).Add(time.Minute * 59).Add(time.Second * 59)

	// Получаем общую сумму трат за указанный месяц
	monthlyExpenses, err := s.expenseRepo.GetTotalAmountByUserIDAndPeriod(ctx, userID, monthStart, monthEnd)
	if err != nil {
		return nil, errors.New("ошибка при получении трат за указанный месяц")
	}

	// Получаем общую сумму накоплений за указанный месяц
	monthlyIncomes, err := s.incomeRepo.GetTotalAmountByUserIDAndPeriod(ctx, userID, monthStart, monthEnd)
	if err != nil {
		return nil, errors.New("ошибка при получении накоплений за указанный месяц")
	}

	// Получаем траты за указанный месяц
	expenses, err := s.expenseRepo.GetByUserIDAndPeriod(ctx, userID, monthStart, monthEnd)
	if err != nil {
		return nil, errors.New("ошибка при получении трат за указанный месяц")
	}

	// Получаем накопления за указанный месяц
	incomes, err := s.incomeRepo.GetByUserIDAndPeriod(ctx, userID, monthStart, monthEnd)
	if err != nil {
		return nil, errors.New("ошибка при получении накоплений за указанный месяц")
	}

	// Получаем сводку трат по категориям за указанный месяц
	expensesByCategory, err := s.expenseRepo.GetCategorySummaryByUserIDAndPeriod(ctx, userID, monthStart, monthEnd)
	if err != nil {
		return nil, errors.New("ошибка при получении сводки трат по категориям")
	}

	// Получаем сводку накоплений по источникам за указанный месяц
	incomesBySource, err := s.incomeRepo.GetSourceSummaryByUserIDAndPeriod(ctx, userID, monthStart, monthEnd)
	if err != nil {
		return nil, errors.New("ошибка при получении сводки накоплений по источникам")
	}

	// Получаем пользователя для получения лимитов
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("пользователь не найден")
	}

	// Формируем ответ
	result := map[string]interface{}{
		"period": map[string]interface{}{
			"year":       year,
			"month":      month,
			"start_date": monthStart,
			"end_date":   monthEnd,
		},
		"summary": map[string]interface{}{
			"expenses":         monthlyExpenses,
			"incomes":          monthlyIncomes,
			"balance":          monthlyIncomes - monthlyExpenses,
			"expenses_percent": calculatePercentage(monthlyExpenses, user.MonthlyLimit),
			"savings_percent":  calculatePercentage(monthlyIncomes, user.SavingsGoal),
		},
		"expenses_by_category": expensesByCategory,
		"incomes_by_source":    incomesBySource,
		"expenses":             expenses,
		"incomes":              incomes,
	}

	return result, nil
}

// GetYearlyStats получает статистику за указанный год
func (s *DashboardServiceImpl) GetYearlyStats(ctx context.Context, userID int64, year int) (map[string]interface{}, error) {
	// Определяем начало и конец указанного года
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.Now().Location())
	yearEnd := time.Date(year, 12, 31, 23, 59, 59, 0, time.Now().Location())

	// Получаем общую сумму трат за указанный год
	yearlyExpenses, err := s.expenseRepo.GetTotalAmountByUserIDAndPeriod(ctx, userID, yearStart, yearEnd)
	if err != nil {
		return nil, errors.New("ошибка при получении трат за указанный год")
	}

	// Получаем общую сумму накоплений за указанный год
	yearlyIncomes, err := s.incomeRepo.GetTotalAmountByUserIDAndPeriod(ctx, userID, yearStart, yearEnd)
	if err != nil {
		return nil, errors.New("ошибка при получении накоплений за указанный год")
	}

	// Получаем сводку трат по категориям за указанный год
	expensesByCategory, err := s.expenseRepo.GetCategorySummaryByUserIDAndPeriod(ctx, userID, yearStart, yearEnd)
	if err != nil {
		return nil, errors.New("ошибка при получении сводки трат по категориям")
	}

	// Получаем сводку накоплений по источникам за указанный год
	incomesBySource, err := s.incomeRepo.GetSourceSummaryByUserIDAndPeriod(ctx, userID, yearStart, yearEnd)
	if err != nil {
		return nil, errors.New("ошибка при получении сводки накоплений по источникам")
	}

	// Получаем данные по месяцам
	monthlyData := make([]map[string]interface{}, 12)
	for month := 1; month <= 12; month++ {
		monthStart := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Now().Location())
		monthEnd := monthStart.AddDate(0, 1, -1).Add(time.Hour * 23).Add(time.Minute * 59).Add(time.Second * 59)

		// Расходы за месяц
		monthExpenses, err := s.expenseRepo.GetTotalAmountByUserIDAndPeriod(ctx, userID, monthStart, monthEnd)
		if err != nil {
			monthExpenses = 0
		}

		// Доходы за месяц
		monthIncomes, err := s.incomeRepo.GetTotalAmountByUserIDAndPeriod(ctx, userID, monthStart, monthEnd)
		if err != nil {
			monthIncomes = 0
		}

		monthlyData[month-1] = map[string]interface{}{
			"month":     month,
			"expenses":  monthExpenses,
			"incomes":   monthIncomes,
			"balance":   monthIncomes - monthExpenses,
			"month_str": monthStart.Month().String(),
		}
	}

	// Формируем ответ
	result := map[string]interface{}{
		"period": map[string]interface{}{
			"year":       year,
			"start_date": yearStart,
			"end_date":   yearEnd,
		},
		"summary": map[string]interface{}{
			"expenses": yearlyExpenses,
			"incomes":  yearlyIncomes,
			"balance":  yearlyIncomes - yearlyExpenses,
			"avg_month": map[string]float64{
				"expenses": yearlyExpenses / 12,
				"incomes":  yearlyIncomes / 12,
			},
		},
		"expenses_by_category": expensesByCategory,
		"incomes_by_source":    incomesBySource,
		"monthly_data":         monthlyData,
	}

	return result, nil
}

// calculatePercentage вычисляет процент от значения
func calculatePercentage(value, total float64) float64 {
	if total == 0 {
		return 0
	}
	return (value / total) * 100
}

// GetBudgetGoals получает бюджетные цели пользователя
func (s *DashboardServiceImpl) GetBudgetGoals(ctx context.Context, userID int64) (map[string]interface{}, error) {
	// Получаем пользователя для проверки
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("пользователь не найден")
	}

	// Получаем текущий месяц для расчёта прогресса
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Now().Location())
	monthEnd := monthStart.AddDate(0, 1, -1).Add(time.Hour * 23).Add(time.Minute * 59).Add(time.Second * 59)

	// Получаем сводку трат по категориям за текущий месяц
	expensesByCategory, err := s.expenseRepo.GetCategorySummaryByUserIDAndPeriod(ctx, userID, monthStart, monthEnd)
	if err != nil {
		return nil, errors.New("ошибка при получении сводки трат по категориям")
	}

	// Получаем бюджетные цели из репозитория (заглушка, так как у нас пока нет отдельной таблицы)
	// В реальном приложении здесь будет запрос к специальному репозиторию бюджетных целей
	budgetGoals := make(map[string]interface{})

	// Для демонстрации создаем несколько стандартных категорий бюджетов
	// Используем параметры пользователя из БД
	defaultCategories := []string{
		"Продукты", "Транспорт", "Жильё", "Коммунальные", "Покупки",
		"Развлечения", "Здоровье", "Образование", "Путешествия", "Другое",
	}

	// Заполняем бюджетные цели - в реальном приложении данные брались бы из БД
	for _, category := range defaultCategories {
		amount := float64(0)

		// Устанавливаем дефолтные значения
		switch category {
		case "Продукты":
			amount = 15000
		case "Транспорт":
			amount = 5000
		case "Жильё":
			amount = 20000
		case "Коммунальные":
			amount = 8000
		case "Развлечения":
			amount = 10000
		default:
			amount = 5000
		}

		// Получаем фактические расходы по категории
		spent := float64(0)
		if val, exists := expensesByCategory[category]; exists {
			// Простое присваивание, так как знаем что категории содержат float64
			spent = val
		}

		// Формируем данные о цели
		budgetGoals[category] = map[string]interface{}{
			"amount": amount,
			"spent":  spent,
		}
	}

	return budgetGoals, nil
}

// SetBudgetGoal устанавливает бюджетную цель для категории
func (s *DashboardServiceImpl) SetBudgetGoal(ctx context.Context, userID int64, category string, amount float64) error {
	// Получаем пользователя для проверки
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return errors.New("пользователь не найден")
	}

	// Проверяем корректность суммы
	if amount < 0 {
		return errors.New("сумма бюджета не может быть отрицательной")
	}

	return nil
}
