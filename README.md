# IT-X

### Структура проекта

Монорепо из 4 компонентов:
- `backend` - Go API бэкенд
- `landing-frontend` - лендинг проекта
- `admin-frontend` - фронт админ-панели
- `platform-frontend` - фронт платформы

### Подготовка к запуску

```
git clone git@github.com:ithozyaeva/itx.git
```

#### После клонирования:
1) Убедитесь, что у Вас установлен Docker и Docker Compose.
2) Скопируйте и переименуйте файлы `.env.template` в `.env` в корневом и `backend`. При необходимости измените значения некоторых ключей.
3) В `backend/migrations/001_initial_schema.sql` на последней строчке можно изменить данные первого пользователя админки.

### Запуск проекта

Для запуска проекта необходимо выполнить следующую команду: 

<code>docker-compose up --build -d</code>

После успешного запуска Вы сможете получить доступы к: 
- Лендингу: http://localhost
- Админ-панели: http://localhost/admin
- Платформе: http://localhost/platform
- API серверу: http://localhost:3000
- Базе данных: localhost:5432

>Креды админки admin/3_lgWfY0yI

### Запуск отдельных компонентов проекта

#### Запуск бэкенда

Для запуска только бэкенда выполните команду:

```
cd backend
docker-compose up --build -d
```

Бэкенд будет доступен по адресу: http://localhost:3000

#### Запуск лендинга

Для запуска только лендинга выполните команду:

```
cd landing-frontend
npm install
npm run start
```

Лендинг будет доступен по адресу: http://localhost:5173

#### Запуск админ-панели

Для запуска только админ-панели выполните команду:

```
cd admin-frontend
npm install
npm run start
```

Админ-панель будет доступна по адресу: http://localhost:5174

#### Запуск платформы

Для запуска только платформы выполните команду:

```
cd platform-frontend
npm install
npm run start
```

Платформа будет доступна по адресу: http://localhost:5175

### Обновление репозитория

```
git pull
```