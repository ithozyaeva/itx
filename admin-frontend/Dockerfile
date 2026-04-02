FROM node:20-alpine

WORKDIR /app

ARG VITE_TELEGRAM_BOT_NAME
ENV VITE_TELEGRAM_BOT_NAME=${VITE_TELEGRAM_BOT_NAME}

COPY package*.json ./

RUN npm install

COPY . .

RUN npm run build

# Создаем директорию для статики
RUN mkdir -p /app/dist

# Очищаем старую статику и копируем свежую сборку
CMD ["sh", "-c", "rm -rf /static/admin-frontend/* && cp -r /app/dist/* /static/admin-frontend/"]