import { describe, expect, it } from 'vitest'
import { cn, dateFormatter, wrapLinks } from '@/lib/utils'

describe('utils', () => {
  describe('cn()', () => {
    it('merges simple class names', () => {
      expect(cn('foo', 'bar')).toBe('foo bar')
    })

    it('handles conditional classes via clsx', () => {
      expect(cn('base', false && 'hidden', 'visible')).toBe('base visible')
    })

    it('merges tailwind classes correctly (later wins)', () => {
      expect(cn('p-4', 'p-2')).toBe('p-2')
    })

    it('merges conflicting tailwind utilities', () => {
      expect(cn('text-red-500', 'text-blue-500')).toBe('text-blue-500')
    })

    it('keeps non-conflicting tailwind classes', () => {
      const result = cn('p-4 text-red-500', 'mt-2')
      expect(result).toContain('p-4')
      expect(result).toContain('text-red-500')
      expect(result).toContain('mt-2')
    })

    it('handles empty input', () => {
      expect(cn()).toBe('')
    })

    it('handles undefined and null values', () => {
      expect(cn('a', undefined, null, 'b')).toBe('a b')
    })

    it('handles array input', () => {
      expect(cn(['foo', 'bar'])).toBe('foo bar')
    })

    it('handles object input', () => {
      expect(cn({ foo: true, bar: false, baz: true })).toBe('foo baz')
    })
  })

  describe('wrapLinks()', () => {
    it('wraps HTTP URLs in anchor tags', () => {
      const result = wrapLinks('Visit http://example.com for more')
      expect(result).toContain('<a href="http://example.com"')
      expect(result).toContain('target="_blank"')
      expect(result).toContain('rel="noopener noreferrer"')
    })

    it('wraps HTTPS URLs in anchor tags', () => {
      const result = wrapLinks('Visit https://example.com for more')
      expect(result).toContain('<a href="https://example.com"')
    })

    it('leaves text without URLs unchanged', () => {
      expect(wrapLinks('no links here')).toBe('no links here')
    })

    it('handles multiple URLs', () => {
      const result = wrapLinks('See https://a.com and https://b.com')
      expect(result).toContain('href="https://a.com"')
      expect(result).toContain('href="https://b.com"')
    })

    it('adds line break before link', () => {
      const result = wrapLinks('Check https://example.com')
      expect(result).toContain('<br />')
    })
  })

  describe('dateFormatter', () => {
    it('is an instance of Intl.DateTimeFormat', () => {
      expect(dateFormatter).toBeInstanceOf(Intl.DateTimeFormat)
    })

    it('formats dates in Russian locale', () => {
      const options = dateFormatter.resolvedOptions()
      expect(options.locale).toContain('ru')
    })
  })
})
