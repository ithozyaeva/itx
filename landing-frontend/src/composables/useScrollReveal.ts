import type { Ref } from 'vue'
import { useIntersectionObserver } from '@vueuse/core'
import { ref } from 'vue'

export function useScrollReveal(options?: { threshold?: number, rootMargin?: string }) {
  const containerRef = ref<HTMLElement | null>(null) as Ref<HTMLElement | null>
  const isVisible = ref(false)

  useIntersectionObserver(
    containerRef,
    ([entry]) => {
      if (entry?.isIntersecting) {
        isVisible.value = true
      }
    },
    {
      threshold: options?.threshold ?? 0.1,
      rootMargin: options?.rootMargin ?? '0px',
    },
  )

  return { containerRef, isVisible }
}
