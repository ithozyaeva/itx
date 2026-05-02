import { onUnmounted, ref } from 'vue'

export type SSEEventType = 'notifications' | 'tasks' | 'raffles' | 'kudos' | 'points' | 'quests' | 'dailies' | 'streak' | 'challenges' | 'connected'

interface SSEMessage {
  type: SSEEventType
  data?: unknown
}

const listeners = new Map<SSEEventType, Set<() => void>>()
let eventSource: EventSource | null = null
let reconnectTimer: ReturnType<typeof setTimeout> | null = null
let consecutiveErrors = 0
const connected = ref(false)

function connect() {
  if (eventSource)
    return

  const token = localStorage.getItem('tg_token')
  if (!token)
    return

  eventSource = new EventSource(`/api/platform/sse?token=${encodeURIComponent(token)}`)

  eventSource.onmessage = (event) => {
    consecutiveErrors = 0
    try {
      const msg: SSEMessage = JSON.parse(event.data)
      if (msg.type === 'connected') {
        connected.value = true
        return
      }
      const callbacks = listeners.get(msg.type)
      if (callbacks) {
        callbacks.forEach(cb => cb())
      }
    }
    catch {
      // ignore parse errors
    }
  }

  eventSource.onerror = () => {
    consecutiveErrors++
    connected.value = false
    eventSource?.close()
    eventSource = null
    const token = localStorage.getItem('tg_token')
    if (token) {
      const delay = Math.min(1000 * 2 ** consecutiveErrors, 60000)
      reconnectTimer = setTimeout(connect, delay)
    }
  }
}

function disconnect() {
  if (reconnectTimer) {
    clearTimeout(reconnectTimer)
    reconnectTimer = null
  }
  if (eventSource) {
    eventSource.close()
    eventSource = null
  }
  consecutiveErrors = 0
  connected.value = false
}

export function startSSE() {
  connect()
}

export function stopSSE() {
  disconnect()
  listeners.clear()
}

export function useSSE(eventType: SSEEventType, callback: () => void) {
  if (!listeners.has(eventType))
    listeners.set(eventType, new Set())

  listeners.get(eventType)!.add(callback)

  // Auto-connect on first use
  connect()

  onUnmounted(() => {
    listeners.get(eventType)?.delete(callback)
  })

  return { connected }
}
