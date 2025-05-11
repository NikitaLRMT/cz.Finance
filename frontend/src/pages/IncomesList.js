import React, { useState, useEffect } from 'react';
import { 
  Box, 
  Typography, 
  Paper, 
  Table, 
  TableBody, 
  TableCell, 
  TableContainer, 
  TableHead, 
  TableRow, 
  Button,
  IconButton,
  TextField,
  InputAdornment,
  Chip,
  Tooltip,
  CircularProgress,
  Fab,
  Alert,
  Snackbar,
  TablePagination
} from '@mui/material';
import { 
  Add as AddIcon, 
  Edit as EditIcon, 
  Delete as DeleteIcon,
  FilterList as FilterListIcon,
  Search as SearchIcon
} from '@mui/icons-material';
import { format, parseISO } from 'date-fns';
import { ru } from 'date-fns/locale';
import incomesService, { incomeSources } from '../services/incomes';
import AddIncomeForm from '../components/incomes/AddIncomeForm';
import EditIncomeForm from '../components/incomes/EditIncomeForm';

// Логируем, какие компоненты и данные доступны
console.log('IncomesList: Модуль загружен, incomeSources:', incomeSources);
console.log('IncomesList: Модуль загружен, incomesService:', incomesService);
console.log('IncomesList: Доступны компоненты AddIncomeForm:', !!AddIncomeForm, 'EditIncomeForm:', !!EditIncomeForm);

// Дополнительные категории для демонстрации, пока API не готов
const incomeCategories = {
  1: { name: 'Зарплата', color: '#4caf50' },
  2: { name: 'Фриланс', color: '#2196f3' },
  3: { name: 'Подарок', color: '#9c27b0' },
  4: { name: 'Инвестиции', color: '#ff9800' },
  5: { name: 'Проценты по вкладу', color: '#f44336' },
  6: { name: 'Возврат долга', color: '#607d8b' },
  7: { name: 'Продажа имущества', color: '#3f51b5' },
  8: { name: 'Прочее', color: '#795548' }
};

// Моковые данные для демонстрации
const mockIncomes = [
  { id: 1, amount: 80000, description: 'Зарплата за август', date: '2023-08-10', category_id: 1 },
  { id: 2, amount: 15000, description: 'Проект по фрилансу', date: '2023-08-08', category_id: 2 },
  { id: 3, amount: 5000, description: 'Подарок на день рождения', date: '2023-08-06', category_id: 3 },
  { id: 4, amount: 3200, description: 'Дивиденды по акциям', date: '2023-08-05', category_id: 4 },
  { id: 5, amount: 1800, description: 'Проценты по вкладу', date: '2023-08-03', category_id: 5 },
  { id: 6, amount: 10000, description: 'Возврат долга от Ивана', date: '2023-08-01', category_id: 6 },
  { id: 7, amount: 45000, description: 'Продажа старого ноутбука', date: '2023-07-29', category_id: 7 },
  { id: 8, amount: 2500, description: 'Кэшбэк с покупок', date: '2023-07-27', category_id: 8 }
];

export default function IncomesList() {
  console.log('IncomesList: Компонент рендерится');
  
  const [incomes, setIncomes] = useState([]);
  const [loading, setLoading] = useState(true);
  const [openAddDialog, setOpenAddDialog] = useState(false);
  const [openEditDialog, setOpenEditDialog] = useState(false);
  const [currentIncome, setCurrentIncome] = useState(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [notification, setNotification] = useState({ open: false, message: '', severity: 'success' });
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  // Добавляем состояние для отслеживания ошибок рендеринга
  const [renderError, setRenderError] = useState(null);

  useEffect(() => {
    console.log('IncomesList: useEffect запущен');
    try {
      fetchIncomes();
    } catch (error) {
      console.error('IncomesList: Ошибка в useEffect:', error);
      setRenderError(error.message);
    }
  }, []);

  const fetchIncomes = async () => {
    console.log('IncomesList: Начало загрузки доходов');
    setLoading(true);
    try {
      console.log('IncomesList: Вызов API для получения доходов');
      const data = await incomesService.getIncomes();
      console.log('IncomesList: Получены данные о доходах:', data);
      setIncomes(data || []);
      setLoading(false);
    } catch (error) {
      console.error('IncomesList: Ошибка при загрузке доходов:', error);
      setNotification({
        open: true,
        message: 'Ошибка при загрузке доходов: ' + (error.message || 'Неизвестная ошибка'),
        severity: 'error'
      });
      setRenderError(error.message || 'Ошибка при загрузке доходов');
      setLoading(false);
    }
  };

  const handleAddIncome = (newIncome) => {
    // После успешного добавления обновляем список
    fetchIncomes();
    
    setNotification({
      open: true,
      message: 'Доход успешно добавлен',
      severity: 'success'
    });
  };

  const handleEditIncome = (income) => {
    setCurrentIncome(income);
    setOpenEditDialog(true);
  };

  const handleIncomeUpdated = () => {
    // После успешного обновления обновляем список
    fetchIncomes();
    
    setNotification({
      open: true,
      message: 'Доход успешно обновлен',
      severity: 'success'
    });
  };

  const handleDeleteIncome = async (id) => {
    try {
      await incomesService.deleteIncome(id);
      
      // После успешного удаления обновляем список
      fetchIncomes();
      
      setNotification({
        open: true,
        message: 'Доход успешно удален',
        severity: 'success'
      });
    } catch (error) {
      console.error(`Ошибка при удалении дохода с id ${id}:`, error);
      setNotification({
        open: true,
        message: 'Ошибка при удалении дохода',
        severity: 'error'
      });
    }
  };

  const handleCloseNotification = () => {
    setNotification({ ...notification, open: false });
  };

  const handleChangePage = (event, newPage) => {
    setPage(newPage);
  };

  const handleChangeRowsPerPage = (event) => {
    setRowsPerPage(parseInt(event.target.value, 10));
    setPage(0);
  };

  // Фильтрация по поисковому запросу
  const filteredIncomes = incomes.filter(income => 
    income.description?.toLowerCase().includes(searchTerm.toLowerCase()) ||
    (incomeSources[income.source]?.name.toLowerCase().includes(searchTerm.toLowerCase()))
  );

  // Получаем данные для текущей страницы
  const paginatedIncomes = filteredIncomes.slice(
    page * rowsPerPage,
    page * rowsPerPage + rowsPerPage
  );

  // Общая сумма доходов
  const totalAmount = filteredIncomes.reduce((sum, income) => sum + income.amount, 0);

  // Форматируем дату
  const formatDate = (dateString) => {
    try {
      return format(
        typeof dateString === 'string' ? parseISO(dateString) : new Date(dateString), 
        'dd MMM yyyy', 
        { locale: ru }
      );
    } catch (error) {
      console.error('Ошибка форматирования даты:', error);
      return 'Неизвестная дата';
    }
  };

  // Оборачиваем весь рендеринг в try-catch
  try {
    console.log('IncomesList: Начало рендеринга, loading:', loading);
    console.log('IncomesList: Количество доходов:', filteredIncomes?.length);
    
    return (
      <Box className="page-container">
        {renderError && (
          <Alert severity="error" sx={{ mb: 2 }}>
            Ошибка при загрузке страницы: {renderError}
          </Alert>
        )}
      
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
          <Typography variant="h4" component="h1">
            Список доходов
          </Typography>
          <Fab 
            color="primary" 
            aria-label="add" 
            onClick={() => setOpenAddDialog(true)}
            size="medium"
          >
            <AddIcon />
          </Fab>
        </Box>

        <Paper sx={{ mb: 3, p: 2 }}>
          <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
            <TextField
              label="Поиск по описанию или источнику"
              variant="outlined"
              size="small"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              fullWidth
              sx={{ mr: 2 }}
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <SearchIcon />
                  </InputAdornment>
                ),
              }}
            />
            <Tooltip title="Фильтры">
              <IconButton>
                <FilterListIcon />
              </IconButton>
            </Tooltip>
          </Box>
          
          <Box sx={{ display: 'flex', justifyContent: 'flex-end', mb: 1 }}>
            <Typography variant="subtitle1" fontWeight="bold">
              Общая сумма: {typeof totalAmount === 'number' ? totalAmount.toLocaleString() : '0'} ₽
            </Typography>
          </Box>
        </Paper>

        {loading ? (
          <Box sx={{ display: 'flex', justifyContent: 'center', my: 4 }}>
            <CircularProgress />
          </Box>
        ) : Array.isArray(paginatedIncomes) && paginatedIncomes.length > 0 ? (
          <Paper>
            <TableContainer>
              <Table>
                <TableHead>
                  <TableRow>
                    <TableCell>Дата</TableCell>
                    <TableCell>Источник</TableCell>
                    <TableCell>Описание</TableCell>
                    <TableCell align="right">Сумма</TableCell>
                    <TableCell align="center">Действия</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {paginatedIncomes.map((income) => (
                    <TableRow key={income.id}>
                      <TableCell>
                        {formatDate(income.date)}
                      </TableCell>
                      <TableCell>
                        <Chip 
                          label={incomeSources[income.source]?.name || incomeCategories[income.category_id]?.name || 'Прочее'} 
                          size="small"
                          style={{ 
                            backgroundColor: incomeSources[income.source]?.color || incomeCategories[income.category_id]?.color || '#757575',
                            color: 'white' 
                          }}
                        />
                      </TableCell>
                      <TableCell>{income.description}</TableCell>
                      <TableCell align="right">{income.amount.toLocaleString()} ₽</TableCell>
                      <TableCell align="center">
                        <IconButton 
                          size="small" 
                          color="primary"
                          onClick={() => handleEditIncome(income)}
                        >
                          <EditIcon fontSize="small" />
                        </IconButton>
                        <IconButton 
                          size="small" 
                          color="error"
                          onClick={() => handleDeleteIncome(income.id)}
                        >
                          <DeleteIcon fontSize="small" />
                        </IconButton>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </TableContainer>
            
            <TablePagination
              rowsPerPageOptions={[5, 10, 25]}
              component="div"
              count={filteredIncomes.length}
              rowsPerPage={rowsPerPage}
              page={page}
              onPageChange={handleChangePage}
              onRowsPerPageChange={handleChangeRowsPerPage}
              labelRowsPerPage="Строк на странице:"
              labelDisplayedRows={({ from, to, count }) => `${from}-${to} из ${count}`}
            />
          </Paper>
        ) : (
          <Alert severity="info" sx={{ mt: 2 }}>
            {searchTerm 
              ? 'Доходы по вашему запросу не найдены' 
              : 'У вас пока нет добавленных доходов'}
          </Alert>
        )}

        <AddIncomeForm 
          open={openAddDialog} 
          handleClose={() => setOpenAddDialog(false)}
          onIncomeAdded={handleAddIncome}
        />

        <EditIncomeForm
          open={openEditDialog}
          handleClose={() => setOpenEditDialog(false)}
          onIncomeUpdated={handleIncomeUpdated}
          income={currentIncome}
        />

        <Snackbar 
          open={notification.open} 
          autoHideDuration={6000} 
          onClose={handleCloseNotification}
          anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
        >
          <Alert onClose={handleCloseNotification} severity={notification.severity}>
            {notification.message}
          </Alert>
        </Snackbar>
      </Box>
    );
  } catch (error) {
    console.error('IncomesList: Ошибка при рендеринге компонента:', error);
    return (
      <Box className="page-container">
        <Alert severity="error" sx={{ mb: 2 }}>
          Произошла ошибка при отображении страницы: {error.message || 'Неизвестная ошибка'}
        </Alert>
      </Box>
    );
  }
} 