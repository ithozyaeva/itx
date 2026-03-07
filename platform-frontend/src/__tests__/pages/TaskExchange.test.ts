import { describe, expect, it } from 'vitest'

describe('TaskExchange logic', () => {
  const statusTabs = [
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
    it('has 5 tabs', () => {
      expect(statusTabs).toHaveLength(5)
    })

    it('has correct keys', () => {
      expect(statusTabs.map(t => t.key)).toEqual(['all', 'OPEN', 'IN_PROGRESS', 'DONE', 'APPROVED'])
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
        id: 1, status: 'OPEN', creatorId: 10, maxAssignees: 3,
        assignees: [{ id: 5, firstName: 'A', lastName: 'B', tg: 'ab' }],
      }
      expect(isAssignee(task, 5)).toBe(true)
    })

    it('returns false when user is not in assignees', () => {
      const task: Task = {
        id: 1, status: 'OPEN', creatorId: 10, maxAssignees: 3,
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
        id: 1, status: 'OPEN', creatorId: 10, maxAssignees: 2,
        assignees: [{ id: 20, firstName: '', lastName: '', tg: '' }],
      }
      expect(canTakeTask(task, 20)).toBe(false)
    })

    it('disallows when max assignees reached', () => {
      const task: Task = {
        id: 1, status: 'OPEN', creatorId: 10, maxAssignees: 1,
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
