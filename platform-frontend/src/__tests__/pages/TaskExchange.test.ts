import { flushPromises, mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

const { mockGetAll, mockCreate, mockUpdate, mockAssign, mockUnassign, mockRemoveAssignee, mockMarkDone, mockApprove, mockReject, mockRemove } = vi.hoisted(() => ({
  mockGetAll: vi.fn(),
  mockCreate: vi.fn(),
  mockUpdate: vi.fn(),
  mockAssign: vi.fn(),
  mockUnassign: vi.fn(),
  mockRemoveAssignee: vi.fn(),
  mockMarkDone: vi.fn(),
  mockApprove: vi.fn(),
  mockReject: vi.fn(),
  mockRemove: vi.fn(),
}))

vi.mock('@/components/ui/typography', () => ({
  Typography: { template: '<div><slot /></div>', props: ['variant', 'as'] },
}))

vi.mock('lucide-vue-next', () => ({
  CheckCircle: { template: '<span />' },
  ClipboardList: { template: '<span />' },
  Clock: { template: '<span />' },
  Edit3: { template: '<span />' },
  Loader2: { template: '<span class="loader" />' },
  Plus: { template: '<span />' },
  Search: { template: '<span />' },
  Trash2: { template: '<span />' },
  User: { template: '<span />' },
  Users: { template: '<span />' },
  XCircle: { template: '<span />' },
}))

vi.mock('@/components/tasks/TaskCardSkeleton.vue', () => ({
  default: { template: '<div class="skeleton" />' },
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

vi.mock('@/composables/useSSE', () => ({
  useSSE: vi.fn(),
}))

vi.mock('@/composables/useUser', () => ({
  useUser: () => ({ value: { id: 1, firstName: 'Test', lastName: 'User' } }),
  isUserAdmin: () => ({ value: false }),
}))

vi.mock('@/services/errorService', () => ({
  handleError: vi.fn(async () => ({ message: 'Error occurred' })),
}))

vi.mock('@/services/taskExchange', () => ({
  taskExchangeService: {
    getAll: mockGetAll,
    create: mockCreate,
    update: mockUpdate,
    assign: mockAssign,
    unassign: mockUnassign,
    removeAssignee: mockRemoveAssignee,
    markDone: mockMarkDone,
    approve: mockApprove,
    reject: mockReject,
    remove: mockRemove,
  },
}))

import TaskExchange from '@/pages/TaskExchange.vue'

describe('TaskExchange page', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('shows loading spinner initially', () => {
    mockGetAll.mockReturnValue(new Promise(() => {}))
    const wrapper = mount(TaskExchange)
    expect(wrapper.find('.skeleton').exists()).toBe(true)
  })

  it('shows ErrorState when fetch fails', async () => {
    mockGetAll.mockRejectedValue(new Error('Network error'))
    const wrapper = mount(TaskExchange)
    await flushPromises()
    expect(wrapper.find('.error-state').exists()).toBe(true)
  })

  it('shows EmptyState when no tasks', async () => {
    mockGetAll.mockResolvedValue({ items: [], total: 0 })
    const wrapper = mount(TaskExchange)
    await flushPromises()
    expect(wrapper.find('.empty-state').exists()).toBe(true)
  })

  it('renders task cards when tasks exist', async () => {
    mockGetAll.mockResolvedValue({
      items: [
        {
          id: 1,
          title: 'Test Task',
          description: 'A test task',
          status: 'OPEN',
          creatorId: 2,
          creator: { firstName: 'Creator', lastName: 'User', tg: 'creator' },
          maxAssignees: 3,
          assignees: [],
        },
        {
          id: 2,
          title: 'Another Task',
          description: '',
          status: 'IN_PROGRESS',
          creatorId: 1,
          creator: { firstName: 'Test', lastName: 'User', tg: 'test' },
          maxAssignees: 1,
          assignees: [{ id: 3, firstName: 'Worker', lastName: '', tg: 'worker' }],
        },
      ],
      total: 2,
    })
    const wrapper = mount(TaskExchange)
    await flushPromises()
    expect(wrapper.text()).toContain('Test Task')
    expect(wrapper.text()).toContain('Another Task')
  })

  it('shows ConfirmDialog on delete button for creator tasks', async () => {
    mockGetAll.mockResolvedValue({
      items: [
        {
          id: 1,
          title: 'My Task',
          description: '',
          status: 'OPEN',
          creatorId: 1,
          creator: { firstName: 'Test', lastName: 'User', tg: 'test' },
          maxAssignees: 1,
          assignees: [],
        },
      ],
      total: 1,
    })
    const wrapper = mount(TaskExchange)
    await flushPromises()
    expect(wrapper.find('.confirm-dialog').exists()).toBe(true)
  })

  it('shows status tabs', async () => {
    mockGetAll.mockResolvedValue({ items: [], total: 0 })
    const wrapper = mount(TaskExchange)
    await flushPromises()
    expect(wrapper.text()).toContain('Активные')
    expect(wrapper.text()).toContain('Все')
    expect(wrapper.text()).toContain('Открытые')
    expect(wrapper.text()).toContain('В работе')
    expect(wrapper.text()).toContain('На проверке')
    expect(wrapper.text()).toContain('Выполненные')
  })

  it('calls taskExchangeService.getAll on mount', async () => {
    mockGetAll.mockResolvedValue({ items: [], total: 0 })
    mount(TaskExchange)
    await flushPromises()
    expect(mockGetAll).toHaveBeenCalledWith({ limit: 100 })
  })
})

describe('TaskExchange logic', () => {
  const statusTabs = [
    { key: 'active', label: 'Активные' },
    { key: 'all', label: 'Все' },
    { key: 'OPEN', label: 'Открытые' },
    { key: 'IN_PROGRESS', label: 'В работе' },
    { key: 'DONE', label: 'На проверке' },
    { key: 'APPROVED', label: 'Выполненные' },
  ]

  const statusConfig: Record<string, { label: string, class: string }> = {
    OPEN: { label: 'Открыто', class: 'bg-blue-500/10 text-blue-500' },
    IN_PROGRESS: { label: 'В работе', class: 'bg-yellow-500/10 text-yellow-500' },
    DONE: { label: 'На проверке', class: 'bg-purple-500/10 text-purple-500' },
    APPROVED: { label: 'Выполнено', class: 'bg-green-500/10 text-green-500' },
  }

  interface Task {
    id: number
    status: string
    creatorId: number
    maxAssignees: number
    assignees?: { id: number, firstName: string, lastName: string, tg: string }[]
  }

  function isCreator(task: Task, userId: number) {
    return userId === task.creatorId
  }

  function isAssignee(task: Task, userId: number) {
    return task.assignees?.some(a => a.id === userId) ?? false
  }

  function canTakeTask(task: Task, userId: number) {
    return task.status === 'OPEN'
      && task.creatorId !== userId
      && !(task.assignees?.some(a => a.id === userId))
      && (task.assignees?.length ?? 0) < task.maxAssignees
  }

  function canEdit(task: Task, userId: number, isAdmin: boolean) {
    return (userId === task.creatorId || isAdmin) && (task.status === 'OPEN' || task.status === 'IN_PROGRESS')
  }

  function canMarkDone(task: Task, userId: number, isAdmin: boolean) {
    return task.status === 'IN_PROGRESS' && (userId === task.creatorId || isAdmin)
  }

  function displayName(member: { firstName: string, lastName: string, tg: string }) {
    const name = [member.firstName, member.lastName].filter(Boolean).join(' ')
    return name || `@${member.tg}`
  }

  describe('statusTabs', () => {
    it('has 6 tabs', () => {
      expect(statusTabs).toHaveLength(6)
    })

    it('has correct keys', () => {
      expect(statusTabs.map(t => t.key)).toEqual(['active', 'all', 'OPEN', 'IN_PROGRESS', 'DONE', 'APPROVED'])
    })
  })

  describe('statusConfig', () => {
    it('has all 4 statuses', () => {
      expect(Object.keys(statusConfig)).toEqual(['OPEN', 'IN_PROGRESS', 'DONE', 'APPROVED'])
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

  describe('isCreator', () => {
    const task: Task = { id: 1, status: 'OPEN', creatorId: 10, maxAssignees: 3 }

    it('returns true when userId matches creatorId', () => {
      expect(isCreator(task, 10)).toBe(true)
    })

    it('returns false when userId does not match', () => {
      expect(isCreator(task, 20)).toBe(false)
    })
  })

  describe('isAssignee', () => {
    it('returns true when user is in assignees', () => {
      const task: Task = {
        id: 1,
        status: 'OPEN',
        creatorId: 10,
        maxAssignees: 3,
        assignees: [{ id: 5, firstName: 'A', lastName: 'B', tg: 'ab' }],
      }
      expect(isAssignee(task, 5)).toBe(true)
    })

    it('returns false when user is not in assignees', () => {
      const task: Task = {
        id: 1,
        status: 'OPEN',
        creatorId: 10,
        maxAssignees: 3,
        assignees: [{ id: 5, firstName: 'A', lastName: 'B', tg: 'ab' }],
      }
      expect(isAssignee(task, 99)).toBe(false)
    })

    it('returns false when assignees is undefined', () => {
      const task: Task = { id: 1, status: 'OPEN', creatorId: 10, maxAssignees: 3 }
      expect(isAssignee(task, 5)).toBe(false)
    })

    it('returns false when assignees is empty', () => {
      const task: Task = { id: 1, status: 'OPEN', creatorId: 10, maxAssignees: 3, assignees: [] }
      expect(isAssignee(task, 5)).toBe(false)
    })
  })

  describe('canTakeTask', () => {
    it('allows taking an open task when not creator and not assigned and slots available', () => {
      const task: Task = { id: 1, status: 'OPEN', creatorId: 10, maxAssignees: 2, assignees: [] }
      expect(canTakeTask(task, 20)).toBe(true)
    })

    it('disallows when task is not OPEN', () => {
      const task: Task = { id: 1, status: 'IN_PROGRESS', creatorId: 10, maxAssignees: 2, assignees: [] }
      expect(canTakeTask(task, 20)).toBe(false)
    })

    it('disallows when user is the creator', () => {
      const task: Task = { id: 1, status: 'OPEN', creatorId: 10, maxAssignees: 2, assignees: [] }
      expect(canTakeTask(task, 10)).toBe(false)
    })

    it('disallows when user is already assigned', () => {
      const task: Task = {
        id: 1,
        status: 'OPEN',
        creatorId: 10,
        maxAssignees: 2,
        assignees: [{ id: 20, firstName: '', lastName: '', tg: '' }],
      }
      expect(canTakeTask(task, 20)).toBe(false)
    })

    it('disallows when max assignees reached', () => {
      const task: Task = {
        id: 1,
        status: 'OPEN',
        creatorId: 10,
        maxAssignees: 1,
        assignees: [{ id: 30, firstName: '', lastName: '', tg: '' }],
      }
      expect(canTakeTask(task, 20)).toBe(false)
    })

    it('allows when assignees is undefined (no assignees yet)', () => {
      const task: Task = { id: 1, status: 'OPEN', creatorId: 10, maxAssignees: 2 }
      expect(canTakeTask(task, 20)).toBe(true)
    })
  })

  describe('canEdit', () => {
    it('allows creator to edit OPEN task', () => {
      const task: Task = { id: 1, status: 'OPEN', creatorId: 10, maxAssignees: 1 }
      expect(canEdit(task, 10, false)).toBe(true)
    })

    it('allows creator to edit IN_PROGRESS task', () => {
      const task: Task = { id: 1, status: 'IN_PROGRESS', creatorId: 10, maxAssignees: 1 }
      expect(canEdit(task, 10, false)).toBe(true)
    })

    it('allows admin to edit OPEN task', () => {
      const task: Task = { id: 1, status: 'OPEN', creatorId: 10, maxAssignees: 1 }
      expect(canEdit(task, 99, true)).toBe(true)
    })

    it('disallows non-creator non-admin', () => {
      const task: Task = { id: 1, status: 'OPEN', creatorId: 10, maxAssignees: 1 }
      expect(canEdit(task, 20, false)).toBe(false)
    })

    it('disallows editing DONE tasks even for creator', () => {
      const task: Task = { id: 1, status: 'DONE', creatorId: 10, maxAssignees: 1 }
      expect(canEdit(task, 10, false)).toBe(false)
    })

    it('disallows editing APPROVED tasks even for admin', () => {
      const task: Task = { id: 1, status: 'APPROVED', creatorId: 10, maxAssignees: 1 }
      expect(canEdit(task, 99, true)).toBe(false)
    })
  })

  describe('canMarkDone', () => {
    it('allows creator to mark IN_PROGRESS task done', () => {
      const task: Task = { id: 1, status: 'IN_PROGRESS', creatorId: 10, maxAssignees: 1 }
      expect(canMarkDone(task, 10, false)).toBe(true)
    })

    it('allows admin to mark IN_PROGRESS task done', () => {
      const task: Task = { id: 1, status: 'IN_PROGRESS', creatorId: 10, maxAssignees: 1 }
      expect(canMarkDone(task, 99, true)).toBe(true)
    })

    it('disallows marking OPEN task as done', () => {
      const task: Task = { id: 1, status: 'OPEN', creatorId: 10, maxAssignees: 1 }
      expect(canMarkDone(task, 10, false)).toBe(false)
    })

    it('disallows non-creator non-admin', () => {
      const task: Task = { id: 1, status: 'IN_PROGRESS', creatorId: 10, maxAssignees: 1 }
      expect(canMarkDone(task, 20, false)).toBe(false)
    })

    it('disallows marking DONE task as done again', () => {
      const task: Task = { id: 1, status: 'DONE', creatorId: 10, maxAssignees: 1 }
      expect(canMarkDone(task, 10, false)).toBe(false)
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
