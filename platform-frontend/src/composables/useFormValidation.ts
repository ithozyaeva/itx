import { reactive } from 'vue'

type Validator = (value: unknown) => string | null

export function required(message = 'Обязательное поле'): Validator {
  return (value: unknown) => {
    if (typeof value === 'string' && !value.trim())
      return message
    if (value === null || value === undefined)
      return message
    return null
  }
}

export function minLength(min: number, message?: string): Validator {
  return (value: unknown) => {
    if (typeof value === 'string' && value.trim().length < min)
      return message || `Минимум ${min} символов`
    return null
  }
}

export function maxLength(max: number, message?: string): Validator {
  return (value: unknown) => {
    if (typeof value === 'string' && value.trim().length > max)
      return message || `Максимум ${max} символов`
    return null
  }
}

interface FieldRules {
  [field: string]: Validator[]
}

export function useFormValidation(rules: FieldRules) {
  const errors = reactive<Record<string, string | undefined>>(
    Object.fromEntries(Object.keys(rules).map(k => [k, undefined])),
  )

  function validateField(field: string, value: unknown): boolean {
    const fieldRules = rules[field]
    if (!fieldRules)
      return true

    for (const rule of fieldRules) {
      const error = rule(value)
      if (error) {
        errors[field] = error
        return false
      }
    }
    errors[field] = undefined
    return true
  }

  function validateAll(values: Record<string, any>): boolean {
    let valid = true
    for (const field of Object.keys(rules)) {
      if (!validateField(field, values[field]))
        valid = false
    }
    return valid
  }

  function clearErrors() {
    for (const field of Object.keys(rules)) {
      errors[field] = undefined
    }
  }

  return { errors, validateField, validateAll, clearErrors }
}
