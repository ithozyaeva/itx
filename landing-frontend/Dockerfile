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

COPY . .

RUN npm run build

# Создаем директорию для статики
RUN mkdir -p /app/dist

# Копируем собранные файлы в dist
CMD ["sh", "-c", "cp -r /app/dist/* /static/landing-frontend/"] 