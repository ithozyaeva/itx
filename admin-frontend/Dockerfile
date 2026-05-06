FROM node:22-alpine

WORKDIR /app

ARG VITE_TELEGRAM_BOT_NAME
ENV VITE_TELEGRAM_BOT_NAME=${VITE_TELEGRAM_BOT_NAME}

COPY package*.json ./

RUN npm install

ARG CACHEBUST=1
COPY . .

# Shared UI package (copied by deploy script before build)
RUN if [ -d "_packages_ui" ]; then mkdir -p /packages/ui && cp -r _packages_ui/* /packages/ui/; fi

RUN npm run build

CMD ["sh", "-c", "mkdir -p /static/admin-frontend && find /static/admin-frontend -mindepth 1 -delete && cp -a /app/dist/. /static/admin-frontend/"]
