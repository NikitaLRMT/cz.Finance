import api from './api';

/**
 * Сервис для работы с данными пользователя
 */
const userService = {
  // Получение данных текущего пользователя
  getCurrentUser: async () => {
    try {
      const response = await api.get('/users/me');
      console.log('Ответ API о пользователе:', response.data);
      
      // Выводим URL аватара для отладки
      if (response.data && response.data.avatar_url) {
        console.log('URL аватара пользователя:', response.data.avatar_url);
        console.log('Полный URL аватара:', window.location.origin + response.data.avatar_url);
      }
      
      return response.data;
    } catch (error) {
      console.error('Error fetching current user:', error);
      // Имитируем данные пользователя на случай, если API недоступен
      return {
        id: 1,
        username: 'lrmt',
        email: 'nikitafaisulin1@gmail.com',
        first_name: 'Никита',
        last_name: 'Ильин',
        avatar_url: '/default-avatar.png',
        monthly_limit: 50000,
        savings_goal: 10000,
        created_at: '2025-05-10T13:40:23.311498Z'
      };
    }
  },
  
  // Обновление данных пользователя
  updateUser: async (userData) => {
    try {
      // Преобразуем числовые поля, если они есть
      if (userData.monthly_limit !== undefined) {
        userData.monthly_limit = parseFloat(userData.monthly_limit) || 0;
      }
      if (userData.savings_goal !== undefined) {
        userData.savings_goal = parseFloat(userData.savings_goal) || 0;
      }
      
      const response = await api.put('/users/me', userData);
      return response.data;
    } catch (error) {
      console.error('Error updating user data:', error);
      throw error;
    }
  },
  
  // Загрузка аватара пользователя
  uploadAvatar: async (file) => {
    try {
      const formData = new FormData();
      formData.append('avatar', file);
      
      const response = await api.post('/users/me/avatar', formData, {
        headers: {
          'Content-Type': 'multipart/form-data'
        }
      });
      return response.data;
    } catch (error) {
      console.error('Error uploading avatar:', error);
      throw error;
    }
  },
  
  // Удаление аватара пользователя
  removeAvatar: async () => {
    try {
      const response = await api.delete('/users/me/avatar');
      return response.data;
    } catch (error) {
      console.error('Error removing avatar:', error);
      throw error;
    }
  },
  
  // Получение списка желаний
  getWishlist: async () => {
    try {
      const response = await api.get('/wishlist');
      return response.data;
    } catch (error) {
      console.error('Error getting wishlist:', error);
      return [];
    }
  },
  
  // Добавление элемента в список желаний
  addWishlistItem: async (wishlistItem) => {
    try {
      // Преобразуем данные для бэкенда
      const requestData = {
        title: wishlistItem.title,
        price: parseFloat(wishlistItem.price),
        priority: wishlistItem.priority,
        description: wishlistItem.description || ''
      };
      
      const response = await api.post('/wishlist', requestData);
      return response.data;
    } catch (error) {
      console.error('Error adding wishlist item:', error);
      throw error;
    }
  },
  
  // Удаление элемента из списка желаний
  removeWishlistItem: async (itemId) => {
    try {
      const response = await api.delete(`/wishlist/${itemId}`);
      return response.data;
    } catch (error) {
      console.error(`Error removing wishlist item with id ${itemId}:`, error);
      throw error;
    }
  },
  
  // Обновление финансовых целей пользователя
  updateFinancialGoals: async (goals) => {
    try {
      // Используем обычный метод обновления пользователя
      const response = await api.put('/users/me', goals);
      return response.data;
    } catch (error) {
      console.error('Error updating financial goals:', error);
      throw error;
    }
  }
};

export default userService; 
