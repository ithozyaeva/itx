import { mount } from '@vue/test-utils'
import { defineComponent, ref } from 'vue'
import { describe, expect, it, vi } from 'vitest'
import { useCardReveal } from '../composables/useCardReveal'

describe('useCardReveal', () => {
  function stubObservers() {
    const ioObserve = vi.fn()
    const ioDisconnect = vi.fn()
    const moObserve = vi.fn()
    const moDisconnect = vi.fn()

    vi.stubGlobal('IntersectionObserver', vi.fn(() => ({
      observe: ioObserve,
      unobserve: vi.fn(),
      disconnect: ioDisconnect,
    })))

    vi.stubGlobal('MutationObserver', vi.fn(() => ({
      observe: moObserve,
      disconnect: moDisconnect,
    })))

    return { ioObserve, ioDisconnect, moObserve, moDisconnect }
  }

  it('creates observers on mount and disconnects on unmount', () => {
    const { ioDisconnect, moDisconnect } = stubObservers()

    const TestComponent = defineComponent({
      setup() {
        const containerRef = ref<HTMLElement | null>(null)
        useCardReveal(containerRef)
        return { containerRef }
      },
      template: '<div ref="containerRef"><div data-reveal>Card</div></div>',
    })

    const wrapper = mount(TestComponent)
    expect(IntersectionObserver).toHaveBeenCalled()
    expect(MutationObserver).toHaveBeenCalled()

    wrapper.unmount()
    expect(ioDisconnect).toHaveBeenCalled()
    expect(moDisconnect).toHaveBeenCalled()

    vi.unstubAllGlobals()
  })

  it('observes elements with data-reveal attribute', () => {
    const { ioObserve } = stubObservers()

    const TestComponent = defineComponent({
      setup() {
        const containerRef = ref<HTMLElement | null>(null)
        useCardReveal(containerRef)
        return { containerRef }
      },
      template: `
        <div ref="containerRef">
          <div data-reveal>One</div>
          <div data-reveal>Two</div>
          <div>Not revealed</div>
        </div>
      `,
    })

    mount(TestComponent)
    expect(ioObserve).toHaveBeenCalledTimes(2)

    vi.unstubAllGlobals()
  })

  it('sets up MutationObserver for dynamic content', () => {
    const { moObserve } = stubObservers()

    const TestComponent = defineComponent({
      setup() {
        const containerRef = ref<HTMLElement | null>(null)
        useCardReveal(containerRef)
        return { containerRef }
      },
      template: '<div ref="containerRef"></div>',
    })

    mount(TestComponent)
    expect(moObserve).toHaveBeenCalledWith(
      expect.any(HTMLElement),
      { childList: true, subtree: true },
    )

    vi.unstubAllGlobals()
  })
})
