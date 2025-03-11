FROM node:20-alpine

WORKDIR /app

COPY package*.json ./

RUN npm install

# Копируем файлы проекта
COPY . .

RUN npm run build

# Устанавливаем nginx для раздачи статики
FROM nginx:alpine

# Копируем собранное приложение в директорию nginx
COPY --from=0 /app/dist /usr/share/nginx/html

# Открываем порт 80
EXPOSE 80

# Запускаем nginx
CMD ["nginx", "-g", "daemon off;"] 