FROM node:20-alpine

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY . .

RUN npm run build

# Создаем директорию для статики
RUN mkdir -p /app/dist

# Копируем собранные файлы в dist
CMD ["sh", "-c", "cp -r /app/dist/* /static/admin-frontend/"]