FROM node:20-alpine

WORKDIR /app

ARG VITE_TELEGRAM_BOT_NAME
ENV VITE_TELEGRAM_BOT_NAME=${VITE_TELEGRAM_BOT_NAME}

COPY packages/ui/ /packages/ui/

COPY admin-frontend/package*.json ./

RUN npm install

ARG CACHEBUST=1
COPY admin-frontend/ .

RUN npm run build

CMD ["sh", "-c", "rm -rf /static/admin-frontend/* && cp -r /app/dist/* /static/admin-frontend/"]
