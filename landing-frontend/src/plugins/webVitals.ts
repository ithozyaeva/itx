import type { Metric } from 'web-vitals'
import { onCLS, onINP, onLCP } from 'web-vitals'

declare global {
  interface Window {
    ym?: (id: number | string, action: string, ...args: unknown[]) => void
  }
}

// CLS приходит в долях единицы (0..1+); нормализуем под int-значение Метрики
function toMetrikaValue(metric: Metric): number {
  return Math.round(metric.name === 'CLS' ? metric.value * 1000 : metric.value)
}

function send(metric: Metric) {
  const id = import.meta.env.VITE_YANDEX_METRIKA_ID
  if (id && typeof window.ym === 'function') {
    window.ym(Number(id), 'params', {
      webVitals: {
        [metric.name]: {
          value: toMetrikaValue(metric),
          rating: metric.rating,
        },
      },
    })
  }
}

export function initWebVitals() {
  onLCP(send)
  onINP(send)
  onCLS(send)
}
