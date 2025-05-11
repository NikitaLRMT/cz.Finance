import React, { useState } from 'react';
import { 
  Box, 
  Typography, 
  TextField, 
  Button, 
  Grid, 
  Paper, 
  Divider,
  InputAdornment,
  Card,
  CardContent,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Tabs,
  Tab
} from '@mui/material';
import { useParams } from 'react-router-dom';

// Компонент калькулятора сложного процента
const CompoundInterestCalculator = () => {
  const [initialAmount, setInitialAmount] = useState('10000');
  const [monthlyContribution, setMonthlyContribution] = useState('500');
  const [annualRate, setAnnualRate] = useState('7');
  const [years, setYears] = useState('10');
  const [results, setResults] = useState(null);

  const calculate = () => {
    const principal = parseFloat(initialAmount) || 0;
    const monthly = parseFloat(monthlyContribution) || 0;
    const rate = (parseFloat(annualRate) || 0) / 100;
    const period = parseInt(years) || 0;
    
    const yearlyResults = [];
    let currentBalance = principal;
    
    for (let year = 1; year <= period; year++) {
      const yearlyContributions = monthly * 12;
      const interestForYear = (currentBalance + yearlyContributions / 2) * rate;
      currentBalance += yearlyContributions + interestForYear;
      
      yearlyResults.push({
        year,
        contributions: principal + yearlyContributions * year,
        interest: currentBalance - (principal + yearlyContributions * year),
        balance: currentBalance
      });
    }
    
    setResults(yearlyResults);
  };

  return (
    <Box>
      <Typography variant="h5" component="h2" gutterBottom>
        Калькулятор сложного процента
      </Typography>
      <Typography variant="body2" color="text.secondary" paragraph>
        Рассчитайте, как ваши инвестиции будут расти со временем с учетом сложного процента
      </Typography>
      
      <Paper elevation={3} sx={{ p: 3, mb: 4 }}>
        <Grid container spacing={3}>
          <Grid item xs={12} sm={6}>
            <TextField
              label="Начальная сумма"
              fullWidth
              value={initialAmount}
              onChange={(e) => setInitialAmount(e.target.value)}
              InputProps={{
                startAdornment: <InputAdornment position="start">₽</InputAdornment>,
              }}
            />
          </Grid>
          <Grid item xs={12} sm={6}>
            <TextField
              label="Ежемесячное пополнение"
              fullWidth
              value={monthlyContribution}
              onChange={(e) => setMonthlyContribution(e.target.value)}
              InputProps={{
                startAdornment: <InputAdornment position="start">₽</InputAdornment>,
              }}
            />
          </Grid>
          <Grid item xs={12} sm={6}>
            <TextField
              label="Годовая процентная ставка"
              fullWidth
              value={annualRate}
              onChange={(e) => setAnnualRate(e.target.value)}
              InputProps={{
                endAdornment: <InputAdornment position="end">%</InputAdornment>,
              }}
            />
          </Grid>
          <Grid item xs={12} sm={6}>
            <TextField
              label="Количество лет"
              fullWidth
              value={years}
              onChange={(e) => setYears(e.target.value)}
            />
          </Grid>
          <Grid item xs={12}>
            <Button 
              variant="contained" 
              fullWidth 
              onClick={calculate}
              size="large"
            >
              Рассчитать
            </Button>
          </Grid>
        </Grid>
      </Paper>
      
      {results && (
        <Paper elevation={3} sx={{ p: 3 }}>
          <Typography variant="h6" gutterBottom>Результаты расчета</Typography>
          
          <Box sx={{ mb: 3 }}>
            <Typography variant="subtitle1">
              Итоговая сумма через {years} лет: 
              <Typography component="span" fontWeight="bold" color="primary" sx={{ ml: 1 }}>
                {results[results.length - 1].balance.toLocaleString()} ₽
              </Typography>
            </Typography>
            <Typography variant="subtitle1">
              Общая сумма вложений: 
              <Typography component="span" fontWeight="bold" sx={{ ml: 1 }}>
                {results[results.length - 1].contributions.toLocaleString()} ₽
              </Typography>
            </Typography>
            <Typography variant="subtitle1">
              Полученный процент: 
              <Typography component="span" fontWeight="bold" color="success.main" sx={{ ml: 1 }}>
                {results[results.length - 1].interest.toLocaleString()} ₽
              </Typography>
            </Typography>
          </Box>
          
          <TableContainer component={Paper} variant="outlined">
            <Table size="small">
              <TableHead>
                <TableRow>
                  <TableCell>Год</TableCell>
                  <TableCell align="right">Сумма взносов</TableCell>
                  <TableCell align="right">Процентный доход</TableCell>
                  <TableCell align="right">Итоговый баланс</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {results.map((row) => (
                  <TableRow key={row.year}>
                    <TableCell component="th" scope="row">{row.year}</TableCell>
                    <TableCell align="right">{row.contributions.toLocaleString()} ₽</TableCell>
                    <TableCell align="right">{row.interest.toLocaleString()} ₽</TableCell>
                    <TableCell align="right">{row.balance.toLocaleString()} ₽</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        </Paper>
      )}
    </Box>
  );
};

// Компонент ипотечного калькулятора
const MortgageCalculator = () => {
  const [loanAmount, setLoanAmount] = useState('5000000');
  const [interestRate, setInterestRate] = useState('7.5');
  const [loanTerm, setLoanTerm] = useState('20');
  const [downPayment, setDownPayment] = useState('1000000');
  const [results, setResults] = useState(null);

  const calculate = () => {
    const principal = parseFloat(loanAmount) - parseFloat(downPayment);
    const interest = parseFloat(interestRate) / 100 / 12;
    const payments = parseFloat(loanTerm) * 12;
    
    // Расчет ежемесячного платежа
    const x = Math.pow(1 + interest, payments);
    const monthly = (principal * x * interest) / (x - 1);
    
    const totalPayment = monthly * payments;
    const totalInterest = totalPayment - principal;
    
    setResults({
      monthlyPayment: monthly,
      totalPayment: totalPayment,
      totalInterest: totalInterest,
      principal: principal
    });
  };

  return (
    <Box>
      <Typography variant="h5" component="h2" gutterBottom>
        Ипотечный калькулятор
      </Typography>
      <Typography variant="body2" color="text.secondary" paragraph>
        Рассчитайте ежемесячный платеж и общую стоимость ипотечного кредита
      </Typography>
      
      <Paper elevation={3} sx={{ p: 3, mb: 4 }}>
        <Grid container spacing={3}>
          <Grid item xs={12} sm={6}>
            <TextField
              label="Стоимость недвижимости"
              fullWidth
              value={loanAmount}
              onChange={(e) => setLoanAmount(e.target.value)}
              InputProps={{
                startAdornment: <InputAdornment position="start">₽</InputAdornment>,
              }}
            />
          </Grid>
          <Grid item xs={12} sm={6}>
            <TextField
              label="Первоначальный взнос"
              fullWidth
              value={downPayment}
              onChange={(e) => setDownPayment(e.target.value)}
              InputProps={{
                startAdornment: <InputAdornment position="start">₽</InputAdornment>,
              }}
            />
          </Grid>
          <Grid item xs={12} sm={6}>
            <TextField
              label="Процентная ставка"
              fullWidth
              value={interestRate}
              onChange={(e) => setInterestRate(e.target.value)}
              InputProps={{
                endAdornment: <InputAdornment position="end">%</InputAdornment>,
              }}
            />
          </Grid>
          <Grid item xs={12} sm={6}>
            <TextField
              label="Срок кредита (лет)"
              fullWidth
              value={loanTerm}
              onChange={(e) => setLoanTerm(e.target.value)}
            />
          </Grid>
          <Grid item xs={12}>
            <Button 
              variant="contained" 
              fullWidth 
              onClick={calculate}
              size="large"
            >
              Рассчитать
            </Button>
          </Grid>
        </Grid>
      </Paper>
      
      {results && (
        <Paper elevation={3} sx={{ p: 3 }}>
          <Typography variant="h6" gutterBottom>Результаты расчета</Typography>
          
          <Grid container spacing={3}>
            <Grid item xs={12} sm={6}>
              <Card variant="outlined">
                <CardContent>
                  <Typography color="text.secondary" gutterBottom>
                    Ежемесячный платеж
                  </Typography>
                  <Typography variant="h5" component="div" color="primary">
                    {results.monthlyPayment.toLocaleString(undefined, {maximumFractionDigits: 0})} ₽
                  </Typography>
                </CardContent>
              </Card>
            </Grid>
            <Grid item xs={12} sm={6}>
              <Card variant="outlined">
                <CardContent>
                  <Typography color="text.secondary" gutterBottom>
                    Сумма кредита
                  </Typography>
                  <Typography variant="h5" component="div">
                    {results.principal.toLocaleString()} ₽
                  </Typography>
                </CardContent>
              </Card>
            </Grid>
            <Grid item xs={12} sm={6}>
              <Card variant="outlined">
                <CardContent>
                  <Typography color="text.secondary" gutterBottom>
                    Общая сумма выплат
                  </Typography>
                  <Typography variant="h5" component="div">
                    {results.totalPayment.toLocaleString(undefined, {maximumFractionDigits: 0})} ₽
                  </Typography>
                </CardContent>
              </Card>
            </Grid>
            <Grid item xs={12} sm={6}>
              <Card variant="outlined">
                <CardContent>
                  <Typography color="text.secondary" gutterBottom>
                    Переплата по кредиту
                  </Typography>
                  <Typography variant="h5" component="div" color="error">
                    {results.totalInterest.toLocaleString(undefined, {maximumFractionDigits: 0})} ₽
                  </Typography>
                </CardContent>
              </Card>
            </Grid>
          </Grid>
        </Paper>
      )}
    </Box>
  );
};

// Главный компонент калькуляторов
export default function Calculators({ type }) {
  const params = useParams();
  const calculatorType = type || params.type || 'default';
  const [activeTab, setActiveTab] = useState(calculatorType);
  
  const handleTabChange = (event, newValue) => {
    setActiveTab(newValue);
  };
  
  // Если калькулятор открыт в защищенной зоне
  if (calculatorType === 'default') {
    return (
      <Box className="page-container">
        <Typography variant="h4" component="h1" gutterBottom fontWeight="bold">
          Финансовые калькуляторы
        </Typography>
        
        {/* Стилизованные карточки для выбора калькулятора */}
        {activeTab === 'default' ? (
          <Grid container spacing={3} sx={{ mb: 4 }}>
            <Grid item xs={12} md={6}>
              <Card 
                sx={{ 
                  cursor: 'pointer',
                  height: '100%',
                  transition: 'transform 0.2s, box-shadow 0.2s',
                  '&:hover': {
                    transform: 'translateY(-4px)',
                    boxShadow: '0 8px 16px rgba(0, 0, 0, 0.2)'
                  }
                }}
                onClick={() => setActiveTab('compound-interest')}
              >
                <CardContent sx={{ p: 4 }}>
                  <Typography variant="h5" component="h2" gutterBottom color="primary">
                    Калькулятор сложного процента
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    Рассчитайте, как ваши инвестиции будут расти со временем с учетом сложного процента.
                    Планируйте своё финансовое будущее с точностью до рубля.
                  </Typography>
                  <Button 
                    variant="outlined" 
                    color="primary" 
                    sx={{ mt: 2 }}
                  >
                    Открыть калькулятор
                  </Button>
                </CardContent>
              </Card>
            </Grid>
            
            <Grid item xs={12} md={6}>
              <Card 
                sx={{ 
                  cursor: 'pointer',
                  height: '100%',
                  transition: 'transform 0.2s, box-shadow 0.2s',
                  '&:hover': {
                    transform: 'translateY(-4px)',
                    boxShadow: '0 8px 16px rgba(0, 0, 0, 0.2)'
                  }
                }}
                onClick={() => setActiveTab('mortgage')}
              >
                <CardContent sx={{ p: 4 }}>
                  <Typography variant="h5" component="h2" gutterBottom color="primary">
                    Ипотечный калькулятор
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    Рассчитайте ежемесячный платеж и общую стоимость ипотечного кредита.
                    Оцените свои финансовые возможности перед покупкой недвижимости.
                  </Typography>
                  <Button 
                    variant="outlined" 
                    color="primary" 
                    sx={{ mt: 2 }}
                  >
                    Открыть калькулятор
                  </Button>
                </CardContent>
              </Card>
            </Grid>
          </Grid>
        ) : (
          <Box sx={{ mb: 4, display: 'flex', alignItems: 'center', gap: 2 }}>
            <Button 
              variant="outlined" 
              onClick={() => setActiveTab('default')}
              startIcon={<span>←</span>}
            >
              Назад к списку
            </Button>
            
            <Typography variant="h6" component="h2" sx={{ ml: 2 }}>
              {activeTab === 'compound-interest' ? 'Калькулятор сложного процента' : 'Ипотечный калькулятор'}
            </Typography>
          </Box>
        )}
        
        {activeTab === 'compound-interest' && <CompoundInterestCalculator />}
        {activeTab === 'mortgage' && <MortgageCalculator />}
      </Box>
    );
  }
  
  // Если открыт конкретный калькулятор (для неавторизованных пользователей)
  return (
    <Box sx={{ maxWidth: 900, mx: 'auto', p: 3 }}>
      {calculatorType === 'compound-interest' && <CompoundInterestCalculator />}
      {calculatorType === 'mortgage' && <MortgageCalculator />}
    </Box>
  );
} 