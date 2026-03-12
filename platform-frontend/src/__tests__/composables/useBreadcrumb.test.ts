import { describe, expect, it } from 'vitest'
import { useBreadcrumb } from '@/composables/useBreadcrumb'

describe('useBreadcrumb', () => {
  it('dynamicLabel starts as null', () => {
    const { dynamicLabel } = useBreadcrumb()
    expect(dynamicLabel.value).toBeNull()
  })

  it('setLabel sets the label', () => {
    const { dynamicLabel, setLabel } = useBreadcrumb()

    setLabel('Профиль')
    expect(dynamicLabel.value).toBe('Профиль')
  })

  it('clearLabel resets to null', () => {
    const { dynamicLabel, setLabel, clearLabel } = useBreadcrumb()

    setLabel('Профиль')
    expect(dynamicLabel.value).toBe('Профиль')

    clearLabel()
    expect(dynamicLabel.value).toBeNull()
  })

  it('multiple calls share the same singleton state', () => {
    const first = useBreadcrumb()
    const second = useBreadcrumb()

    expect(first.dynamicLabel).toBe(second.dynamicLabel)

    first.setLabel('Тест')
    expect(second.dynamicLabel.value).toBe('Тест')

    second.clearLabel()
    expect(first.dynamicLabel.value).toBeNull()
  })
})
