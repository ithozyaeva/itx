import type { Ref } from 'vue'
import { onMounted, onUnmounted, ref } from 'vue'

export function useMagneticHover(elRef: Ref<HTMLElement | null>, strength = 0.35) {
  const x = ref(0)
  const y = ref(0)

  function onMove(e: MouseEvent) {
    const el = elRef.value
    if (!el)
      return
    const rect = el.getBoundingClientRect()
    const cx = rect.left + rect.width / 2
    const cy = rect.top + rect.height / 2
    x.value = (e.clientX - cx) * strength
    y.value = (e.clientY - cy) * strength
  }

  function onLeave() {
    x.value = 0
    y.value = 0
  }

  onMounted(() => {
    elRef.value?.addEventListener('mousemove', onMove)
    elRef.value?.addEventListener('mouseleave', onLeave)
  })

  onUnmounted(() => {
    elRef.value?.removeEventListener('mousemove', onMove)
    elRef.value?.removeEventListener('mouseleave', onLeave)
  })

  return { x, y }
}
