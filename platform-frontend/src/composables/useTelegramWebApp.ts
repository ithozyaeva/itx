// Тонкий враппер над window.Telegram.WebApp. Здесь только то, что нужно для
// MVP miniapp-входа: определить, что мы открыты внутри Telegram, и отдать
// initData бэкенду. BackButton/MainButton/themeParams/haptics — намеренно
// не трогаем, см. план tg-cached-newt.md.

import type { Ref } from 'vue'
import { onUnmounted, watch } from 'vue'

interface TelegramWebApp {
  initData: string
  ready: () => void
  expand: () => void
  close: () => void
  isExpanded?: boolean
  viewportHeight?: number
  viewportStableHeight?: number
  // Появилось в Bot API 7.7. В старых клиентах метода нет — вызывать
  // строго через optional chaining, иначе TypeError свалит initTelegramWebApp.
  disableVerticalSwipes?: () => void
  // Открывает t.me/* ссылку прямо в Telegram-клиенте (без выпадения в
  // системный браузер). Bot API 6.1+.
  openTelegramLink?: (url: string) => void
  // Открывает произвольную внешнюю ссылку во встроенном in-app браузере
  // Telegram. Bot API 6.1+.
  openLink?: (url: string, options?: { try_instant_view?: boolean }) => void
  // Включает confirm-диалог при попытке закрыть miniapp (свайпом или
  // close-кнопкой). Используется для форм с unsaved-данными. Bot API 6.2+.
  enableClosingConfirmation?: () => void
  disableClosingConfirmation?: () => void
  // Цвет шапки/фона miniapp в Telegram-клиенте. Принимает hex (#aabbcc),
  // 'bg_color' или 'secondary_bg_color'. Bot API 6.1+ / 6.9+.
  setHeaderColor?: (color: string) => void
  setBackgroundColor?: (color: string) => void
  onEvent?: (eventType: string, eventHandler: () => void) => void
  offEvent?: (eventType: string, eventHandler: () => void) => void
}

interface TelegramNamespace {
  WebApp?: TelegramWebApp
}

declare global {
  interface Window {
    Telegram?: TelegramNamespace
  }
}

export function getTelegramWebApp(): TelegramWebApp | null {
  return window.Telegram?.WebApp ?? null
}

// isMiniApp — true только если Telegram-клиент действительно открыл нас как
// WebApp и передал initData. Просто наличие window.Telegram.WebApp недоста-
// точно: SDK иногда подгружается и в обычном браузере, но initData там пустой,
// и подписать его нечем — на бэке такая попытка отвалится 401.
export function isMiniApp(): boolean {
  const tg = getTelegramWebApp()
  return !!tg && typeof tg.initData === 'string' && tg.initData.length > 0
}

const RGB_RE = /^rgba?\((\d+),\s*(\d+),\s*(\d+)/i

// rgbToHex — Telegram.setHeaderColor принимает только hex (#rrggbb), а
// computed-style браузер отдаёт rgb(r, g, b). Конвертим. На неожиданных
// форматах (transparent, hsl) возвращаем null — не вызываем сеттер.
function rgbToHex(rgb: string): string | null {
  const m = RGB_RE.exec(rgb)
  if (!m)
    return null
  const [, r, g, b] = m
  return `#${[r, g, b].map(n => Number(n).toString(16).padStart(2, '0')).join('')}`
}

// syncTelegramColors — приводим цвет шапки и фона Telegram-клиента к фону
// нашего приложения, чтобы между body и шапкой/нижней полосой Телеги не
// было контрастной полоски. Берём фактический backgroundColor у body после
// рендера темы (учитывает класс dark).
function syncTelegramColors(tg: TelegramWebApp) {
  const hex = rgbToHex(getComputedStyle(document.body).backgroundColor)
  if (!hex)
    return
  try {
    tg.setHeaderColor?.(hex)
    tg.setBackgroundColor?.(hex)
  }
  catch {
    // Старый клиент без поддержки — игнорируем.
  }
}

// initTelegramWebApp — вызвать один раз при старте приложения, если мы
// внутри Telegram. ready() сообщает клиенту, что UI готов отрисоваться (без
// этого мобильный TG показывает чёрный экран до первого fetch). expand()
// разворачивает miniapp на полную высоту — иначе TG открывает её половиной
// экрана, и юзер думает, что это ошибка вёрстки.
// disableVerticalSwipes выключает свайп-вниз-чтобы-закрыть: внутри прило-
// жения постоянно скроллят вертикально, и без этого юзер случайно гасит
// miniapp на каждой второй прокрутке.
// setHeaderColor/setBackgroundColor — синхронизируем с фоном приложения
// сразу и при каждом изменении класса dark/light на <html>.
//
// TODO: переподключить viewportChanged → --tg-viewport-stable-height,
// когда понадобится корректная высота модалок при выезде клавиатуры в
// iOS Mini App (#354 завёл setter, но консьюмеров так и не написали —
// модалки всё ещё на inset-0 / 100dvh).
export function initTelegramWebApp() {
  const tg = getTelegramWebApp()
  if (!tg)
    return
  try {
    tg.ready()
    tg.expand()
    tg.disableVerticalSwipes?.()
    syncTelegramColors(tg)
    new MutationObserver(() => syncTelegramColors(tg))
      .observe(document.documentElement, { attributes: true, attributeFilter: ['class'] })
  }
  catch {
    // SDK странно повёл себя в старом TG-клиенте — не критично, дальше идём.
  }
}

const TELEGRAM_LINK_RE = /^https?:\/\/(?:t\.me|telegram\.me)\//i

// openLink — открыть произвольную ссылку с учётом окружения. В Mini App
// t.me/* идут через openTelegramLink (юзер остаётся в Телеге, открывается
// нужный чат), внешние — через openLink (in-app браузер Телеги). В обычном
// браузере — window.open в новой вкладке. Использовать вместо <a target=
// "_blank"> и window.open везде, где есть переход на внешний ресурс.
export function openLink(url: string) {
  const tg = getTelegramWebApp()
  if (tg) {
    try {
      if (TELEGRAM_LINK_RE.test(url) && tg.openTelegramLink) {
        tg.openTelegramLink(url)
        return
      }
      if (tg.openLink) {
        tg.openLink(url)
        return
      }
    }
    catch {
      // SDK ругнулся (например, кривой URL) — fallback на window.open.
    }
  }
  window.open(url, '_blank', 'noopener,noreferrer')
}

// installMockTelegram — для локальной разработки без реального TG-клиента.
// Вызывать ОДИН раз в main.ts при VITE_MOCK_TELEGRAM=true до createApp(),
// чтобы getTelegramWebApp() в App.vue setup увидел мок. initData оставляем
// пустым — тогда isMiniApp() возвращает false и App.vue не пытается обмен-
// ять (заведомо невалидную) подпись на сессию через /api/auth/telegram-
// webapp. Все методы стабают console.info + window.open для openLink, чтобы
// в DevTools было видно, что клиентский код их вызывает.
export function installMockTelegram() {
  if (window.Telegram?.WebApp)
    return
  // eslint-disable-next-line no-console
  const log = (msg: string, ...args: unknown[]) => console.info(`[tg-mock] ${msg}`, ...args)
  const mock: TelegramWebApp = {
    initData: '',
    ready: () => log('ready'),
    expand: () => log('expand'),
    close: () => log('close'),
    isExpanded: true,
    viewportHeight: window.innerHeight,
    viewportStableHeight: window.innerHeight,
    disableVerticalSwipes: () => log('disableVerticalSwipes'),
    openTelegramLink: (url) => {
      log('openTelegramLink', url)
      window.open(url, '_blank', 'noopener,noreferrer')
    },
    openLink: (url, opts) => {
      log('openLink', url, opts)
      window.open(url, '_blank', 'noopener,noreferrer')
    },
    enableClosingConfirmation: () => log('enableClosingConfirmation'),
    disableClosingConfirmation: () => log('disableClosingConfirmation'),
    setHeaderColor: color => log('setHeaderColor', color),
    setBackgroundColor: color => log('setBackgroundColor', color),
    onEvent: (eventType, handler) => log('onEvent', eventType, handler),
    offEvent: (eventType, handler) => log('offEvent', eventType, handler),
  }
  window.Telegram = { WebApp: mock }
  log('installed (VITE_MOCK_TELEGRAM=true)')
}

// useClosingConfirmation — на формах с unsaved-данными. Пока dirty=true,
// Telegram при попытке закрыть miniapp покажет диалог «Точно закрыть?».
// При dirty=false и при unmount флаг снимаем, чтобы не остался висеть на
// других экранах.
export function useClosingConfirmation(dirty: Ref<boolean>) {
  const tg = getTelegramWebApp()
  if (!tg)
    return
  watch(
    dirty,
    (isDirty) => {
      try {
        if (isDirty)
          tg.enableClosingConfirmation?.()
        else
          tg.disableClosingConfirmation?.()
      }
      catch {
        // no-op: старый клиент без поддержки.
      }
    },
    { immediate: true },
  )
  onUnmounted(() => {
    try {
      tg.disableClosingConfirmation?.()
    }
    catch {
      // no-op
    }
  })
}
