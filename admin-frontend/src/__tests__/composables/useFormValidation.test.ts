import { describe, expect, it } from 'vitest'
import { requiredArrayRule, requiredRule, useFormValidation } from '@/composables/useFormValidation'

describe('requiredRule', () => {
  it('returns false for empty string', () => {
    expect(requiredRule.validate('')).toBe(false)
  })

  it('returns false for whitespace-only string', () => {
    expect(requiredRule.validate('   ')).toBe(false)
  })

  it('returns true for non-empty string', () => {
    expect(requiredRule.validate('hello')).toBe(true)
  })
})

describe('requiredArrayRule', () => {
  it('returns false for empty array', () => {
    expect(requiredArrayRule.validate([])).toBe(false)
  })

  it('returns false for non-array value', () => {
    expect(requiredArrayRule.validate(null as any)).toBe(false)
    expect(requiredArrayRule.validate(undefined as any)).toBe(false)
  })

  it('returns true for array with items', () => {
    expect(requiredArrayRule.validate([1, 2])).toBe(true)
  })
})

describe('useFormValidation', () => {
  const initialValues = { name: '', email: '' }

  it('has initial values matching initialValues', () => {
    const { values, errors, touched } = useFormValidation(initialValues)

    expect(values.value).toEqual({ name: '', email: '' })
    expect(errors.value.name).toBeNull()
    expect(errors.value.email).toBeNull()
    expect(touched.value.name).toBe(false)
    expect(touched.value.email).toBe(false)
  })

  it('validateField returns true for field without rules', async () => {
    const { validateField } = useFormValidation(initialValues)

    const result = await validateField('name')

    expect(result).toBe(true)
  })

  it('validateField sets error for invalid field', async () => {
    const { validateField, errors } = useFormValidation(initialValues, {
      name: [requiredRule],
    })

    const result = await validateField('name')

    expect(result).toBe(false)
    expect(errors.value.name).toBe('Это поле обязательно для заполнения')
  })

  it('validate validates all fields and marks them as touched', async () => {
    const { validate, touched } = useFormValidation(initialValues, {
      name: [requiredRule],
    })

    await validate()

    expect(touched.value.name).toBe(true)
    expect(touched.value.email).toBe(true)
  })

  it('validate returns false if any field is invalid', async () => {
    const { validate } = useFormValidation(initialValues, {
      name: [requiredRule],
    })

    const result = await validate()

    expect(result).toBe(false)
  })

  it('validate returns true if all fields are valid', async () => {
    const { validate, values } = useFormValidation({ name: '', email: '' }, {
      name: [requiredRule],
    })

    values.value.name = 'John'
    const result = await validate()

    expect(result).toBe(true)
  })

  it('handleBlur marks field as touched', () => {
    const { handleBlur, touched } = useFormValidation(initialValues)

    handleBlur('name')

    expect(touched.value.name).toBe(true)
    expect(touched.value.email).toBe(false)
  })

  it('resetForm resets to initial values', async () => {
    const { values, errors, touched, resetForm, validateField } = useFormValidation(initialValues, {
      name: [requiredRule],
    })

    values.value.name = 'John'
    touched.value.name = true
    await validateField('name')

    resetForm()

    expect(values.value).toEqual({ name: '', email: '' })
    expect(errors.value.name).toBeNull()
    expect(touched.value.name).toBe(false)
  })

  it('setValues updates values', () => {
    const { values, setValues } = useFormValidation(initialValues)

    setValues({ name: 'Jane' })

    expect(values.value.name).toBe('Jane')
    expect(values.value.email).toBe('')
  })

  it('isValid computed is true when no errors', () => {
    const { isValid } = useFormValidation(initialValues)

    expect(isValid.value).toBe(true)
  })

  it('isValid computed is false when there are errors', async () => {
    const { isValid, validateField } = useFormValidation(initialValues, {
      name: [requiredRule],
    })

    await validateField('name')

    expect(isValid.value).toBe(false)
  })
})
