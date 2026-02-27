/// <reference types="vite/client" />
import 'vite-svg-loader'

interface ImportMetaEnv {
  readonly VITE_YANDEX_METRIKA_ID?: string
  readonly VITE_YANDEX_METRIKA_ENABLED?: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
