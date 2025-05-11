package routes

import (
	"database/sql"
	"net/http"

	"cz.Finance/backend/configs"
	"cz.Finance/backend/handlers"
	"cz.Finance/backend/middleware"
	"cz.Finance/backend/repositories"
	"cz.Finance/backend/services"

	"github.com/gorilla/mux"
)

// RegisterRoutes регистрирует все маршруты приложения
func RegisterRoutes(router *mux.Router, db *sql.DB, config *configs.Config) {
	// Инициализация репозиториев
	userRepo := repositories.NewUserRepository(db)
	expenseRepo := repositories.NewExpenseRepository(db)
	incomeRepo := repositories.NewIncomeRepository(db)
	wishlistRepo := repositories.NewWishlistRepository(db)
	telegramRepo := repositories.NewTelegramUserRepository(db)

	// Инициализация сервисов
	authService := services.NewAuthService(config.JWT)
	userService := services.NewUserService(userRepo, authService)
	expenseService := services.NewExpenseService(expenseRepo, userRepo)
	incomeService := services.NewIncomeService(incomeRepo, userRepo)
	dashboardService := services.NewDashboardService(expenseRepo, incomeRepo, userRepo)
	wishlistService := services.NewWishlistService(wishlistRepo, userRepo)
	telegramService := services.NewTelegramService(telegramRepo, userRepo)
	calculatorHandler := handlers.NewCalculatorHandler()

	// Инициализация обработчиков
	userHandler := handlers.NewUserHandler(userService)
	expenseHandler := handlers.NewExpenseHandler(expenseService)
	incomeHandler := handlers.NewIncomeHandler(incomeService)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)
	wishlistHandler := handlers.NewWishlistHandler(wishlistService)

	// Настройка маршрутов для публичных API
	public := router.PathPrefix("/api").Subrouter()

	// Авторизация и регистрация
	public.HandleFunc("/auth/signup", userHandler.SignUp).Methods("POST")
	public.HandleFunc("/auth/login", userHandler.Login).Methods("POST")

	// Маршруты для калькуляторов (доступны без авторизации)
	public.HandleFunc("/calculators/compound-interest", calculatorHandler.CompoundInterestCalculator).Methods("POST")
	public.HandleFunc("/calculators/mortgage", calculatorHandler.MortgageCalculator).Methods("POST")

	// Настройка маршрутов для статических файлов
	public.PathPrefix("/uploads/").Handler(http.StripPrefix("/api/uploads/", http.FileServer(http.Dir("./uploads"))))

	// Настройка маршрутов для приватных API (требуют авторизации)
	private := router.PathPrefix("/api").Subrouter()
	private.Use(middleware.AuthMiddleware(config.JWT))

	// Маршруты пользователя
	private.HandleFunc("/users/me", userHandler.GetUser).Methods("GET")
	private.HandleFunc("/users/me", userHandler.UpdateUser).Methods("PUT")
	private.HandleFunc("/users/me", userHandler.DeleteUser).Methods("DELETE")
	private.HandleFunc("/users/me/avatar", userHandler.UploadAvatar).Methods("POST")
	private.HandleFunc("/users/me/avatar", userHandler.RemoveAvatar).Methods("DELETE")

	// Маршруты для трат
	private.HandleFunc("/expenses", expenseHandler.CreateExpense).Methods("POST")
	private.HandleFunc("/expenses", expenseHandler.GetUserExpenses).Methods("GET")
	private.HandleFunc("/expenses/{id:[0-9]+}", expenseHandler.GetExpense).Methods("GET")
	private.HandleFunc("/expenses/{id:[0-9]+}", expenseHandler.UpdateExpense).Methods("PUT")
	private.HandleFunc("/expenses/{id:[0-9]+}", expenseHandler.DeleteExpense).Methods("DELETE")
	private.HandleFunc("/expenses/summary", expenseHandler.GetExpenseSummary).Methods("GET")

	// Маршруты для накоплений/доходов
	private.HandleFunc("/incomes", incomeHandler.CreateIncome).Methods("POST")
	private.HandleFunc("/incomes", incomeHandler.GetUserIncomes).Methods("GET")
	private.HandleFunc("/incomes/{id:[0-9]+}", incomeHandler.GetIncome).Methods("GET")
	private.HandleFunc("/incomes/{id:[0-9]+}", incomeHandler.UpdateIncome).Methods("PUT")
	private.HandleFunc("/incomes/{id:[0-9]+}", incomeHandler.DeleteIncome).Methods("DELETE")
	private.HandleFunc("/incomes/summary", incomeHandler.GetIncomeSummary).Methods("GET")

	// Маршруты для списка желаний
	private.HandleFunc("/wishlist", wishlistHandler.CreateWishlistItem).Methods("POST")
	private.HandleFunc("/wishlist", wishlistHandler.GetUserWishlist).Methods("GET")
	private.HandleFunc("/wishlist/{id:[0-9]+}", wishlistHandler.GetWishlistItem).Methods("GET")
	private.HandleFunc("/wishlist/{id:[0-9]+}", wishlistHandler.UpdateWishlistItem).Methods("PUT")
	private.HandleFunc("/wishlist/{id:[0-9]+}", wishlistHandler.DeleteWishlistItem).Methods("DELETE")

	// Маршруты для информационной панели
	private.HandleFunc("/dashboard", dashboardHandler.GetDashboardSummary).Methods("GET")
	private.HandleFunc("/dashboard/monthly/{year:[0-9]+}/{month:[0-9]+}", dashboardHandler.GetMonthlyStats).Methods("GET")
	private.HandleFunc("/dashboard/yearly/{year:[0-9]+}", dashboardHandler.GetYearlyStats).Methods("GET")

	// Маршруты для бюджетных целей
	private.HandleFunc("/budget/goals", dashboardHandler.GetBudgetGoals).Methods("GET")
	private.HandleFunc("/budget/goals", dashboardHandler.SetBudgetGoal).Methods("POST")

	// Регистрируем обработчики телеграма
	telegramHandler := handlers.NewTelegramHandler(telegramService, userService)
	registerTelegramRoutes(router, telegramHandler)
}

// registerTelegramRoutes регистрирует маршруты для Telegram
func registerTelegramRoutes(router *mux.Router, handler handlers.TelegramHandler) {
	// Маршруты для Telegram
	telegram := router.PathPrefix("/api/telegram").Subrouter()

	telegram.HandleFunc("/link", handler.LinkTelegramAccount).Methods("POST")
	telegram.HandleFunc("/user/{telegram_id:[0-9]+}", handler.GetUserByTelegramID).Methods("GET")
	telegram.HandleFunc("/unlink/{telegram_id:[0-9]+}", handler.UnlinkTelegramAccount).Methods("DELETE")
}
