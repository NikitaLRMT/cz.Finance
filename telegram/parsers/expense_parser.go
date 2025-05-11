package parsers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"cz.Finance/backend/models"
)

// ParseExpense парсит текстовое сообщение из Telegram и преобразует его в запрос на создание траты
func ParseExpense(message string) (*models.CreateExpenseRequest, error) {
	// Формат сообщения: "Категория Наименование Сумма [Описание]"
	parts := strings.SplitN(message, " ", 4)

	if len(parts) < 3 {
		return nil, fmt.Errorf("неверный формат сообщения. Используйте: 'Категория Наименование Сумма [Описание]'")
	}

	// Парсим категорию
	category := strings.ToLower(parts[0])
	validCategory := false
	switch category {
	case "продукты":
		category = string(models.CategoryFood)
		validCategory = true
	case "транспорт":
		category = string(models.CategoryTransport)
		validCategory = true
	case "жильё":
		category = string(models.CategoryHousing)
		validCategory = true
	case "коммунальные":
		category = string(models.CategoryUtilities)
		validCategory = true
	case "покупки":
		category = string(models.CategoryShopping)
		validCategory = true
	case "развлечения":
		category = string(models.CategoryEntertainment)
		validCategory = true
	case "здоровье":
		category = string(models.CategoryHealthcare)
		validCategory = true
	case "образование":
		category = string(models.CategoryEducation)
		validCategory = true
	case "путешествия":
		category = string(models.CategoryTravel)
		validCategory = true
	case "другое":
		category = string(models.CategoryOther)
		validCategory = true
	}

	if !validCategory {
		return nil, fmt.Errorf("неверная категория. Доступные категории: продукты, транспорт, жильё, коммунальные, покупки, развлечения, здоровье, образование, путешествия, другое")
	}

	// Парсим наименование
	title := parts[1]
	if len(title) < 2 {
		return nil, fmt.Errorf("название должно содержать не менее 2 символов")
	}

	// Парсим сумму
	amountStr := strings.Replace(parts[2], ",", ".", -1)
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return nil, fmt.Errorf("неверный формат суммы: %s", parts[2])
	}

	if amount <= 0 {
		return nil, fmt.Errorf("сумма должна быть положительным числом")
	}

	// Парсим описание (если есть)
	description := ""
	if len(parts) > 3 {
		description = parts[3]
	}

	// Создаем запрос
	request := &models.CreateExpenseRequest{
		Title:       title,
		Amount:      amount,
		Category:    models.ExpenseCategory(category),
		Date:        time.Now(),
		Description: description,
	}

	return request, nil
}
