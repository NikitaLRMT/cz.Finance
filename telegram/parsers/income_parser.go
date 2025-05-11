package parsers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"cz.Finance/backend/models"
)

// ParseIncome парсит текстовое сообщение из Telegram и преобразует его в запрос на создание дохода
func ParseIncome(message string) (*models.CreateIncomeRequest, error) {
	// Формат сообщения: "Источник Сумма [Описание]"
	parts := strings.SplitN(message, " ", 3)

	if len(parts) < 2 {
		return nil, fmt.Errorf("неверный формат сообщения. Используйте: 'Источник Сумма [Описание]'")
	}

	// Парсим источник
	source := strings.ToLower(parts[0])
	validSource := false
	switch source {
	case "зарплата":
		source = string(models.SourceSalary)
		validSource = true
	case "фриланс":
		source = string(models.SourceFreelance)
		validSource = true
	case "инвестиции":
		source = string(models.SourceInvestment)
		validSource = true
	case "подарок":
		source = string(models.SourceGift)
		validSource = true
	case "аренда":
		source = string(models.SourceRental)
		validSource = true
	case "другое":
		source = string(models.SourceOther)
		validSource = true
	}

	if !validSource {
		return nil, fmt.Errorf("неверный источник. Доступные источники: зарплата, фриланс, инвестиции, подарок, аренда, другое")
	}

	// Парсим сумму
	amountStr := strings.Replace(parts[1], ",", ".", -1)
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return nil, fmt.Errorf("неверный формат суммы: %s", parts[1])
	}

	if amount <= 0 {
		return nil, fmt.Errorf("сумма должна быть положительным числом")
	}

	// Парсим описание (если есть)
	description := ""
	if len(parts) > 2 {
		description = parts[2]
	}

	// Создаем запрос
	request := &models.CreateIncomeRequest{
		Amount:      amount,
		Source:      models.IncomeSource(source),
		Date:        time.Now(),
		Description: description,
	}

	return request, nil
}
