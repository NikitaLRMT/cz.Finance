import React, { useState, useEffect } from 'react';
import { 
  Box, 
  Typography, 
  Paper, 
  Grid, 
  TextField, 
  Button, 
  IconButton, 
  Divider, 
  List, 
  ListItem, 
  ListItemText, 
  ListItemSecondaryAction, 
  Card, 
  CardContent, 
  CardHeader,
  InputAdornment,
  CircularProgress,
  Snackbar,
  Alert,
  Tooltip,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  SvgIcon
} from '@mui/material';
import { 
  CameraAlt as CameraAltIcon, 
  Delete as DeleteIcon, 
  Edit as EditIcon, 
  Add as AddIcon,
  Save as SaveIcon,
  RemoveCircle as RemoveCircleIcon
} from '@mui/icons-material';
import userService from '../services/user';
import Avatar from '../components/Avatar';

// Компонент иконки Telegram
const TelegramIcon = (props) => (
  <SvgIcon {...props} viewBox="0 0 1000 1000">
    <circle fill="#0088CC" cx="500" cy="500" r="500" />
    <path fill="#FFFFFF" d="M226.328419,494.722069 C372.088573,431.216685 469.284839,389.350049 517.917216,369.122161 C656.772535,311.36743 685.625481,301.334277 704.431427,301.003289 C708.567621,300.93067 717.815839,301.955722 723.806446,306.816968 C728.864797,310.92031 730.256552,316.46052 730.922551,320.357329 C731.588551,324.254139 732.417879,333.113018 731.758626,340.040661 C724.234007,419.102522 691.675292,610.964683 675.110952,699.515209 C668.10208,736.984878 654.301336,749.547465 640.940618,750.777479 C611.904684,753.448806 589.698552,731.458052 561.856476,713.158904 C517.117172,684.306305 492.293994,666.826643 448.669137,638.634524 C398.707288,606.11397 431.420903,588.117079 460.368587,558.329476 C467.713445,550.752113 591.360613,438.177924 593.709877,427.427302 C594.063058,425.478474 594.376003,417.940126 590.115686,414.143673 C585.855368,410.34722 579.68841,411.649557 575.11525,412.5936 C568.508897,413.9091 503.866839,457.688231 381.189198,544.929483 C360.363678,559.346517 341.55393,566.316028 324.75994,565.837519 C306.309602,565.311225 270.981465,555.553296 244.621375,547.252234 C212.270553,537.106332 186.598001,531.70654 189.05669,513.696279 C190.333088,504.411671 203.156872,494.900653 226.328419,494.722069 Z" />
  </SvgIcon>
);

export default function Profile() {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [editMode, setEditMode] = useState(false);
  const [editedUser, setEditedUser] = useState(null);
  const [notification, setNotification] = useState({ open: false, message: '', severity: 'success' });
  const [openWishlistDialog, setOpenWishlistDialog] = useState(false);
  const [newWishlistItem, setNewWishlistItem] = useState({ title: '', price: '', priority: 'medium' });
  const [fileInput, setFileInput] = useState(null);
  const [saving, setSaving] = useState(false);
  
  // Список желаний - теперь пустой массив по умолчанию
  const [wishlist, setWishlist] = useState([]);
  
  useEffect(() => {
    fetchUserData();
  }, []);
  
  const fetchUserData = async () => {
    try {
      setLoading(true);
      // Получаем данные пользователя
      const userData = await userService.getCurrentUser();
      console.log('Получены данные пользователя:', userData);
      setUser(userData);
      setEditedUser(userData);
      
      // Загружаем список желаний
      const wishlistData = await userService.getWishlist();
      console.log('Загружен список желаний:', wishlistData);
      setWishlist(wishlistData || []);
      
      setLoading(false);
    } catch (error) {
      console.error('Ошибка при загрузке данных пользователя:', error);
      setNotification({
        open: true,
        message: 'Не удалось загрузить данные пользователя',
        severity: 'error'
      });
      setLoading(false);
    }
  };
  
  const handleEditToggle = () => {
    if (editMode) {
      // Отменяем изменения, возвращаем пользователя к исходному состоянию
      setEditedUser(user);
    }
    setEditMode(!editMode);
  };
  
  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setEditedUser({
      ...editedUser,
      [name]: value
    });
  };
  
  const handleSaveChanges = async () => {
    try {
      setSaving(true);
      
      // Создаем копию объекта для изменений
      const userData = {
        first_name: editedUser.first_name,
        last_name: editedUser.last_name
      };
      
      // Отправляем данные на сервер
      const updatedUser = await userService.updateUser(userData);
      
      setUser(updatedUser);
      setEditMode(false);
      setSaving(false);
      
      setNotification({
        open: true,
        message: 'Изменения сохранены',
        severity: 'success'
      });
    } catch (error) {
      console.error('Ошибка при сохранении данных пользователя:', error);
      setSaving(false);
      
      let errorMessage = 'Не удалось сохранить изменения';
      
      // Проверяем, содержит ли объект ошибки сообщение от сервера
      if (error.response && error.response.data && error.response.data.message) {
        errorMessage = error.response.data.message;
      } else if (error.message) {
        errorMessage = error.message;
      }
      
      setNotification({
        open: true,
        message: errorMessage,
        severity: 'error'
      });
    }
  };
  
  const handleSaveFinancialGoals = async () => {
    try {
      setSaving(true);
      
      // Создаем объект только с финансовыми целями
      const goalsData = {
        monthly_limit: parseFloat(editedUser.monthly_limit) || 0,
        savings_goal: parseFloat(editedUser.savings_goal) || 0
      };
      
      // Используем метод updateUser
      const updatedUser = await userService.updateUser(goalsData);
      
      // Обновляем локальное состояние
      setUser({
        ...user,
        monthly_limit: goalsData.monthly_limit,
        savings_goal: goalsData.savings_goal
      });
      
      setSaving(false);
      
      setNotification({
        open: true,
        message: 'Финансовые цели обновлены',
        severity: 'success'
      });
    } catch (error) {
      console.error('Ошибка при сохранении финансовых целей:', error);
      setSaving(false);
      
      let errorMessage = 'Не удалось обновить финансовые цели';
      
      // Проверяем, содержит ли объект ошибки сообщение от сервера
      if (error.response && error.response.data && error.response.data.message) {
        errorMessage = error.response.data.message;
      } else if (error.message) {
        errorMessage = error.message;
      }
      
      setNotification({
        open: true,
        message: errorMessage,
        severity: 'error'
      });
    }
  };
  
  const handleAvatarChange = async (e) => {
    if (e.target.files && e.target.files[0]) {
      const file = e.target.files[0];
      
      try {
        setSaving(true);
        
        // Проверяем размер файла (максимум 5MB)
        if (file.size > 5 * 1024 * 1024) {
          setNotification({
            open: true,
            message: 'Размер файла не должен превышать 5MB',
            severity: 'error'
          });
          setSaving(false);
          return;
        }
        
        // Проверяем тип файла
        const validTypes = ['image/jpeg', 'image/png', 'image/gif'];
        if (!validTypes.includes(file.type)) {
          setNotification({
            open: true,
            message: 'Поддерживаются только изображения формата JPEG, PNG и GIF',
            severity: 'error'
          });
          setSaving(false);
          return;
        }
        
        // Загружаем файл на сервер
        const response = await userService.uploadAvatar(file);
        
        // Обновляем данные пользователя с новым аватаром
        setUser({
          ...user,
          avatar_url: response.avatar_url
        });
        
        setSaving(false);
        
        setNotification({
          open: true,
          message: 'Аватар успешно обновлен',
          severity: 'success'
        });
      } catch (error) {
        console.error('Ошибка при загрузке аватара:', error);
        setSaving(false);
        
        let errorMessage = 'Не удалось загрузить аватар';
        
        // Проверяем, содержит ли объект ошибки сообщение от сервера
        if (error.response && error.response.data && error.response.data.message) {
          errorMessage = error.response.data.message;
        } else if (error.message) {
          errorMessage = error.message;
        }
        
        setNotification({
          open: true,
          message: errorMessage,
          severity: 'error'
        });
      }
    }
  };
  
  const handleOpenFileInput = () => {
    if (fileInput) {
      fileInput.click();
    }
  };
  
  const handleRemoveAvatar = async () => {
    try {
      setSaving(true);
      
      // Отправляем запрос на удаление аватара
      const response = await userService.removeAvatar();
      
      setUser({
        ...user,
        avatar_url: response.avatar_url
      });
      
      setSaving(false);
      
      setNotification({
        open: true,
        message: 'Аватар успешно удален',
        severity: 'success'
      });
    } catch (error) {
      console.error('Ошибка при удалении аватара:', error);
      setSaving(false);
      setNotification({
        open: true,
        message: 'Не удалось удалить аватар',
        severity: 'error'
      });
    }
  };
  
  const handleWishlistItemChange = (e) => {
    const { name, value } = e.target;
    setNewWishlistItem({
      ...newWishlistItem,
      [name]: name === 'price' ? (value === '' ? '' : parseFloat(value)) : value
    });
  };
  
  const handleAddWishlistItem = async () => {
    try {
      // Добавляем элемент в список желаний
      const newItem = {
        title: newWishlistItem.title,
        price: parseFloat(newWishlistItem.price),
        priority: newWishlistItem.priority
      };
      
      const addedItem = await userService.addWishlistItem(newItem);
      
      // Обновляем локальный список
      setWishlist([...wishlist, addedItem]);
      setNewWishlistItem({ title: '', price: '', priority: 'medium' });
      setOpenWishlistDialog(false);
      
      setNotification({
        open: true,
        message: 'Элемент добавлен в список желаний',
        severity: 'success'
      });
    } catch (error) {
      console.error('Ошибка при добавлении элемента в список желаний:', error);
      setNotification({
        open: true,
        message: 'Не удалось добавить элемент',
        severity: 'error'
      });
    }
  };
  
  const handleRemoveWishlistItem = async (id) => {
    try {
      // Удаляем элемент на сервере
      await userService.removeWishlistItem(id);
      
      // Обновляем локальный список
      setWishlist(wishlist.filter(item => item.id !== id));
      
      setNotification({
        open: true,
        message: 'Элемент удален из списка желаний',
        severity: 'success'
      });
    } catch (error) {
      console.error('Ошибка при удалении элемента из списка желаний:', error);
      setNotification({
        open: true,
        message: 'Не удалось удалить элемент',
        severity: 'error'
      });
    }
  };
  
  const handleCloseNotification = () => {
    setNotification({ ...notification, open: false });
  };
  
  if (loading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '50vh' }}>
        <CircularProgress />
      </Box>
    );
  }
  
  return (
    <Box className="page-container">
      <Typography variant="h4" component="h1" gutterBottom>
        Профиль пользователя
      </Typography>
      
      <Grid container spacing={3}>
        {/* Основная информация о пользователе */}
        <Grid item xs={12} md={6}>
          <Paper sx={{ p: 3, position: 'relative' }}>
            <Box sx={{ position: 'absolute', top: 16, right: 16 }}>
              <Button 
                startIcon={editMode ? <SaveIcon /> : <EditIcon />}
                variant={editMode ? "contained" : "outlined"}
                color={editMode ? "primary" : "inherit"}
                onClick={editMode ? handleSaveChanges : handleEditToggle}
                disabled={saving}
              >
                {editMode ? 'Сохранить' : 'Редактировать'}
              </Button>
            </Box>
            
            <Box sx={{ display: 'flex', flexDirection: { xs: 'column', sm: 'row' }, alignItems: 'center', mb: 3 }}>
              <Box sx={{ position: 'relative', mr: { xs: 0, sm: 3 }, mb: { xs: 3, sm: 0 } }}>
                <Avatar
                  src={user.avatar_url}
                  alt={`${user.first_name} ${user.last_name}`}
                  sx={{ width: 100, height: 100 }}
                />
                {/* Отладочная информация о URL аватара */}
                {process.env.NODE_ENV === 'development' && (
                  <Box sx={{ position: 'absolute', top: 110, width: 300, left: -100, fontSize: '10px', color: 'gray' }}>
                    <div>URL: {user.avatar_url}</div>
                    <div>Origin: {window.location.origin}</div>
                    <div>Full URL: {`${window.location.origin}${user.avatar_url}`}</div>
                    <Button 
                      variant="outlined" 
                      size="small"
                      sx={{ mt: 1, fontSize: '8px' }}
                      onClick={() => window.open(`${window.location.origin}${user.avatar_url}`, '_blank')}
                    >
                      Открыть изображение
                    </Button>
                  </Box>
                )}
                {editMode && (
                  <Box sx={{ position: 'absolute', bottom: -10, right: -10 }}>
                    <input
                      type="file"
                      accept="image/*"
                      style={{ display: 'none' }}
                      onChange={handleAvatarChange}
                      ref={el => setFileInput(el)}
                    />
                    <Tooltip title="Изменить аватар">
                      <IconButton 
                        color="primary" 
                        aria-label="upload picture" 
                        component="span"
                        onClick={handleOpenFileInput}
                        disabled={saving}
                      >
                        <CameraAltIcon />
                      </IconButton>
                    </Tooltip>
                    <Tooltip title="Удалить аватар">
                      <IconButton 
                        color="error" 
                        aria-label="remove picture" 
                        component="span"
                        onClick={handleRemoveAvatar}
                        disabled={saving}
                      >
                        <DeleteIcon />
                      </IconButton>
                    </Tooltip>
                  </Box>
                )}
              </Box>
              
              <Box>
                <Typography variant="h5" gutterBottom>
                  {user.first_name} {user.last_name}
                </Typography>
                <Typography variant="body1" color="text.secondary">
                  @{user.username}
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  {user.email}
                </Typography>
                <Typography variant="body2" sx={{ mt: 1 }}>
                  <a href="https://t.me/czfinancebot" target="_blank" rel="noopener noreferrer" 
                    style={{ color: '#0088cc', textDecoration: 'none', display: 'flex', alignItems: 'center' }}>
                    <TelegramIcon sx={{ fontSize: 16, mr: 0.5, color: '#0088cc' }} />
                    Управляйте финансами через наш Telegram-бот
                  </a>
                </Typography>
              </Box>
            </Box>
            
            <Divider sx={{ my: 3 }} />
            
            {editMode ? (
              <Grid container spacing={2}>
                <Grid item xs={12} sm={6}>
                  <TextField
                    name="first_name"
                    label="Имя"
                    value={editedUser.first_name}
                    onChange={handleInputChange}
                    fullWidth
                    variant="outlined"
                    margin="normal"
                    disabled={saving}
                  />
                </Grid>
                <Grid item xs={12} sm={6}>
                  <TextField
                    name="last_name"
                    label="Фамилия"
                    value={editedUser.last_name}
                    onChange={handleInputChange}
                    fullWidth
                    variant="outlined"
                    margin="normal"
                    disabled={saving}
                  />
                </Grid>
                <Grid item xs={12}>
                  <TextField
                    name="email"
                    label="Email"
                    value={editedUser.email}
                    onChange={handleInputChange}
                    fullWidth
                    variant="outlined"
                    margin="normal"
                    type="email"
                    disabled={saving}
                  />
                </Grid>
                <Grid item xs={12}>
                  <TextField
                    name="username"
                    label="Имя пользователя"
                    value={editedUser.username}
                    onChange={handleInputChange}
                    fullWidth
                    variant="outlined"
                    margin="normal"
                    disabled={saving}
                  />
                </Grid>
              </Grid>
            ) : (
              <Box>
                <Typography variant="body1" gutterBottom>
                  <strong>Имя:</strong> {user.first_name}
                </Typography>
                <Typography variant="body1" gutterBottom>
                  <strong>Фамилия:</strong> {user.last_name}
                </Typography>
                <Typography variant="body1" gutterBottom>
                  <strong>Email:</strong> {user.email}
                </Typography>
                <Typography variant="body1" gutterBottom>
                  <strong>Имя пользователя:</strong> {user.username}
                </Typography>
                <Typography variant="body1" gutterBottom>
                  <strong>Telegram-бот:</strong> <a href="https://t.me/czfinancebot" target="_blank" rel="noopener noreferrer" 
                    style={{ color: '#0088cc', textDecoration: 'none', display: 'inline-flex', alignItems: 'center' }}>
                    <TelegramIcon sx={{ fontSize: 16, mr: 0.5, color: '#0088cc' }} />
                    @czfinancebot
                  </a>
                </Typography>
                <Typography variant="body2" color="text.secondary" sx={{ mt: 2 }}>
                  Аккаунт создан: {new Date(user.created_at).toLocaleDateString()}
                </Typography>
              </Box>
            )}
          </Paper>
        </Grid>
        
        {/* Финансовые цели */}
        <Grid item xs={12} md={6}>
          <Paper sx={{ p: 3 }}>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
              <Typography variant="h6">
                Финансовые цели
              </Typography>
              <Button 
                startIcon={<SaveIcon />}
                variant="contained"
                color="primary"
                onClick={handleSaveFinancialGoals}
                disabled={saving}
              >
                Сохранить
              </Button>
            </Box>
            
            <Grid container spacing={3}>
              <Grid item xs={12}>
                <TextField
                  name="monthly_limit"
                  label="Месячный лимит расходов"
                  value={editedUser.monthly_limit}
                  onChange={handleInputChange}
                  fullWidth
                  variant="outlined"
                  type="number"
                  InputProps={{
                    startAdornment: <InputAdornment position="start">₽</InputAdornment>,
                  }}
                  disabled={saving}
                />
                <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
                  Установите максимальную сумму, которую вы планируете тратить в месяц
                </Typography>
              </Grid>
              
              <Grid item xs={12}>
                <TextField
                  name="savings_goal"
                  label="Цель по накоплениям в месяц"
                  value={editedUser.savings_goal}
                  onChange={handleInputChange}
                  fullWidth
                  variant="outlined"
                  type="number"
                  InputProps={{
                    startAdornment: <InputAdornment position="start">₽</InputAdornment>,
                  }}
                  disabled={saving}
                />
                <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
                  Установите сумму, которую вы планируете откладывать каждый месяц
                </Typography>
              </Grid>
            </Grid>
          </Paper>
        </Grid>
        
        {/* Список желаний */}
        <Grid item xs={12}>
          <Card>
            <CardHeader 
              title="Список желаний" 
              action={
                <Button 
                  startIcon={<AddIcon />}
                  variant="outlined"
                  onClick={() => setOpenWishlistDialog(true)}
                >
                  Добавить
                </Button>
              }
            />
            <Divider />
            <CardContent>
              {wishlist.length > 0 ? (
                <List>
                  {wishlist.map((item) => (
                    <ListItem key={item.id} divider sx={{ py: 2 }}>
                      <ListItemText
                        primary={item.title}
                        secondary={
                          <Box sx={{ display: 'flex', alignItems: 'center', mt: 1 }}>
                            <Typography
                              variant="body2"
                              sx={{
                                backgroundColor: 
                                  item.priority === 'high' ? '#ffcdd2' :
                                  item.priority === 'medium' ? '#fff9c4' :
                                  '#c8e6c9',
                                color: 
                                  item.priority === 'high' ? '#c62828' : 
                                  item.priority === 'medium' ? '#f57f17' : 
                                  '#2e7d32',
                                px: 1,
                                py: 0.5,
                                borderRadius: 1,
                                mr: 2,
                                fontSize: '0.7rem'
                              }}
                            >
                              {item.priority === 'high' ? 'Высокий' : 
                               item.priority === 'medium' ? 'Средний' : 
                               'Низкий'} приоритет
                            </Typography>
                            <Typography variant="body2" color="text.secondary">
                              Примерная стоимость: {item.price.toLocaleString()} ₽
                            </Typography>
                          </Box>
                        }
                      />
                      <ListItemSecondaryAction>
                        <IconButton 
                          edge="end" 
                          aria-label="delete"
                          onClick={() => handleRemoveWishlistItem(item.id)}
                        >
                          <RemoveCircleIcon color="error" />
                        </IconButton>
                      </ListItemSecondaryAction>
                    </ListItem>
                  ))}
                </List>
              ) : (
                <Typography variant="body1" color="text.secondary" align="center" sx={{ py: 3 }}>
                  Ваш список желаний пуст. Добавьте что-то, о чем вы мечтаете!
                </Typography>
              )}
            </CardContent>
          </Card>
        </Grid>
      </Grid>
      
      {/* Диалог добавления элемента в список желаний */}
      <Dialog open={openWishlistDialog} onClose={() => setOpenWishlistDialog(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Добавить в список желаний</DialogTitle>
        <DialogContent>
          <Grid container spacing={2} sx={{ mt: 1 }}>
            <Grid item xs={12}>
              <TextField
                name="title"
                label="Название"
                value={newWishlistItem.title}
                onChange={handleWishlistItemChange}
                fullWidth
                required
              />
            </Grid>
            
            <Grid item xs={12}>
              <TextField
                name="price"
                label="Примерная стоимость"
                value={newWishlistItem.price}
                onChange={handleWishlistItemChange}
                fullWidth
                type="number"
                InputProps={{
                  startAdornment: <InputAdornment position="start">₽</InputAdornment>,
                }}
              />
            </Grid>
            
            <Grid item xs={12}>
              <TextField
                select
                name="priority"
                label="Приоритет"
                value={newWishlistItem.priority}
                onChange={handleWishlistItemChange}
                fullWidth
                SelectProps={{
                  native: true,
                }}
              >
                <option value="high">Высокий</option>
                <option value="medium">Средний</option>
                <option value="low">Низкий</option>
              </TextField>
            </Grid>
          </Grid>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenWishlistDialog(false)} color="inherit">
            Отмена
          </Button>
          <Button 
            onClick={handleAddWishlistItem} 
            color="primary" 
            variant="contained"
            disabled={!newWishlistItem.title}
          >
            Добавить
          </Button>
        </DialogActions>
      </Dialog>
      
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
} 