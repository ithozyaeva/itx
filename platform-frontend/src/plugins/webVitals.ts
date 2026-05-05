import type { Metric } from 'web-vitals'
import { onCLS, onINP, onLCP } from 'web-vitals'

// RUM-бэкенда нет — пишем в console, чтобы видеть метрики в DevTools.
// Эмитим CustomEvent, чтобы при необходимости подцепиться внешним скриптом.
// TODO: заменить на отправку в реальный RUM (Я.Метрика/own backend),
// когда такой эндпоинт появится.
function send(metric: Metric) {
  // eslint-disable-next-line no-console
  console.info(`[web-vitals] ${metric.name}=${metric.value.toFixed(2)} (${metric.rating})`)
  window.dispatchEvent(new CustomEvent('web-vitals', { detail: metric }))
}

export function initWebVitals() {
  onLCP(send)
  onINP(send)
  onCLS(send)
}
