# GitHub Actions для автоматического обновления сабмодулей

## Обзор

Данная настройка GitHub Actions позволяет автоматически обновлять родительский репозиторий (itx) при изменениях в сабмодулях landing-frontend и admin-frontend.

## Как это работает

1. При пуше изменений в ветку master сабмодуля landing-frontend запускается GitHub Action из файла `.github/workflows/notify-main-repo.yml`, которая отправляет событие "repository_dispatch" в родительский репозиторий.

2. В родительском репозитории настроен workflow `.github/workflows/update-landing-submodule.yml`, который реагирует на это событие и автоматически обновляет ссылку на сабмодуль до последнего коммита.

3. Аналогично настроено и для admin-frontend:
   - В сабмодуле: `.github/workflows/notify-main-repo.yml` отправляет событие типа "admin-frontend-update"
   - В родительском репозитории: `.github/workflows/update-admin-submodule.yml` обрабатывает это событие

## Требования

Для корректной работы необходимо:

1. Создать Personal Access Token (PAT) с правами доступа ко всем репозиториям.
2. Добавить этот токен в секреты всех репозиториев (родительского и сабмодулей) под именем `REPO_ACCESS_TOKEN`.

## Ручной запуск

Workflow в родительском репозитории можно запустить вручную через интерфейс GitHub Actions (опция `workflow_dispatch`). 