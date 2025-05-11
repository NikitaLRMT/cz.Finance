package services

import (
	"context"
	"errors"
	"time"

	"cz.Finance/backend/models"
	"cz.Finance/backend/repositories"
)

// ExpenseServiceImpl представляет реализацию сервиса трат
type ExpenseServiceImpl struct {
	expenseRepo repositories.ExpenseRepository
	userRepo    repositories.UserRepository
}

// NewExpenseService создает новый экземпляр сервиса трат
func NewExpenseService(expenseRepo repositories.ExpenseRepository, userRepo repositories.UserRepository) ExpenseService {
	return &ExpenseServiceImpl{
		expenseRepo: expenseRepo,
		userRepo:    userRepo,
	}
}

// CreateExpense создает новую трату
func (s *ExpenseServiceImpl) CreateExpense(ctx context.Context, userID int64, request *models.CreateExpenseRequest) (*models.Expense, error) {
	// Проверяем существование пользователя
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("пользователь не найден")
	}

	// Создаем новую трату
	expense := &models.Expense{
		UserID:      userID,
		Title:       request.Title,
		Amount:      request.Amount,
		Category:    request.Category,
		Date:        request.Date,
		Description: request.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Если дата не указана, используем текущую
	if expense.Date.IsZero() {
		expense.Date = time.Now()
	}

	// Сохраняем трату в базе данных
	expenseID, err := s.expenseRepo.Create(ctx, expense)
	if err != nil {
		return nil, errors.New("ошибка при создании траты")
	}

	// Устанавливаем ID траты
	expense.ID = expenseID

	return expense, nil
}

// GetExpense получает трату по ID
func (s *ExpenseServiceImpl) GetExpense(ctx context.Context, id int64, userID int64) (*models.Expense, error) {
	expense, err := s.expenseRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("трата не найдена")
	}

	// Проверяем, что трата принадлежит пользователю
	if expense.UserID != userID {
		return nil, errors.New("у вас нет прав на просмотр этой траты")
	}

	return expense, nil
}

// GetUserExpenses получает список трат пользователя с пагинацией
func (s *ExpenseServiceImpl) GetUserExpenses(ctx context.Context, userID int64, limit, offset int) ([]models.Expense, error) {
	// Ограничиваем лимит выдачи
	if limit <= 0 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}

	return s.expenseRepo.GetByUserID(ctx, userID, limit, offset)
}

// GetUserExpensesByPeriod получает список трат пользователя за определенный период
func (s *ExpenseServiceImpl) GetUserExpensesByPeriod(ctx context.Context, userID int64, startDate, endDate time.Time) ([]models.Expense, error) {
	return s.expenseRepo.GetByUserIDAndPeriod(ctx, userID, startDate, endDate)
}

// GetExpenseSummary получает сводку по тратам пользователя за период
func (s *ExpenseServiceImpl) GetExpenseSummary(ctx context.Context, userID int64, startDate, endDate time.Time) (*models.ExpenseSummary, error) {
	// Получаем информацию о пользователе для получения месячного лимита
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("пользователь не найден")
	}

	// Получаем общую сумму трат за период
	totalAmount, err := s.expenseRepo.GetTotalAmountByUserIDAndPeriod(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, errors.New("ошибка при получении общей суммы трат")
	}

	// Получаем сводку по категориям
	categorySummary, err := s.expenseRepo.GetCategorySummaryByUserIDAndPeriod(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, errors.New("ошибка при получении сводки по категориям")
	}

	// Получаем последние 5 трат
	recentExpenses, err := s.expenseRepo.GetByUserIDAndPeriod(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, errors.New("ошибка при получении последних трат")
	}

	// Ограничиваем список последних трат до 5 элементов
	if len(recentExpenses) > 5 {
		recentExpenses = recentExpenses[:5]
	}

	// Вычисляем процент от месячного лимита
	limitPercentage := 0.0
	if user.MonthlyLimit > 0 {
		limitPercentage = (totalAmount / user.MonthlyLimit) * 100
	}

	return &models.ExpenseSummary{
		TotalAmount:     totalAmount,
		MonthlyLimit:    user.MonthlyLimit,
		LimitPercentage: limitPercentage,
		CategorySummary: categorySummary,
		RecentExpenses:  recentExpenses,
	}, nil
}

// UpdateExpense обновляет информацию о трате
func (s *ExpenseServiceImpl) UpdateExpense(ctx context.Context, id int64, userID int64, request *models.UpdateExpenseRequest) (*models.Expense, error) {
	// Получаем текущую трату
	expense, err := s.expenseRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("трата не найдена")
	}

	// Проверяем, что трата принадлежит пользователю
	if expense.UserID != userID {
		return nil, errors.New("у вас нет прав на изменение этой траты")
	}

	// Обновляем поля, если они указаны в запросе
	if request.Title != nil {
		expense.Title = *request.Title
	}
	if request.Amount != nil {
		expense.Amount = *request.Amount
	}
	if request.Category != nil {
		expense.Category = *request.Category
	}
	if request.Date != nil {
		expense.Date = *request.Date
	}
	if request.Description != nil {
		expense.Description = *request.Description
	}

	// Обновляем время изменения
	expense.UpdatedAt = time.Now()

	// Сохраняем изменения в базе данных
	err = s.expenseRepo.Update(ctx, expense)
	if err != nil {
		return nil, errors.New("ошибка при обновлении траты")
	}

	return expense, nil
}

// DeleteExpense удаляет трату
func (s *ExpenseServiceImpl) DeleteExpense(ctx context.Context, id int64, userID int64) error {
	// Проверяем существование траты
	expense, err := s.expenseRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("трата не найдена")
	}

	// Проверяем, что трата принадлежит пользователю
	if expense.UserID != userID {
		return errors.New("у вас нет прав на удаление этой траты")
	}

	return s.expenseRepo.Delete(ctx, id, userID)
}
