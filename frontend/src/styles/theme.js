import { createTheme } from '@mui/material/styles';

const theme = createTheme({
  palette: {
    mode: 'dark',
    primary: {
      main: '#3861FB', // CoinMarketCap синий
      light: '#6C87FF',
      dark: '#0035C7',
      contrastText: '#ffffff',
    },
    secondary: {
      main: '#16C784', // CoinMarketCap зеленый
      light: '#58fbb8',
      dark: '#009653',
      contrastText: '#ffffff',
    },
    error: {
      main: '#EA3943', // CoinMarketCap красный
      light: '#ff6f6f',
      dark: '#b0001a',
    },
    warning: {
      main: '#FFB528', // CoinMarketCap оранжево-желтый
      light: '#ffe75a',
      dark: '#c78500',
    },
    info: {
      main: '#A6B0C3', // CoinMarketCap серо-голубой
      light: '#d7e0f7',
      dark: '#778293',
    },
    success: {
      main: '#16C784', // CoinMarketCap зеленый
      light: '#58fbb8',
      dark: '#009653',
    },
    background: {
      default: '#171924', // CoinMarketCap тёмный фон
      paper: '#222531', // CoinMarketCap цвет карточек в темной теме
      dark: '#13151F', // CoinMarketCap более темный фон
      light: '#2B2F3E', // CoinMarketCap светлый фон для тёмной темы
      card: '#222531', // Цвет карточек в темной теме
    },
    text: {
      primary: '#F8FAFD', // Основной текст в темной теме
      secondary: '#A6B0C3', // Вторичный текст в темной теме
      disabled: '#616E85',
    },
    divider: '#2B2F3E', // Разделители в темной теме
  },
  typography: {
    fontFamily: [
      'Inter',
      'Roboto',
      'Arial',
      'sans-serif',
    ].join(','),
    h1: {
      fontSize: '2.5rem',
      fontWeight: 600,
      lineHeight: 1.2,
    },
    h2: {
      fontSize: '2rem',
      fontWeight: 600,
      lineHeight: 1.3,
    },
    h3: {
      fontSize: '1.75rem',
      fontWeight: 600,
      lineHeight: 1.3,
    },
    h4: {
      fontSize: '1.5rem',
      fontWeight: 600,
      lineHeight: 1.4,
    },
    h5: {
      fontSize: '1.25rem',
      fontWeight: 500,
      lineHeight: 1.4,
    },
    h6: {
      fontSize: '1rem',
      fontWeight: 500,
      lineHeight: 1.5,
    },
    body1: {
      fontSize: '1rem',
      lineHeight: 1.5,
    },
    body2: {
      fontSize: '0.875rem',
      lineHeight: 1.6,
    },
    button: {
      fontWeight: 500,
      fontSize: '0.875rem',
      textTransform: 'none',
    },
  },
  shape: {
    borderRadius: 8,
  },
  components: {
    MuiCssBaseline: {
      styleOverrides: {
        body: {
          backgroundColor: '#171924',
          scrollbarWidth: 'thin',
          '&::-webkit-scrollbar': {
            width: '8px',
            height: '8px',
          },
          '&::-webkit-scrollbar-track': {
            background: '#13151F',
          },
          '&::-webkit-scrollbar-thumb': {
            background: '#2B2F3E',
            borderRadius: '4px',
          },
          '&::-webkit-scrollbar-thumb:hover': {
            background: '#616E85',
          },
        },
      },
    },
    MuiButton: {
      styleOverrides: {
        root: {
          textTransform: 'none',
          borderRadius: 8,
          fontWeight: 500,
          boxShadow: 'none',
          '&:hover': {
            boxShadow: '0 4px 12px rgba(0, 0, 0, 0.2)',
          },
        },
        contained: {
          '&:hover': {
            boxShadow: '0 4px 12px rgba(0, 0, 0, 0.25)',
          },
        },
        containedPrimary: {
          '&:hover': {
            backgroundColor: '#4471FF',
            boxShadow: '0 4px 12px rgba(56, 97, 251, 0.35)',
          },
        },
        containedSecondary: {
          '&:hover': {
            backgroundColor: '#1DDF96',
            boxShadow: '0 4px 12px rgba(22, 199, 132, 0.35)',
          },
        },
      },
    },
    MuiCard: {
      styleOverrides: {
        root: {
          borderRadius: 12,
          boxShadow: '0 2px 8px rgba(0, 0, 0, 0.2)',
          border: '1px solid #2B2F3E',
          transition: 'transform 0.2s ease-in-out, box-shadow 0.2s ease-in-out',
          backgroundColor: '#222531',
          '&:hover': {
            boxShadow: '0 4px 16px rgba(0, 0, 0, 0.3)',
            transform: 'translateY(-2px)',
          },
        },
      },
    },
    MuiPaper: {
      styleOverrides: {
        root: {
          boxShadow: '0 2px 8px rgba(0, 0, 0, 0.2)',
          backgroundImage: 'none',
          backgroundColor: '#222531',
        },
        outlined: {
          border: '1px solid #2B2F3E',
        },
      },
    },
    MuiAppBar: {
      styleOverrides: {
        root: {
          boxShadow: '0 1px 2px rgba(0, 0, 0, 0.2)',
          backgroundImage: 'none',
          backgroundColor: '#171924',
        },
      },
    },
    MuiTableCell: {
      styleOverrides: {
        root: {
          borderBottom: '1px solid #2B2F3E',
          padding: '16px',
        },
        head: {
          fontWeight: 600,
          backgroundColor: '#13151F',
          color: '#A6B0C3',
        },
      },
    },
    MuiChip: {
      styleOverrides: {
        root: {
          fontWeight: 500,
        },
        colorPrimary: {
          backgroundColor: 'rgba(56, 97, 251, 0.15)',
          color: '#6C87FF',
        },
        colorSecondary: {
          backgroundColor: 'rgba(22, 199, 132, 0.15)',
          color: '#58fbb8',
        },
      },
    },
    MuiListItemButton: {
      styleOverrides: {
        root: {
          borderRadius: 8,
          '&.Mui-selected': {
            backgroundColor: 'rgba(56, 97, 251, 0.2)',
            color: '#6C87FF',
            '&:hover': {
              backgroundColor: 'rgba(56, 97, 251, 0.25)',
            },
            '& .MuiListItemIcon-root': {
              color: '#6C87FF',
            },
          },
          '&:hover': {
            backgroundColor: 'rgba(43, 47, 62, 0.8)',
          },
        },
      },
    },
    MuiDrawer: {
      styleOverrides: {
        paper: {
          backgroundColor: '#171924',
          borderRight: '1px solid #2B2F3E',
        },
      },
    },
    MuiDivider: {
      styleOverrides: {
        root: {
          borderColor: '#2B2F3E',
        },
      },
    },
    MuiAvatar: {
      styleOverrides: {
        root: {
          backgroundColor: 'rgba(56, 97, 251, 0.2)',
          color: '#6C87FF',
        },
      },
    },
    MuiTooltip: {
      styleOverrides: {
        tooltip: {
          backgroundColor: '#13151F',
          fontSize: '0.75rem',
          padding: '8px 12px',
          borderRadius: 6,
        },
      },
    },
    MuiInputBase: {
      styleOverrides: {
        root: {
          backgroundColor: '#2B2F3E',
          color: '#F8FAFD',
          '&.Mui-focused': {
            backgroundColor: '#2B2F3E',
          },
          '&:hover': {
            backgroundColor: '#333846',
          },
        },
      },
    },
    MuiOutlinedInput: {
      styleOverrides: {
        root: {
          '& .MuiOutlinedInput-notchedOutline': {
            borderColor: '#2B2F3E',
          },
          '&:hover .MuiOutlinedInput-notchedOutline': {
            borderColor: '#616E85',
          },
          '&.Mui-focused .MuiOutlinedInput-notchedOutline': {
            borderColor: '#3861FB',
          },
        },
      },
    },
  },
});

export default theme; 