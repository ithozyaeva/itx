FROM node:20-alpine

WORKDIR /app

ARG VITE_YANDEX_METRIKA_ID
ARG VITE_YANDEX_METRIKA_ENABLED
ARG VITE_TELEGRAM_BOT_NAME

ENV VITE_YANDEX_METRIKA_ID=${VITE_YANDEX_METRIKA_ID}
ENV VITE_YANDEX_METRIKA_ENABLED=${VITE_YANDEX_METRIKA_ENABLED}
ENV VITE_TELEGRAM_BOT_NAME=${VITE_TELEGRAM_BOT_NAME}

COPY package*.json ./

RUN npm install

ARG CACHEBUST=1
COPY . .

RUN npm run build

# Очищаем старую статику и копируем свежую сборку
CMD ["sh", "-c", "rm -rf /static/landing-frontend/* && cp -r /app/dist/* /static/landing-frontend/"]
