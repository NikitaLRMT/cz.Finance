import React from 'react';
import { Outlet } from 'react-router-dom';
import { Box, Container, Paper, Typography } from '@mui/material';

export default function AuthLayout() {
  return (
    <Box className="auth-form-container">
      <Container maxWidth="sm">
        <Paper elevation={6} sx={{ p: 4, borderRadius: 2 }}>
          <Box sx={{ mb: 3, textAlign: 'center' }}>
            <Typography variant="h4" component="h1" gutterBottom>
              Finance App
            </Typography>
            <Typography variant="subtitle1" color="text.secondary">
              Управление личными финансами
            </Typography>
          </Box>
          <Outlet />
        </Paper>
      </Container>
    </Box>
  );
} 