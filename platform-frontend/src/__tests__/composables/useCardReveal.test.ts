import { beforeEach, describe, expect, it, vi } from 'vitest'
import { ref } from 'vue'
import { withSetup } from '../helpers'

// Mock IntersectionObserver
class MockIntersectionObserver {
  static instances: MockIntersectionObserver[] = []
  callback: IntersectionObserverCallback
  options: IntersectionObserverInit | undefined
  observedElements: Element[] = []
  unobservedElements: Element[] = []
  disconnected = false

  constructor(callback: IntersectionObserverCallback, options?: IntersectionObserverInit) {
    this.callback = callback
    this.options = options
    MockIntersectionObserver.instances.push(this)
  }

  observe(el: Element) {
    this.observedElements.push(el)
  }

  unobserve(el: Element) {
    this.unobservedElements.push(el)
  }

  disconnect() {
    this.disconnected = true
  }

  // Simulate entries intersecting
  simulateIntersection(entries: Partial<IntersectionObserverEntry>[]) {
    this.callback(entries as IntersectionObserverEntry[], this as unknown as IntersectionObserver)
  }
}

// Mock MutationObserver
class MockMutationObserver {
  static instances: MockMutationObserver[] = []
  callback: MutationCallback
  observedNode: Node | null = null
  observeOptions: MutationObserverInit | undefined
  disconnected = false

  constructor(callback: MutationCallback) {
    this.callback = callback
    MockMutationObserver.instances.push(this)
  }

  observe(node: Node, options?: MutationObserverInit) {
    this.observedNode = node
    this.observeOptions = options
  }

  disconnect() {
    this.disconnected = true
  }
}

describe('useCardReveal', () => {
  beforeEach(() => {
    MockIntersectionObserver.instances = []
    MockMutationObserver.instances = []
    vi.stubGlobal('IntersectionObserver', MockIntersectionObserver)
    vi.stubGlobal('MutationObserver', MockMutationObserver)
  })

  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('creates IntersectionObserver on mount with threshold 0.1', async () => {
    const { useCardReveal } = await import('@/composables/useCardReveal')
    const container = document.createElement('div')
    const containerRef = ref(container)

    withSetup(() => useCardReveal(containerRef))

    expect(MockIntersectionObserver.instances).toHaveLength(1)
    expect(MockIntersectionObserver.instances[0].options).toEqual({ threshold: 0.1 })
  })

  it('observes existing [data-reveal] elements in container', async () => {
    const { useCardReveal } = await import('@/composables/useCardReveal')
    const container = document.createElement('div')
    const card1 = document.createElement('div')
    card1.setAttribute('data-reveal', '')
    const card2 = document.createElement('div')
    card2.setAttribute('data-reveal', '')
    container.appendChild(card1)
    container.appendChild(card2)

    const containerRef = ref(container)
    withSetup(() => useCardReveal(containerRef))

    const observer = MockIntersectionObserver.instances[0]
    expect(observer.observedElements).toHaveLength(2)
  })

  it('adds animate-card-reveal class when element intersects', async () => {
    const { useCardReveal } = await import('@/composables/useCardReveal')
    const container = document.createElement('div')
    const card = document.createElement('div')
    card.setAttribute('data-reveal', '')
    container.appendChild(card)

    const containerRef = ref(container)
    withSetup(() => useCardReveal(containerRef))

    const observer = MockIntersectionObserver.instances[0]
    observer.simulateIntersection([
      { isIntersecting: true, target: card },
    ])

    expect(card.classList.contains('animate-card-reveal')).toBe(true)
  })

  it('unobserves element after it has been revealed', async () => {
    const { useCardReveal } = await import('@/composables/useCardReveal')
    const container = document.createElement('div')
    const card = document.createElement('div')
    card.setAttribute('data-reveal', '')
    container.appendChild(card)

    const containerRef = ref(container)
    withSetup(() => useCardReveal(containerRef))

    const observer = MockIntersectionObserver.instances[0]
    observer.simulateIntersection([
      { isIntersecting: true, target: card },
    ])

    expect(observer.unobservedElements).toContain(card)
  })

  it('does not add class when element is not intersecting', async () => {
    const { useCardReveal } = await import('@/composables/useCardReveal')
    const container = document.createElement('div')
    const card = document.createElement('div')
    card.setAttribute('data-reveal', '')
    container.appendChild(card)

    const containerRef = ref(container)
    withSetup(() => useCardReveal(containerRef))

    const observer = MockIntersectionObserver.instances[0]
    observer.simulateIntersection([
      { isIntersecting: false, target: card },
    ])

    expect(card.classList.contains('animate-card-reveal')).toBe(false)
  })

  it('sets up MutationObserver for dynamic elements', async () => {
    const { useCardReveal } = await import('@/composables/useCardReveal')
    const container = document.createElement('div')
    const containerRef = ref(container)

    withSetup(() => useCardReveal(containerRef))

    expect(MockMutationObserver.instances).toHaveLength(1)
    const mutObs = MockMutationObserver.instances[0]
    expect(mutObs.observedNode).toBe(container)
    expect(mutObs.observeOptions).toEqual({ childList: true, subtree: true })
  })

  it('observes dynamically added [data-reveal] elements via MutationObserver', async () => {
    const { useCardReveal } = await import('@/composables/useCardReveal')
    const container = document.createElement('div')
    const containerRef = ref(container)

    withSetup(() => useCardReveal(containerRef))

    const mutObs = MockMutationObserver.instances[0]
    const intObs = MockIntersectionObserver.instances[0]

    // Simulate a dynamically added node with data-reveal
    const newCard = document.createElement('div')
    newCard.setAttribute('data-reveal', '')

    mutObs.callback(
      [
        {
          addedNodes: [newCard] as unknown as NodeList,
          removedNodes: [] as unknown as NodeList,
          type: 'childList',
        } as unknown as MutationRecord,
      ],
      mutObs as unknown as MutationObserver,
    )

    expect(intObs.observedElements).toContain(newCard)
  })

  it('handles null containerRef gracefully', async () => {
    const { useCardReveal } = await import('@/composables/useCardReveal')
    const containerRef = ref<HTMLElement | null>(null)

    // Should not throw
    expect(() => withSetup(() => useCardReveal(containerRef))).not.toThrow()

    // IntersectionObserver is still created, but MutationObserver is not set up
    expect(MockIntersectionObserver.instances).toHaveLength(1)
    expect(MockMutationObserver.instances).toHaveLength(0)
  })

  it('disconnects observers on unmount', async () => {
    const { useCardReveal } = await import('@/composables/useCardReveal')
    const container = document.createElement('div')
    const containerRef = ref(container)

    const { app } = withSetup(() => useCardReveal(containerRef))

    const intObs = MockIntersectionObserver.instances[0]
    const mutObs = MockMutationObserver.instances[0]

    app.unmount()

    expect(intObs.disconnected).toBe(true)
    expect(mutObs.disconnected).toBe(true)
  })
})
