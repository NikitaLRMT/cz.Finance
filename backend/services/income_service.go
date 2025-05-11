package services

import (
	"context"
	"errors"
	"time"

	"cz.Finance/backend/models"
	"cz.Finance/backend/repositories"
)

// IncomeServiceImpl представляет реализацию сервиса накоплений
type IncomeServiceImpl struct {
	incomeRepo repositories.IncomeRepository
	userRepo   repositories.UserRepository
}

// NewIncomeService создает новый экземпляр сервиса накоплений
func NewIncomeService(incomeRepo repositories.IncomeRepository, userRepo repositories.UserRepository) IncomeService {
	return &IncomeServiceImpl{
		incomeRepo: incomeRepo,
		userRepo:   userRepo,
	}
}

// CreateIncome создает новое накопление
func (s *IncomeServiceImpl) CreateIncome(ctx context.Context, userID int64, request *models.CreateIncomeRequest) (*models.Income, error) {
	// Проверяем существование пользователя
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("пользователь не найден")
	}

	// Создаем новое накопление
	income := &models.Income{
		UserID:      userID,
		Amount:      request.Amount,
		Source:      request.Source,
		Date:        request.Date,
		Description: request.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Если дата не указана, используем текущую
	if income.Date.IsZero() {
		income.Date = time.Now()
	}

	// Сохраняем накопление в базе данных
	incomeID, err := s.incomeRepo.Create(ctx, income)
	if err != nil {
		return nil, errors.New("ошибка при создании накопления")
	}

	// Устанавливаем ID накопления
	income.ID = incomeID

	return income, nil
}

// GetIncome получает накопление по ID
func (s *IncomeServiceImpl) GetIncome(ctx context.Context, id int64, userID int64) (*models.Income, error) {
	income, err := s.incomeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("накопление не найдено")
	}

	// Проверяем, что накопление принадлежит пользователю
	if income.UserID != userID {
		return nil, errors.New("у вас нет прав на просмотр этого накопления")
	}

	return income, nil
}

// GetUserIncomes получает список накоплений пользователя с пагинацией
func (s *IncomeServiceImpl) GetUserIncomes(ctx context.Context, userID int64, limit, offset int) ([]models.Income, error) {
	// Ограничиваем лимит выдачи
	if limit <= 0 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}

	return s.incomeRepo.GetByUserID(ctx, userID, limit, offset)
}

// GetUserIncomesByPeriod получает список накоплений пользователя за определенный период
func (s *IncomeServiceImpl) GetUserIncomesByPeriod(ctx context.Context, userID int64, startDate, endDate time.Time) ([]models.Income, error) {
	return s.incomeRepo.GetByUserIDAndPeriod(ctx, userID, startDate, endDate)
}

// GetIncomeSummary получает сводку по накоплениям пользователя за период
func (s *IncomeServiceImpl) GetIncomeSummary(ctx context.Context, userID int64, startDate, endDate time.Time) (*models.IncomeSummary, error) {
	// Получаем информацию о пользователе для получения цели накоплений
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("пользователь не найден")
	}

	// Получаем общую сумму накоплений за период
	totalAmount, err := s.incomeRepo.GetTotalAmountByUserIDAndPeriod(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, errors.New("ошибка при получении общей суммы накоплений")
	}

	// Получаем сводку по источникам
	sourceSummary, err := s.incomeRepo.GetSourceSummaryByUserIDAndPeriod(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, errors.New("ошибка при получении сводки по источникам")
	}

	// Получаем последние 5 накоплений
	recentIncomes, err := s.incomeRepo.GetByUserIDAndPeriod(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, errors.New("ошибка при получении последних накоплений")
	}

	// Ограничиваем список последних накоплений до 5 элементов
	if len(recentIncomes) > 5 {
		recentIncomes = recentIncomes[:5]
	}

	// Вычисляем процент от цели накоплений
	goalPercentage := 0.0
	if user.SavingsGoal > 0 {
		goalPercentage = (totalAmount / user.SavingsGoal) * 100
	}

	return &models.IncomeSummary{
		TotalAmount:    totalAmount,
		SavingsGoal:    user.SavingsGoal,
		GoalPercentage: goalPercentage,
		SourceSummary:  sourceSummary,
		RecentIncomes:  recentIncomes,
	}, nil
}

// UpdateIncome обновляет информацию о накоплении
func (s *IncomeServiceImpl) UpdateIncome(ctx context.Context, id int64, userID int64, request *models.UpdateIncomeRequest) (*models.Income, error) {
	// Получаем текущее накопление
	income, err := s.incomeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("накопление не найдено")
	}

	// Проверяем, что накопление принадлежит пользователю
	if income.UserID != userID {
		return nil, errors.New("у вас нет прав на изменение этого накопления")
	}

	// Обновляем поля, если они указаны в запросе
	if request.Amount != nil {
		income.Amount = *request.Amount
	}
	if request.Source != nil {
		income.Source = *request.Source
	}
	if request.Date != nil {
		income.Date = *request.Date
	}
	if request.Description != nil {
		income.Description = *request.Description
	}

	// Обновляем время изменения
	income.UpdatedAt = time.Now()

	// Сохраняем изменения в базе данных
	err = s.incomeRepo.Update(ctx, income)
	if err != nil {
		return nil, errors.New("ошибка при обновлении накопления")
	}

	return income, nil
}

// DeleteIncome удаляет накопление
func (s *IncomeServiceImpl) DeleteIncome(ctx context.Context, id int64, userID int64) error {
	// Проверяем существование накопления
	income, err := s.incomeRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("накопление не найдено")
	}

	// Проверяем, что накопление принадлежит пользователю
	if income.UserID != userID {
		return errors.New("у вас нет прав на удаление этого накопления")
	}

	return s.incomeRepo.Delete(ctx, id, userID)
}
