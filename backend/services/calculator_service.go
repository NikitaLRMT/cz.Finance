package services

import (
	"math"
)

// CalculatorService интерфейс для калькуляторов финансовых расчетов
type CalculatorService interface {
	CalculateCompoundInterest(principal float64, rate float64, time float64, frequency int) map[string]interface{}
	CalculateMortgage(principal float64, rate float64, years int) map[string]interface{}
}

// CalculatorServiceImpl представляет реализацию сервиса калькуляторов
type CalculatorServiceImpl struct{}

// NewCalculatorService создает новый экземпляр сервиса калькуляторов
func NewCalculatorService() CalculatorService {
	return &CalculatorServiceImpl{}
}

// CalculateCompoundInterest вычисляет сложный процент
func (s *CalculatorServiceImpl) CalculateCompoundInterest(principal float64, rate float64, time float64, frequency int) map[string]interface{} {
	// Переводим процентную ставку из процентов в доли
	rate = rate / 100

	// Вычисляем количество периодов
	n := frequency * int(math.Floor(time))
	extraMonths := int((time - math.Floor(time)) * float64(frequency))

	// Вычисляем процентную ставку за период
	r := rate / float64(frequency)

	// Вычисляем итоговую сумму
	amount := principal * math.Pow(1+r, float64(n))

	// Если есть дополнительные месяцы, добавляем их
	if extraMonths > 0 {
		amount = amount * math.Pow(1+r, float64(extraMonths))
	}

	// Вычисляем заработанные проценты
	interest := amount - principal

	// Подготавливаем детальную информацию по годам
	yearlyDetails := make([]map[string]interface{}, int(math.Ceil(time)))
	currentAmount := principal

	for year := 0; year < int(math.Ceil(time)); year++ {
		startAmount := currentAmount
		periodsInYear := frequency

		// Для последнего неполного года
		if year == int(math.Floor(time)) && extraMonths > 0 {
			periodsInYear = extraMonths
		}

		// Расчет для каждого периода в году
		for period := 0; period < periodsInYear; period++ {
			currentAmount = currentAmount * (1 + r)
		}

		yearlyInterest := currentAmount - startAmount

		yearlyDetails[year] = map[string]interface{}{
			"year":            year + 1,
			"start_amount":    startAmount,
			"end_amount":      currentAmount,
			"yearly_interest": yearlyInterest,
			"total_interest":  currentAmount - principal,
		}
	}

	return map[string]interface{}{
		"principal":      principal,
		"rate":           rate * 100,
		"time":           time,
		"frequency":      frequency,
		"final_amount":   amount,
		"total_interest": interest,
		"yearly_details": yearlyDetails,
	}
}

// CalculateMortgage вычисляет параметры ипотечного кредита
func (s *CalculatorServiceImpl) CalculateMortgage(principal float64, rate float64, years int) map[string]interface{} {
	// Переводим процентную ставку из процентов в доли
	rate = rate / 100

	// Количество месяцев
	months := years * 12

	// Месячная процентная ставка
	monthlyRate := rate / 12

	// Расчет ежемесячного платежа
	monthlyPayment := principal * (monthlyRate * math.Pow(1+monthlyRate, float64(months))) / (math.Pow(1+monthlyRate, float64(months)) - 1)

	// Общая выплаченная сумма
	totalPayment := monthlyPayment * float64(months)

	// Общая сумма процентов
	totalInterest := totalPayment - principal

	// Подготавливаем детальную информацию по месяцам
	amortizationSchedule := make([]map[string]interface{}, months)
	remainingPrincipal := principal

	for month := 0; month < months; month++ {
		// Расчет процентов за месяц
		interestPayment := remainingPrincipal * monthlyRate

		// Расчет платежа в счет основного долга
		principalPayment := monthlyPayment - interestPayment

		// Обновление оставшегося основного долга
		remainingPrincipal -= principalPayment

		amortizationSchedule[month] = map[string]interface{}{
			"month":               month + 1,
			"payment":             monthlyPayment,
			"principal_payment":   principalPayment,
			"interest_payment":    interestPayment,
			"remaining_principal": remainingPrincipal,
		}
	}

	// Группировка по годам для более удобного отображения
	yearlyDetails := make([]map[string]interface{}, years)
	for year := 0; year < years; year++ {
		yearStart := year * 12
		yearEnd := (year+1)*12 - 1

		yearlyPrincipal := 0.0
		yearlyInterest := 0.0

		for month := yearStart; month <= yearEnd; month++ {
			yearlyPrincipal += amortizationSchedule[month]["principal_payment"].(float64)
			yearlyInterest += amortizationSchedule[month]["interest_payment"].(float64)
		}

		remainingPrincipalAtYearEnd := amortizationSchedule[yearEnd]["remaining_principal"].(float64)

		yearlyDetails[year] = map[string]interface{}{
			"year":                      year + 1,
			"yearly_principal_payment":  yearlyPrincipal,
			"yearly_interest_payment":   yearlyInterest,
			"yearly_total_payment":      yearlyPrincipal + yearlyInterest,
			"remaining_principal":       remainingPrincipalAtYearEnd,
			"paid_principal_percentage": (principal - remainingPrincipalAtYearEnd) / principal * 100,
		}
	}

	return map[string]interface{}{
		"principal":             principal,
		"rate":                  rate * 100,
		"years":                 years,
		"months":                months,
		"monthly_payment":       monthlyPayment,
		"total_payment":         totalPayment,
		"total_interest":        totalInterest,
		"yearly_details":        yearlyDetails,
		"amortization_schedule": amortizationSchedule,
	}
}
