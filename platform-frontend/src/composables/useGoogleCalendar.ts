import type { CommunityEvent } from '@/models/event'

export function useGoogleCalendar() {
  const formatDateForGoogle = (date: Date): string => {
    return `${date
      .toISOString()
      .split('.')[0]
      .replace(/[-:]/g, '')}Z`
  }

  const buildGoogleCalendarUrl = (event: CommunityEvent, durationMinutes: number = 60): string => {
    const start = new Date(event.date)
    const end = new Date(start.getTime() + durationMinutes * 60 * 1000)

    const startStr = formatDateForGoogle(start)
    const endStr = formatDateForGoogle(end)

    const url = new URL('https://calendar.google.com/calendar/render')
    url.searchParams.set('action', 'TEMPLATE')
    url.searchParams.set('text', event.title || '')
    url.searchParams.set('dates', `${startStr}/${endStr}`)

    let details = event.description || ''
    if (event.timezone && event.timezone !== 'UTC') {
      details += `\n\nВремя указано для таймзоны: ${event.timezone}`
    }
    url.searchParams.set('details', details)
    url.searchParams.set('location', `${event.videoLink || ''} ${event.place || ''}`.trim())

    return url.toString()
  }

  const openInGoogleCalendar = (event: CommunityEvent, durationMinutes: number = 60) => {
    const link = buildGoogleCalendarUrl(event, durationMinutes)
    window.open(link, '_blank')
  }

  return {
    openInGoogleCalendar,
    buildGoogleCalendarUrl,
    formatDateForGoogle,
  }
}
