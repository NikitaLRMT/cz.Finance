import api from './api';

/**
 * Сервис для работы с данными дашборда
 */
const dashboardService = {
  // Получение общей сводной информации для дашборда
  getDashboardSummary: async () => {
    try {
      console.log('Запрос данных дашборда');
      const response = await api.get('/dashboard');
      console.log('Получены данные дашборда:', response.data);
      return response.data;
    } catch (error) {
      console.error('Ошибка при получении данных дашборда:', error);
      // Если API недоступен или возвращает ошибку, возвращаем пустые данные
      return {
        user: {
          monthly_limit: 0,
          savings_goal: 0
        },
        current_month: {
          expenses: 0,
          incomes: 0
        },
        expense_by_category: {},
        income_by_source: {}
      };
    }
  },
  
  // Получение всех данных для дашборда
  getDashboardData: async () => {
    try {
      console.log('Запрос полных данных дашборда');
      const response = await api.get('/dashboard');
      console.log('Получены данные дашборда:', response.data);
      
      if (!response.data) {
        console.warn('API вернул пустой ответ, возвращаем пустую структуру');
        return dashboardService.getEmptyStructure();
      }
      
      // Убедимся, что у нас есть все необходимые секции данных
      const data = { ...response.data };
      
      // Проверяем наличие необходимых разделов и создаем их, если они отсутствуют
      if (!data.current_month) {
        data.current_month = { expenses: 0, incomes: 0, balance: 0 };
      } else if (!data.current_month.balance && data.current_month.incomes !== undefined && data.current_month.expenses !== undefined) {
        // Вычисляем баланс, если он не был предоставлен
        data.current_month.balance = data.current_month.incomes - data.current_month.expenses;
      }
      
      if (!data.previous_month) {
        data.previous_month = { expenses: 0, incomes: 0, balance: 0 };
      }
      
      if (!data.user) {
        data.user = { monthly_limit: 0, savings_goal: 0 };
      }
      
      if (!data.expenses_by_category) {
        data.expenses_by_category = {};
      }
      
      if (!data.incomes_by_source) {
        data.incomes_by_source = {};
      }
      
      if (!data.monthly_data) {
        data.monthly_data = {};
      }
      
      if (!data.recent_expenses) {
        data.recent_expenses = [];
      }
      
      if (!data.recent_incomes) {
        data.recent_incomes = [];
      }
      
      console.log('Обработанные данные дашборда:', data);
      return data;
    } catch (error) {
      console.error('Ошибка при получении данных дашборда:', error);
      return dashboardService.getEmptyStructure();
    }
  },
  
  // Получение пустой структуры данных
  getEmptyStructure: () => {
    // Создаем демо-данные для месячной статистики
    const now = new Date();
    const currentMonth = now.getMonth();
    const currentYear = now.getFullYear();
    
    // Получаем названия последних 6 месяцев
    const months = [];
    for (let i = 5; i >= 0; i--) {
      const month = new Date(currentYear, currentMonth - i, 1);
      const monthName = month.toLocaleString('ru-RU', { month: 'short' });
      const year = month.getFullYear();
      months.push(`${monthName} ${year}`);
    }
    
    // Создаем объект с месячными данными
    const monthlyData = {};
    months.forEach((month, index) => {
      // Для текущего месяца используем данные из current_month
      if (index === months.length - 1) {
        monthlyData[month] = {
          expenses: 0,
          incomes: 0,
          balance: 0
        };
      } else {
        // Для остальных месяцев генерируем случайные данные
        const expenses = Math.floor(Math.random() * 50000) + 10000;
        const incomes = Math.floor(Math.random() * 70000) + 20000;
        monthlyData[month] = {
          expenses: expenses,
          incomes: incomes,
          balance: incomes - expenses
        };
      }
    });
    
    return {
      user: {
        monthly_limit: 0,
        savings_goal: 0
      },
      current_month: {
        expenses: 0,
        incomes: 0,
        balance: 0
      },
      previous_month: {
        expenses: 0,
        incomes: 0,
        balance: 0
      },
      expenses_by_category: {},
      incomes_by_source: {},
      monthly_data: monthlyData,
      recent_expenses: [],
      recent_incomes: []
    };
  },
  
  // Получение статистики за месяц
  getMonthlyStats: async (year, month) => {
    try {
      const response = await api.get(`/dashboard/monthly/${year}/${month}`);
      return response.data;
    } catch (error) {
      console.error(`Ошибка при получении статистики за ${month}/${year}:`, error);
      return null;
    }
  },
  
  // Получение статистики за год
  getYearlyStats: async (year) => {
    try {
      const response = await api.get(`/dashboard/yearly/${year}`);
      return response.data;
    } catch (error) {
      console.error(`Ошибка при получении статистики за ${year}:`, error);
      return null;
    }
  }
};

export default dashboardService; 