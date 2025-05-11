import React from 'react';
import { 
  Box, 
  Button, 
  Container, 
  Typography, 
  Grid, 
  Card, 
  CardContent, 
  CardActions,
  CardMedia,
  Avatar,
  Stack,
  Chip
} from '@mui/material';
import { Link as RouterLink } from 'react-router-dom';
import { 
  AccountBalanceWallet as WalletIcon,
  TrendingUp as TrendingUpIcon,
  CreditCard as CreditCardIcon,
  ShowChart as ChartIcon,
  Savings as SavingsIcon,
  PieChart as PieChartIcon,
  BarChart as BarChartIcon,
  Timeline as TimelineIcon,
  CheckCircle as CheckCircleIcon
} from '@mui/icons-material';

export default function Landing() {
  const features = [
    {
      title: 'Умный учет расходов',
      description: 'Автоматическая категоризация транзакций и отслеживание всех финансовых потоков',
      icon: <WalletIcon fontSize="large" />
    },
    {
      title: 'Финансовые цели',
      description: 'Устанавливайте цели накоплений с конкретными сроками и следите за прогрессом',
      icon: <TrendingUpIcon fontSize="large" />
    },
    {
      title: 'Интеллектуальное бюджетирование',
      description: 'Создавайте бюджеты по категориям и получайте уведомления при приближении к лимитам',
      icon: <CreditCardIcon fontSize="large" />
    },
    {
      title: 'Продвинутая аналитика',
      description: 'Анализируйте структуру расходов и доходов с помощью наглядных графиков и отчетов',
      icon: <ChartIcon fontSize="large" />
    }
  ];

  const benefits = [
    "Контроль над расходами",
    "Достижение финансовых целей",
    "Снижение ненужных трат",
    "Накопление сбережений",
    "Финансовая безопасность",
    "Уверенность в будущем"
  ];

  return (
    <Box className="landing-container">
      {/* Hero секция */}
      <Box className="landing-hero">
        <Container maxWidth="md" sx={{ position: 'relative', zIndex: 2 }}>
          <Grid container spacing={4} alignItems="center">
            <Grid item xs={12} md={6} className="animate-fadeIn">
              <Typography variant="h2" component="h1" gutterBottom fontWeight="bold">
                Finance App
              </Typography>
              <Typography variant="h5" component="h2" gutterBottom sx={{ mb: 3, color: 'text.secondary' }}>
                Управляйте финансами <span style={{ color: '#3861FB' }}>умно</span> и <span style={{ color: '#16C784' }}>эффективно</span>
              </Typography>
              <Typography variant="body1" sx={{ mb: 4, color: 'text.secondary' }}>
                Современный сервис для учета личных финансов, управления бюджетом и достижения финансовых целей.
              </Typography>
              
              <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2}>
                <Button 
                  component={RouterLink} 
                  to="/register" 
                  variant="contained" 
                  size="large"
                  sx={{ py: 1.5, px: 4, fontWeight: 'bold' }}
                >
                  Начать бесплатно
                </Button>
                <Button 
                  component={RouterLink} 
                  to="/login" 
                  variant="outlined" 
                  size="large"
                  sx={{ py: 1.5, px: 4 }}
                >
                  Войти
                </Button>
              </Stack>
              
              <Box sx={{ mt: 4, display: 'flex', alignItems: 'center' }}>
                <CheckCircleIcon color="success" sx={{ mr: 1 }} />
                <Typography variant="body2" color="text.secondary">
                  14-дневный пробный период. Регистрация за 2 минуты.
                </Typography>
              </Box>
            </Grid>
            
            <Grid item xs={12} md={6} sx={{ display: { xs: 'none', md: 'block' } }}>
              <Box sx={{
                width: '100%',
                height: '350px',
                backgroundColor: 'background.paper',
                borderRadius: '16px',
                boxShadow: '0 20px 40px rgba(0, 0, 0, 0.3)',
                overflow: 'hidden',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                background: 'linear-gradient(135deg, #1a1f2c 0%, #222531 100%)',
                padding: '20px'
              }}>
                <img 
                  src="/assets/dashboard-dark.svg" 
                  alt="Finance App Dashboard" 
                  style={{ 
                    width: '100%', 
                    height: 'auto',
                    objectFit: 'contain',
                    borderRadius: '8px'
                  }}
                  onError={(e) => {
                    e.target.onerror = null;
                    e.target.src = '/assets/fallback-dashboard.png';
                  }}
                />
              </Box>
            </Grid>
          </Grid>
        </Container>
      </Box>
      
      {/* Преимущества */}
      <Container maxWidth="md" sx={{ py: 8 }}>
        <Typography variant="h4" component="h2" align="center" gutterBottom fontWeight="bold">
          Почему стоит выбрать Finance App?
        </Typography>
        <Typography variant="body1" align="center" color="text.secondary" paragraph sx={{ mb: 6, maxWidth: '700px', mx: 'auto' }}>
          Finance App поможет вам не только отслеживать деньги, но и принимать более обоснованные финансовые решения
        </Typography>
        
        <Grid container spacing={3} justifyContent="center">
          {benefits.map((benefit, index) => (
            <Grid item key={index} xs={6} sm={4}>
              <Chip 
                label={index === 1 ? "Достижение целей" : benefit}
                icon={<CheckCircleIcon />} 
                color="primary" 
                variant="outlined" 
                sx={{ 
                  width: '100%', 
                  py: 2.5,
                  justifyContent: 'flex-start',
                  '& .MuiChip-label': { 
                    fontWeight: 500,
                    fontSize: '0.9rem' 
                  }
                }}
              />
            </Grid>
          ))}
        </Grid>
      </Container>

      {/* Возможности приложения */}
      <Box className="landing-features">
        <Container maxWidth="lg">
          <Typography variant="h4" component="h2" align="center" gutterBottom fontWeight="bold">
            Мощные возможности Finance App
          </Typography>
          <Typography variant="body1" align="center" color="text.secondary" paragraph sx={{ mb: 6, maxWidth: '700px', mx: 'auto' }}>
            Получите полный контроль над своими финансами благодаря инновационным функциям
          </Typography>
          
          <Grid container spacing={4}>
            {features.map((feature, index) => (
              <Grid item key={index} xs={12} sm={6} md={3}>
                <Card className="landing-feature-card">
                  <Box className="landing-feature-icon">
                    {feature.icon}
                  </Box>
                  <Typography variant="h6" component="h3" align="center" gutterBottom sx={{ mt: 1, fontWeight: 'bold' }}>
                    {feature.title}
                  </Typography>
                  <Typography variant="body2" align="center" color="text.secondary">
                    {feature.description}
                  </Typography>
                </Card>
              </Grid>
            ))}
          </Grid>
        </Container>
      </Box>

      {/* Секция с калькуляторами */}
      <Container maxWidth="md" sx={{ py: 8 }}>
        <Typography variant="h4" component="h2" align="center" gutterBottom fontWeight="bold">
          Финансовые калькуляторы
        </Typography>
        <Typography variant="body1" align="center" color="text.secondary" paragraph sx={{ mb: 6, maxWidth: '700px', mx: 'auto' }}>
          Планируйте будущее с помощью наших специализированных финансовых калькуляторов
        </Typography>
        
        <Grid container spacing={4}>
          <Grid item xs={12} sm={6}>
            <Card sx={{ height: '100%', display: 'flex', flexDirection: 'column', overflow: 'hidden' }}>
              <Box
                sx={{
                  bgcolor: 'primary.dark',
                  pt: '56.25%', // 16:9 aspect ratio
                  position: 'relative',
                  overflow: 'hidden'
                }}
              >
                <SavingsIcon 
                  sx={{ 
                    position: 'absolute', 
                    top: '50%', 
                    left: '50%', 
                    transform: 'translate(-50%, -50%)',
                    fontSize: 80,
                    color: 'rgba(255, 255, 255, 0.2)'
                  }} 
                />
              </Box>
              <CardContent sx={{ flexGrow: 1 }}>
                <Typography gutterBottom variant="h5" component="h2" fontWeight="bold">
                  Калькулятор сложного процента
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  Рассчитайте рост ваших инвестиций с учетом сложного процента и узнайте потенциальную прибыль за выбранный период времени.
                </Typography>
              </CardContent>
              <CardActions sx={{ p: 2, pt: 0 }}>
                <Button 
                  component={RouterLink} 
                  to="/calculators/compound-interest" 
                  size="large" 
                  fullWidth
                  variant="contained"
                >
                  Открыть калькулятор
                </Button>
              </CardActions>
            </Card>
          </Grid>
          
          <Grid item xs={12} sm={6}>
            <Card sx={{ height: '100%', display: 'flex', flexDirection: 'column', overflow: 'hidden' }}>
              <Box
                sx={{
                  bgcolor: 'secondary.dark',
                  pt: '56.25%', // 16:9 aspect ratio
                  position: 'relative',
                  overflow: 'hidden'
                }}
              >
                <TimelineIcon 
                  sx={{ 
                    position: 'absolute', 
                    top: '50%', 
                    left: '50%', 
                    transform: 'translate(-50%, -50%)',
                    fontSize: 80,
                    color: 'rgba(255, 255, 255, 0.2)'
                  }} 
                />
              </Box>
              <CardContent sx={{ flexGrow: 1 }}>
                <Typography gutterBottom variant="h5" component="h2" fontWeight="bold">
                  Ипотечный калькулятор
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  Рассчитайте ежемесячный платеж по ипотеке, общую стоимость кредита и составьте график погашения при различных условиях.
                </Typography>
              </CardContent>
              <CardActions sx={{ p: 2, pt: 0 }}>
                <Button 
                  component={RouterLink} 
                  to="/calculators/mortgage" 
                  size="large" 
                  fullWidth
                  variant="contained"
                  color="secondary"
                >
                  Открыть калькулятор
                </Button>
              </CardActions>
            </Card>
          </Grid>
        </Grid>
      </Container>

      {/* CTA секция */}
      <Box sx={{ py: 8, background: 'linear-gradient(135deg, #3861FB, #0035C7)', position: 'relative', overflow: 'hidden' }}>
        <Container maxWidth="md" sx={{ position: 'relative', zIndex: 2, textAlign: 'center' }}>
          <Typography variant="h3" component="h2" gutterBottom color="white" fontWeight="bold">
            Готовы взять контроль над финансами?
          </Typography>
          <Typography variant="h6" paragraph color="white" sx={{ opacity: 0.8, mb: 4, maxWidth: '700px', mx: 'auto' }}>
            Присоединитесь к тысячам пользователей, которые уже управляют своими финансами с помощью Finance App
          </Typography>
          <Button 
            component={RouterLink} 
            to="/register" 
            variant="contained" 
            size="large"
            sx={{ 
              py: 1.5, 
              px: 5, 
              bgcolor: 'white', 
              color: 'primary.main',
              '&:hover': {
                bgcolor: 'rgba(255, 255, 255, 0.9)',
              }
            }}
          >
            Начать бесплатно
          </Button>
        </Container>
        <Box
          sx={{
            position: 'absolute',
            top: 0,
            left: 0,
            right: 0,
            bottom: 0,
            opacity: 0.1,
            background: 'url("data:image/svg+xml,%3Csvg width=\'100\' height=\'100\' viewBox=\'0 0 100 100\' xmlns=\'http://www.w3.org/2000/svg\'%3E%3Cpath d=\'M11 18c3.866 0 7-3.134 7-7s-3.134-7-7-7-7 3.134-7 7 3.134 7 7 7zm48 25c3.866 0 7-3.134 7-7s-3.134-7-7-7-7 3.134-7 7 3.134 7 7 7zm-43-7c1.657 0 3-1.343 3-3s-1.343-3-3-3-3 1.343-3 3 1.343 3 3 3zm63 31c1.657 0 3-1.343 3-3s-1.343-3-3-3-3 1.343-3 3 1.343 3 3 3zM34 90c1.657 0 3-1.343 3-3s-1.343-3-3-3-3 1.343-3 3 1.343 3 3 3zm56-76c1.657 0 3-1.343 3-3s-1.343-3-3-3-3 1.343-3 3 1.343 3 3 3zM12 86c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm28-65c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm23-11c2.76 0 5-2.24 5-5s-2.24-5-5-5-5 2.24-5 5 2.24 5 5 5zm-6 60c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm29 22c2.76 0 5-2.24 5-5s-2.24-5-5-5-5 2.24-5 5 2.24 5 5 5zM32 63c2.76 0 5-2.24 5-5s-2.24-5-5-5-5 2.24-5 5 2.24 5 5 5zm57-13c2.76 0 5-2.24 5-5s-2.24-5-5-5-5 2.24-5 5 2.24 5 5 5zm-9-21c1.105 0 2-.895 2-2s-.895-2-2-2-2 .895-2 2 .895 2 2 2zM60 91c1.105 0 2-.895 2-2s-.895-2-2-2-2 .895-2 2 .895 2 2 2zM35 41c1.105 0 2-.895 2-2s-.895-2-2-2-2 .895-2 2 .895 2 2 2zM12 60c1.105 0 2-.895 2-2s-.895-2-2-2-2 .895-2 2 .895 2 2 2z\' fill=\'%23ffffff\' fill-opacity=\'1\' fill-rule=\'evenodd\'/%3E%3C/svg%3E")'
          }}
        />
      </Box>
    </Box>
  );
} 