import api from './api';

// Константы для источников доходов
export const incomeSources = {
  salary: {
    name: 'Зарплата',
    color: '#16C784'
  },
  freelance: {
    name: 'Фриланс',
    color: '#3861FB'
  },
  investments: {
    name: 'Инвестиции',
    color: '#F3BA2F'
  },
  rental: {
    name: 'Аренда',
    color: '#627EEA'
  },
  business: {
    name: 'Бизнес',
    color: '#8A63D2'
  },
  gifts: {
    name: 'Подарки',
    color: '#2775CA'
  },
  other: {
    name: 'Прочие доходы',
    color: '#42A5F5'
  }
};

/**
 * Сервис для работы с доходами
 */
const incomeService = {
  // Получение всех доходов
  getIncomes: async () => {
    try {
      const response = await api.get('/incomes');
      return response.data;
    } catch (error) {
      console.error('Ошибка при получении доходов:', error);
      return [];
    }
  },
  
  // Получение дохода по ID
  getIncome: async (id) => {
    try {
      const response = await api.get(`/incomes/${id}`);
      return response.data;
    } catch (error) {
      console.error(`Ошибка при получении дохода с ID ${id}:`, error);
      return null;
    }
  },
  
  // Создание нового дохода
  createIncome: async (incomeData) => {
    try {
      // Форматирование даты
      let formattedDate;
      if (incomeData.date) {
        // Проверяем формат даты и преобразуем его соответствующим образом
        if (incomeData.date instanceof Date) {
          formattedDate = incomeData.date.toISOString();
        } else if (typeof incomeData.date === 'string') {
          // Если дата представлена как строка, создаем объект Date
          formattedDate = new Date(incomeData.date).toISOString();
        } else {
          // По умолчанию используем текущую дату
          formattedDate = new Date().toISOString();
        }
      } else {
        // Если дата не указана, используем текущую
        formattedDate = new Date().toISOString();
      }
      
      // Обработка суммы
      let amount;
      if (incomeData.amount) {
        // Убираем все не цифровые символы, кроме точки (для десятичной части)
        const cleanAmount = String(incomeData.amount).replace(/[^\d.]/g, '');
        amount = parseFloat(cleanAmount);
        
        // Проверка на корректность суммы
        if (isNaN(amount) || amount <= 0) {
          throw new Error('Сумма должна быть положительным числом');
        }
      } else {
        throw new Error('Сумма не указана');
      }
      
      // Проверка и установка источника дохода
      const source = incomeData.category || 'other';
      
      // Проверяем, что источник входит в список доступных источников
      if (!Object.keys(incomeSources).includes(source)) {
        console.warn(`Неизвестный источник дохода: ${source}. Используется источник "other"`);
      }
      
      // Адаптация данных в формат, ожидаемый бэкендом
      const adaptedData = {
        amount: amount,
        source: source,
        date: formattedDate,
        description: incomeData.description || ''
      };
      
      console.log('Отправляем данные дохода:', adaptedData);
      
      const response = await api.post('/incomes', adaptedData);
      console.log('Ответ сервера при создании дохода:', response.data);
      return response.data;
    } catch (error) {
      console.error('Ошибка при создании дохода:', error);
      throw error;
    }
  },
  
  // Обновление дохода
  updateIncome: async (id, incomeData) => {
    try {
      // Адаптация данных в формат, ожидаемый бэкендом
      const adaptedData = {};
      
      if (incomeData.amount) adaptedData.amount = parseFloat(incomeData.amount);
      if (incomeData.category) adaptedData.source = incomeData.category;
      if (incomeData.date) {
        adaptedData.date = incomeData.date instanceof Date ? 
          incomeData.date.toISOString() : 
          new Date(incomeData.date).toISOString();
      }
      if (incomeData.description !== undefined) adaptedData.description = incomeData.description;
      
      const response = await api.put(`/incomes/${id}`, adaptedData);
      return response.data;
    } catch (error) {
      console.error(`Ошибка при обновлении дохода с ID ${id}:`, error);
      throw error;
    }
  },
  
  // Удаление дохода
  deleteIncome: async (id) => {
    try {
      await api.delete(`/incomes/${id}`);
      return true;
    } catch (error) {
      console.error(`Ошибка при удалении дохода с ID ${id}:`, error);
      throw error;
    }
  },
  
  // Получение сводки доходов
  getIncomeSummary: async () => {
    try {
      const response = await api.get('/incomes/summary');
      return response.data;
    } catch (error) {
      console.error('Error fetching income summary:', error);
      throw error;
    }
  },
  
  // Преобразование источников дохода в массив для селектов
  getSourcesArray: () => {
    return Object.entries(incomeSources).map(([value, data]) => ({
      id: value,
      name: data.name
    }));
  }
};

export default incomeService; 