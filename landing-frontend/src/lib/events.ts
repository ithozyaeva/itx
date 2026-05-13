// Build-time fetch событий с прод-API.
// На билде запрашиваем next/old, обрабатываем и возвращаем массив для рендера.
// Не падаем — если API недоступен, возвращаем пустой массив.

export interface EventTag {
  id: number
  name: string
}

export interface EventHost {
  id: number
  telegramID: number
  tg: string
  firstName: string
  lastName: string
  avatarUrl?: string
}

export interface CommunityEvent {
  id: number
  title: string
  description: string
  date: string
  endDate?: string
  timezone: string
  placeType: 'ONLINE' | 'OFFLINE' | 'HYBRID'
  place: string
  customPlaceType: string
  eventType: string
  open: boolean
  isRepeating: boolean
  repeatPeriod?: string
  repeatInterval?: number
  repeatEndDate?: string
  hosts: EventHost[]
  eventTags: EventTag[]
  videoLink: string
}

const API_BASE = (import.meta.env.PUBLIC_API_BASE as string | undefined) ?? 'https://ithozyaeva.ru'

async function fetchEvents(path: string): Promise<CommunityEvent[]> {
  try {
    const res = await fetch(`${API_BASE}${path}`, {
      headers: { Accept: 'application/json' },
    })
    if (!res.ok) {
      console.warn(`[events] ${path} returned ${res.status}, skipping`)
      return []
    }
    const data = await res.json() as { items?: CommunityEvent[] }
    return data.items ?? []
  }
  catch (e) {
    console.warn(`[events] ${path} fetch failed:`, e)
    return []
  }
}

export async function getNextEvents(): Promise<CommunityEvent[]> {
  return fetchEvents('/api/events/next')
}

export async function getOldEvents(): Promise<CommunityEvent[]> {
  return fetchEvents('/api/events/old')
}

const SLASH_RE = /\//g

export function formatEventStamp(dateString: string): string {
  const d = new Date(dateString)
  const fmt = new Intl.DateTimeFormat('en-GB', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
    timeZone: 'Europe/Moscow',
  }).format(d)
  return fmt.replace(',', '').replace(SLASH_RE, '-')
}

export function formatEventDate(dateString: string): string {
  const d = new Date(dateString)
  return new Intl.DateTimeFormat('ru-RU', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
    timeZone: 'Europe/Moscow',
  }).format(d)
}

const GC_DATE_RE = /[-:]/g

function formatGcDate(d: Date): string {
  return `${d.toISOString().split('.')[0].replace(GC_DATE_RE, '')}Z`
}

export function buildGoogleCalendarUrl(event: CommunityEvent, durationMinutes = 60): string {
  const start = new Date(event.date)
  const end = event.endDate
    ? new Date(event.endDate)
    : new Date(start.getTime() + durationMinutes * 60 * 1000)

  const url = new URL('https://calendar.google.com/calendar/render')
  url.searchParams.set('action', 'TEMPLATE')
  url.searchParams.set('text', event.title || '')
  url.searchParams.set('dates', `${formatGcDate(start)}/${formatGcDate(end)}`)

  let details = event.description || ''
  if (event.timezone && event.timezone !== 'UTC')
    details += `\n\n⏰ Время указано для таймзоны: ${event.timezone}`
  url.searchParams.set('details', details)
  url.searchParams.set('location', `${event.videoLink || ''} ${event.place || ''}`.trim())

  return url.toString()
}
