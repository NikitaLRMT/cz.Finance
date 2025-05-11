import React from 'react';
import { useNavigate } from 'react-router-dom';
import { Box, Button, Typography, Container } from '@mui/material';
import { Error as ErrorIcon } from '@mui/icons-material';

export default function NotFound() {
  const navigate = useNavigate();

  return (
    <Container>
      <Box
        sx={{
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          justifyContent: 'center',
          minHeight: '100vh',
          textAlign: 'center',
        }}
      >
        <ErrorIcon sx={{ fontSize: 100, color: 'error.main', mb: 2 }} />
        <Typography variant="h1" component="h1" gutterBottom>
          404
        </Typography>
        <Typography variant="h4" component="h2" gutterBottom>
          Страница не найдена
        </Typography>
        <Typography variant="body1" color="text.secondary" paragraph>
          Запрашиваемая страница не существует или была удалена.
        </Typography>
        <Button
          variant="contained"
          color="primary"
          onClick={() => navigate('/')}
          sx={{ mt: 2 }}
        >
          На главную
        </Button>
      </Box>
    </Container>
  );
} 