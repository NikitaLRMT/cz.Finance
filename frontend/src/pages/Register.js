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
  Grid,
  Tabs,
  Tab,
  Divider,
  IconButton,
  InputAdornment
} from '@mui/material';
import {
  AccountBalanceWallet as WalletIcon,
  Google as GoogleIcon,
  Facebook as FacebookIcon,
  Apple as AppleIcon,
  Visibility as VisibilityIcon,
  VisibilityOff as VisibilityOffIcon
} from '@mui/icons-material';
import { useAuth } from '../hooks/useAuth';

export default function Register() {
  const [formData, setFormData] = useState({
    email: '',
    username: '',
    password: '',
    first_name: '',
    last_name: ''
  });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [tabValue, setTabValue] = useState(1);
  const [showPassword, setShowPassword] = useState(false);
  const { register } = useAuth();
  const navigate = useNavigate();

  const handleTabChange = (event, newValue) => {
    setTabValue(newValue);
  };

  const handleClickShowPassword = () => {
    setShowPassword(!showPassword);
  };

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      await register(formData);
      navigate('/');
    } catch (err) {
      setError(err.response?.data?.message || 'Ошибка регистрации. Пожалуйста, попробуйте снова.');
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
              onClick={() => navigate('/login')}
            />
            <Tab 
              label="Регистрация" 
              sx={{
                fontWeight: 'bold',
                '&.Mui-selected': {
                  color: 'primary.main',
                }
              }}
            />
          </Tabs>
        </div>
        
        <Typography variant="h5" sx={{ mb: 2, fontWeight: 'bold' }}>
          Создайте аккаунт
        </Typography>
        
        <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
          Зарегистрируйтесь, чтобы начать управлять своими финансами эффективно
        </Typography>

        <Box component="form" onSubmit={handleSubmit} noValidate className="animate-fadeIn">
          {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}
          
          <Grid container spacing={2}>
            <Grid item xs={12} sm={6}>
              <TextField
                required
                fullWidth
                id="first_name"
                label="Имя"
                name="first_name"
                autoComplete="given-name"
                value={formData.first_name}
                onChange={handleChange}
                disabled={loading}
                sx={{ mb: 1 }}
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                required
                fullWidth
                id="last_name"
                label="Фамилия"
                name="last_name"
                autoComplete="family-name"
                value={formData.last_name}
                onChange={handleChange}
                disabled={loading}
                sx={{ mb: 1 }}
              />
            </Grid>
          </Grid>
          
          <TextField
            margin="normal"
            required
            fullWidth
            id="username"
            label="Имя пользователя"
            name="username"
            autoComplete="username"
            value={formData.username}
            onChange={handleChange}
            disabled={loading}
            sx={{ mb: 2 }}
          />
          
          <TextField
            margin="normal"
            required
            fullWidth
            id="email"
            label="Email"
            name="email"
            autoComplete="email"
            value={formData.email}
            onChange={handleChange}
            disabled={loading}
            sx={{ mb: 2 }}
          />
          
          <TextField
            margin="normal"
            required
            fullWidth
            name="password"
            label="Пароль"
            type={showPassword ? 'text' : 'password'}
            id="password"
            autoComplete="new-password"
            value={formData.password}
            onChange={handleChange}
            disabled={loading}
            sx={{ mb: 2 }}
            InputProps={{
              endAdornment: (
                <InputAdornment position="end">
                  <IconButton
                    aria-label="toggle password visibility"
                    onClick={handleClickShowPassword}
                    edge="end"
                  >
                    {showPassword ? <VisibilityOffIcon /> : <VisibilityIcon />}
                  </IconButton>
                </InputAdornment>
              )
            }}
          />
          
          <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
            Нажимая кнопку "Зарегистрироваться", вы соглашаетесь с нашими условиями использования и политикой конфиденциальности.
          </Typography>
          
          <Button
            type="submit"
            fullWidth
            variant="contained"
            disabled={loading}
            className="auth-submit-button"
          >
            {loading ? <CircularProgress size={24} /> : 'Зарегистрироваться'}
          </Button>
          
          <Divider sx={{ my: 3 }}>
            <Typography variant="body2" color="text.secondary">
              или зарегистрируйтесь с помощью
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
              Уже есть аккаунт?{' '}
              <Link component={RouterLink} to="/login" variant="body2" color="primary.main" fontWeight="bold">
                Войти
              </Link>
            </Typography>
          </Box>
        </Box>
      </div>
    </div>
  );
} 