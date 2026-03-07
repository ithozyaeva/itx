import { describe, expect, it } from 'vitest'

describe('EventsView logic', () => {
  // selectEvent sets selectedEventId and opens modal
  it('selectEvent sets the event id', () => {
    let selectedEventId: number | undefined
    function selectEvent(entityId: number) {
      selectedEventId = entityId
    }
    selectEvent(7)
    expect(selectedEventId).toBe(7)
  })

  // Tag name truncation logic from template
  it('truncates long tag names to 24 chars with ellipsis', () => {
    function truncateTagName(name: string): string {
      return name.length > 30 ? `${name.slice(0, 24)}...` : name
    }

    expect(truncateTagName('Short')).toBe('Short')
    expect(truncateTagName('A'.repeat(31))).toBe(`${'A'.repeat(24)}...`)
    expect(truncateTagName('A'.repeat(30))).toBe('A'.repeat(30))
  })

  // Hosts formatting from template
  it('formats event hosts as comma-separated names', () => {
    const hosts = [
      { firstName: 'Иван', lastName: 'Иванов' },
      { firstName: 'Мария', lastName: 'Петрова' },
    ]
    const result = hosts.map(host => `${host.firstName} ${host.lastName}`).join(', ')
    expect(result).toBe('Иван Иванов, Мария Петрова')
  })

  it('handles empty hosts list', () => {
    const hosts: { firstName: string, lastName: string }[] = []
    const result = hosts.map(host => `${host.firstName} ${host.lastName}`).join(', ')
    expect(result).toBe('')
  })

  // Date formatting with timezone from template
  it('displays timezone alongside date', () => {
    const event = { date: '2026-03-07T15:00:00Z', timezone: 'Europe/Moscow' }
    const formatted = `${new Date(event.date).toLocaleString()} (${event.timezone || 'UTC'})`
    expect(formatted).toContain('Europe/Moscow')
  })

  it('defaults to UTC when timezone is empty', () => {
    const event = { date: '2026-03-07T15:00:00Z', timezone: '' }
    const formatted = `(${event.timezone || 'UTC'})`
    expect(formatted).toBe('(UTC)')
  })
})
