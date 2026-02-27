import { useEventListener } from '@vueuse/core'
import { ref } from 'vue'

export function useScrollHeader(threshold = 10) {
  const isScrolled = ref(false)

  useEventListener(
    window,
    'scroll',
    () => {
      isScrolled.value = window.scrollY > threshold
    },
    { passive: true },
  )

  return { isScrolled }
}
