import { computed, ref } from 'vue'

export interface ValidationRule<T = any> {
  validate: (value: T, formValues?: any) => boolean | Promise<boolean>
  message: string
}

export type ValidationRules<T> = {
  [K in keyof T]?: ValidationRule[]
}

export const requiredRule = {
  validate: (value: string) => !!value?.toString()?.trim(),
  message: 'Это поле обязательно для заполнения',
}

export function useFormValidation<T extends Record<string, any>>(
  initialValues: T,
  validationRules: ValidationRules<T> = {},
) {
  const values = ref<T>({ ...initialValues })
  const errors = ref<Record<keyof T, string | null>>({} as any)
  const touched = ref<Record<keyof T, boolean>>({} as any)
  const isValidating = ref(false)

  Object.keys(initialValues).forEach((key) => {
    errors.value[key] = null
    touched.value[key] = false
  })

  async function validateField(field: keyof T): Promise<boolean> {
    const fieldRules = validationRules[field] || []

    if (fieldRules.length === 0) {
      errors.value[field] = null
      return true
    }

    for (const rule of fieldRules) {
      try {
        const isValid = await rule.validate(values.value[field], values.value)

        if (!isValid) {
          errors.value[field] = rule.message
          return false
        }
      }
      catch (error) {
        console.error(`Ошибка при валидации поля ${String(field)}:`, error)
        errors.value[field] = 'Ошибка валидации'
        return false
      }
    }

    errors.value[field] = null
    return true
  }

  async function validate(): Promise<boolean> {
    isValidating.value = true

    try {
      const fields = Object.keys(values.value) as Array<keyof T>
      const results = await Promise.all(fields.map(field => validateField(field)))

      // Отмечаем все поля как touched
      fields.forEach((field) => {
        touched.value[field] = true
      })

      return results.every(Boolean)
    }
    finally {
      isValidating.value = false
    }
  }

  function handleBlur(field: keyof T) {
    touched.value[field] = true
    validateField(field)
  }

  function resetForm() {
    values.value = { ...initialValues }

    Object.keys(initialValues).forEach((key) => {
      errors.value[key] = null
      touched.value[key] = false
    })
  }

  function setValues(newValues: Partial<T>) {
    values.value = { ...values.value, ...newValues }
  }

  const isValid = computed(() => {
    return Object.values(errors.value).every(error => error === null)
  })

  return {
    values,
    errors,
    touched,
    isValidating,
    isValid,
    validate,
    validateField,
    handleBlur,
    resetForm,
    setValues,
  }
}
