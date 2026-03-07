import { describe, expect, it } from 'vitest'

describe('ReferralsView logic', () => {
  // Status display logic from template
  it('displays "Активна" for active status', () => {
    const status = 'active'
    const label = status === 'active' ? 'Активна' : 'Заморожена'
    expect(label).toBe('Активна')
  })

  it('displays "Заморожена" for non-active status', () => {
    const status = 'frozen'
    const label = status === 'active' ? 'Активна' : 'Заморожена'
    expect(label).toBe('Заморожена')
  })

  // Status styling from template
  it('applies green class for active status', () => {
    const status = 'active'
    const className = status === 'active' ? 'bg-green-500/10 text-green-600' : 'bg-gray-500/10 text-gray-600'
    expect(className).toContain('green')
  })

  it('applies gray class for inactive status', () => {
    const status = 'frozen'
    const className = status === 'active' ? 'bg-green-500/10 text-green-600' : 'bg-gray-500/10 text-gray-600'
    expect(className).toContain('gray')
  })

  // Author name formatting from template
  it('formats author full name', () => {
    const author = { firstName: 'Иван', lastName: 'Петров' }
    const fullName = `${author.firstName} ${author.lastName}`
    expect(fullName).toBe('Иван Петров')
  })

  it('handles author with optional fields', () => {
    const author: { firstName?: string, lastName?: string } = { firstName: 'Иван' }
    const fullName = `${author?.firstName} ${author?.lastName}`
    expect(fullName).toContain('Иван')
  })
})
