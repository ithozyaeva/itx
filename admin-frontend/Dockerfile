FROM node:20-alpine

WORKDIR /app

ARG VITE_TELEGRAM_BOT_NAME
ENV VITE_TELEGRAM_BOT_NAME=${VITE_TELEGRAM_BOT_NAME}

COPY package*.json ./

RUN npm install

ARG CACHEBUST=1
COPY . .

RUN npm run build

# Очищаем старую статику и копируем свежую сборку
CMD ["sh", "-c", "rm -rf /static/admin-frontend/* && cp -r /app/dist/* /static/admin-frontend/"]
