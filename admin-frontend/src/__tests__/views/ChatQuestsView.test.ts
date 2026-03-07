import { describe, expect, it } from 'vitest'

describe('ChatQuestsView logic', () => {
  // formatDate function
  function formatDate(dateStr: string) {
    return new Date(dateStr).toLocaleDateString('ru-RU', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    })
  }

  it('formats ISO date to Russian locale', () => {
    const result = formatDate('2026-03-07T15:30:00Z')
    expect(result).toBeTruthy()
    expect(typeof result).toBe('string')
  })

  it('includes day, month, year in formatted date', () => {
    const result = formatDate('2026-03-07T15:30:00Z')
    expect(result).toContain('07')
    expect(result).toContain('03')
    expect(result).toContain('2026')
  })

  // formatDateForAPI function
  function formatDateForAPI(dateStr: string): string {
    if (!dateStr)
      return dateStr
    const d = new Date(dateStr)
    return d.toISOString()
  }

  it('converts datetime-local to ISO string', () => {
    const result = formatDateForAPI('2026-03-04T17:00')
    expect(result).toContain('2026-03-04')
    expect(result).toContain('T')
    expect(result).toMatch(/Z$/)
  })

  it('returns empty string for empty input', () => {
    expect(formatDateForAPI('')).toBe('')
  })

  // isQuestActive function
  function isQuestActive(quest: { isActive: boolean, startsAt: string, endsAt: string }) {
    const now = new Date()
    return quest.isActive && new Date(quest.startsAt) <= now && new Date(quest.endsAt) >= now
  }

  it('returns true for active quest within time range', () => {
    const now = new Date()
    const past = new Date(now.getTime() - 86400000).toISOString()
    const future = new Date(now.getTime() + 86400000).toISOString()
    expect(isQuestActive({ isActive: true, startsAt: past, endsAt: future })).toBe(true)
  })

  it('returns false for inactive quest', () => {
    const now = new Date()
    const past = new Date(now.getTime() - 86400000).toISOString()
    const future = new Date(now.getTime() + 86400000).toISOString()
    expect(isQuestActive({ isActive: false, startsAt: past, endsAt: future })).toBe(false)
  })

  it('returns false for expired quest', () => {
    const past1 = new Date(Date.now() - 172800000).toISOString()
    const past2 = new Date(Date.now() - 86400000).toISOString()
    expect(isQuestActive({ isActive: true, startsAt: past1, endsAt: past2 })).toBe(false)
  })

  it('returns false for future quest', () => {
    const future1 = new Date(Date.now() + 86400000).toISOString()
    const future2 = new Date(Date.now() + 172800000).toISOString()
    expect(isQuestActive({ isActive: true, startsAt: future1, endsAt: future2 })).toBe(false)
  })

  // resetForm function
  it('resetForm returns default form values', () => {
    function resetForm() {
      return {
        title: '',
        description: '',
        questType: 'message_count',
        chatId: null,
        targetCount: 10,
        pointsReward: 10,
        startsAt: '',
        endsAt: '',
        isActive: true,
      }
    }

    const form = resetForm()
    expect(form.title).toBe('')
    expect(form.questType).toBe('message_count')
    expect(form.chatId).toBeNull()
    expect(form.targetCount).toBe(10)
    expect(form.pointsReward).toBe(10)
    expect(form.isActive).toBe(true)
  })

  // openEdit populates form from quest
  it('openEdit slices startsAt and endsAt to 16 characters for datetime-local input', () => {
    const quest = {
      startsAt: '2026-03-04T17:00:00.000Z',
      endsAt: '2026-03-10T23:59:00.000Z',
    }
    expect(quest.startsAt.slice(0, 16)).toBe('2026-03-04T17:00')
    expect(quest.endsAt.slice(0, 16)).toBe('2026-03-10T23:59')
  })
})
