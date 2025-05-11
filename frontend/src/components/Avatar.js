import React, { useState, useEffect } from 'react';
import { Avatar as MuiAvatar } from '@mui/material';

/**
 * Компонент Avatar для корректного отображения аватаров
 * Обрабатывает относительные и абсолютные URL
 */
const Avatar = ({ src, alt, ...props }) => {
  const [imageError, setImageError] = useState(false);
  const [imageSrc, setImageSrc] = useState('');
  
  useEffect(() => {
    // Обработка источника изображения
    if (!src || src === '') {
      setImageSrc('/default-avatar.png');
      return;
    }
    
    // Добавляем временную метку для предотвращения кэширования
    const timestamp = new Date().getTime();
    
    if (src.startsWith('http')) {
      setImageSrc(`${src}?t=${timestamp}`);
    } else if (src.startsWith('/uploads/')) {
      setImageSrc(`${src}?t=${timestamp}`);
      console.log('Загружаем аватар по пути:', `${src}?t=${timestamp}`);
      
      // Проверяем доступность файла через fetch
      fetch(`${window.location.origin}${src}?t=${timestamp}`)
        .then(response => {
          if (!response.ok) {
            console.error('Ошибка загрузки аватара:', response.status);
            setImageError(true);
          } else {
            console.log('Аватар доступен');
            setImageError(false);
          }
        })
        .catch(error => {
          console.error('Ошибка при запросе аватара:', error);
          setImageError(true);
        });
    } else {
      setImageSrc(`${process.env.REACT_APP_API_URL || window.location.origin}${src}?t=${timestamp}`);
    }
  }, [src]);
  
  // Обработка ошибки загрузки
  const handleError = () => {
    console.error('Ошибка загрузки изображения:', imageSrc);
    setImageError(true);
  };
  
  return (
    <MuiAvatar
      src={imageError ? '/default-avatar.png' : imageSrc}
      alt={alt}
      onError={handleError}
      {...props}
    />
  );
};

export default Avatar; 