import { describe, expect, it, vi } from 'vitest'

const mockLoginWithTelegram = vi.fn()
const mockHandleError = vi.fn()

vi.mock('@/services/authService', () => ({
  loginWithTelegram: (...args: unknown[]) => mockLoginWithTelegram(...args),
}))

vi.mock('@/services/errorService', () => ({
  handleError: (...args: unknown[]) => mockHandleError(...args),
}))

vi.mock('@/components/ui/button', () => ({
  Button: { template: '<button :disabled="disabled"><slot /></button>', props: ['disabled'] },
}))

import { mount } from '@vue/test-utils'
import TelegramAuth from '@/components/TelegramAuth.vue'

describe('TelegramAuth', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    // Ensure no token in URL for most tests
    Object.defineProperty(window, 'location', {
      value: { search: '', pathname: '/', href: 'http://localhost/' },
      writable: true,
    })
    window.history.replaceState = vi.fn()
  })

  it('renders the auth button', () => {
    const wrapper = mount(TelegramAuth)
    expect(wrapper.find('button').exists()).toBe(true)
    expect(wrapper.text()).toContain('Зайти через ТГ')
  })

  it('does not show loading state initially', () => {
    const wrapper = mount(TelegramAuth)
    expect(wrapper.find('.loading').exists()).toBe(false)
  })

  it('calls loginWithTelegram when token is present in URL', async () => {
    mockLoginWithTelegram.mockResolvedValue({ token: 'jwt-token' })
    Object.defineProperty(window, 'location', {
      value: { search: '?token=test-token', pathname: '/', href: 'http://localhost/?token=test-token' },
      writable: true,
    })

    mount(TelegramAuth)
    // Wait for onMounted + async handler
    await vi.dynamicImportSettled()

    expect(mockLoginWithTelegram).toHaveBeenCalledWith('test-token')
  })

  it('opens telegram bot link when button is clicked', () => {
    const openSpy = vi.spyOn(window, 'open').mockImplementation(() => null)

    const wrapper = mount(TelegramAuth)
    wrapper.find('button').trigger('click')

    expect(openSpy).toHaveBeenCalledWith(
      expect.stringContaining('https://t.me/'),
      '_blank',
    )
    openSpy.mockRestore()
  })

  it('handles login error gracefully', async () => {
    const error = new Error('Auth failed')
    mockLoginWithTelegram.mockRejectedValue(error)
    Object.defineProperty(window, 'location', {
      value: { search: '?token=bad-token', pathname: '/', href: 'http://localhost/?token=bad-token' },
      writable: true,
    })

    mount(TelegramAuth)
    await vi.dynamicImportSettled()

    expect(mockHandleError).toHaveBeenCalledWith(error)
  })
})
