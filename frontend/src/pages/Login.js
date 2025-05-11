import React, { useState } from 'react';
import { Link as RouterLink, useNavigate } from 'react-router-dom';
import { 
  TextField, 
  Button, 
  Typography, 
  Box, 
  Link,
  Alert,
  CircularProgress,
  Container,
  Paper,
  Tabs,
  Tab,
  Divider,
  IconButton,
  Grid
} from '@mui/material';
import {
  AccountBalanceWallet as WalletIcon,
  Google as GoogleIcon,
  Facebook as FacebookIcon,
  Apple as AppleIcon
} from '@mui/icons-material';
import { useAuth } from '../hooks/useAuth';

export default function Login() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [tabValue, setTabValue] = useState(0);
  const { login } = useAuth();
  const navigate = useNavigate();

  const handleTabChange = (event, newValue) => {
    setTabValue(newValue);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      console.log('Attempting login with:', { email, password });
      
      // DEBUG: Жесткая проверка для тестирования (обход API для диагностики)
      if (email === "testuser@example.com" && password === "password") {
        console.log('DEBUG: Hard-coded auth success!');
        
        // Создаем фиктивный user объект
        const user = {
          id: 999,
          email: email,
          username: "testuser",
          first_name: "Test",
          last_name: "User",
          avatar_url: "",
          monthly_limit: 0,
          savings_goal: 0
        };
        
        // Сохраняем фиктивный токен
        const token = "debug_token_for_testing";
        localStorage.setItem('token', token);
        
        // Эмулируем авторизацию
        navigate('/');
        return;
      }
      
      await login({ email, password });
      navigate('/');
    } catch (err) {
      console.error('Login error:', err);
      console.error('Error details:', err.response?.data);
      setError(err.response?.data?.message || 'Ошибка входа. Проверьте email и пароль.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="auth-form-container">
      <div className="auth-form-card">
        <div className="auth-logo">
          <WalletIcon sx={{ mr: 1, fontSize: 28 }} />
          Finance App
        </div>
        
        <div className="auth-tabs">
          <Tabs
            value={tabValue}
            onChange={handleTabChange}
            variant="fullWidth"
            sx={{
              '& .MuiTabs-indicator': {
                backgroundColor: 'primary.main',
              },
              mb: 3
            }}
          >
            <Tab 
              label="Вход" 
              sx={{
                fontWeight: 'bold',
                '&.Mui-selected': {
                  color: 'primary.main',
                }
              }}
            />
            <Tab 
              label="Регистрация" 
              sx={{
                fontWeight: 'bold',
                '&.Mui-selected': {
                  color: 'primary.main',
                }
              }}
              onClick={() => navigate('/register')}
            />
          </Tabs>
        </div>
        
        <Typography variant="h5" sx={{ mb: 2, fontWeight: 'bold' }}>
          Добро пожаловать!
        </Typography>
        
        <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
          Войдите, чтобы получить доступ к отслеживанию финансов и управлению бюджетом
        </Typography>

        <Box component="form" onSubmit={handleSubmit} noValidate className="animate-fadeIn">
          {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}
          
          <TextField
            margin="normal"
            required
            fullWidth
            id="email"
            label="Email"
            name="email"
            autoComplete="email"
            autoFocus
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            disabled={loading}
            sx={{ mb: 2 }}
          />
          
          <TextField
            margin="normal"
            required
            fullWidth
            name="password"
            label="Пароль"
            type="password"
            id="password"
            autoComplete="current-password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            disabled={loading}
            sx={{ mb: 1 }}
          />
          
          <Box sx={{ textAlign: 'right', mb: 2 }}>
            <Link component={RouterLink} to="/forgot-password" variant="body2" color="primary.main">
              Забыли пароль?
            </Link>
          </Box>
          
          <Button
            type="submit"
            fullWidth
            variant="contained"
            disabled={loading}
            className="auth-submit-button"
          >
            {loading ? <CircularProgress size={24} /> : 'Войти'}
          </Button>
          
          <Divider sx={{ my: 3 }}>
            <Typography variant="body2" color="text.secondary">
              или войти с помощью
            </Typography>
          </Divider>
          
          <Box sx={{ display: 'flex', justifyContent: 'center', gap: 2 }}>
            <IconButton sx={{ bgcolor: 'background.light', color: '#EA4335', '&:hover': { bgcolor: 'background.light' } }}>
              <GoogleIcon />
            </IconButton>
            <IconButton sx={{ bgcolor: 'background.light', color: '#1877F2', '&:hover': { bgcolor: 'background.light' } }}>
              <FacebookIcon />
            </IconButton>
            <IconButton sx={{ bgcolor: 'background.light', color: 'white', '&:hover': { bgcolor: 'background.light' } }}>
              <AppleIcon />
            </IconButton>
          </Box>
          
          <Box sx={{ textAlign: 'center', mt: 3 }}>
            <Typography variant="body2" color="text.secondary">
              Ещё нет аккаунта?{' '}
              <Link component={RouterLink} to="/register" variant="body2" color="primary.main" fontWeight="bold">
                Зарегистрироваться
              </Link>
            </Typography>
          </Box>
        </Box>
      </div>
    </div>
  );
} 