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
import incomesService, { incomeSources } from '../../services/incomes';
import { ru } from 'date-fns/locale';

const EditIncomeForm = ({ open, handleClose, onIncomeUpdated, income }) => {
  const [formData, setFormData] = useState({
    amount: '',
    description: '',
    date: new Date(),
    category: '', // используется для source на бекенде
  });
  
  const [errors, setErrors] = useState({});
  const [sources, setSources] = useState([]);
  const [isSubmitting, setIsSubmitting] = useState(false);
  
  useEffect(() => {
    // Получаем источники из сервиса
    const sourcesArray = incomesService.getSourcesArray();
    setSources(sourcesArray);
    
    // Если передан доход для редактирования, заполняем форму его данными
    if (income) {
      setFormData({
        amount: income.amount.toString(),
        description: income.description || '',
        date: new Date(income.date),
        category: income.source,  // source в income соответствует category в форме
      });
    }
  }, [income]);
  
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
      newErrors.category = 'Выберите источник дохода';
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
      const response = await incomesService.updateIncome(income.id, formData);
      
      setIsSubmitting(false);
      handleClose();
      if (onIncomeUpdated) {
        onIncomeUpdated(response);
      }
    } catch (error) {
      setIsSubmitting(false);
      console.error('Ошибка при обновлении дохода:', error);
      setErrors({
        ...errors,
        submit: 'Ошибка при обновлении дохода. Пожалуйста, попробуйте снова.'
      });
    }
  };
  
  return (
    <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
      <DialogTitle>Редактировать доход</DialogTitle>
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
              <InputLabel id="category-label">Источник дохода</InputLabel>
              <Select
                labelId="category-label"
                name="category"
                value={formData.category}
                onChange={handleChange}
                label="Источник дохода"
              >
                {sources.map((source) => (
                  <MenuItem key={source.id} value={source.id}>
                    {source.name}
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
              label="Описание"
              value={formData.description}
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

export default EditIncomeForm; 