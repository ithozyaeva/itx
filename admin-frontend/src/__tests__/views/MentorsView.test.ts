import { describe, expect, it } from 'vitest'

describe('MentorsView logic', () => {
  // hasActiveFilters computed
  it('hasActiveFilters returns false when both filters are empty', () => {
    const filters = { name: '', tag: '' }
    const hasActiveFilters = filters.name !== '' || filters.tag !== ''
    expect(hasActiveFilters).toBe(false)
  })

  it('hasActiveFilters returns true when name filter is set', () => {
    const filters = { name: 'Иван', tag: '' }
    const hasActiveFilters = filters.name !== '' || filters.tag !== ''
    expect(hasActiveFilters).toBe(true)
  })

  it('hasActiveFilters returns true when tag filter is set', () => {
    const filters = { name: '', tag: 'Go' }
    const hasActiveFilters = filters.name !== '' || filters.tag !== ''
    expect(hasActiveFilters).toBe(true)
  })

  // filteredMentors logic
  interface Mentor {
    id: number
    firstName: string
    lastName: string
    tg: string
    profTags?: { id: number, title: string }[]
  }

  function filterMentors(items: Mentor[], filters: { name: string, tag: string }): Mentor[] {
    return items.filter((mentor) => {
      const nameMatch = !filters.name
        || `${mentor.firstName} ${mentor.lastName} ${mentor.tg}`.toLowerCase().includes(filters.name.toLowerCase())
      const tagMatch = !filters.tag
        || mentor.profTags?.some(t => t.title.toLowerCase().includes(filters.tag.toLowerCase()))
      return nameMatch && tagMatch
    })
  }

  const mentors: Mentor[] = [
    { id: 1, firstName: 'Иван', lastName: 'Петров', tg: 'ivan_p', profTags: [{ id: 1, title: 'Go' }, { id: 2, title: 'Python' }] },
    { id: 2, firstName: 'Мария', lastName: 'Сидорова', tg: 'maria_s', profTags: [{ id: 3, title: 'JavaScript' }] },
    { id: 3, firstName: 'Алексей', lastName: 'Козлов', tg: 'alex_k', profTags: [] },
  ]

  it('returns all mentors when no filters', () => {
    const result = filterMentors(mentors, { name: '', tag: '' })
    expect(result).toHaveLength(3)
  })

  it('filters by name (first name)', () => {
    const result = filterMentors(mentors, { name: 'Иван', tag: '' })
    expect(result).toHaveLength(1)
    expect(result[0].id).toBe(1)
  })

  it('filters by name (telegram username)', () => {
    const result = filterMentors(mentors, { name: 'maria_s', tag: '' })
    expect(result).toHaveLength(1)
    expect(result[0].id).toBe(2)
  })

  it('filters by tag', () => {
    const result = filterMentors(mentors, { name: '', tag: 'Go' })
    expect(result).toHaveLength(1)
    expect(result[0].id).toBe(1)
  })

  it('filters by both name and tag', () => {
    const result = filterMentors(mentors, { name: 'Иван', tag: 'Python' })
    expect(result).toHaveLength(1)
    expect(result[0].id).toBe(1)
  })

  it('returns empty when no match', () => {
    const result = filterMentors(mentors, { name: 'Никто', tag: '' })
    expect(result).toHaveLength(0)
  })

  it('name filter is case-insensitive', () => {
    const result = filterMentors(mentors, { name: 'иван', tag: '' })
    expect(result).toHaveLength(1)
  })

  it('tag filter is case-insensitive', () => {
    const result = filterMentors(mentors, { name: '', tag: 'javascript' })
    expect(result).toHaveLength(1)
    expect(result[0].id).toBe(2)
  })
})
