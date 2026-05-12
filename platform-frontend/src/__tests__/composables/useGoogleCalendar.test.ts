import { describe, expect, it, vi } from 'vitest'
import type { CommunityEvent } from '@/models/event'
import { useGoogleCalendar } from '@/composables/useGoogleCalendar'

function createEvent(overrides: Partial<CommunityEvent> = {}): CommunityEvent {
  return {
    id: 1,
    title: 'Test Event',
    description: 'Test description',
    date: '2024-01-15T10:30:00.000Z',
    timezone: 'UTC',
    placeType: 'ONLINE',
    place: '',
    customPlaceType: '',
    eventType: 'meetup',
    open: true,
    videoLink: 'https://zoom.us/123',
    isRepeating: false,
    recordingUrl: '',
    maxParticipants: 0,
    hosts: [],
    members: [],
    eventTags: [],
    ...overrides,
  }
}

describe('useGoogleCalendar', () => {
  describe('formatDateForGoogle', () => {
    it('formats a Date into Google Calendar date string', () => {
      const { formatDateForGoogle } = useGoogleCalendar()
      const date = new Date('2024-01-15T10:30:00.000Z')

      expect(formatDateForGoogle(date)).toBe('20240115T103000Z')
    })

    it('formats midnight correctly', () => {
      const { formatDateForGoogle } = useGoogleCalendar()
      const date = new Date('2024-12-31T00:00:00.000Z')

      expect(formatDateForGoogle(date)).toBe('20241231T000000Z')
    })
  })

  describe('buildGoogleCalendarUrl', () => {
    it('includes correct params with default 60-minute duration', () => {
      const { buildGoogleCalendarUrl } = useGoogleCalendar()
      const event = createEvent()
      const url = new URL(buildGoogleCalendarUrl(event))

      expect(url.origin + url.pathname).toBe('https://calendar.google.com/calendar/render')
      expect(url.searchParams.get('action')).toBe('TEMPLATE')
      expect(url.searchParams.get('text')).toBe('Test Event')
      expect(url.searchParams.get('dates')).toBe('20240115T103000Z/20240115T113000Z')
      expect(url.searchParams.get('details')).toBe('Test description')
      expect(url.searchParams.get('location')).toBe('https://zoom.us/123')
    })

    it('uses custom duration', () => {
      const { buildGoogleCalendarUrl } = useGoogleCalendar()
      const event = createEvent()
      const url = new URL(buildGoogleCalendarUrl(event, 120))

      expect(url.searchParams.get('dates')).toBe('20240115T103000Z/20240115T123000Z')
    })

    it('appends timezone info when timezone is not UTC', () => {
      const { buildGoogleCalendarUrl } = useGoogleCalendar()
      const event = createEvent({ timezone: 'Europe/Moscow' })
      const url = new URL(buildGoogleCalendarUrl(event))

      const details = url.searchParams.get('details')
      expect(details).toContain('Test description')
      expect(details).toContain('Время указано для таймзоны: Europe/Moscow')
    })

    it('does not append timezone info when timezone is UTC', () => {
      const { buildGoogleCalendarUrl } = useGoogleCalendar()
      const event = createEvent({ timezone: 'UTC' })
      const url = new URL(buildGoogleCalendarUrl(event))

      expect(url.searchParams.get('details')).toBe('Test description')
    })

    it('handles missing optional fields gracefully', () => {
      const { buildGoogleCalendarUrl } = useGoogleCalendar()
      const event = createEvent({
        title: '',
        description: '',
        videoLink: '',
        place: '',
        timezone: '',
      })
      const url = new URL(buildGoogleCalendarUrl(event))

      expect(url.searchParams.get('text')).toBe('')
      expect(url.searchParams.get('details')).toBe('')
      expect(url.searchParams.get('location')).toBe('')
    })

    it('combines videoLink and place in location', () => {
      const { buildGoogleCalendarUrl } = useGoogleCalendar()
      const event = createEvent({
        videoLink: 'https://zoom.us/456',
        place: 'Office #5',
      })
      const url = new URL(buildGoogleCalendarUrl(event))

      expect(url.searchParams.get('location')).toBe('https://zoom.us/456 Office #5')
    })

    it('adds RRULE for weekly repeating event', () => {
      const { buildGoogleCalendarUrl } = useGoogleCalendar()
      const event = createEvent({
        isRepeating: true,
        repeatPeriod: 'WEEKLY',
        repeatInterval: 1,
      })
      const url = new URL(buildGoogleCalendarUrl(event))

      expect(url.searchParams.get('recur')).toBe('RRULE:FREQ=WEEKLY;INTERVAL=1')
    })

    it('adds RRULE with custom interval for repeating event', () => {
      const { buildGoogleCalendarUrl } = useGoogleCalendar()
      const event = createEvent({
        isRepeating: true,
        repeatPeriod: 'DAILY',
        repeatInterval: 3,
      })
      const url = new URL(buildGoogleCalendarUrl(event))

      expect(url.searchParams.get('recur')).toBe('RRULE:FREQ=DAILY;INTERVAL=3')
    })

    it('adds RRULE with UNTIL when repeatEndDate is set', () => {
      const { buildGoogleCalendarUrl } = useGoogleCalendar()
      const event = createEvent({
        isRepeating: true,
        repeatPeriod: 'MONTHLY',
        repeatInterval: 1,
        repeatEndDate: '2024-12-31T23:59:00.000Z',
      })
      const url = new URL(buildGoogleCalendarUrl(event))

      expect(url.searchParams.get('recur')).toBe('RRULE:FREQ=MONTHLY;INTERVAL=1;UNTIL=20241231T235900Z')
    })

    it('defaults repeatInterval to 1 when not set', () => {
      const { buildGoogleCalendarUrl } = useGoogleCalendar()
      const event = createEvent({
        isRepeating: true,
        repeatPeriod: 'YEARLY',
      })
      const url = new URL(buildGoogleCalendarUrl(event))

      expect(url.searchParams.get('recur')).toBe('RRULE:FREQ=YEARLY;INTERVAL=1')
    })

    it('does not add recur param for non-repeating event', () => {
      const { buildGoogleCalendarUrl } = useGoogleCalendar()
      const event = createEvent({ isRepeating: false })
      const url = new URL(buildGoogleCalendarUrl(event))

      expect(url.searchParams.get('recur')).toBeNull()
    })
  })

  describe('openInGoogleCalendar', () => {
    it('opens the calendar URL via openLink (fallback to window.open в браузере)', () => {
      const openSpy = vi.spyOn(window, 'open').mockImplementation(() => null)
      const { openInGoogleCalendar, buildGoogleCalendarUrl } = useGoogleCalendar()
      const event = createEvent()

      openInGoogleCalendar(event, 90)

      const expectedUrl = buildGoogleCalendarUrl(event, 90)
      expect(openSpy).toHaveBeenCalledWith(expectedUrl, '_blank', 'noopener,noreferrer')

      openSpy.mockRestore()
    })
  })
})
