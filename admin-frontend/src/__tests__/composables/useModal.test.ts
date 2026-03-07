import { describe, expect, it } from 'vitest'
import { useModal } from '@/composables/useModal'

describe('useModal', () => {
  it('starts closed', () => {
    const { isOpen } = useModal()

    expect(isOpen.value).toBe(false)
  })

  it('open sets isOpen to true', () => {
    const { isOpen, open } = useModal()

    open()

    expect(isOpen.value).toBe(true)
  })

  it('close sets isOpen to false', () => {
    const { isOpen, open, close } = useModal()

    open()
    close()

    expect(isOpen.value).toBe(false)
  })
})
