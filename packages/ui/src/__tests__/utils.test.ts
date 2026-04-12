import { describe, expect, it } from 'vitest'
import { cn } from '../lib/utils'

describe('cn utility', () => {
  it('merges class names', () => {
    expect(cn('foo', 'bar')).toBe('foo bar')
  })

  it('handles conditional classes', () => {
    expect(cn('base', false && 'hidden', 'end')).toBe('base end')
  })

  it('deduplicates tailwind classes', () => {
    expect(cn('p-4', 'p-2')).toBe('p-2')
  })

  it('handles undefined and null', () => {
    expect(cn('base', undefined, null, 'end')).toBe('base end')
  })

  it('handles empty input', () => {
    expect(cn()).toBe('')
  })

  it('merges conflicting tailwind utilities', () => {
    expect(cn('text-red-500', 'text-blue-500')).toBe('text-blue-500')
  })

  it('preserves non-conflicting classes', () => {
    expect(cn('rounded-sm', 'border', 'bg-card')).toBe('rounded-sm border bg-card')
  })

  it('works with arrays via clsx', () => {
    expect(cn(['foo', 'bar'], 'baz')).toBe('foo bar baz')
  })
})
