import axios from 'axios';

// Создаем экземпляр axios с базовым URL
// Определяем базовый URL в зависимости от окружения
const API_URL = process.env.REACT_APP_API_URL || 
  (window.location.hostname === 'localhost' && window.location.port === '3000' 
    ? 'http://localhost:8080/api'
    : '/api');

console.log('API URL:', API_URL);

const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  }
});

// Функция для проверки валидности JWT токена
function isValidJWT(token) {
  if (!token) return false;
  
  // Проверяем формат токена (должен быть три части, разделенные точками)
  const parts = token.split('.');
  if (parts.length !== 3) {
    console.error('Неверный формат JWT-токена - должно быть 3 части:', parts.length);
    return false;
  }
  
  try {
    // Проверяем, что части токена - валидный base64
    const header = JSON.parse(atob(parts[0]));
    const payload = JSON.parse(atob(parts[1]));
    
    console.log('JWT Header:', header);
    console.log('JWT Payload:', payload);
    
    // Проверяем наличие обязательных полей
    if (!payload.exp) {
      console.error('В токене отсутствует поле exp (expiration time)');
      return false;
    }
    
    // Проверяем, не истек ли токен
    const expirationTime = payload.exp * 1000; // переводим в миллисекунды
    const currentTime = Date.now();
    
    if (currentTime > expirationTime) {
      console.error('Токен истек. Текущее время:', new Date(currentTime).toISOString(), 
                 'Время истечения:', new Date(expirationTime).toISOString());
      return false;
    }
    
    return true;
  } catch (e) {
    console.error('Ошибка при разборе JWT-токена:', e);
    return false;
  }
}

// Добавляем интерцептор для добавления токена к запросам
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      // Проверяем валидность токена перед отправкой
      const isValid = isValidJWT(token);
      console.log('Используется токен:', token);
      console.log('Токен валиден:', isValid);
      
      if (!isValid) {
        console.warn('Токен невалиден, но все равно будет отправлен для проверки на сервере');
      }
      
      config.headers['Authorization'] = `Bearer ${token}`;
    }
    console.log('Исходящий запрос:', {
      url: config.url,
      method: config.method,
      headers: JSON.stringify(config.headers),
      data: config.data
    });
    return config;
  },
  (error) => {
    console.error('Ошибка в запросе:', error);
    return Promise.reject(error);
  }
);

// Добавляем интерцептор для обработки ошибок
api.interceptors.response.use(
  (response) => {
    console.log('Ответ API:', {
      url: response.config.url,
      status: response.status,
      data: response.data
    });
    return response;
  },
  (error) => {
    if (error.response) {
      // Сервер ответил с ошибкой
      console.error('Ошибка API:', {
        url: error.config?.url,
        method: error.config?.method,
        status: error.response.status,
        statusText: error.response.statusText,
        headers: error.response.headers,
        data: error.response.data || 'Нет данных',
        requestData: error.config?.data || 'Нет данных запроса'
      });
    } else if (error.request) {
      // Запрос был сделан, но ответ не получен
      console.error('Ошибка сети:', {
        url: error.config?.url,
        method: error.config?.method,
        request: error.request,
        requestData: error.config?.data || 'Нет данных запроса'
      });
    } else {
      // Что-то пошло не так при настройке запроса
      console.error('Ошибка запроса:', error.message);
    }
    
    if (error.response && error.response.status === 401) {
      // Если ошибка авторизации - разлогиниваем пользователя
      console.warn('Получена ошибка 401 - перенаправление на страницу логина');
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    
    return Promise.reject(error);
  }
);

export default api; 