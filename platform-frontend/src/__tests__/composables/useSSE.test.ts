import { beforeEach, describe, expect, it, vi } from 'vitest'
import { withSetup } from '../helpers'

// Mock EventSource
class MockEventSource {
  static instances: MockEventSource[] = []
  url: string
  onmessage: ((event: MessageEvent) => void) | null = null
  onerror: (() => void) | null = null
  readyState = 0
  close = vi.fn()

  constructor(url: string) {
    this.url = url
    MockEventSource.instances.push(this)
  }

  simulateMessage(data: string) {
    if (this.onmessage) {
      this.onmessage({ data } as MessageEvent)
    }
  }

  simulateError() {
    if (this.onerror) {
      this.onerror()
    }
  }
}

describe('useSSE', () => {
  beforeEach(() => {
    vi.resetModules()
    vi.useFakeTimers()
    localStorage.clear()
    MockEventSource.instances = []
    vi.stubGlobal('EventSource', MockEventSource)
  })

  afterEach(() => {
    vi.useRealTimers()
    vi.unstubAllGlobals()
  })

  it('connects to SSE with token from localStorage', async () => {
    localStorage.setItem('tg_token', 'test-token-123')
    const { useSSE } = await import('@/composables/useSSE')
    const callback = vi.fn()

    withSetup(() => useSSE('notifications', callback))

    expect(MockEventSource.instances).toHaveLength(1)
    expect(MockEventSource.instances[0].url).toBe('/api/platform/sse?token=test-token-123')
  })

  it('does not connect when no token in localStorage', async () => {
    const { useSSE } = await import('@/composables/useSSE')
    const callback = vi.fn()

    withSetup(() => useSSE('notifications', callback))

    expect(MockEventSource.instances).toHaveLength(0)
  })

  it('calls callback when matching event type is received', async () => {
    localStorage.setItem('tg_token', 'token')
    const { useSSE } = await import('@/composables/useSSE')
    const callback = vi.fn()

    withSetup(() => useSSE('notifications', callback))

    const es = MockEventSource.instances[0]
    es.simulateMessage(JSON.stringify({ type: 'notifications' }))

    expect(callback).toHaveBeenCalledOnce()
  })

  it('does not call callback for non-matching event type', async () => {
    localStorage.setItem('tg_token', 'token')
    const { useSSE } = await import('@/composables/useSSE')
    const callback = vi.fn()

    withSetup(() => useSSE('notifications', callback))

    const es = MockEventSource.instances[0]
    es.simulateMessage(JSON.stringify({ type: 'tasks' }))

    expect(callback).not.toHaveBeenCalled()
  })

  it('sets connected to true on "connected" message', async () => {
    localStorage.setItem('tg_token', 'token')
    const { useSSE } = await import('@/composables/useSSE')
    const callback = vi.fn()

    const { result } = withSetup(() => useSSE('notifications', callback))

    const es = MockEventSource.instances[0]
    es.simulateMessage(JSON.stringify({ type: 'connected' }))

    expect(result.connected.value).toBe(true)
    // 'connected' type should not trigger the notifications callback
    expect(callback).not.toHaveBeenCalled()
  })

  it('reconnects after error with 5 second delay', async () => {
    localStorage.setItem('tg_token', 'token')
    const { useSSE } = await import('@/composables/useSSE')
    const callback = vi.fn()

    withSetup(() => useSSE('notifications', callback))

    expect(MockEventSource.instances).toHaveLength(1)
    const firstES = MockEventSource.instances[0]

    // Simulate error
    firstES.simulateError()
    expect(firstES.close).toHaveBeenCalled()

    // Before timer fires, no reconnection
    expect(MockEventSource.instances).toHaveLength(1)

    // Advance timer by 5 seconds
    vi.advanceTimersByTime(5000)

    // A new EventSource should have been created
    expect(MockEventSource.instances).toHaveLength(2)
  })

  it('ignores invalid JSON messages', async () => {
    localStorage.setItem('tg_token', 'token')
    const { useSSE } = await import('@/composables/useSSE')
    const callback = vi.fn()

    withSetup(() => useSSE('notifications', callback))

    const es = MockEventSource.instances[0]
    // Should not throw
    expect(() => es.simulateMessage('invalid json{')).not.toThrow()
    expect(callback).not.toHaveBeenCalled()
  })

  it('stopSSE closes connection and clears listeners', async () => {
    localStorage.setItem('tg_token', 'token')
    const { useSSE, stopSSE } = await import('@/composables/useSSE')
    const callback = vi.fn()

    withSetup(() => useSSE('notifications', callback))

    const es = MockEventSource.instances[0]
    stopSSE()

    expect(es.close).toHaveBeenCalled()
  })

  it('does not create duplicate connections', async () => {
    localStorage.setItem('tg_token', 'token')
    const { useSSE } = await import('@/composables/useSSE')

    withSetup(() => {
      useSSE('notifications', vi.fn())
      useSSE('tasks', vi.fn())
      return {}
    })

    // Only one EventSource should be created
    expect(MockEventSource.instances).toHaveLength(1)
  })
})
