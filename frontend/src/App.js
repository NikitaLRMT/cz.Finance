import React from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import { CssBaseline, ThemeProvider } from '@mui/material';
import { useAuth } from './hooks/useAuth';
import theme from './styles/theme';

// Layouts
import MainLayout from './components/layouts/MainLayout';
import AuthLayout from './components/layouts/AuthLayout';

// Pages
import Dashboard from './pages/Dashboard';
import Login from './pages/Login';
import Register from './pages/Register';
import ExpensesList from './pages/ExpensesList';
import IncomesList from './pages/IncomesList';
import Profile from './pages/Profile';
import Calculators from './pages/Calculators';
import NotFound from './pages/NotFound';
import Landing from './pages/Landing';

// Protected route component
const ProtectedRoute = ({ children }) => {
  const { isAuthenticated, loading } = useAuth();
  
  if (loading) {
    return <div>Загрузка...</div>;
  }
  
  return isAuthenticated ? children : <Navigate to="/login" />;
};

function App() {
  const { isAuthenticated } = useAuth();

  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Routes>
        {/* Лендинг для неаутентифицированных пользователей */}
        <Route path="/" element={isAuthenticated ? <Navigate to="/dashboard" /> : <Landing />} />
        
        {/* Калькуляторы доступны всем */}
        <Route path="/calculators/compound-interest" element={<Calculators type="compound-interest" />} />
        <Route path="/calculators/mortgage" element={<Calculators type="mortgage" />} />
        
        {/* Auth routes */}
        <Route element={<AuthLayout />}>
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
        </Route>
        
        {/* Protected routes */}
        <Route element={<MainLayout />}>
          <Route path="/dashboard" element={
            <ProtectedRoute>
              <Dashboard />
            </ProtectedRoute>
          } />
          <Route path="/expenses" element={
            <ProtectedRoute>
              <ExpensesList />
            </ProtectedRoute>
          } />
          <Route path="/incomes" element={
            <ProtectedRoute>
              <IncomesList />
            </ProtectedRoute>
          } />
          <Route path="/profile" element={
            <ProtectedRoute>
              <Profile />
            </ProtectedRoute>
          } />
          <Route path="/calculators" element={
            <ProtectedRoute>
              <Calculators />
            </ProtectedRoute>
          } />
        </Route>
        
        {/* Not Found */}
        <Route path="*" element={<NotFound />} />
      </Routes>
    </ThemeProvider>
  );
}

export default App; 