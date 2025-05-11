import React, { useState } from 'react';
import { Outlet, useLocation } from 'react-router-dom';
import {
  AppBar,
  Box,
  CssBaseline,
  Divider,
  Drawer,
  IconButton,
  List,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Toolbar,
  Typography,
  Button,
  Avatar,
  Menu,
  MenuItem,
  Badge,
  Chip,
  useMediaQuery,
  useTheme,
  Container
} from '@mui/material';
import {
  Menu as MenuIcon,
  Dashboard as DashboardIcon,
  AttachMoney as ExpenseIcon,
  Savings as IncomeIcon,
  Calculate as CalculatorIcon,
  Person as ProfileIcon,
  Logout as LogoutIcon,
  Notifications as NotificationsIcon,
  AccountBalanceWallet as WalletIcon,
  BarChart as BarChartIcon,
  Settings as SettingsIcon,
  Search as SearchIcon
} from '@mui/icons-material';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';

const drawerWidth = 240;

export default function MainLayout() {
  const [mobileOpen, setMobileOpen] = useState(false);
  const [anchorEl, setAnchorEl] = useState(null);
  const { currentUser, logout } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('md'));
  
  // Активный путь для выделения текущего элемента меню
  const currentPath = location.pathname;

  const handleDrawerToggle = () => {
    setMobileOpen(!mobileOpen);
  };

  const handleProfileMenuOpen = (event) => {
    setAnchorEl(event.currentTarget);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
  };

  const handleLogout = () => {
    handleMenuClose();
    logout();
    navigate('/login');
  };

  const navigateTo = (path) => {
    navigate(path);
    setMobileOpen(false);
  };

  const menuItems = [
    { text: 'Панель управления', icon: <DashboardIcon />, path: '/dashboard' },
    { text: 'Расходы', icon: <ExpenseIcon />, path: '/expenses' },
    { text: 'Доходы', icon: <IncomeIcon />, path: '/incomes' },
    { text: 'Калькуляторы', icon: <CalculatorIcon />, path: '/calculators' },
  ];

  const drawer = (
    <Box sx={{ overflow: 'auto' }}>
      <Box
        sx={{
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          py: 2,
          px: 1,
        }}
      >
        <Typography 
          variant="h6" 
          component="div" 
          sx={{ 
            fontWeight: 'bold', 
            color: 'primary.main',
            display: 'flex',
            alignItems: 'center',
            mb: 2
          }}
        >
          <WalletIcon sx={{ mr: 1 }} />
          Finance App
        </Typography>
      </Box>
      <Divider />
      <List>
        {menuItems.map((item) => (
          <ListItem key={item.text} disablePadding>
            <ListItemButton 
              onClick={() => navigateTo(item.path)}
              selected={currentPath === item.path}
              sx={{ 
                mx: 1, 
                borderRadius: '8px',
                my: 0.5,
                '&.Mui-selected': {
                  backgroundColor: 'primary.main',
                  color: 'white',
                  '&:hover': {
                    backgroundColor: 'primary.dark',
                  },
                  '& .MuiListItemIcon-root': {
                    color: 'white',
                  },
                },
              }}
            >
              <ListItemIcon 
                sx={{ 
                  minWidth: 40,
                  color: currentPath === item.path ? 'white' : 'inherit'
                }}
              >
                {item.icon}
              </ListItemIcon>
              <ListItemText primary={item.text} />
              {item.path === '/dashboard' && (
                <Chip 
                  label="Главная" 
                  size="small" 
                  sx={{ 
                    height: 20, 
                    fontSize: '0.65rem',
                    backgroundColor: currentPath === item.path ? 'rgba(255,255,255,0.2)' : 'primary.light',
                    color: currentPath === item.path ? 'white' : 'primary.contrastText',
                  }} 
                />
              )}
            </ListItemButton>
          </ListItem>
        ))}
      </List>
      <Divider sx={{ my: 2 }} />
      <List>
        <ListItem disablePadding>
          <ListItemButton 
            onClick={() => navigateTo('/profile')}
            selected={currentPath === '/profile'}
            sx={{ 
              mx: 1, 
              borderRadius: '8px',
              my: 0.5,
              '&.Mui-selected': {
                backgroundColor: 'primary.main',
                color: 'white',
                '&:hover': {
                  backgroundColor: 'primary.dark',
                },
                '& .MuiListItemIcon-root': {
                  color: 'white',
                },
              },
            }}
          >
            <ListItemIcon 
              sx={{ 
                minWidth: 40,
                color: currentPath === '/profile' ? 'white' : 'inherit'
              }}
            >
              <ProfileIcon />
            </ListItemIcon>
            <ListItemText primary="Профиль" />
          </ListItemButton>
        </ListItem>
        <ListItem disablePadding>
          <ListItemButton 
            onClick={handleLogout}
            sx={{ 
              mx: 1, 
              borderRadius: '8px',
              my: 0.5,
              color: 'error.main',
              '&:hover': {
                backgroundColor: 'error.light',
                color: 'white',
                '& .MuiListItemIcon-root': {
                  color: 'white',
                },
              },
            }}
          >
            <ListItemIcon sx={{ minWidth: 40, color: 'error.main' }}>
              <LogoutIcon />
            </ListItemIcon>
            <ListItemText primary="Выйти" />
          </ListItemButton>
        </ListItem>
      </List>
    </Box>
  );

  return (
    <Box sx={{ display: 'flex' }}>
      <CssBaseline />
      <AppBar 
        position="fixed" 
        sx={{ 
          zIndex: (theme) => theme.zIndex.drawer + 1,
          backgroundColor: 'background.paper',
          color: 'text.primary',
          boxShadow: '0px 1px 4px rgba(0, 0, 0, 0.05)',
        }}
      >
        <Toolbar>
          <IconButton
            color="inherit"
            aria-label="open drawer"
            edge="start"
            onClick={handleDrawerToggle}
            sx={{ mr: 2, display: { sm: 'none' } }}
          >
            <MenuIcon />
          </IconButton>
          <Typography 
            variant="h6" 
            noWrap 
            component="div" 
            sx={{ 
              display: { xs: 'none', sm: 'flex' }, 
              fontWeight: 'bold',
              color: 'primary.main',
              alignItems: 'center'
            }}
          >
            <WalletIcon sx={{ mr: 1 }} />
            Finance App
          </Typography>
          
          <Box sx={{ flexGrow: 1 }} />
          
          <Box sx={{ display: 'flex', alignItems: 'center' }}>
            {/* Профиль пользователя */}
            <Button 
              onClick={handleProfileMenuOpen}
              sx={{ 
                borderRadius: 2,
                bgcolor: 'background.light',
                py: 0.5,
                '&:hover': {
                  bgcolor: 'background.light',
                }
              }}
              startIcon={
                <Avatar 
                  sx={{ width: 32, height: 32 }}
                  alt={currentUser?.username || 'User'}
                  src={currentUser?.avatar_url || "/default-avatar.png"}
                />
              }
            >
              <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'flex-start', ml: 0.5 }}>
                <Typography variant="body2" fontWeight="medium">
                  {currentUser?.username || 'Пользователь'}
                </Typography>
                <Typography variant="caption" color="text.secondary" sx={{ lineHeight: 1 }}>
                  {currentUser?.email ? currentUser.email.split('@')[0] : 'user'}
                </Typography>
              </Box>
            </Button>
            <Menu
              anchorEl={anchorEl}
              open={Boolean(anchorEl)}
              onClose={handleMenuClose}
              PaperProps={{
                sx: { width: 220, mt: 1.5 }
              }}
            >
              <Box sx={{ py: 1, px: 2, mb: 1 }}>
                <Typography variant="subtitle2" fontWeight="bold">
                  {currentUser?.first_name} {currentUser?.last_name}
                </Typography>
                <Typography variant="caption" color="text.secondary">
                  {currentUser?.email}
                </Typography>
              </Box>
              <Divider />
              <MenuItem onClick={() => { 
                handleMenuClose();
                navigate('/profile');
              }}>
                <ListItemIcon>
                  <ProfileIcon fontSize="small" />
                </ListItemIcon>
                <ListItemText>Профиль</ListItemText>
              </MenuItem>
              <MenuItem onClick={() => { 
                handleMenuClose();
                navigate('/dashboard');
              }}>
                <ListItemIcon>
                  <DashboardIcon fontSize="small" />
                </ListItemIcon>
                <ListItemText>Дашборд</ListItemText>
              </MenuItem>
              <MenuItem onClick={() => { 
                handleMenuClose();
              }}>
                <ListItemIcon>
                  <SettingsIcon fontSize="small" />
                </ListItemIcon>
                <ListItemText>Настройки</ListItemText>
              </MenuItem>
              <Divider />
              <MenuItem onClick={handleLogout}>
                <ListItemIcon>
                  <LogoutIcon fontSize="small" sx={{ color: 'error.main' }} />
                </ListItemIcon>
                <ListItemText sx={{ color: 'error.main' }}>Выйти</ListItemText>
              </MenuItem>
            </Menu>
          </Box>
        </Toolbar>
      </AppBar>
      <Box
        component="nav"
        sx={{ width: { sm: drawerWidth }, flexShrink: { sm: 0 } }}
      >
        <Drawer
          variant="temporary"
          open={mobileOpen}
          onClose={handleDrawerToggle}
          ModalProps={{
            keepMounted: true, // Лучше для мобильной производительности
          }}
          sx={{
            display: { xs: 'block', sm: 'none' },
            '& .MuiDrawer-paper': { boxSizing: 'border-box', width: drawerWidth },
          }}
        >
          {drawer}
        </Drawer>
        <Drawer
          variant="permanent"
          sx={{
            display: { xs: 'none', sm: 'block' },
            '& .MuiDrawer-paper': { 
              boxSizing: 'border-box', 
              width: drawerWidth, 
              borderRight: '1px solid #EFF2F5',
            },
          }}
          open
        >
          {drawer}
        </Drawer>
      </Box>
      <Box
        component="main"
        sx={{ 
          flexGrow: 1, 
          pt: 3, 
          px: 2, 
          width: { sm: `calc(100% - ${drawerWidth}px)` },
          backgroundColor: 'background.default',
          minHeight: '100vh'
        }}
      >
        <Toolbar />
        <Container maxWidth="lg">
          <Outlet />
        </Container>
      </Box>
    </Box>
  );
} 