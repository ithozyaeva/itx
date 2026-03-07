import { describe, expect, it, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { defineComponent, ref } from 'vue'
import { useCardReveal } from '@/composables/useCardReveal'

// Helper to create a wrapper component that uses the composable
function createWrapper(template = '<div ref="containerRef"><div data-reveal>Card</div></div>') {
  return defineComponent({
    template,
    setup() {
      const containerRef = ref<HTMLElement | null>(null)
      useCardReveal(containerRef)
      return { containerRef }
    },
  })
}

function setupMockObservers() {
  const intersectionObserveSpy = vi.fn()
  const intersectionUnobserveSpy = vi.fn()
  const intersectionDisconnectSpy = vi.fn()
  let intersectionCallback: any = null

  const mutationObserveSpy = vi.fn()
  const mutationDisconnectSpy = vi.fn()
  let mutationCallback: any = null

  // Use class-style mocks so `new` works correctly
  vi.stubGlobal('IntersectionObserver', class {
    constructor(cb: any) {
      intersectionCallback = cb
    }

    observe = intersectionObserveSpy
    unobserve = intersectionUnobserveSpy
    disconnect = intersectionDisconnectSpy
  })

  vi.stubGlobal('MutationObserver', class {
    constructor(cb: any) {
      mutationCallback = cb
    }

    observe = mutationObserveSpy
    disconnect = mutationDisconnectSpy
  })

  return {
    intersectionObserveSpy,
    intersectionUnobserveSpy,
    intersectionDisconnectSpy,
    getIntersectionCallback: () => intersectionCallback,
    mutationObserveSpy,
    mutationDisconnectSpy,
    getMutationCallback: () => mutationCallback,
  }
}

describe('useCardReveal', () => {
  it('creates IntersectionObserver on mount', () => {
    const mocks = setupMockObservers()
    const wrapper = mount(createWrapper())

    expect(mocks.getIntersectionCallback()).not.toBeNull()

    wrapper.unmount()
    vi.unstubAllGlobals()
  })

  it('observes existing [data-reveal] elements', () => {
    const mocks = setupMockObservers()
    const wrapper = mount(createWrapper())

    expect(mocks.intersectionObserveSpy).toHaveBeenCalled()

    wrapper.unmount()
    vi.unstubAllGlobals()
  })

  it('sets up MutationObserver for dynamic elements', () => {
    const mocks = setupMockObservers()
    const wrapper = mount(createWrapper())

    expect(mocks.getMutationCallback()).not.toBeNull()
    expect(mocks.mutationObserveSpy).toHaveBeenCalledWith(
      expect.any(Object),
      { childList: true, subtree: true },
    )

    wrapper.unmount()
    vi.unstubAllGlobals()
  })

  it('adds animate-card-reveal class when element is intersecting', () => {
    const mocks = setupMockObservers()
    const wrapper = mount(createWrapper())

    const mockTarget = document.createElement('div')
    mocks.getIntersectionCallback()([
      { isIntersecting: true, target: mockTarget },
    ])

    expect(mockTarget.classList.contains('animate-card-reveal')).toBe(true)

    wrapper.unmount()
    vi.unstubAllGlobals()
  })

  it('unobserves element after it becomes visible', () => {
    const mocks = setupMockObservers()
    const wrapper = mount(createWrapper())

    const mockTarget = document.createElement('div')
    mocks.getIntersectionCallback()([
      { isIntersecting: true, target: mockTarget },
    ])

    expect(mocks.intersectionUnobserveSpy).toHaveBeenCalledWith(mockTarget)

    wrapper.unmount()
    vi.unstubAllGlobals()
  })

  it('does not animate non-intersecting entries', () => {
    const mocks = setupMockObservers()
    const wrapper = mount(createWrapper())

    const mockTarget = document.createElement('div')
    mocks.getIntersectionCallback()([
      { isIntersecting: false, target: mockTarget },
    ])

    expect(mockTarget.classList.contains('animate-card-reveal')).toBe(false)

    wrapper.unmount()
    vi.unstubAllGlobals()
  })

  it('disconnects observers on unmount', () => {
    const mocks = setupMockObservers()
    const wrapper = mount(createWrapper())

    wrapper.unmount()

    expect(mocks.intersectionDisconnectSpy).toHaveBeenCalled()
    expect(mocks.mutationDisconnectSpy).toHaveBeenCalled()

    vi.unstubAllGlobals()
  })

  it('handles dynamically added elements via MutationObserver', () => {
    const mocks = setupMockObservers()
    const wrapper = mount(createWrapper())

    mocks.intersectionObserveSpy.mockClear()

    const newElement = document.createElement('div')
    newElement.setAttribute('data-reveal', '')
    mocks.getMutationCallback()([
      { addedNodes: [newElement] },
    ])

    expect(mocks.intersectionObserveSpy).toHaveBeenCalledWith(newElement)

    wrapper.unmount()
    vi.unstubAllGlobals()
  })

  it('handles nested [data-reveal] in dynamically added elements', () => {
    const mocks = setupMockObservers()
    const wrapper = mount(createWrapper())

    mocks.intersectionObserveSpy.mockClear()

    const parent = document.createElement('div')
    const child = document.createElement('div')
    child.setAttribute('data-reveal', '')
    parent.appendChild(child)

    mocks.getMutationCallback()([
      { addedNodes: [parent] },
    ])

    expect(mocks.intersectionObserveSpy).toHaveBeenCalledWith(child)

    wrapper.unmount()
    vi.unstubAllGlobals()
  })
})
