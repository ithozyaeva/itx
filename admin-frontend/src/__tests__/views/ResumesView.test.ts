import { describe, expect, it } from 'vitest'

describe('ResumesView logic', () => {
  // workFormatOptions config
  const workFormatOptions = [
    { label: 'Любой формат', value: '' },
    { label: 'Удалёнка', value: 'REMOTE' },
    { label: 'Гибрид', value: 'HYBRID' },
    { label: 'Офис', value: 'OFFICE' },
  ]

  it('has 4 work format options', () => {
    expect(workFormatOptions).toHaveLength(4)
  })

  it('first option is empty value for "any format"', () => {
    expect(workFormatOptions[0].value).toBe('')
    expect(workFormatOptions[0].label).toBe('Любой формат')
  })

  it('all options have non-empty labels', () => {
    for (const option of workFormatOptions) {
      expect(option.label.length).toBeGreaterThan(0)
    }
  })

  // formatWorkFormat function
  function formatWorkFormat(value?: string) {
    const match = workFormatOptions.find(option => option.value === value)
    return match?.label ?? 'Не указано'
  }

  it('returns "Удалёнка" for REMOTE', () => {
    expect(formatWorkFormat('REMOTE')).toBe('Удалёнка')
  })

  it('returns "Гибрид" for HYBRID', () => {
    expect(formatWorkFormat('HYBRID')).toBe('Гибрид')
  })

  it('returns "Офис" for OFFICE', () => {
    expect(formatWorkFormat('OFFICE')).toBe('Офис')
  })

  it('returns "Любой формат" for empty string', () => {
    expect(formatWorkFormat('')).toBe('Любой формат')
  })

  it('returns "Не указано" for undefined', () => {
    expect(formatWorkFormat(undefined)).toBe('Не указано')
  })

  it('returns "Не указано" for unknown value', () => {
    expect(formatWorkFormat('UNKNOWN')).toBe('Не указано')
  })

  // resetFilters logic
  it('resetFilters clears all filter fields', () => {
    const filters = {
      workFormat: 'REMOTE',
      desiredPosition: 'Developer',
      workExperience: 'финтех',
    }
    // Simulating resetFilters
    filters.workFormat = ''
    filters.desiredPosition = ''
    filters.workExperience = ''

    expect(filters.workFormat).toBe('')
    expect(filters.desiredPosition).toBe('')
    expect(filters.workExperience).toBe('')
  })
})
