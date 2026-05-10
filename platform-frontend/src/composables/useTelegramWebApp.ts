// Тонкий враппер над window.Telegram.WebApp. Здесь только то, что нужно для
// MVP miniapp-входа: определить, что мы открыты внутри Telegram, и отдать
// initData бэкенду. BackButton/MainButton/themeParams/haptics — намеренно
// не трогаем, см. план tg-cached-newt.md.

interface TelegramWebApp {
  initData: string
  ready: () => void
  expand: () => void
  close: () => void
  // Появилось в Bot API 7.7. В старых клиентах метода нет — вызывать
  // строго через optional chaining, иначе TypeError свалит initTelegramWebApp.
  disableVerticalSwipes?: () => void
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

// initTelegramWebApp — вызвать один раз при старте приложения, если мы
// внутри Telegram. ready() сообщает клиенту, что UI готов отрисоваться (без
// этого мобильный TG показывает чёрный экран до первого fetch). expand()
// разворачивает miniapp на полную высоту — иначе TG открывает её половиной
// экрана, и юзер думает, что это ошибка вёрстки.
// disableVerticalSwipes выключает свайп-вниз-чтобы-закрыть: внутри прило-
// жения постоянно скроллят вертикально, и без этого юзер случайно гасит
// miniapp на каждой второй прокрутке.
export function initTelegramWebApp() {
  const tg = getTelegramWebApp()
  if (!tg)
    return
  try {
    tg.ready()
    tg.expand()
    tg.disableVerticalSwipes?.()
  }
  catch {
    // SDK странно повёл себя в старом TG-клиенте — не критично, дальше идём.
  }
}
