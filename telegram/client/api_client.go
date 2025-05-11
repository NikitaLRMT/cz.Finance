package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"cz.Finance/backend/models"
)

// APIClient представляет HTTP клиент для работы с API
type APIClient struct {
	baseURL    string
	httpClient *http.Client
	tokenCache map[int64]string // Кэш JWT токенов по Telegram ID
}

// NewAPIClient создает новый экземпляр API клиента
func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		tokenCache: make(map[int64]string),
	}
}

// LinkTelegramAccount связывает аккаунт Telegram с аккаунтом пользователя
func (c *APIClient) LinkTelegramAccount(telegramID int64, username, firstName, lastName string, request *models.TelegramLinkRequest) (*models.User, error) {
	// Создаем данные для запроса
	data := map[string]interface{}{
		"telegram_id": telegramID,
		"username":    username,
		"first_name":  firstName,
		"last_name":   lastName,
		"email":       request.Email,
		"password":    request.Password,
	}

	// Кодируем данные в JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("ошибка при кодировании JSON: %v", err)
	}

	fmt.Printf("Отправляем запрос на связывание аккаунта для Telegram ID %d\n", telegramID)

	// Отправляем запрос
	resp, err := c.doRequest("POST", "/telegram/link", jsonData, 0) // 0 означает что не используем токен авторизации
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	// Читаем тело ответа полностью
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка при чтении ответа: %v", err)
	}

	fmt.Printf("Получен ответ от сервера: %s\n", string(body))

	// Пытаемся распарсить как TokenResponse
	var tokenResponse struct {
		Token     string      `json:"token"`
		ExpiresAt time.Time   `json:"expires_at"`
		User      models.User `json:"user"`
	}

	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		fmt.Printf("Ошибка при разборе ответа как TokenResponse: %v\n", err)
		// Если не удалось распарсить как TokenResponse, пробуем как User
		var user models.User
		if err := json.Unmarshal(body, &user); err != nil {
			return nil, fmt.Errorf("ошибка при декодировании ответа: %v", err)
		}
		fmt.Printf("Успешно получен пользователь ID=%d без токена\n", user.ID)
		return &user, nil
	}

	// Если получили токен, сохраняем его в кэше
	if tokenResponse.Token != "" {
		c.tokenCache[telegramID] = tokenResponse.Token
		fmt.Printf("Сохранен токен для пользователя Telegram ID %d, срок действия до %v\n",
			telegramID, tokenResponse.ExpiresAt)
	} else {
		fmt.Printf("Предупреждение: токен отсутствует в ответе для Telegram ID %d\n", telegramID)
	}

	return &tokenResponse.User, nil
}

// GetUserByTelegramID получает пользователя по Telegram ID
func (c *APIClient) GetUserByTelegramID(telegramID int64) (*models.User, error) {
	// Отправляем запрос
	resp, err := c.doRequest("GET", fmt.Sprintf("/telegram/user/%d", telegramID), nil, int(telegramID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	// Декодируем ответ
	var user models.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("ошибка при декодировании ответа: %v", err)
	}

	return &user, nil
}

// CreateExpense создает новую трату
func (c *APIClient) CreateExpense(userID int64, request *models.CreateExpenseRequest, telegramID int64) (*models.Expense, error) {
	// Кодируем данные в JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("ошибка при кодировании JSON: %v", err)
	}

	// Отправляем запрос
	resp, err := c.doRequest("POST", "/expenses", jsonData, int(telegramID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusCreated {
		return nil, c.handleErrorResponse(resp)
	}

	// Декодируем ответ
	var expense models.Expense
	if err := json.NewDecoder(resp.Body).Decode(&expense); err != nil {
		return nil, fmt.Errorf("ошибка при декодировании ответа: %v", err)
	}

	return &expense, nil
}

// CreateIncome создает новое поступление
func (c *APIClient) CreateIncome(userID int64, request *models.CreateIncomeRequest, telegramID int64) (*models.Income, error) {
	// Кодируем данные в JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("ошибка при кодировании JSON: %v", err)
	}

	// Отправляем запрос
	resp, err := c.doRequest("POST", "/incomes", jsonData, int(telegramID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusCreated {
		return nil, c.handleErrorResponse(resp)
	}

	// Декодируем ответ
	var income models.Income
	if err := json.NewDecoder(resp.Body).Decode(&income); err != nil {
		return nil, fmt.Errorf("ошибка при декодировании ответа: %v", err)
	}

	return &income, nil
}

// GetMonthlyStats получает статистику за месяц
func (c *APIClient) GetMonthlyStats(userID int64, year, month int, telegramID int64) (map[string]interface{}, error) {
	// Отправляем запрос
	resp, err := c.doRequest("GET", fmt.Sprintf("/dashboard/monthly/%d/%d", year, month), nil, int(telegramID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	// Декодируем ответ
	var stats map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, fmt.Errorf("ошибка при декодировании ответа: %v", err)
	}

	return stats, nil
}

// GetRecentTransactions получает последние транзакции пользователя
func (c *APIClient) GetRecentTransactions(limit int, telegramID int64) (map[string]interface{}, error) {
	// Отправляем запрос
	resp, err := c.doRequest("GET", fmt.Sprintf("/dashboard?limit=%d", limit), nil, int(telegramID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	// Декодируем ответ
	var transactions map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&transactions); err != nil {
		return nil, fmt.Errorf("ошибка при декодировании ответа: %v", err)
	}

	return transactions, nil
}

// GetExpensesByCategory получает расходы по определенной категории
func (c *APIClient) GetExpensesByCategory(category string, telegramID int64) ([]models.Expense, error) {
	// Отправляем запрос
	resp, err := c.doRequest("GET", fmt.Sprintf("/expenses?category=%s", category), nil, int(telegramID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	// Декодируем ответ
	var expenses []models.Expense
	if err := json.NewDecoder(resp.Body).Decode(&expenses); err != nil {
		return nil, fmt.Errorf("ошибка при декодировании ответа: %v", err)
	}

	return expenses, nil
}

// SetBudgetGoal устанавливает бюджетную цель для категории
func (c *APIClient) SetBudgetGoal(category string, amount float64, telegramID int64) error {
	// Создаем данные для запроса
	data := map[string]interface{}{
		"category": category,
		"amount":   amount,
	}

	// Кодируем данные в JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("ошибка при кодировании JSON: %v", err)
	}

	// Отправляем запрос
	resp, err := c.doRequest("POST", "/budget/goals", jsonData, int(telegramID))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return c.handleErrorResponse(resp)
	}

	return nil
}

// GetBudgetGoals получает все бюджетные цели
func (c *APIClient) GetBudgetGoals(telegramID int64) (map[string]interface{}, error) {
	url := c.baseURL + "/budget/goals"
	fmt.Printf("Отправляем запрос на получение бюджетных целей: %s\n", url)

	// Отправляем запрос
	resp, err := c.doRequest("GET", "/budget/goals", nil, int(telegramID))
	if err != nil {
		fmt.Printf("Ошибка при отправке запроса: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		errResp := c.handleErrorResponse(resp)
		fmt.Printf("Ошибка от сервера: %v (статус %d)\n", errResp, resp.StatusCode)
		return nil, errResp
	}

	// Декодируем ответ
	var goals map[string]interface{}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Ошибка при чтении тела ответа: %v\n", err)
		return nil, fmt.Errorf("ошибка при чтении ответа: %v", err)
	}

	fmt.Printf("Получен ответ: %s\n", string(body))

	if err := json.Unmarshal(body, &goals); err != nil {
		fmt.Printf("Ошибка при декодировании JSON: %v\n", err)
		return nil, fmt.Errorf("ошибка при декодировании ответа: %v", err)
	}

	return goals, nil
}

// UnlinkAccount отвязывает аккаунт Telegram от аккаунта пользователя
func (c *APIClient) UnlinkAccount(telegramID int64) error {
	fmt.Printf("Отправляем запрос на отвязку аккаунта для Telegram ID %d\n", telegramID)

	// Отправляем запрос
	resp, err := c.doRequest("DELETE", fmt.Sprintf("/telegram/unlink/%d", telegramID), nil, int(telegramID))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return c.handleErrorResponse(resp)
	}

	// Удаляем токен из кэша
	delete(c.tokenCache, telegramID)
	fmt.Printf("Аккаунт Telegram ID %d успешно отвязан, токен удален из кэша\n", telegramID)

	return nil
}

// doRequest выполняет HTTP запрос
func (c *APIClient) doRequest(method, path string, body []byte, telegramID int) (*http.Response, error) {
	url := c.baseURL + path
	var req *http.Request
	var err error

	if body != nil {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(body))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return nil, fmt.Errorf("ошибка при создании запроса: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Устанавливаем заголовок авторизации, если указан ID телеграма
	if telegramID != 0 {
		telegramID64 := int64(telegramID)
		token, exists := c.tokenCache[telegramID64]
		if exists && token != "" {
			fmt.Printf("Добавляем заголовок авторизации для Telegram ID %d: Bearer %s\n", telegramID64, token)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		} else {
			fmt.Printf("Предупреждение: токен не найден для Telegram ID %d\n", telegramID64)
		}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %v", err)
	}

	return resp, nil
}

// handleErrorResponse обрабатывает ответ с ошибкой
func (c *APIClient) handleErrorResponse(resp *http.Response) error {
	var errResp struct {
		Message string `json:"message"`
		Error   string `json:"error"`
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ошибка HTTP %d: не удалось прочитать тело ответа", resp.StatusCode)
	}

	if err := json.Unmarshal(body, &errResp); err != nil {
		return fmt.Errorf("ошибка HTTP %d: %s", resp.StatusCode, string(body))
	}

	if errResp.Error != "" {
		return fmt.Errorf("%s: %s", errResp.Message, errResp.Error)
	}

	return fmt.Errorf("%s", errResp.Message)
}
