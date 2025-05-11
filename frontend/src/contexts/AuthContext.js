import React, { createContext, useState, useEffect } from 'react';
import axios from 'axios';

export const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
  const [currentUser, setCurrentUser] = useState(null);
  const [token, setToken] = useState(localStorage.getItem('token') || null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // Настраиваем axios
  useEffect(() => {
    if (token) {
      axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
    } else {
      delete axios.defaults.headers.common['Authorization'];
    }
  }, [token]);

  // Загружаем данные пользователя при наличии токена
  useEffect(() => {
    const loadUser = async () => {
      if (!token) {
        setLoading(false);
        return;
      }

      try {
        const response = await axios.get('/api/users/me');
        setCurrentUser(response.data);
        setLoading(false);
      } catch (err) {
        console.error('Ошибка при загрузке пользователя:', err);
        logout();
        setLoading(false);
      }
    };

    loadUser();
  }, [token]);

  // Регистрация пользователя
  const register = async (userData) => {
    setError(null);
    
    try {
      const response = await axios.post('/api/auth/signup', userData);
      const { token, user } = response.data;
      
      localStorage.setItem('token', token);
      setToken(token);
      setCurrentUser(user);
      
      return user;
    } catch (err) {
      setError(err.response?.data?.message || 'Ошибка регистрации');
      throw err;
    }
  };

  // Вход пользователя
  const login = async (credentials) => {
    setError(null);
    
    try {
      console.log('Попытка входа с данными:', credentials);
      const response = await axios.post('/api/auth/login', credentials);
      console.log('Ответ сервера при входе:', response.data);
      
      const { token, user } = response.data;
      
      localStorage.setItem('token', token);
      setToken(token);
      setCurrentUser(user);
      
      return user;
    } catch (err) {
      console.error('Ошибка входа:', err);
      console.error('Детали ошибки:', err.response?.data);
      setError(err.response?.data?.message || 'Ошибка входа');
      throw err;
    }
  };

  // Выход пользователя
  const logout = () => {
    localStorage.removeItem('token');
    setToken(null);
    setCurrentUser(null);
  };

  // Обновление данных пользователя
  const updateUser = async (userData) => {
    setError(null);
    
    try {
      const response = await axios.put('/api/users/me', userData);
      setCurrentUser(response.data);
      return response.data;
    } catch (err) {
      setError(err.response?.data?.message || 'Ошибка обновления профиля');
      throw err;
    }
  };

  return (
    <AuthContext.Provider
      value={{
        currentUser,
        token,
        loading,
        error,
        isAuthenticated: !!token,
        register,
        login,
        logout,
        updateUser
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}; 