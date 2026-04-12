import { ref } from 'vue'

export function useCountUp(target: number, duration = 1800) {
  const value = ref(0)
  const started = ref(false)

  function start() {
    if (started.value)
      return
    started.value = true
    const startTime = performance.now()
    const ease = (t: number) => 1 - (1 - t) ** 3

    function tick(now: number) {
      const elapsed = now - startTime
      const progress = Math.min(elapsed / duration, 1)
      value.value = Math.round(ease(progress) * target)
      if (progress < 1)
        requestAnimationFrame(tick)
    }
    requestAnimationFrame(tick)
  }

  return { value, start }
}
