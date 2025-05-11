package models

// TelegramUser представляет связь между пользователем системы и пользователем Telegram
type TelegramUser struct {
	ID         int64  `json:"id" db:"id"`
	UserID     int64  `json:"user_id" db:"user_id"`
	TelegramID int64  `json:"telegram_id" db:"telegram_id"`
	Username   string `json:"username" db:"username"`
	FirstName  string `json:"first_name" db:"first_name"`
	LastName   string `json:"last_name" db:"last_name"`
}

// TelegramLinkRequest представляет запрос на связывание аккаунта Telegram с аккаунтом пользователя
type TelegramLinkRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
