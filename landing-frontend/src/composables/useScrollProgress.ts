import { onMounted, onUnmounted, ref } from 'vue'

export function useScrollProgress() {
  const progress = ref(0)

  function update() {
    const h = document.documentElement.scrollHeight - window.innerHeight
    progress.value = h > 0 ? Math.min(window.scrollY / h, 1) : 0
  }

  onMounted(() => window.addEventListener('scroll', update, { passive: true }))
  onUnmounted(() => window.removeEventListener('scroll', update))

  return { progress }
}
