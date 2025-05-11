import React, { useState, useEffect } from 'react';
import { 
  Button, 
  TextField, 
  Dialog, 
  DialogActions, 
  DialogContent, 
  DialogTitle, 
  FormControl, 
  InputLabel, 
  Select, 
  MenuItem,
  InputAdornment,
  Grid,
  FormHelperText
} from '@mui/material';
import { DatePicker } from '@mui/x-date-pickers/DatePicker';
import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import expensesService, { expenseCategories } from '../../services/expenses';
import { ru } from 'date-fns/locale';

const AddExpenseForm = ({ open, handleClose, onExpenseAdded }) => {
  const [formData, setFormData] = useState({
    amount: '',
    description: '',
    date: new Date(),
    category: '',
    notes: ''
  });
  
  const [errors, setErrors] = useState({});
  const [categories, setCategories] = useState([]);
  const [isSubmitting, setIsSubmitting] = useState(false);
  
  useEffect(() => {
    // Получаем категории из сервиса
    const categoriesArray = expensesService.getCategoriesArray();
    setCategories(categoriesArray);
  }, []);
  
  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData({
      ...formData,
      [name]: value
    });
    
    // Очищаем ошибку поля при изменении
    if (errors[name]) {
      setErrors({
        ...errors,
        [name]: null
      });
    }
  };
  
  const handleDateChange = (date) => {
    setFormData({
      ...formData,
      date
    });
  };
  
  const validate = () => {
    const newErrors = {};
    
    if (!formData.amount || formData.amount <= 0) {
      newErrors.amount = 'Введите корректную сумму';
    }
    
    if (!formData.category) {
      newErrors.category = 'Выберите категорию';
    }
    
    if (!formData.description.trim()) {
      newErrors.description = 'Введите описание';
    }
    
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };
  
  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!validate()) {
      return;
    }
    
    setIsSubmitting(true);
    
    try {
      const response = await expensesService.createExpense(formData);
      
      setIsSubmitting(false);
      handleClose();
      if (onExpenseAdded) {
        onExpenseAdded(response);
      }
      
      // Очищаем форму
      setFormData({
        amount: '',
        description: '',
        date: new Date(),
        category: '',
        notes: ''
      });
    } catch (error) {
      setIsSubmitting(false);
      console.error('Ошибка при добавлении расхода:', error);
      setErrors({
        ...errors,
        submit: 'Ошибка при добавлении расхода. Пожалуйста, попробуйте снова.'
      });
    }
  };
  
  return (
    <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
      <DialogTitle>Добавить новый расход</DialogTitle>
      <DialogContent>
        <Grid container spacing={2} sx={{ mt: 1 }}>
          <Grid item xs={12} sm={6}>
            <TextField
              name="amount"
              label="Сумма"
              value={formData.amount}
              onChange={handleChange}
              fullWidth
              required
              type="number"
              InputProps={{
                startAdornment: <InputAdornment position="start">₽</InputAdornment>,
              }}
              error={!!errors.amount}
              helperText={errors.amount}
            />
          </Grid>
          
          <Grid item xs={12} sm={6}>
            <FormControl fullWidth required error={!!errors.category}>
              <InputLabel id="category-label">Категория</InputLabel>
              <Select
                labelId="category-label"
                name="category"
                value={formData.category}
                onChange={handleChange}
                label="Категория"
              >
                {categories.map((category) => (
                  <MenuItem key={category.id} value={category.id}>
                    {category.name}
                  </MenuItem>
                ))}
              </Select>
              {errors.category && (
                <FormHelperText>{errors.category}</FormHelperText>
              )}
            </FormControl>
          </Grid>
          
          <Grid item xs={12}>
            <TextField
              name="description"
              label="Наименование"
              value={formData.description}
              onChange={handleChange}
              fullWidth
              required
              error={!!errors.description}
              helperText={errors.description}
            />
          </Grid>
          
          <Grid item xs={12}>
            <TextField
              name="notes"
              label="Дополнительное описание"
              value={formData.notes}
              onChange={handleChange}
              fullWidth
              multiline
              rows={2}
            />
          </Grid>
          
          <Grid item xs={12}>
            <LocalizationProvider dateAdapter={AdapterDateFns} adapterLocale={ru}>
              <DatePicker
                label="Дата"
                value={formData.date}
                onChange={handleDateChange}
                slotProps={{ textField: { fullWidth: true } }}
                maxDate={new Date()}
              />
            </LocalizationProvider>
          </Grid>
          
          {errors.submit && (
            <Grid item xs={12}>
              <FormHelperText error>{errors.submit}</FormHelperText>
            </Grid>
          )}
        </Grid>
      </DialogContent>
      <DialogActions>
        <Button onClick={handleClose} color="inherit">
          Отмена
        </Button>
        <Button 
          onClick={handleSubmit} 
          color="primary" 
          variant="contained"
          disabled={isSubmitting}
        >
          {isSubmitting ? 'Сохранение...' : 'Сохранить'}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default AddExpenseForm; 