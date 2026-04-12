import { onMounted, onUnmounted } from 'vue'

export function useRevealObserver() {
  let observer: IntersectionObserver | null = null

  onMounted(() => {
    observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting)
            entry.target.classList.add('revealed')
        })
      },
      { threshold: 0.08, rootMargin: '0px 0px -40px 0px' },
    )

    document.querySelectorAll('.reveal, .reveal-left, .reveal-right, .reveal-stagger').forEach((el) => {
      observer!.observe(el)
    })
  })

  onUnmounted(() => {
    observer?.disconnect()
  })
}
