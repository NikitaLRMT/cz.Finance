# Используем многоступенчатую сборку для оптимизации размера образа
FROM node:18-alpine as builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем только package.json и package-lock.json для эффективного кеширования
COPY ./frontend/package*.json ./

# Устанавливаем зависимости
RUN npm install

# Затем копируем остальной код фронтенда
COPY ./frontend/ ./

# Собираем приложение
RUN npm run build

# Используем Nginx для раздачи собранного приложения
FROM nginx:1.24-alpine

# Копируем собранные файлы из стадии builder
COPY --from=builder /app/build /usr/share/nginx/html

# Копируем кастомную конфигурацию Nginx
COPY ./docker/nginx.conf /etc/nginx/conf.d/default.conf

# Создаем директорию для загруженных файлов
RUN mkdir -p /usr/share/nginx/html/uploads

# Expose порта
EXPOSE 80

# Запускаем Nginx
CMD ["nginx", "-g", "daemon off;"] 