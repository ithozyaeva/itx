# CLAUDE.md — IT-X Project

## Project Overview

IT-X — монорепо с полноценной платформой сообщества: Go-бэкенд, три Vue 3 фронтенда, Telegram-бот, PostgreSQL. Всё контейнеризировано через Docker.

## Repository Structure

```
itx/
├── backend/              # Go API
├── landing-frontend/     # Лендинг
├── admin-frontend/       # Админ-панель
├── platform-frontend/    # Основная платформа
├── _nginx/               # Конфигурация Nginx
├── docker/               # Скрипты certbot и nginx
├── docs/                 # Архитектурная документация (C4 диаграммы)
├── static/               # Собранные фронтенды
├── docker-compose.yml    # Локальная разработка
└── docker-compose.prod.yml  # Продакшн (SSL)
```

Все четыре подпроекта находятся в монорепо как обычные директории.

## Tech Stack

**Backend:** Go 1.24, Fiber v2, GORM, PostgreSQL 15, JWT auth, AWS S3, Telegram Bot API
**Frontend (все три):** Vue 3, TypeScript, Vite, Tailwind CSS, Pinia, Vue Router, TanStack Query, Reka UI, Lucide icons
**DevOps:** Docker, Nginx, Let's Encrypt/Certbot, GitHub Actions

## Common Commands

### Запуск всего проекта (Docker)
```bash
docker-compose up --build -d
```
- Landing: http://localhost
- Admin: http://localhost/admin
- Platform: http://localhost/platform
- API: http://localhost:3000

### Фронтенд отдельно
```bash
cd landing-frontend && npm install && npm run dev    # :5173
cd admin-frontend && npm install && npm run dev      # :5174
cd platform-frontend && npm install && npm run dev   # :5175
```

### Бэкенд отдельно
```bash
cd backend && docker-compose up --build -d
```

### Линтинг и проверки (внутри фронтенд-директорий)
```bash
npm run lint          # Проверка
npm run lint:fix      # Автоисправление
npm run type-check    # Проверка типов
npm run build         # Продакшн-сборка
```

## Backend Architecture

```
backend/
├── cmd/main.go              # Entry point, порт 3000
├── config/                  # Конфигурация
├── database/                # Подключение к БД, миграции (18+ файлов)
├── internal/
│   ├── handler/             # HTTP хендлеры (13 файлов)
│   ├── service/             # Бизнес-логика
│   ├── repository/          # Слой доступа к данным
│   ├── middleware/           # Auth middleware
│   ├── models/              # Модели данных
│   ├── bot/                 # Telegram бот
│   └── utils/               # Утилиты
```

## Frontend Architecture

Все фронтенды используют одинаковую структуру:
- `src/components/` — компоненты
- `src/composables/` — переиспользуемая логика
- `src/services/` — API-клиенты
- `src/router/` — маршрутизация
- `src/sections/` — секции страниц
- `src/assets/` — статика

Path alias: `@/*` → `./src/*`

## Code Style

- ESLint с @antfu/eslint-config
- Husky pre-commit хуки: `npm run lint`
- Vue: max 1 атрибут на строку
- TypeScript strict mode

## Environment Variables

**Backend:** DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, PORT, JWT_SECRET, JWT_EXPIRATION, TELEGRAM_BOT_TOKEN, TELEGRAM_MAIN_CHAT_ID, S3_ENDPOINT, S3_REGION, S3_ACCESS_KEY, S3_SECRET_KEY, S3_BUCKET

**Frontend:** VITE_YANDEX_METRIKA_ID, VITE_YANDEX_METRIKA_ENABLED, VITE_TELEGRAM_BOT_NAME

## Git Workflow

- **Никогда не пушить напрямую в master.** Для любых изменений создавать feature-ветку и pull request.
- Формат веток: `feature/<описание>`, `fix/<описание>`, `chore/<описание>`
- PR должен содержать описание изменений

## Deployment

- Dev деплой автоматически при push в master
- Продакшн деплой через GitHub Actions (manual trigger с подтверждением)
- SSH на сервер → git pull → docker-compose rebuild
- SSL через Certbot с автообновлением

## Key Domain Concepts

- Members / Users — участники сообщества
- Mentors — менторы с профессиональными тегами
- Events — мероприятия (поддержка повторяющихся, таймзоны)
- Services — услуги
- Reviews — отзывы
- Resumes — резюме
- Referral system — реферальная система
- Telegram OAuth — авторизация через Telegram
