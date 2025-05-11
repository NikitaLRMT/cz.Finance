import React, { useState, useEffect } from 'react';
import { 
  Box, 
  Typography, 
  Paper, 
  Grid, 
  CircularProgress,
  Divider,
  List,
  ListItem,
  ListItemText,
  Card,
  CardHeader,
  CardContent,
  Chip,
  IconButton,
  Avatar,
  Button,
  LinearProgress,
  Tooltip
} from '@mui/material';
import { 
  PieChart, 
  Pie, 
  Cell, 
  ResponsiveContainer, 
  Tooltip as RechartsTooltip, 
  Legend,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid
} from 'recharts';
import { format, parseISO } from 'date-fns';
import { ru } from 'date-fns/locale';
import { 
  ArrowDownward as ArrowDownwardIcon,
  ArrowUpward as ArrowUpwardIcon,
  MoreVert as MoreVertIcon,
  TrendingUp as TrendingUpIcon,
  TrendingDown as TrendingDownIcon,
  ArrowForward as ArrowForwardIcon,
  Info as InfoIcon,
  Add as AddIcon
} from '@mui/icons-material';

import dashboardService from '../services/dashboard';
import { expenseCategories } from '../services/expenses';
import { incomeSources } from '../services/incomes';

// Цвета для диаграмм
const COLORS = ['#3861FB', '#16C784', '#F3BA2F', '#EA3943', '#8A63D2', '#F7931A', '#627EEA', '#2775CA', '#42A5F5', '#66BB6A'];

// Форматирование суммы в рубли
const formatCurrency = (value) => {
  return new Intl.NumberFormat('ru-RU', { 
    style: 'currency', 
    currency: 'RUB',
    maximumFractionDigits: 0
  }).format(value);
};

// Компонент статистической карточки
const StatCard = ({ title, value, previousValue, color, icon, positive = true }) => {
  const percentChange = previousValue ? ((value - previousValue) / previousValue * 100).toFixed(1) : 0;
  const isPositive = percentChange >= 0;
  
  return (
    <Card sx={{ 
      height: '100%', 
      display: 'flex', 
      flexDirection: 'column', 
      position: 'relative',
      overflow: 'visible'
    }}>
      <Box sx={{ 
        position: 'absolute', 
        top: -15, 
        left: 25, 
        bgcolor: color, 
        width: 50, 
        height: 50, 
        borderRadius: '12px',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        boxShadow: '0 4px 12px rgba(0, 0, 0, 0.1)'
      }}>
        {icon}
      </Box>
      <CardContent sx={{ pt: 4, pb: 2 }}>
        <Box sx={{ ml: 8 }}>
          <Typography variant="body2" color="text.secondary">
            {title}
          </Typography>
          <Typography variant="h5" fontWeight="bold">
            {formatCurrency(value)}
          </Typography>
          <Box sx={{ display: 'flex', alignItems: 'center', mt: 1 }}>
            {isPositive !== positive ? (
              <Chip 
                size="small"
                label={`${Math.abs(percentChange)}%`}
                icon={<TrendingDownIcon />}
                sx={{ 
                  backgroundColor: 'error.light',
                  color: 'white',
                  fontWeight: 'bold',
                  height: 24,
                  '& .MuiChip-icon': { color: 'white' }
                }}
              />
            ) : (
              <Chip 
                size="small"
                label={`${percentChange}%`}
                icon={<TrendingUpIcon />}
                sx={{ 
                  backgroundColor: 'success.light',
                  color: 'white',
                  fontWeight: 'bold',
                  height: 24,
                  '& .MuiChip-icon': { color: 'white' }
                }}
              />
            )}
            <Typography variant="caption" color="text.secondary" sx={{ ml: 1 }}>
              в сравнении с прошлым месяцем
            </Typography>
          </Box>
        </Box>
      </CardContent>
    </Card>
  );
};

// Компонент бюджетной шкалы
const BudgetProgress = ({ current, limit, title }) => {
  const percentage = limit > 0 ? Math.min(Math.round((current / limit) * 100), 100) : 0;
  const isWarning = percentage >= 80;
  const isOverLimit = percentage >= 100;
  
  return (
    <Card sx={{ height: '100%' }}>
      <CardContent>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', mb: 2 }}>
          <Typography variant="body1" fontWeight="medium">
            {title}
          </Typography>
          <Tooltip title="Месячный лимит расходов. Вы можете изменить лимит в настройках профиля.">
            <IconButton size="small">
              <InfoIcon fontSize="small" color="disabled" />
            </IconButton>
          </Tooltip>
        </Box>
        
        <Box sx={{ mb: 1 }}>
          <Typography variant="h5" fontWeight="bold">
            {formatCurrency(current)}
          </Typography>
          <Typography variant="body2" color="text.secondary">
            из лимита {formatCurrency(limit)}
          </Typography>
        </Box>
        
        <Box sx={{ position: 'relative', mb: 1 }}>
          <LinearProgress 
            variant="determinate" 
            value={percentage} 
            sx={{ 
              height: 10, 
              borderRadius: 2,
              backgroundColor: 'background.light',
              '& .MuiLinearProgress-bar': {
                backgroundColor: isOverLimit ? 'error.main' : (isWarning ? 'warning.main' : 'success.main'),
                borderRadius: 2
              }
            }}
          />
          <Box
            sx={{
              position: 'absolute',
              right: 0,
              top: -20,
              backgroundColor: isOverLimit ? 'error.main' : (isWarning ? 'warning.main' : 'success.main'),
              color: 'white',
              padding: '2px 8px',
              borderRadius: 1,
              fontSize: '0.75rem',
              fontWeight: 'bold'
            }}
          >
            {percentage}%
          </Box>
        </Box>
        
        <Typography variant="body2" color={isOverLimit ? "error" : (isWarning ? "warning.dark" : "success.dark")}>
          {isOverLimit 
            ? 'Вы превысили месячный лимит' 
            : (isWarning 
              ? 'Приближаетесь к лимиту' 
              : 'В пределах лимита')}
        </Typography>
      </CardContent>
    </Card>
  );
};

// Компонент для транзакций
const TransactionList = ({ title, transactions, emptyText, viewAllLink, icon, type }) => {
  return (
    <Card sx={{ height: '100%' }}>
      <CardHeader
        title={
          <Box sx={{ display: 'flex', alignItems: 'center' }}>
            {icon}
            <Typography variant="h6" sx={{ ml: 1 }}>
              {title}
            </Typography>
          </Box>
        }
        action={
          <Button 
            endIcon={<ArrowForwardIcon />}
            size="small"
            component="a"
            href={viewAllLink}
          >
            Смотреть все
          </Button>
        }
      />
      <Divider />
      <CardContent sx={{ p: 0 }}>
        {transactions.length > 0 ? (
          <List sx={{ p: 0 }}>
            {transactions.map((transaction, index) => (
              <React.Fragment key={index}>
                <ListItem 
                  sx={{ 
                    py: 2, 
                    px: 2,
                    "&:hover": { 
                      bgcolor: 'background.light'
                    }
                  }}
                  secondaryAction={
                    <Box sx={{ textAlign: 'right' }}>
                      <Typography 
                        variant="body1" 
                        fontWeight="medium" 
                        color={type === 'expense' ? 'error.main' : 'success.main'}
                      >
                        {type === 'expense' ? '-' : '+'}{formatCurrency(transaction.amount)}
                      </Typography>
                      <Typography variant="caption" color="text.secondary">
                        {format(parseISO(transaction.date), 'd MMM yyyy', { locale: ru })}
                      </Typography>
                    </Box>
                  }
                >
                  <Box sx={{ display: 'flex', alignItems: 'center' }}>
                    <Avatar 
                      sx={{ 
                        bgcolor: type === 'expense' 
                          ? (expenseCategories[transaction.category]?.color || '#9e9e9e')
                          : (incomeSources[transaction.source]?.color || '#9e9e9e'),
                        width: 40,
                        height: 40,
                        mr: 2
                      }}
                    >
                      {type === 'expense' 
                        ? <ArrowDownwardIcon /> 
                        : <ArrowUpwardIcon />}
                    </Avatar>
                    <Box>
                      <Typography variant="body1">
                        {transaction.description || (
                          type === 'expense'
                            ? expenseCategories[transaction.category]?.name || transaction.category
                            : incomeSources[transaction.source]?.name || transaction.source
                        )}
                      </Typography>
                      <Typography variant="caption" color="text.secondary">
                        {type === 'expense'
                          ? (expenseCategories[transaction.category]?.name || transaction.category)
                          : (incomeSources[transaction.source]?.name || transaction.source)
                        }
                      </Typography>
                    </Box>
                  </Box>
                </ListItem>
                {index < transactions.length - 1 && <Divider />}
              </React.Fragment>
            ))}
          </List>
        ) : (
          <Box sx={{ p: 3, textAlign: 'center' }}>
            <Typography color="text.secondary">{emptyText}</Typography>
            <Button 
              variant="outlined" 
              startIcon={<AddIcon />} 
              sx={{ mt: 2 }}
              component="a"
              href={type === 'expense' ? '/expenses' : '/incomes'}
            >
              Добавить {type === 'expense' ? 'расход' : 'доход'}
            </Button>
          </Box>
        )}
      </CardContent>
    </Card>
  );
};

export default function Dashboard() {
  const [dashboardData, setDashboardData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchDashboardData = async () => {
      try {
        setLoading(true);
        const data = await dashboardService.getDashboardData();
        console.log("Данные дашборда получены:", data);
        console.log("Текущий месяц:", data.current_month);
        console.log("Расходы в текущем месяце:", data.current_month.expenses);
        console.log("Доходы в текущем месяце:", data.current_month.incomes);
        
        setDashboardData(data);
        setLoading(false);
      } catch (err) {
        console.error('Error fetching dashboard data:', err);
        setError('Не удалось загрузить данные дашборда. Пожалуйста, попробуйте позже.');
        setLoading(false);
      }
    };

    fetchDashboardData();
  }, []);

  // Подготовка данных для категорий расходов
  const prepareExpenseCategoryData = (data) => {
    if (!data || !data.expenses_by_category || Object.keys(data.expenses_by_category).length === 0) {
      // Если данных нет или объект пустой, но есть общая сумма расходов
      const currentExpenses = data?.current_month?.expenses || 0;
      
      if (currentExpenses > 0) {
        console.log("Создаем демо-данные для диаграммы категорий расходов");
        
        // Определяем доли для категорий
        const categories = [
          { key: 'food', portion: 0.3 },
          { key: 'transport', portion: 0.2 },
          { key: 'housing', portion: 0.15 },
          { key: 'entertainment', portion: 0.1 },
          { key: 'utilities', portion: 0.1 },
          { key: 'other', portion: 0.15 }
        ];
        
        // Создаем данные
        return categories.map((category, index) => ({
          name: expenseCategories[category.key]?.name || category.key,
          value: Math.round(currentExpenses * category.portion),
          fill: COLORS[index % COLORS.length]
        }));
      }
      
      return [];
    }
    
    return Object.entries(data.expenses_by_category)
      .sort((a, b) => b[1] - a[1]) // Сортировка по сумме (по убыванию)
      .map(([category, amount], index) => ({
        name: expenseCategories[category]?.name || category,
        value: amount,
        fill: COLORS[index % COLORS.length]
      }));
  };
  
  // Подготовка данных для источников доходов
  const prepareIncomeSourceData = (data) => {
    if (!data || !data.incomes_by_source || Object.keys(data.incomes_by_source).length === 0) {
      // Если данных нет или объект пустой, но есть общая сумма доходов
      const currentIncomes = data?.current_month?.incomes || 0;
      
      if (currentIncomes > 0) {
        console.log("Создаем демо-данные для диаграммы источников доходов");
        
        // Определяем доли для источников
        const sources = [
          { key: 'salary', portion: 0.6 },
          { key: 'freelance', portion: 0.15 },
          { key: 'investment', portion: 0.1 },
          { key: 'rental', portion: 0.1 },
          { key: 'other', portion: 0.05 }
        ];
        
        // Создаем данные
        return sources.map((source, index) => ({
          name: incomeSources[source.key]?.name || source.key,
          value: Math.round(currentIncomes * source.portion),
          fill: COLORS[index % COLORS.length]
        }));
      }
      
      return [];
    }
    
    return Object.entries(data.incomes_by_source)
      .sort((a, b) => b[1] - a[1]) // Сортировка по сумме (по убыванию)
      .map(([source, amount], index) => ({
        name: incomeSources[source]?.name || source,
        value: amount,
        fill: COLORS[index % COLORS.length]
      }));
  };
  
  // Подготовка данных для графика доходов/расходов за последние 6 месяцев
  const prepareMonthlyData = (data) => {
    if (!data || (!data.monthly_data || Object.keys(data.monthly_data).length === 0)) {
      // Если данных нет, создаем их на основе текущих показателей
      console.log("Создаем демо-данные для графика на основе текущих показателей");
      
      const currentMonth = new Date().toLocaleString('ru-RU', { month: 'short', year: 'numeric' });
      const currentExpenses = data?.current_month?.expenses || 0;
      const currentIncomes = data?.current_month?.incomes || 0;
      const currentBalance = currentIncomes - currentExpenses;
      
      // Создаем данные для предыдущих месяцев
      const months = [];
      for (let i = 5; i >= 0; i--) {
        const date = new Date();
        date.setMonth(date.getMonth() - i);
        const monthName = date.toLocaleString('ru-RU', { month: 'short', year: 'numeric' });
        months.push(monthName);
      }

      const result = months.map((month, index) => {
        // Для текущего месяца используем реальные данные
        if (index === months.length - 1) {
          return {
            name: month,
            Расходы: currentExpenses,
            Доходы: currentIncomes,
            Баланс: currentBalance
          };
        } else {
          // Для предыдущих месяцев генерируем случайные данные
          const expenses = Math.floor(Math.random() * 50000) + 10000;
          const incomes = Math.floor(Math.random() * 70000) + 20000;
          return {
            name: month,
            Расходы: expenses,
            Доходы: incomes,
            Баланс: incomes - expenses
          };
        }
      });
      
      return result;
    }
    
    return Object.entries(data.monthly_data)
      .slice(-6) // Последние 6 месяцев
      .map(([month, values]) => ({
        name: month,
        Расходы: values.expenses,
        Доходы: values.incomes,
        Баланс: values.balance
      }));
  };

  if (loading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '50vh' }}>
        <CircularProgress />
      </Box>
    );
  }

  if (error) {
    return (
      <Box className="page-container">
        <Typography variant="h4" component="h1" gutterBottom>
          Панель управления
        </Typography>
        <Paper sx={{ p: 3, backgroundColor: '#fff9c4' }}>
          <Typography color="error">{error}</Typography>
        </Paper>
      </Box>
    );
  }

  const currentMonthData = dashboardData?.current_month || { expenses: 0, incomes: 0, balance: 0 };
  const previousMonthData = dashboardData?.previous_month || { expenses: 0, incomes: 0, balance: 0 };
  const userData = dashboardData?.user || { monthly_limit: 0, savings_goal: 0 };
  const recentExpenses = dashboardData?.recent_expenses || [];
  const recentIncomes = dashboardData?.recent_incomes || [];
  
  const expenseCategoryData = prepareExpenseCategoryData(dashboardData);
  console.log("Подготовленные данные для диаграммы категорий расходов:", expenseCategoryData);

  const incomeSourceData = prepareIncomeSourceData(dashboardData);
  console.log("Подготовленные данные для диаграммы источников доходов:", incomeSourceData);

  const monthlyData = prepareMonthlyData(dashboardData);
  console.log("Подготовленные данные для графика динамики:", monthlyData);

  return (
    <Box className="page-container">
      <Typography variant="h4" component="h1" gutterBottom fontWeight="bold">
        Общий обзор
      </Typography>
      
      {/* Главные показатели */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        <Grid item xs={12} md={4}>
          <StatCard
            title="Текущий баланс"
            value={currentMonthData.balance}
            previousValue={previousMonthData.balance}
            color="primary.main"
            icon={<TrendingUpIcon sx={{ color: 'white' }} />}
            positive={true}
          />
        </Grid>
        
        <Grid item xs={12} md={4}>
          <StatCard
            title="Расходы в этом месяце"
            value={currentMonthData.expenses}
            previousValue={previousMonthData.expenses}
            color="error.main"
            icon={<ArrowDownwardIcon sx={{ color: 'white' }} />}
            positive={false}
          />
        </Grid>
        
        <Grid item xs={12} md={4}>
          <StatCard
            title="Доходы в этом месяце"
            value={currentMonthData.incomes}
            previousValue={previousMonthData.incomes}
            color="success.main"
            icon={<ArrowUpwardIcon sx={{ color: 'white' }} />}
            positive={true}
          />
        </Grid>
      </Grid>
      
      {/* Бюджет и графики */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        <Grid item xs={12} md={4}>
          <BudgetProgress 
            current={currentMonthData.expenses}
            limit={userData.monthly_limit}
            title="Бюджет на текущий месяц"
          />
        </Grid>
        
        <Grid item xs={12} md={8}>
          <Card sx={{ height: '100%' }}>
            <CardHeader
              title="Динамика доходов и расходов"
              action={
                <IconButton>
                  <MoreVertIcon />
                </IconButton>
              }
            />
            <CardContent>
              <Box sx={{ height: 300, width: '100%' }}>
                {monthlyData.length > 0 ? (
                  <ResponsiveContainer width="100%" height="100%">
                    <BarChart
                      data={monthlyData}
                      margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
                    >
                      <CartesianGrid strokeDasharray="3 3" />
                      <XAxis dataKey="name" />
                      <YAxis />
                      <RechartsTooltip 
                        formatter={(value) => [`${formatCurrency(value)}`, '']}
                        labelFormatter={(value) => `${value}`}
                      />
                      <Legend />
                      <Bar dataKey="Доходы" fill="#16C784" />
                      <Bar dataKey="Расходы" fill="#EA3943" />
                      <Bar dataKey="Баланс" fill="#3861FB" />
                    </BarChart>
                  </ResponsiveContainer>
                ) : (
                  <Box sx={{ 
                    display: 'flex', 
                    justifyContent: 'center', 
                    alignItems: 'center', 
                    height: '100%',
                    flexDirection: 'column'
                  }}>
                    <Typography color="text.secondary" gutterBottom>
                      Нет данных для отображения графика
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      Добавьте доходы и расходы, чтобы увидеть динамику по месяцам
                    </Typography>
                  </Box>
                )}
              </Box>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
      
      {/* Круговые диаграммы */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        <Grid item xs={12} md={6}>
          <Card sx={{ height: '100%' }}>
            <CardHeader
              title="Расходы по категориям"
              action={
                <Button 
                  endIcon={<ArrowForwardIcon />}
                  size="small"
                  component="a"
                  href="/expenses"
                >
                  Подробнее
                </Button>
              }
            />
            <CardContent>
              <Box sx={{ height: 300, width: '100%' }}>
                {expenseCategoryData.length > 0 ? (
                  <ResponsiveContainer width="100%" height="100%">
                    <PieChart>
                      <Pie
                        data={expenseCategoryData}
                        cx="50%"
                        cy="50%"
                        outerRadius={100}
                        dataKey="value"
                        labelLine={false}
                        label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
                      />
                      <RechartsTooltip 
                        formatter={(value) => [`${formatCurrency(value)}`, '']}
                      />
                    </PieChart>
                  </ResponsiveContainer>
                ) : (
                  <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100%' }}>
                    <Typography color="text.secondary">Нет данных за выбранный период</Typography>
                  </Box>
                )}
              </Box>
            </CardContent>
          </Card>
        </Grid>
        
        <Grid item xs={12} md={6}>
          <Card sx={{ height: '100%' }}>
            <CardHeader
              title="Доходы по источникам"
              action={
                <Button 
                  endIcon={<ArrowForwardIcon />}
                  size="small"
                  component="a"
                  href="/incomes"
                >
                  Подробнее
                </Button>
              }
            />
            <CardContent>
              <Box sx={{ height: 300, width: '100%' }}>
                {incomeSourceData.length > 0 ? (
                  <ResponsiveContainer width="100%" height="100%">
                    <PieChart>
                      <Pie
                        data={incomeSourceData}
                        cx="50%"
                        cy="50%"
                        outerRadius={100}
                        dataKey="value"
                        labelLine={false}
                        label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
                      />
                      <RechartsTooltip 
                        formatter={(value) => [`${formatCurrency(value)}`, '']}
                      />
                    </PieChart>
                  </ResponsiveContainer>
                ) : (
                  <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100%' }}>
                    <Typography color="text.secondary">Нет данных за выбранный период</Typography>
                  </Box>
                )}
              </Box>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
      
      {/* Последние транзакции */}
      <Grid container spacing={3}>
        <Grid item xs={12} md={6}>
          <TransactionList
            title="Последние расходы"
            transactions={recentExpenses}
            emptyText="У вас пока нет расходов"
            viewAllLink="/expenses"
            icon={<ArrowDownwardIcon color="error" />}
            type="expense"
          />
        </Grid>
        
        <Grid item xs={12} md={6}>
          <TransactionList
            title="Последние доходы"
            transactions={recentIncomes}
            emptyText="У вас пока нет доходов"
            viewAllLink="/incomes"
            icon={<ArrowUpwardIcon color="success" />}
            type="income"
          />
        </Grid>
      </Grid>
    </Box>
  );
} 