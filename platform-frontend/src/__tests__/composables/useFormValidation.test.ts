import { describe, expect, it } from 'vitest'
import { maxLength, minLength, required, useFormValidation } from '@/composables/useFormValidation'

describe('required', () => {
  it('returns error message for empty string', () => {
    const validate = required()
    expect(validate('')).toBe('Обязательное поле')
  })

  it('returns error message for whitespace-only string', () => {
    const validate = required()
    expect(validate('   ')).toBe('Обязательное поле')
  })

  it('returns error message for null', () => {
    const validate = required()
    expect(validate(null)).toBe('Обязательное поле')
  })

  it('returns error message for undefined', () => {
    const validate = required()
    expect(validate(undefined)).toBe('Обязательное поле')
  })

  it('returns null for non-empty string', () => {
    const validate = required()
    expect(validate('hello')).toBeNull()
  })

  it('uses custom message when provided', () => {
    const validate = required('Поле обязательно')
    expect(validate('')).toBe('Поле обязательно')
  })
})

describe('minLength', () => {
  it('returns error for strings shorter than n after trim', () => {
    const validate = minLength(5)
    expect(validate('abc')).toBe('Минимум 5 символов')
  })

  it('returns error for strings shorter than n when padded with spaces', () => {
    const validate = minLength(5)
    expect(validate('  ab  ')).toBe('Минимум 5 символов')
  })

  it('returns null for strings with length >= n', () => {
    const validate = minLength(3)
    expect(validate('abc')).toBeNull()
    expect(validate('abcd')).toBeNull()
  })

  it('uses custom message when provided', () => {
    const validate = minLength(5, 'Слишком коротко')
    expect(validate('ab')).toBe('Слишком коротко')
  })
})

describe('maxLength', () => {
  it('returns error for strings longer than n after trim', () => {
    const validate = maxLength(3)
    expect(validate('abcd')).toBe('Максимум 3 символов')
  })

  it('returns null for strings with length <= n', () => {
    const validate = maxLength(5)
    expect(validate('abc')).toBeNull()
    expect(validate('abcde')).toBeNull()
  })

  it('returns null for strings that are long only due to whitespace', () => {
    const validate = maxLength(3)
    expect(validate('  ab  ')).toBeNull()
  })

  it('uses custom message when provided', () => {
    const validate = maxLength(3, 'Слишком длинно')
    expect(validate('abcd')).toBe('Слишком длинно')
  })
})

describe('useFormValidation', () => {
  it('errors starts with all fields undefined', () => {
    const { errors } = useFormValidation({
      name: [required()],
      email: [required()],
    })

    expect(errors.name).toBeUndefined()
    expect(errors.email).toBeUndefined()
  })

  it('validateField sets error when validation fails', () => {
    const { errors, validateField } = useFormValidation({
      name: [required()],
    })

    validateField('name', '')
    expect(errors.name).toBe('Обязательное поле')
  })

  it('validateField clears error when validation passes', () => {
    const { errors, validateField } = useFormValidation({
      name: [required()],
    })

    validateField('name', '')
    expect(errors.name).toBe('Обязательное поле')

    validateField('name', 'John')
    expect(errors.name).toBeUndefined()
  })

  it('validateAll validates all fields and returns false if any fail', () => {
    const { errors, validateAll } = useFormValidation({
      name: [required()],
      email: [required()],
    })

    const result = validateAll({ name: 'John', email: '' })

    expect(result).toBe(false)
    expect(errors.name).toBeUndefined()
    expect(errors.email).toBe('Обязательное поле')
  })

  it('validateAll returns true when all pass', () => {
    const { validateAll } = useFormValidation({
      name: [required()],
      email: [required()],
    })

    const result = validateAll({ name: 'John', email: 'john@example.com' })
    expect(result).toBe(true)
  })

  it('clearErrors resets all errors to undefined', () => {
    const { errors, validateAll, clearErrors } = useFormValidation({
      name: [required()],
      email: [required()],
    })

    validateAll({ name: '', email: '' })
    expect(errors.name).toBe('Обязательное поле')
    expect(errors.email).toBe('Обязательное поле')

    clearErrors()
    expect(errors.name).toBeUndefined()
    expect(errors.email).toBeUndefined()
  })

  it('works with multiple fields and multiple rules per field', () => {
    const { errors, validateField } = useFormValidation({
      name: [required(), minLength(3), maxLength(20)],
      bio: [maxLength(100)],
    })

    validateField('name', 'ab')
    expect(errors.name).toBe('Минимум 3 символов')

    validateField('name', 'abc')
    expect(errors.name).toBeUndefined()

    validateField('bio', 'a'.repeat(101))
    expect(errors.bio).toBe('Максимум 100 символов')
  })
})
