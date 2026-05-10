/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_TELEGRAM_BOT_NAME?: string
  readonly VITE_YANDEX_METRIKA_ID?: string
  readonly VITE_YANDEX_METRIKA_ENABLED?: string
  // Если 'true' и DEV — main.ts ставит мок window.Telegram.WebApp,
  // чтобы тестировать openLink/closingConfirmation/setHeaderColor без
  // реального TG-клиента. Все методы логируют в console.
  readonly VITE_MOCK_TELEGRAM?: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
