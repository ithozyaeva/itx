import { flushPromises, mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

const { mockGetAll, mockCreate, mockRequestPurchase, mockCancelPurchase, mockMarkSold, mockRemove } = vi.hoisted(() => ({
  mockGetAll: vi.fn(),
  mockCreate: vi.fn(),
  mockRequestPurchase: vi.fn(),
  mockCancelPurchase: vi.fn(),
  mockMarkSold: vi.fn(),
  mockRemove: vi.fn(),
}))

vi.mock('itx-ui-kit', () => ({
  Typography: { template: '<div><slot /></div>', props: ['variant', 'as'] },
}))

vi.mock('lucide-vue-next', () => ({
  Loader2: { template: '<span class="loader" />' },
  Package: { template: '<span />' },
  Plus: { template: '<span />' },
  Trash2: { template: '<span />' },
  User: { template: '<span />' },
}))

vi.mock('@/components/common/EmptyState.vue', () => ({
  default: { template: '<div class="empty-state"><slot /></div>', props: ['icon', 'title', 'description', 'actionLabel'] },
}))

vi.mock('@/components/common/ErrorState.vue', () => ({
  default: { template: '<div class="error-state"><slot /></div>', props: ['message'], emits: ['retry'] },
}))

vi.mock('@/components/common/FormField.vue', () => ({
  default: { template: '<div class="form-field"><slot /></div>', props: ['label', 'error', 'htmlFor', 'required'] },
}))

vi.mock('@/components/ConfirmDialog.vue', () => ({
  default: { template: '<div class="confirm-dialog"><slot name="trigger" /></div>', props: ['title', 'description', 'confirmLabel'], emits: ['confirm'] },
}))

vi.mock('@/components/ui/dialog', () => ({
  Dialog: { template: '<div class="dialog"><slot /></div>', props: ['open'] },
  DialogContent: { template: '<div class="dialog-content"><slot /></div>' },
  DialogScrollContent: { template: '<div class="dialog-scroll-content"><slot /></div>' },
  DialogHeader: { template: '<div class="dialog-header"><slot /></div>' },
  DialogTitle: { template: '<div class="dialog-title"><slot /></div>' },
  DialogDescription: { template: '<div class="dialog-description"><slot /></div>' },
  DialogFooter: { template: '<div class="dialog-footer"><slot /></div>' },
}))

vi.mock('@/composables/useFormValidation', () => ({
  required: () => (v: string) => (!v ? 'Required' : ''),
  useFormValidation: () => ({
    errors: {},
    validateAll: vi.fn(() => true),
    validateField: vi.fn(() => true),
    clearErrors: vi.fn(),
  }),
}))

vi.mock('@/composables/useUser', () => ({
  useUser: () => ({ value: { id: 1, firstName: 'Test', lastName: 'User', tg: 'testuser' } }),
  isUserAdmin: () => ({ value: false }),
}))

vi.mock('@/services/errorService', () => ({
  handleError: vi.fn(async () => ({ message: 'Error occurred' })),
}))

vi.mock('@/services/marketplace', () => ({
  marketplaceService: {
    getAll: mockGetAll,
    create: mockCreate,
    requestPurchase: mockRequestPurchase,
    cancelPurchase: mockCancelPurchase,
    markSold: mockMarkSold,
    remove: mockRemove,
  },
}))

import Marketplace from '@/pages/Marketplace.vue'

describe('Marketplace page', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('shows loading spinner initially', () => {
    mockGetAll.mockReturnValue(new Promise(() => {}))
    const wrapper = mount(Marketplace)
    expect(wrapper.find('.loader').exists()).toBe(true)
  })

  it('shows ErrorState when fetch fails', async () => {
    mockGetAll.mockRejectedValue(new Error('Network error'))
    const wrapper = mount(Marketplace)
    await flushPromises()
    expect(wrapper.find('.error-state').exists()).toBe(true)
  })

  it('shows EmptyState when no items', async () => {
    mockGetAll.mockResolvedValue({ items: [], total: 0 })
    const wrapper = mount(Marketplace)
    await flushPromises()
    expect(wrapper.find('.empty-state').exists()).toBe(true)
  })

  it('renders item cards when items exist', async () => {
    mockGetAll.mockResolvedValue({
      items: [
        {
          id: 1,
          title: 'MacBook Pro',
          description: 'Like new',
          price: '100000',
          city: 'Москва',
          canShip: true,
          condition: 'USED',
          defects: 'Minor scratch',
          packageContents: 'Laptop, charger',
          contactTelegram: 'seller1',
          contactEmail: '',
          contactPhone: '',
          imagePath: '',
          sellerId: 2,
          seller: { firstName: 'Seller', lastName: 'One', tg: 'seller1' },
          buyerId: null,
          buyer: null,
          status: 'ACTIVE',
          createdAt: '2026-01-01',
          updatedAt: '2026-01-01',
        },
        {
          id: 2,
          title: 'iPhone 15',
          description: 'Brand new',
          price: '80000',
          city: 'СПб',
          canShip: false,
          condition: 'NEW',
          defects: '',
          packageContents: 'Phone, cable',
          contactTelegram: 'seller2',
          contactEmail: '',
          contactPhone: '',
          imagePath: 'https://example.com/img.jpg',
          sellerId: 1,
          seller: { firstName: 'Test', lastName: 'User', tg: 'testuser' },
          buyerId: null,
          buyer: null,
          status: 'ACTIVE',
          createdAt: '2026-01-02',
          updatedAt: '2026-01-02',
        },
      ],
      total: 2,
    })
    const wrapper = mount(Marketplace)
    await flushPromises()
    expect(wrapper.text()).toContain('MacBook Pro')
    expect(wrapper.text()).toContain('iPhone 15')
  })

  it('shows ConfirmDialog on delete button for seller items', async () => {
    mockGetAll.mockResolvedValue({
      items: [
        {
          id: 1,
          title: 'My Item',
          description: '',
          price: '500',
          city: '',
          canShip: false,
          condition: 'NEW',
          defects: '',
          packageContents: '',
          contactTelegram: 'testuser',
          contactEmail: '',
          contactPhone: '',
          imagePath: '',
          sellerId: 1,
          seller: { firstName: 'Test', lastName: 'User', tg: 'testuser' },
          buyerId: null,
          buyer: null,
          status: 'ACTIVE',
          createdAt: '2026-01-01',
          updatedAt: '2026-01-01',
        },
      ],
      total: 1,
    })
    const wrapper = mount(Marketplace)
    await flushPromises()
    expect(wrapper.find('.confirm-dialog').exists()).toBe(true)
  })

  it('shows status tabs', async () => {
    mockGetAll.mockResolvedValue({ items: [], total: 0 })
    const wrapper = mount(Marketplace)
    await flushPromises()
    expect(wrapper.text()).toContain('Все')
    expect(wrapper.text()).toContain('Активные')
    expect(wrapper.text()).toContain('Забронированные')
    expect(wrapper.text()).toContain('Проданные')
  })

  it('calls marketplaceService.getAll on mount', async () => {
    mockGetAll.mockResolvedValue({ items: [], total: 0 })
    mount(Marketplace)
    await flushPromises()
    expect(mockGetAll).toHaveBeenCalledWith({ limit: 100 })
  })
})

describe('Marketplace logic', () => {
  const statusTabs = [
    { key: 'all', label: 'Все' },
    { key: 'ACTIVE', label: 'Активные' },
    { key: 'RESERVED', label: 'Забронированные' },
    { key: 'SOLD', label: 'Проданные' },
  ]

  const statusConfig: Record<string, { label: string, class: string }> = {
    ACTIVE: { label: 'Активно', class: 'bg-blue-500/10 text-blue-500' },
    RESERVED: { label: 'Забронировано', class: 'bg-yellow-500/10 text-yellow-500' },
    SOLD: { label: 'Продано', class: 'bg-green-500/10 text-green-500' },
    ARCHIVED: { label: 'В архиве', class: 'bg-muted text-muted-foreground' },
  }

  interface Item {
    id: number
    status: string
    sellerId: number
    buyerId: number | null
  }

  function isSeller(item: Item, userId: number) {
    return userId === item.sellerId
  }

  function isBuyer(item: Item, userId: number) {
    return userId === item.buyerId
  }

  function displayName(member: { firstName: string, lastName: string, tg: string }) {
    const name = [member.firstName, member.lastName].filter(Boolean).join(' ')
    return name || `@${member.tg}`
  }

  describe('statusTabs', () => {
    it('has 4 tabs', () => {
      expect(statusTabs).toHaveLength(4)
    })

    it('has correct keys', () => {
      expect(statusTabs.map(t => t.key)).toEqual(['all', 'ACTIVE', 'RESERVED', 'SOLD'])
    })
  })

  describe('statusConfig', () => {
    it('has all 4 statuses', () => {
      expect(Object.keys(statusConfig)).toEqual(['ACTIVE', 'RESERVED', 'SOLD', 'ARCHIVED'])
    })

    it('each status has label and class', () => {
      for (const status of Object.values(statusConfig)) {
        expect(status).toHaveProperty('label')
        expect(status).toHaveProperty('class')
        expect(typeof status.label).toBe('string')
        expect(typeof status.class).toBe('string')
      }
    })
  })

  describe('isSeller', () => {
    it('returns true when userId matches sellerId', () => {
      const item: Item = { id: 1, status: 'ACTIVE', sellerId: 10, buyerId: null }
      expect(isSeller(item, 10)).toBe(true)
    })

    it('returns false when userId does not match', () => {
      const item: Item = { id: 1, status: 'ACTIVE', sellerId: 10, buyerId: null }
      expect(isSeller(item, 20)).toBe(false)
    })
  })

  describe('isBuyer', () => {
    it('returns true when userId matches buyerId', () => {
      const item: Item = { id: 1, status: 'RESERVED', sellerId: 10, buyerId: 20 }
      expect(isBuyer(item, 20)).toBe(true)
    })

    it('returns false when buyerId is null', () => {
      const item: Item = { id: 1, status: 'ACTIVE', sellerId: 10, buyerId: null }
      expect(isBuyer(item, 20)).toBe(false)
    })

    it('returns false when userId does not match buyerId', () => {
      const item: Item = { id: 1, status: 'RESERVED', sellerId: 10, buyerId: 30 }
      expect(isBuyer(item, 20)).toBe(false)
    })
  })

  describe('displayName', () => {
    it('joins first and last name', () => {
      expect(displayName({ firstName: 'John', lastName: 'Doe', tg: 'johndoe' })).toBe('John Doe')
    })

    it('returns first name only when last name is empty', () => {
      expect(displayName({ firstName: 'John', lastName: '', tg: 'johndoe' })).toBe('John')
    })

    it('falls back to @tg when both names are empty', () => {
      expect(displayName({ firstName: '', lastName: '', tg: 'johndoe' })).toBe('@johndoe')
    })

    it('returns last name only when first name is empty', () => {
      expect(displayName({ firstName: '', lastName: 'Doe', tg: 'johndoe' })).toBe('Doe')
    })
  })
})
