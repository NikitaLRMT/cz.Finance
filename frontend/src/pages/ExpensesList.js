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
import expensesService, { expenseCategories } from '../services/expenses';
import AddExpenseForm from '../components/expenses/AddExpenseForm';
import EditExpenseForm from '../components/expenses/EditExpenseForm';

export default function ExpensesList() {
  // Логируем монтирование компонента
  console.log('ExpensesList компонент рендерится');
  
  const [expenses, setExpenses] = useState([]);
  const [loading, setLoading] = useState(true);
  const [openAddDialog, setOpenAddDialog] = useState(false);
  const [openEditDialog, setOpenEditDialog] = useState(false);
  const [currentExpense, setCurrentExpense] = useState(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [notification, setNotification] = useState({ open: false, message: '', severity: 'success' });
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  // Добавляем состояние для отслеживания ошибок рендеринга
  const [renderError, setRenderError] = useState(null);

  useEffect(() => {
    console.log('ExpensesList: useEffect запустился');
    try {
      fetchExpenses();
    } catch (error) {
      console.error('Ошибка в useEffect:', error);
      setRenderError(error.message);
    }
  }, []);

  const fetchExpenses = async () => {
    console.log('ExpensesList: Начинаем загрузку расходов');
    setLoading(true);
    try {
      console.log('ExpensesList: Вызов API для получения расходов');
      const data = await expensesService.getExpenses();
      console.log('ExpensesList: Данные получены:', data);
      setExpenses(data || []); // Защита от null
      setLoading(false);
    } catch (error) {
      console.error('Ошибка при загрузке расходов:', error);
      setNotification({
        open: true,
        message: 'Ошибка при загрузке расходов: ' + (error.message || 'Неизвестная ошибка'),
        severity: 'error'
      });
      setRenderError(error.message || 'Ошибка при загрузке расходов');
      setLoading(false);
    }
  };

  const handleAddExpense = (newExpense) => {
    // После успешного добавления обновляем список
    fetchExpenses();
    
    setNotification({
      open: true,
      message: 'Расход успешно добавлен',
      severity: 'success'
    });
  };

  const handleDeleteExpense = async (id) => {
    try {
      await expensesService.deleteExpense(id);
      
      // После успешного удаления обновляем список
      fetchExpenses();
      
      setNotification({
        open: true,
        message: 'Расход успешно удален',
        severity: 'success'
      });
    } catch (error) {
      console.error(`Ошибка при удалении расхода с id ${id}:`, error);
      setNotification({
        open: true,
        message: 'Ошибка при удалении расхода',
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

  const handleEditExpense = (expense) => {
    setCurrentExpense(expense);
    setOpenEditDialog(true);
  };

  const handleExpenseUpdated = () => {
    // После успешного обновления обновляем список
    fetchExpenses();
    
    setNotification({
      open: true,
      message: 'Расход успешно обновлен',
      severity: 'success'
    });
  };

  // Фильтрация по поисковому запросу
  const filteredExpenses = expenses.filter(expense => 
    expense.title?.toLowerCase().includes(searchTerm.toLowerCase()) ||
    expense.description?.toLowerCase().includes(searchTerm.toLowerCase()) ||
    (expenseCategories[expense.category]?.name.toLowerCase().includes(searchTerm.toLowerCase()))
  );

  // Получаем данные для текущей страницы
  const paginatedExpenses = filteredExpenses.slice(
    page * rowsPerPage,
    page * rowsPerPage + rowsPerPage
  );

  // Общая сумма расходов
  const totalAmount = filteredExpenses.reduce((sum, expense) => sum + expense.amount, 0);

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

  // Обновляем секцию return с проверкой ошибок
  try {
    console.log('ExpensesList: Рендеринг компонента, состояние loading:', loading);
    console.log('ExpensesList: Количество расходов:', filteredExpenses?.length);
    
    return (
      <Box className="page-container">
        {renderError && (
          <Alert severity="error" sx={{ mb: 2 }}>
            Ошибка при загрузке страницы: {renderError}
          </Alert>
        )}
        
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
          <Typography variant="h4" component="h1">
            Список расходов
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
              label="Поиск по описанию или категории"
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
        ) : Array.isArray(paginatedExpenses) && paginatedExpenses.length > 0 ? (
          <Paper>
            <TableContainer>
              <Table>
                <TableHead>
                  <TableRow>
                    <TableCell>Дата</TableCell>
                    <TableCell>Категория</TableCell>
                    <TableCell>Наименование</TableCell>
                    <TableCell>Описание</TableCell>
                    <TableCell align="right">Сумма</TableCell>
                    <TableCell align="center">Действия</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {paginatedExpenses.map((expense) => (
                    <TableRow key={expense.id}>
                      <TableCell>
                        {formatDate(expense.date)}
                      </TableCell>
                      <TableCell>
                        <Chip 
                          label={expenseCategories[expense.category]?.name || 'Прочее'} 
                          size="small"
                          style={{ 
                            backgroundColor: expenseCategories[expense.category]?.color || '#757575',
                            color: 'white' 
                          }}
                        />
                      </TableCell>
                      <TableCell>{expense.title}</TableCell>
                      <TableCell>{expense.description}</TableCell>
                      <TableCell align="right">{expense.amount.toLocaleString()} ₽</TableCell>
                      <TableCell align="center">
                        <IconButton 
                          size="small" 
                          color="primary"
                          onClick={() => handleEditExpense(expense)}
                        >
                          <EditIcon fontSize="small" />
                        </IconButton>
                        <IconButton 
                          size="small" 
                          color="error"
                          onClick={() => handleDeleteExpense(expense.id)}
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
              count={filteredExpenses.length}
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
              ? 'Расходы по вашему запросу не найдены' 
              : 'У вас пока нет добавленных расходов'}
          </Alert>
        )}

        <AddExpenseForm 
          open={openAddDialog} 
          handleClose={() => setOpenAddDialog(false)}
          onExpenseAdded={handleAddExpense}
        />

        <EditExpenseForm
          open={openEditDialog}
          handleClose={() => setOpenEditDialog(false)}
          onExpenseUpdated={handleExpenseUpdated}
          expense={currentExpense}
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
    console.error('Ошибка при рендеринге компонента ExpensesList:', error);
    return (
      <Box className="page-container">
        <Alert severity="error" sx={{ mb: 2 }}>
          Произошла ошибка при отображении страницы: {error.message || 'Неизвестная ошибка'}
        </Alert>
      </Box>
    );
  }
} 