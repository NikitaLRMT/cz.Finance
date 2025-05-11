package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"cz.Finance/backend/models"
	"cz.Finance/telegram/client"
	"cz.Finance/telegram/parsers"
	"gopkg.in/telebot.v3"
)

// BotHandlers структура для обработчиков бота
type BotHandlers struct {
	apiClient *client.APIClient
}

// NewBotHandlers создает новые обработчики для бота
func NewBotHandlers(apiClient *client.APIClient) *BotHandlers {
	return &BotHandlers{
		apiClient: apiClient,
	}
}

// RegisterHandlers регистрирует обработчики команд
func (h *BotHandlers) RegisterHandlers(bot *telebot.Bot) {
	// Обработчик команды /start
	bot.Handle("/start", h.HandleStart)

	// Обработчик команды /link
	bot.Handle("/link", h.HandleLink)

	// Обработчик команды /unlink
	bot.Handle("/unlink", h.HandleUnlink)

	// Обработчик команды /expense
	bot.Handle("/expense", h.HandleExpenseCommand)

	// Обработчик команды /income
	bot.Handle("/income", h.HandleIncomeCommand)

	// Обработчик команды /balance
	bot.Handle("/balance", h.HandleBalance)

	// Обработчик команды /transactions
	bot.Handle("/transactions", h.HandleTransactions)

	// Обработчик команды /category
	bot.Handle("/category", h.HandleCategory)

	// Обработчики для бюджетных целей
	bot.Handle("/budget", h.HandleBudgetGoals)
	bot.Handle("/setbudget", h.HandleSetBudgetGoal)

	// Обработчик для добавления траты
	bot.Handle(telebot.OnText, h.HandleMessage)
}

// HandleStart обрабатывает команду /start
func (h *BotHandlers) HandleStart(c telebot.Context) error {
	welcomeMessage := `
Привет! Я бот для учета финансов.

Доступные команды:
/link - Связать ваш аккаунт Telegram с аккаунтом в приложении
/unlink - Отвязать ваш аккаунт Telegram от аккаунта в приложении
/expense - Добавить трату
/income - Добавить поступление
/balance - Посмотреть баланс
/transactions - Последние транзакции
/category - Расходы по категории
/budget - Посмотреть бюджетные цели
/setbudget - Установить бюджетную цель

Чтобы связать аккаунт, используйте команду /link и введите ваш email и пароль в формате:
/link email@example.com password
`
	return c.Send(welcomeMessage)
}

// HandleLink обрабатывает команду /link для связывания аккаунтов
func (h *BotHandlers) HandleLink(c telebot.Context) error {
	args := c.Args()
	if len(args) != 2 {
		return c.Send("Неверный формат команды. Используйте: /link email пароль")
	}

	email := args[0]
	password := args[1]

	// Создаем запрос на связывание
	request := &models.TelegramLinkRequest{
		Email:    email,
		Password: password,
	}

	// Получаем информацию о пользователе Telegram
	telegramID := c.Sender().ID
	username := c.Sender().Username
	firstName := c.Sender().FirstName
	lastName := c.Sender().LastName

	// Связываем аккаунты через API
	user, err := h.apiClient.LinkTelegramAccount(telegramID, username, firstName, lastName, request)
	if err != nil {
		return c.Send(fmt.Sprintf("Ошибка при связывании аккаунтов: %s", err.Error()))
	}

	return c.Send(fmt.Sprintf("Аккаунт успешно связан с пользователем %s", user.Username))
}

// HandleUnlink обрабатывает команду /unlink для отвязки аккаунта
func (h *BotHandlers) HandleUnlink(c telebot.Context) error {
	telegramID := c.Sender().ID

	// Проверяем, связан ли аккаунт
	_, err := h.apiClient.GetUserByTelegramID(telegramID)
	if err != nil {
		return c.Send("Ваш аккаунт Telegram не связан с аккаунтом в приложении.")
	}

	// Отвязываем аккаунт
	err = h.apiClient.UnlinkAccount(telegramID)
	if err != nil {
		return c.Send(fmt.Sprintf("Ошибка при отвязке аккаунта: %s", err.Error()))
	}

	return c.Send("Ваш аккаунт Telegram успешно отвязан от аккаунта в приложении.")
}

// HandleExpenseCommand обрабатывает команду /expense
func (h *BotHandlers) HandleExpenseCommand(c telebot.Context) error {
	// Проверяем, связан ли аккаунт
	_, err := h.apiClient.GetUserByTelegramID(c.Sender().ID)
	if err != nil {
		return c.Send("Вы не связали аккаунт. Используйте команду /link")
	}

	expenseTemplate := `
Чтобы добавить трату, отправьте сообщение в формате:
Категория Наименование Сумма [Описание]

Например:
Продукты Пятерочка 1300 Еженедельная закупка

Доступные категории:
- Продукты
- Транспорт
- Жильё
- Коммунальные
- Покупки
- Развлечения
- Здоровье
- Образование
- Путешествия
- Другое
`
	return c.Send(expenseTemplate)
}

// HandleIncomeCommand обрабатывает команду /income
func (h *BotHandlers) HandleIncomeCommand(c telebot.Context) error {
	// Проверяем, связан ли аккаунт
	_, err := h.apiClient.GetUserByTelegramID(c.Sender().ID)
	if err != nil {
		return c.Send("Вы не связали аккаунт. Используйте команду /link")
	}

	incomeTemplate := `
Чтобы добавить поступление, отправьте сообщение в формате:
Источник Сумма [Описание]

Например:
Зарплата 50000 Аванс

Доступные источники:
- Зарплата
- Фриланс
- Инвестиции
- Подарок
- Аренда
- Другое
`
	return c.Send(incomeTemplate)
}

// HandleBalance обрабатывает команду /balance
func (h *BotHandlers) HandleBalance(c telebot.Context) error {
	telegramID := c.Sender().ID

	// Проверяем, связан ли аккаунт
	user, err := h.apiClient.GetUserByTelegramID(telegramID)
	if err != nil {
		return c.Send("Вы не связали аккаунт. Используйте команду /link")
	}

	// Получаем текущий месяц и год
	now := time.Now()
	year, month, _ := now.Date()

	// Получаем статистику за текущий месяц через API
	stats, err := h.apiClient.GetMonthlyStats(user.ID, year, int(month), telegramID)
	if err != nil {
		return c.Send(fmt.Sprintf("Ошибка при получении баланса: %s", err.Error()))
	}

	// Форматируем сообщение
	summary := stats["summary"].(map[string]interface{})
	balance := summary["balance"].(float64)
	expenses := summary["expenses"].(float64)
	incomes := summary["incomes"].(float64)

	balanceMessage := fmt.Sprintf(`
Ваш баланс:

🗓 Период: %s %d

💰 Поступления: %.2f руб.
💸 Расходы: %.2f руб.
📊 Баланс: %.2f руб.
`, now.Month().String(), year, incomes, expenses, balance)

	return c.Send(balanceMessage)
}

// HandleTransactions обрабатывает команду /transactions
func (h *BotHandlers) HandleTransactions(c telebot.Context) error {
	telegramID := c.Sender().ID

	// Проверяем, связан ли аккаунт
	user, err := h.apiClient.GetUserByTelegramID(telegramID)
	if err != nil {
		return c.Send("Вы не связали аккаунт. Используйте команду /link")
	}

	// Получаем последние 5 транзакций через API
	transactions, err := h.apiClient.GetRecentTransactions(5, telegramID)
	if err != nil {
		return c.Send(fmt.Sprintf("Ошибка при получении транзакций: %s", err.Error()))
	}

	// Добавляем отладочный вывод содержимого ответа
	fmt.Printf("Получены транзакции: %+v\n", transactions)

	// Форматируем сообщение
	message := fmt.Sprintf("Последние транзакции для %s:\n\n", user.Username)

	// Обработка расходов
	hasTransactions := false

	// Проверяем наличие расходов в разных возможных форматах
	if recentExpenses, ok := transactions["recent_expenses"]; ok {
		// Пытаемся получить расходы из поля recent_expenses
		hasTransactions = handleExpenses(recentExpenses, &message)
	} else if expenses, ok := transactions["expenses"]; ok {
		// Пытаемся получить расходы из поля expenses
		hasTransactions = handleExpenses(expenses, &message)
	}

	// Обработка доходов
	if recentIncomes, ok := transactions["recent_incomes"]; ok {
		// Пытаемся получить доходы из поля recent_incomes
		hasTransactions = handleIncomes(recentIncomes, &message) || hasTransactions
	} else if incomes, ok := transactions["incomes"]; ok {
		// Пытаемся получить доходы из поля incomes
		hasTransactions = handleIncomes(incomes, &message) || hasTransactions
	}

	if !hasTransactions {
		message += "Транзакций не найдено."
	}

	return c.Send(message)
}

// handleExpenses обрабатывает список расходов для вывода в сообщении
func handleExpenses(expenses interface{}, message *string) bool {
	fmt.Printf("Обработка расходов: %+v\n", expenses)

	// Проверяем, есть ли расходы
	switch expList := expenses.(type) {
	case []interface{}:
		if len(expList) == 0 {
			return false
		}
		*message += "💸 Расходы:\n"
		for i, exp := range expList {
			if i >= 5 {
				break
			}

			// Пытаемся извлечь данные
			var date, title, category string
			var amount float64

			switch expense := exp.(type) {
			case map[string]interface{}:
				// Извлекаем данные из map
				if dateVal, ok := expense["date"]; ok {
					date = fmt.Sprintf("%v", dateVal)
				} else if dateVal, ok := expense["created_at"]; ok {
					date = fmt.Sprintf("%v", dateVal)
				}

				if titleVal, ok := expense["title"]; ok {
					title = fmt.Sprintf("%v", titleVal)
				}

				if categoryVal, ok := expense["category"]; ok {
					category = fmt.Sprintf("%v", categoryVal)
				}

				if amountVal, ok := expense["amount"]; ok {
					switch a := amountVal.(type) {
					case float64:
						amount = a
					case int:
						amount = float64(a)
					case string:
						amount, _ = strconv.ParseFloat(a, 64)
					}
				}
			}

			// Формируем строку
			if date != "" {
				if len(date) > 10 {
					date = date[:10]
				}
				*message += fmt.Sprintf("- %s | %s | %s | %.2f руб.\n", date, category, title, amount)
			} else {
				*message += fmt.Sprintf("- %s | %s | %.2f руб.\n", category, title, amount)
			}
		}
		*message += "\n"
		return true

	default:
		fmt.Printf("Неизвестный формат расходов: %T\n", expenses)
		return false
	}
}

// handleIncomes обрабатывает список доходов для вывода в сообщении
func handleIncomes(incomes interface{}, message *string) bool {
	fmt.Printf("Обработка доходов: %+v\n", incomes)

	// Проверяем, есть ли доходы
	switch incList := incomes.(type) {
	case []interface{}:
		if len(incList) == 0 {
			return false
		}
		*message += "💰 Поступления:\n"
		for i, inc := range incList {
			if i >= 5 {
				break
			}

			// Пытаемся извлечь данные
			var date, source string
			var amount float64

			switch income := inc.(type) {
			case map[string]interface{}:
				// Извлекаем данные из map
				if dateVal, ok := income["date"]; ok {
					date = fmt.Sprintf("%v", dateVal)
				} else if dateVal, ok := income["created_at"]; ok {
					date = fmt.Sprintf("%v", dateVal)
				}

				if sourceVal, ok := income["source"]; ok {
					source = fmt.Sprintf("%v", sourceVal)
				}

				if amountVal, ok := income["amount"]; ok {
					switch a := amountVal.(type) {
					case float64:
						amount = a
					case int:
						amount = float64(a)
					case string:
						amount, _ = strconv.ParseFloat(a, 64)
					}
				}
			}

			// Формируем строку
			if date != "" {
				if len(date) > 10 {
					date = date[:10]
				}
				*message += fmt.Sprintf("- %s | %s | %.2f руб.\n", date, source, amount)
			} else {
				*message += fmt.Sprintf("- %s | %.2f руб.\n", source, amount)
			}
		}
		*message += "\n"
		return true

	default:
		fmt.Printf("Неизвестный формат доходов: %T\n", incomes)
		return false
	}
}

// HandleCategory обрабатывает команду /category для фильтрации расходов по категориям
func (h *BotHandlers) HandleCategory(c telebot.Context) error {
	telegramID := c.Sender().ID

	// Проверяем, связан ли аккаунт
	_, err := h.apiClient.GetUserByTelegramID(telegramID)
	if err != nil {
		return c.Send("Вы не связали аккаунт. Используйте команду /link")
	}

	args := c.Args()
	if len(args) == 0 {
		// Если категория не указана, выводим список доступных категорий
		categories := []string{
			"Продукты",
			"Транспорт",
			"Жильё",
			"Коммунальные",
			"Покупки",
			"Развлечения",
			"Здоровье",
			"Образование",
			"Путешествия",
			"Другое",
		}

		message := "Укажите категорию для фильтрации расходов.\nНапример: /category Продукты\n\nДоступные категории:\n"
		for _, category := range categories {
			message += fmt.Sprintf("- %s\n", category)
		}

		return c.Send(message)
	}

	// Получаем указанную категорию
	category := args[0]

	// Получаем расходы по указанной категории
	expenses, err := h.apiClient.GetExpensesByCategory(category, telegramID)
	if err != nil {
		return c.Send(fmt.Sprintf("Ошибка при получении расходов: %s", err.Error()))
	}

	if len(expenses) == 0 {
		return c.Send(fmt.Sprintf("Расходы по категории '%s' не найдены.", category))
	}

	// Форматируем сообщение
	message := fmt.Sprintf("Расходы по категории '%s':\n\n", category)

	// Ограничиваем количество выводимых расходов до 10
	limit := 10
	if len(expenses) < limit {
		limit = len(expenses)
	}

	// Сортируем расходы по дате (от новых к старым)
	// Здесь предполагается, что модель Expense имеет поле CreatedAt или Date

	totalAmount := 0.0
	for i := 0; i < limit; i++ {
		date := expenses[i].CreatedAt.Format("2006-01-02")
		message += fmt.Sprintf("- %s | %s | %.2f руб.\n",
			date, expenses[i].Title, expenses[i].Amount)
		totalAmount += expenses[i].Amount
	}

	// Добавляем итоговую сумму
	message += fmt.Sprintf("\nИтого: %.2f руб.", totalAmount)

	// Если расходов больше, чем выведено
	if len(expenses) > limit {
		message += fmt.Sprintf("\n\nПоказано %d из %d расходов.", limit, len(expenses))
	}

	return c.Send(message)
}

// HandleBudgetGoals обрабатывает команду /budget
func (h *BotHandlers) HandleBudgetGoals(c telebot.Context) error {
	telegramID := c.Sender().ID

	// Проверяем, связан ли аккаунт
	_, err := h.apiClient.GetUserByTelegramID(telegramID)
	if err != nil {
		return c.Send("Вы не связали аккаунт. Используйте команду /link")
	}

	// Получаем бюджетные цели
	goals, err := h.apiClient.GetBudgetGoals(telegramID)
	if err != nil {
		return c.Send(fmt.Sprintf("Ошибка при получении бюджетных целей: %s", err.Error()))
	}

	// Проверяем наличие целей
	if goals == nil || len(goals) == 0 {
		return c.Send("У вас нет установленных бюджетных целей. Используйте /setbudget для установки.")
	}

	// Форматируем сообщение
	message := "Ваши бюджетные цели:\n\n"

	// Обрабатываем цели
	for category, data := range goals {
		goalData := data.(map[string]interface{})
		amount := goalData["amount"].(float64)
		spent := goalData["spent"].(float64)
		remaining := amount - spent
		percentage := (spent / amount) * 100

		// Добавляем эмодзи-индикатор
		var emoji string
		if percentage >= 90 {
			emoji = "🔴" // красный - превышение бюджета
		} else if percentage >= 70 {
			emoji = "🟠" // оранжевый - приближается к лимиту
		} else {
			emoji = "🟢" // зеленый - в пределах бюджета
		}

		message += fmt.Sprintf("%s %s:\n", emoji, category)
		message += fmt.Sprintf("   Бюджет: %.2f руб.\n", amount)
		message += fmt.Sprintf("   Потрачено: %.2f руб. (%.1f%%)\n", spent, percentage)
		message += fmt.Sprintf("   Осталось: %.2f руб.\n\n", remaining)
	}

	return c.Send(message)
}

// HandleSetBudgetGoal обрабатывает команду /setbudget
func (h *BotHandlers) HandleSetBudgetGoal(c telebot.Context) error {
	telegramID := c.Sender().ID

	// Проверяем, связан ли аккаунт
	_, err := h.apiClient.GetUserByTelegramID(telegramID)
	if err != nil {
		return c.Send("Вы не связали аккаунт. Используйте команду /link")
	}

	args := c.Args()
	if len(args) != 2 {
		// Если аргументы не указаны, выводим подсказку
		categories := []string{
			"Продукты",
			"Транспорт",
			"Жильё",
			"Коммунальные",
			"Покупки",
			"Развлечения",
			"Здоровье",
			"Образование",
			"Путешествия",
			"Другое",
		}

		message := "Укажите категорию и сумму для установки бюджетной цели.\nНапример: /setbudget Продукты 10000\n\nДоступные категории:\n"
		for _, category := range categories {
			message += fmt.Sprintf("- %s\n", category)
		}

		return c.Send(message)
	}

	// Получаем категорию и сумму
	category := args[0]
	amountStr := args[1]

	// Парсим сумму
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return c.Send("Неверный формат суммы. Введите число, например: 10000")
	}

	// Устанавливаем бюджетную цель
	err = h.apiClient.SetBudgetGoal(category, amount, telegramID)
	if err != nil {
		return c.Send(fmt.Sprintf("Ошибка при установке бюджетной цели: %s", err.Error()))
	}

	return c.Send(fmt.Sprintf("Бюджетная цель для категории '%s' установлена: %.2f руб.", category, amount))
}

// HandleMessage обрабатывает текстовые сообщения
func (h *BotHandlers) HandleMessage(c telebot.Context) error {
	// Пропускаем команды
	if c.Message().Text[0] == '/' {
		return nil
	}

	telegramID := c.Sender().ID

	// Получаем пользователя через API
	user, err := h.apiClient.GetUserByTelegramID(telegramID)
	if err != nil {
		return c.Send("Вы не связали аккаунт. Используйте команду /link")
	}

	// В зависимости от контекста парсим сообщение как трату или поступление
	// Здесь мы используем простую эвристику: если первое слово - известная категория траты,
	// то это трата, иначе - поступление
	firstWord := getFirstWord(c.Message().Text)

	// Проверяем, является ли первое слово категорией траты
	switch strings.ToLower(firstWord) {
	case "продукты", "транспорт", "жильё", "коммунальные", "покупки", "развлечения", "здоровье", "образование", "путешествия", "другое":
		// Парсим сообщение как трату
		return h.handleExpense(c, user)
	case "зарплата", "фриланс", "инвестиции", "подарок", "аренда":
		// Парсим сообщение как поступление
		return h.handleIncome(c, user)
	default:
		// Если не удалось определить, спрашиваем пользователя
		return c.Send("Не удалось определить тип операции. Пожалуйста, используйте команды /expense или /income для добавления трат или поступлений.")
	}
}

// handleExpense обрабатывает добавление траты
func (h *BotHandlers) handleExpense(c telebot.Context, user *models.User) error {
	telegramID := c.Sender().ID

	// Парсим сообщение
	expenseRequest, err := parsers.ParseExpense(c.Message().Text)
	if err != nil {
		return c.Send(fmt.Sprintf("Ошибка при парсинге траты: %s", err.Error()))
	}

	// Добавляем трату через API
	expense, err := h.apiClient.CreateExpense(user.ID, expenseRequest, telegramID)
	if err != nil {
		return c.Send(fmt.Sprintf("Ошибка при добавлении траты: %s", err.Error()))
	}

	return c.Send(fmt.Sprintf("Трата успешно добавлена:\n- Категория: %s\n- Наименование: %s\n- Сумма: %.2f руб.", expense.Category, expense.Title, expense.Amount))
}

// handleIncome обрабатывает добавление поступления
func (h *BotHandlers) handleIncome(c telebot.Context, user *models.User) error {
	telegramID := c.Sender().ID

	// Парсим сообщение
	incomeRequest, err := parsers.ParseIncome(c.Message().Text)
	if err != nil {
		return c.Send(fmt.Sprintf("Ошибка при парсинге поступления: %s", err.Error()))
	}

	// Добавляем поступление через API
	income, err := h.apiClient.CreateIncome(user.ID, incomeRequest, telegramID)
	if err != nil {
		return c.Send(fmt.Sprintf("Ошибка при добавлении поступления: %s", err.Error()))
	}

	return c.Send(fmt.Sprintf("Поступление успешно добавлено:\n- Источник: %s\n- Сумма: %.2f руб.", income.Source, income.Amount))
}

// getFirstWord возвращает первое слово из строки
func getFirstWord(text string) string {
	words := strings.Fields(text)
	if len(words) > 0 {
		return words[0]
	}
	return ""
}
