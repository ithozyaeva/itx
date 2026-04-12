import type { Ref } from 'vue'
import { onMounted, onUnmounted, ref } from 'vue'

export function useCardReveal(containerRef: Ref<HTMLElement | null>) {
  const intersectionObserver = ref<IntersectionObserver | null>(null)
  const mutationObserver = ref<MutationObserver | null>(null)

  function observeElement(el: Element) {
    intersectionObserver.value?.observe(el)
  }

  onMounted(() => {
    intersectionObserver.value = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            entry.target.classList.add('animate-card-reveal')
            intersectionObserver.value?.unobserve(entry.target)
          }
        })
      },
      { threshold: 0.1 },
    )

    if (containerRef.value) {
      // Observe existing elements
      const cards = containerRef.value.querySelectorAll('[data-reveal]')
      cards.forEach(card => observeElement(card))

      // Watch for dynamically added elements (async data loading)
      mutationObserver.value = new MutationObserver((mutations) => {
        mutations.forEach((mutation) => {
          mutation.addedNodes.forEach((node) => {
            if (node instanceof HTMLElement) {
              if (node.hasAttribute('data-reveal')) {
                observeElement(node)
              }
              node.querySelectorAll('[data-reveal]').forEach(el => observeElement(el))
            }
          })
        })
      })

      mutationObserver.value.observe(containerRef.value, {
        childList: true,
        subtree: true,
      })
    }
  })

  onUnmounted(() => {
    intersectionObserver.value?.disconnect()
    mutationObserver.value?.disconnect()
  })
}
