import { describe, expect, it } from 'vitest'
import { cleanParams, datetimeLocalToISO, formatDateToInput, parseTimezoneOffsetMinutes, toDatetimeLocal } from '@/lib/utils'

describe('utils', () => {
  describe('formatDateToInput', () => {
    it('extracts date portion from ISO string', () => {
      expect(formatDateToInput('2026-03-07T15:30:00Z')).toBe('2026-03-07')
    })

    it('returns empty string for undefined', () => {
      expect(formatDateToInput()).toBe('')
    })

    it('returns empty string for empty string', () => {
      expect(formatDateToInput('')).toBe('')
    })

    it('handles date-only strings', () => {
      expect(formatDateToInput('2026-12-25')).toBe('2026-12-25')
    })
  })

  describe('parseTimezoneOffsetMinutes', () => {
    it('returns 0 for UTC', () => {
      expect(parseTimezoneOffsetMinutes('UTC')).toBe(0)
    })

    it('returns 0 for empty string', () => {
      expect(parseTimezoneOffsetMinutes('')).toBe(0)
    })

    it('parses positive offset', () => {
      expect(parseTimezoneOffsetMinutes('UTC+3')).toBe(180)
    })

    it('parses negative offset', () => {
      expect(parseTimezoneOffsetMinutes('UTC-5')).toBe(-300)
    })

    it('parses large offset', () => {
      expect(parseTimezoneOffsetMinutes('UTC+12')).toBe(720)
    })

    it('returns 0 for invalid format', () => {
      expect(parseTimezoneOffsetMinutes('EST')).toBe(0)
      expect(parseTimezoneOffsetMinutes('GMT+3')).toBe(0)
    })
  })

  describe('toDatetimeLocal', () => {
    it('converts UTC ISO to datetime-local in UTC', () => {
      const result = toDatetimeLocal('2026-03-07T15:30:00Z')
      expect(result).toBe('2026-03-07T15:30')
    })

    it('converts UTC ISO to datetime-local in UTC+3', () => {
      const result = toDatetimeLocal('2026-03-07T15:30:00Z', 'UTC+3')
      expect(result).toBe('2026-03-07T18:30')
    })

    it('converts UTC ISO to datetime-local in UTC-5', () => {
      const result = toDatetimeLocal('2026-03-07T15:30:00Z', 'UTC-5')
      expect(result).toBe('2026-03-07T10:30')
    })

    it('handles day boundary crossing', () => {
      const result = toDatetimeLocal('2026-03-07T23:00:00Z', 'UTC+3')
      expect(result).toBe('2026-03-08T02:00')
    })
  })

  describe('datetimeLocalToISO', () => {
    it('converts datetime-local in UTC to ISO', () => {
      const result = datetimeLocalToISO('2026-03-07T15:30', 'UTC')
      expect(result).toBe('2026-03-07T15:30:00.000Z')
    })

    it('converts datetime-local in UTC+3 to ISO (subtracts offset)', () => {
      const result = datetimeLocalToISO('2026-03-07T18:30', 'UTC+3')
      expect(result).toBe('2026-03-07T15:30:00.000Z')
    })

    it('converts datetime-local in UTC-5 to ISO (adds offset)', () => {
      const result = datetimeLocalToISO('2026-03-07T10:30', 'UTC-5')
      expect(result).toBe('2026-03-07T15:30:00.000Z')
    })

    it('is inverse of toDatetimeLocal', () => {
      const original = '2026-06-15T12:00:00.000Z'
      const timezone = 'UTC+3'
      const local = toDatetimeLocal(original, timezone)
      const backToISO = datetimeLocalToISO(local, timezone)
      expect(backToISO).toBe(original)
    })
  })

  describe('cleanParams', () => {
    it('removes undefined values', () => {
      expect(cleanParams({ a: 'hello', b: undefined })).toEqual({ a: 'hello' })
    })

    it('removes null values', () => {
      expect(cleanParams({ a: 'hello', b: null })).toEqual({ a: 'hello' })
    })

    it('removes empty string values', () => {
      expect(cleanParams({ a: 'hello', b: '' })).toEqual({ a: 'hello' })
    })

    it('removes empty arrays', () => {
      expect(cleanParams({ a: 'hello', b: [] })).toEqual({ a: 'hello' })
    })

    it('keeps non-empty arrays', () => {
      expect(cleanParams({ a: [1, 2] })).toEqual({ a: [1, 2] })
    })

    it('keeps zero values', () => {
      expect(cleanParams({ a: 0 })).toEqual({ a: 0 })
    })

    it('keeps false values', () => {
      expect(cleanParams({ a: false })).toEqual({ a: false })
    })

    it('handles all clean params', () => {
      expect(cleanParams({ a: 'x', b: 1, c: true })).toEqual({ a: 'x', b: 1, c: true })
    })

    it('handles all dirty params', () => {
      expect(cleanParams({ a: undefined, b: null, c: '', d: [] })).toEqual({})
    })
  })
})
