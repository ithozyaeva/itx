/// <reference path="../.astro/types.d.ts" />

interface ImportMetaEnv {
  readonly PUBLIC_YANDEX_METRIKA_ID?: string
  readonly PUBLIC_YANDEX_METRIKA_ENABLED?: string
  readonly PUBLIC_TELEGRAM_BOT_NAME?: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
