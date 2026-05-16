FROM node:22-alpine

WORKDIR /app

# CI всё ещё проксирует переменные окружения через VITE_* имена (из секретов GitHub
# и серверного окружения). Astro читает PUBLIC_* — поэтому маппим внутри Dockerfile,
# чтобы не править deploy-workflow и серверные .env при переезде с Vue.
ARG VITE_YANDEX_METRIKA_ID
ARG VITE_YANDEX_METRIKA_ENABLED
ARG VITE_TELEGRAM_BOT_NAME

ENV PUBLIC_YANDEX_METRIKA_ID=${VITE_YANDEX_METRIKA_ID}
ENV PUBLIC_YANDEX_METRIKA_ENABLED=${VITE_YANDEX_METRIKA_ENABLED}
ENV PUBLIC_TELEGRAM_BOT_NAME=${VITE_TELEGRAM_BOT_NAME}

COPY package*.json ./

RUN npm install

ARG CACHEBUST=1
COPY . .

RUN npm run build

# Очищаем старую статику и копируем свежую сборку.
# mkdir + find вместо `rm -rf /static/.../*`: glob не раскрывается в пустой директории
# и rm падает exit 1, из-за чего `&& cp` не выполняется и volume остаётся старым.
# /app/dist/. — копирует содержимое директории включая dot-файлы, независимо от shell.
CMD ["sh", "-c", "mkdir -p /static/landing-frontend && find /static/landing-frontend -mindepth 1 -delete && cp -a /app/dist/. /static/landing-frontend/"]
