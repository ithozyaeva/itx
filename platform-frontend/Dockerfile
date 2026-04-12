FROM node:20-alpine

WORKDIR /app

ARG VITE_TELEGRAM_BOT_NAME
ENV VITE_TELEGRAM_BOT_NAME=${VITE_TELEGRAM_BOT_NAME}

COPY packages/ui/ /packages/ui/

COPY platform-frontend/package*.json ./

RUN npm install

ARG CACHEBUST=1
COPY platform-frontend/ .

RUN npm run build

CMD ["sh", "-c", "rm -rf /static/platform-frontend/* && cp -r /app/dist/* /static/platform-frontend/"]
