import api from './api';
import axios from 'axios';

// Константы для категорий расходов
export const expenseCategories = {
  food: {
    name: 'Продукты',
    color: '#16C784'
  },
  transport: {
    name: 'Транспорт',
    color: '#3861FB'
  },
  housing: {
    name: 'Жильё',
    color: '#F3BA2F'
  },
  entertainment: {
    name: 'Развлечения',
    color: '#F3BA2F'
  },
  utilities: {
    name: 'Коммунальные услуги',
    color: '#627EEA'
  },
  health: {
    name: 'Здоровье',
    color: '#8A63D2'
  },
  education: {
    name: 'Образование',
    color: '#2775CA'
  },
  shopping: {
    name: 'Покупки',
    color: '#F7931A'
  },
  travel: {
    name: 'Путешествия',
    color: '#42A5F5'
  },
  other: {
    name: 'Прочие расходы',
    color: '#EA3943'
  }
};

/**
 * Сервис для работы с расходами
 */
const expenseService = {
  // Получение всех расходов
  getExpenses: async () => {
    try {
      const response = await api.get('/expenses');
      return response.data;
    } catch (error) {
      console.error('Ошибка при получении расходов:', error);
      return [];
    }
  },
  
  // Получение расхода по ID
  getExpense: async (id) => {
    try {
      const response = await api.get(`/expenses/${id}`);
      return response.data;
    } catch (error) {
      console.error(`Ошибка при получении расхода с ID ${id}:`, error);
      return null;
    }
  },
  
  // Создание нового расхода
  createExpense: async (expenseData) => {
    try {
      console.log('=== НАЧАЛО СОЗДАНИЯ РАСХОДА ===');
      console.log('Исходные данные:', JSON.stringify(expenseData));
      
      // Форматирование даты
      let formattedDate;
      if (expenseData.date) {
        // Проверяем формат даты и преобразуем его соответствующим образом
        if (expenseData.date instanceof Date) {
          formattedDate = expenseData.date.toISOString();
        } else if (typeof expenseData.date === 'string') {
          // Если дата представлена как строка, создаем объект Date
          formattedDate = new Date(expenseData.date).toISOString();
        } else {
          // По умолчанию используем текущую дату
          formattedDate = new Date().toISOString();
        }
      } else {
        // Если дата не указана, используем текущую
        formattedDate = new Date().toISOString();
      }
      
      console.log('Форматированная дата:', formattedDate);
      
      // Обработка суммы
      let amount;
      if (expenseData.amount) {
        // Убираем все не цифровые символы, кроме точки (для десятичной части)
        const cleanAmount = String(expenseData.amount).replace(/[^\d.]/g, '');
        amount = parseFloat(cleanAmount);
        
        // Проверка на корректность суммы
        if (isNaN(amount) || amount <= 0) {
          throw new Error('Сумма должна быть положительным числом');
        }
      } else {
        throw new Error('Сумма не указана');
      }
      
      console.log('Обработанная сумма:', amount);
      
      // Проверка и установка категории
      const category = expenseData.category || 'other';
      
      // Проверяем, что категория входит в список доступных категорий
      if (!Object.keys(expenseCategories).includes(category)) {
        console.warn(`Неизвестная категория: ${category}. Используется категория "other"`);
      }
      
      // Адаптация данных в формат, ожидаемый бэкендом
      const adaptedData = {
        title: expenseData.description || "Расход", // Гарантируем, что title всегда существует
        amount: amount,
        category: category,
        date: formattedDate,
        description: expenseData.notes || ''
      };
      
      console.log('Отправляем данные расхода:', adaptedData);
      
      // Проверяем текущий JWT токен
      const token = localStorage.getItem('token');
      console.log('Текущий токен:', token);
      
      // Выводим текущие заголовки axios
      console.log('Глобальные заголовки axios:', axios.defaults.headers.common);
      
      try {
        const response = await api.post('/expenses', adaptedData);
        console.log('Ответ сервера при создании расхода:', response.data);
        console.log('=== КОНЕЦ СОЗДАНИЯ РАСХОДА (УСПЕХ) ===');
        return response.data;
      } catch (apiError) {
        console.error('Ошибка API при создании расхода:', apiError);
        console.log('Ответ сервера при ошибке:', apiError.response?.data);
        console.log('=== КОНЕЦ СОЗДАНИЯ РАСХОДА (ОШИБКА API) ===');
        throw apiError;
      }
    } catch (error) {
      console.error('Общая ошибка при создании расхода:', error);
      console.log('=== КОНЕЦ СОЗДАНИЯ РАСХОДА (ОБЩАЯ ОШИБКА) ===');
      throw error;
    }
  },
  
  // Обновление расхода
  updateExpense: async (id, expenseData) => {
    try {
      // Адаптация данных в формат, ожидаемый бэкендом
      const adaptedData = {};
      
      if (expenseData.description) adaptedData.title = expenseData.description;
      if (expenseData.amount) adaptedData.amount = parseFloat(expenseData.amount);
      if (expenseData.category) adaptedData.category = expenseData.category;
      if (expenseData.date) {
        adaptedData.date = expenseData.date instanceof Date ? 
          expenseData.date.toISOString() : 
          new Date(expenseData.date).toISOString();
      }
      if (expenseData.notes !== undefined) adaptedData.description = expenseData.notes;
      
      const response = await api.put(`/expenses/${id}`, adaptedData);
      return response.data;
    } catch (error) {
      console.error(`Ошибка при обновлении расхода с ID ${id}:`, error);
      throw error;
    }
  },
  
  // Удаление расхода
  deleteExpense: async (id) => {
    try {
      await api.delete(`/expenses/${id}`);
      return true;
    } catch (error) {
      console.error(`Ошибка при удалении расхода с ID ${id}:`, error);
      throw error;
    }
  },
  
  // Получение сводки расходов
  getExpenseSummary: async () => {
    try {
      const response = await api.get('/expenses/summary');
      return response.data;
    } catch (error) {
      console.error('Error fetching expense summary:', error);
      throw error;
    }
  },
  
  // Преобразование категорий в массив для селектов
  getCategoriesArray: () => {
    return Object.entries(expenseCategories).map(([value, data]) => ({
      id: value,
      name: data.name
    }));
  }
};

export default expenseService; 