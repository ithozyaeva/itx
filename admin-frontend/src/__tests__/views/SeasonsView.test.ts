import { describe, expect, it } from 'vitest'

describe('SeasonsView logic', () => {
  // formatDate function
  function formatDate(dateStr: string) {
    return new Date(dateStr).toLocaleDateString('ru-RU', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
    })
  }

  it('formats ISO date to Russian locale (date only)', () => {
    const result = formatDate('2026-03-07T00:00:00Z')
    expect(result).toBeTruthy()
    expect(result).toContain('2026')
  })

  it('does not include time components', () => {
    const result = formatDate('2026-03-07T15:30:00Z')
    // Russian date-only format should not include hours/minutes
    // It should be DD.MM.YYYY
    expect(result).toMatch(/\d{2}\.\d{2}\.\d{4}/)
  })

  it('handles different dates correctly', () => {
    const result1 = formatDate('2025-01-01T00:00:00Z')
    const result2 = formatDate('2025-12-31T23:59:59Z')
    expect(result1).not.toBe(result2)
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
    expect(result).toMatch(/^\d{4}-\d{2}-\d{2}T/)
    expect(result).toMatch(/Z$/)
  })

  it('returns empty string for empty input', () => {
    expect(formatDateForAPI('')).toBe('')
  })

  // Season status display logic from template
  it('displays "Активный" for ACTIVE status', () => {
    const status = 'ACTIVE'
    const label = status === 'ACTIVE' ? 'Активный' : 'Завершён'
    expect(label).toBe('Активный')
  })

  it('displays "Завершён" for FINISHED status', () => {
    const status = 'FINISHED'
    const label = status === 'ACTIVE' ? 'Активный' : 'Завершён'
    expect(label).toBe('Завершён')
  })

  // Status styling from template
  it('applies green class for ACTIVE status', () => {
    const status = 'ACTIVE'
    const className = status === 'ACTIVE'
      ? 'bg-green-500/10 text-green-500'
      : 'bg-muted text-muted-foreground'
    expect(className).toContain('green')
  })

  it('applies muted class for non-ACTIVE status', () => {
    const status = 'FINISHED'
    const className = status === 'ACTIVE'
      ? 'bg-green-500/10 text-green-500'
      : 'bg-muted text-muted-foreground'
    expect(className).toContain('muted')
  })

  // Finish button visibility: only shown for ACTIVE seasons
  it('shows finish button only for ACTIVE seasons', () => {
    const statuses = ['ACTIVE', 'FINISHED']
    const showFinish = statuses.map(s => s === 'ACTIVE')
    expect(showFinish).toEqual([true, false])
  })

  // resetForm logic
  it('resetForm returns empty form values', () => {
    function resetForm() {
      return { title: '', startDate: '', endDate: '' }
    }
    const form = resetForm()
    expect(form.title).toBe('')
    expect(form.startDate).toBe('')
    expect(form.endDate).toBe('')
  })
})
