import type { TelegramUser } from '@/models/profile'
import { beforeEach, describe, expect, it } from 'vitest'
import { withSetup } from '../helpers'

// Reset the module-level singleton between tests
async function freshImport() {
  const mod = await import('@/composables/useUser')
  return mod
}

// useUser хранит данные обёрнутыми в {data, savedAt} под версионированным
// ключом — пишем через хелпер, чтобы тесты не дублировали логику сериализации.
function setStoredUser(user: TelegramUser | null) {
  localStorage.setItem('tg_user:v2', JSON.stringify({ data: user, savedAt: Date.now() }))
}

describe('useUser', () => {
  beforeEach(() => {
    localStorage.clear()
    // Force re-import to reset the singleton userRef
    vi.resetModules()
  })

  describe('useUser()', () => {
    it('returns null when no user in localStorage', async () => {
      const { useUser } = await freshImport()
      const { result } = withSetup(() => useUser())
      expect(result.value).toBeNull()
    })

    it('reads user from localStorage', async () => {
      const mockUser: TelegramUser = {
        id: 1,
        telegramID: 12345,
        tg: 'testuser',
        birthday: '1990-01-01',
        firstName: 'John',
        lastName: 'Doe',
        bio: 'Test bio',
        grade: 'Senior',
        company: 'TestCo',
        avatarUrl: '',
        roles: ['SUBSCRIBER'],
      }
      setStoredUser(mockUser)

      const { useUser } = await freshImport()
      const { result } = withSetup(() => useUser())
      expect(result.value).toEqual(mockUser)
      expect(result.value!.firstName).toBe('John')
    })

    it('returns the same ref on subsequent calls (singleton)', async () => {
      const { useUser } = await freshImport()
      let ref1: any, ref2: any
      withSetup(() => {
        ref1 = useUser()
        ref2 = useUser()
        return {}
      })
      expect(ref1).toBe(ref2)
    })
  })

  describe('isUserSubscribed()', () => {
    it('returns true when user does not have UNSUBSCRIBER role', async () => {
      const user: TelegramUser = {
        id: 1,
        telegramID: 1,
        tg: 'user',
        birthday: '',
        firstName: '',
        lastName: '',
        bio: '',
        grade: '',
        company: '',
        avatarUrl: '',
        roles: ['SUBSCRIBER'],
      }
      setStoredUser(user)

      const { isUserSubscribed } = await freshImport()
      const { result } = withSetup(() => isUserSubscribed())
      expect(result.value).toBe(true)
    })

    it('returns false when user has UNSUBSCRIBER role', async () => {
      const user: TelegramUser = {
        id: 1,
        telegramID: 1,
        tg: 'user',
        birthday: '',
        firstName: '',
        lastName: '',
        bio: '',
        grade: '',
        company: '',
        avatarUrl: '',
        roles: ['UNSUBSCRIBER'],
      }
      setStoredUser(user)

      const { isUserSubscribed } = await freshImport()
      const { result } = withSetup(() => isUserSubscribed())
      expect(result.value).toBe(false)
    })
  })

  describe('isUserMentor()', () => {
    it('returns true when user has MENTOR role', async () => {
      const user: TelegramUser = {
        id: 1,
        telegramID: 1,
        tg: 'user',
        birthday: '',
        firstName: '',
        lastName: '',
        bio: '',
        grade: '',
        company: '',
        avatarUrl: '',
        roles: ['MENTOR'],
      }
      setStoredUser(user)

      const { isUserMentor } = await freshImport()
      const { result } = withSetup(() => isUserMentor())
      expect(result.value).toBe(true)
    })

    it('returns false when user does not have MENTOR role', async () => {
      const user: TelegramUser = {
        id: 1,
        telegramID: 1,
        tg: 'user',
        birthday: '',
        firstName: '',
        lastName: '',
        bio: '',
        grade: '',
        company: '',
        avatarUrl: '',
        roles: ['SUBSCRIBER'],
      }
      setStoredUser(user)

      const { isUserMentor } = await freshImport()
      const { result } = withSetup(() => isUserMentor())
      expect(result.value).toBe(false)
    })
  })

  describe('isUserAdmin()', () => {
    it('returns true when user has ADMIN role', async () => {
      const user: TelegramUser = {
        id: 1,
        telegramID: 1,
        tg: 'user',
        birthday: '',
        firstName: '',
        lastName: '',
        bio: '',
        grade: '',
        company: '',
        avatarUrl: '',
        roles: ['ADMIN'],
      }
      setStoredUser(user)

      const { isUserAdmin } = await freshImport()
      const { result } = withSetup(() => isUserAdmin())
      expect(result.value).toBe(true)
    })

    it('returns false when user has no ADMIN role', async () => {
      const user: TelegramUser = {
        id: 1,
        telegramID: 1,
        tg: 'user',
        birthday: '',
        firstName: '',
        lastName: '',
        bio: '',
        grade: '',
        company: '',
        avatarUrl: '',
        roles: ['SUBSCRIBER'],
      }
      setStoredUser(user)

      const { isUserAdmin } = await freshImport()
      const { result } = withSetup(() => isUserAdmin())
      expect(result.value).toBe(false)
    })

    it('returns false when user is null', async () => {
      const { isUserAdmin } = await freshImport()
      const { result } = withSetup(() => isUserAdmin())
      expect(result.value).toBe(false)
    })
  })

  describe('canViewAdminPanel()', () => {
    it('returns true for ADMIN', async () => {
      const user: TelegramUser = {
        id: 1,
        telegramID: 1,
        tg: 'user',
        birthday: '',
        firstName: '',
        lastName: '',
        bio: '',
        grade: '',
        company: '',
        avatarUrl: '',
        roles: ['ADMIN'],
      }
      setStoredUser(user)

      const { canViewAdminPanel } = await freshImport()
      const { result } = withSetup(() => canViewAdminPanel())
      expect(result.value).toBe(true)
    })

    it('returns true for EVENT_MAKER', async () => {
      const user: TelegramUser = {
        id: 1,
        telegramID: 1,
        tg: 'user',
        birthday: '',
        firstName: '',
        lastName: '',
        bio: '',
        grade: '',
        company: '',
        avatarUrl: '',
        roles: ['EVENT_MAKER'],
      }
      setStoredUser(user)

      const { canViewAdminPanel } = await freshImport()
      const { result } = withSetup(() => canViewAdminPanel())
      expect(result.value).toBe(true)
    })

    it('returns false for regular SUBSCRIBER', async () => {
      const user: TelegramUser = {
        id: 1,
        telegramID: 1,
        tg: 'user',
        birthday: '',
        firstName: '',
        lastName: '',
        bio: '',
        grade: '',
        company: '',
        avatarUrl: '',
        roles: ['SUBSCRIBER'],
      }
      setStoredUser(user)

      const { canViewAdminPanel } = await freshImport()
      const { result } = withSetup(() => canViewAdminPanel())
      expect(result.value).toBe(false)
    })

    it('returns false when user is null', async () => {
      const { canViewAdminPanel } = await freshImport()
      const { result } = withSetup(() => canViewAdminPanel())
      expect(result.value).toBe(false)
    })
  })

  describe('useUserLevel()', () => {
    it('returns Новичок for null user', async () => {
      const { useUserLevel } = await freshImport()
      const { result } = withSetup(() => useUserLevel())
      expect(result.level.value).toBe('Новичок')
      expect(result.levelIndex.value).toBe(0)
    })

    it('returns Новичок for SUBSCRIBER without tier (beginner or null)', async () => {
      const user: TelegramUser = {
        id: 1,
        telegramID: 1,
        tg: 'user',
        birthday: '',
        firstName: '',
        lastName: '',
        bio: '',
        grade: '',
        company: '',
        avatarUrl: '',
        roles: ['SUBSCRIBER'],
        subscriptionTier: { id: 1, slug: 'beginner', name: 'Beginner', level: 1 },
      }
      setStoredUser(user)

      const { useUserLevel } = await freshImport()
      const { result } = withSetup(() => useUserLevel())
      expect(result.level.value).toBe('Новичок')
      expect(result.levelIndex.value).toBe(0)
    })

    it('returns Бригадир for foreman tier', async () => {
      const user: TelegramUser = {
        id: 1,
        telegramID: 1,
        tg: 'user',
        birthday: '',
        firstName: '',
        lastName: '',
        bio: '',
        grade: '',
        company: '',
        avatarUrl: '',
        roles: ['SUBSCRIBER'],
        subscriptionTier: { id: 2, slug: 'foreman', name: 'Foreman', level: 2 },
      }
      setStoredUser(user)

      const { useUserLevel } = await freshImport()
      const { result } = withSetup(() => useUserLevel())
      expect(result.level.value).toBe('Бригадир')
      expect(result.levelIndex.value).toBe(1)
    })

    it('returns Хозяин for master tier', async () => {
      const user: TelegramUser = {
        id: 1,
        telegramID: 1,
        tg: 'user',
        birthday: '',
        firstName: '',
        lastName: '',
        bio: '',
        grade: '',
        company: '',
        avatarUrl: '',
        roles: ['SUBSCRIBER'],
        subscriptionTier: { id: 3, slug: 'master', name: 'Master', level: 3 },
      }
      setStoredUser(user)

      const { useUserLevel } = await freshImport()
      const { result } = withSetup(() => useUserLevel())
      expect(result.level.value).toBe('Хозяин')
      expect(result.levelIndex.value).toBe(2)
    })

    it('returns King for king tier', async () => {
      const user: TelegramUser = {
        id: 1,
        telegramID: 1,
        tg: 'user',
        birthday: '',
        firstName: '',
        lastName: '',
        bio: '',
        grade: '',
        company: '',
        avatarUrl: '',
        roles: ['SUBSCRIBER'],
        subscriptionTier: { id: 4, slug: 'king', name: 'King', level: 4 },
      }
      setStoredUser(user)

      const { useUserLevel } = await freshImport()
      const { result } = withSetup(() => useUserLevel())
      expect(result.level.value).toBe('King')
      expect(result.levelIndex.value).toBe(3)
    })

    it('returns Хозяин for MENTOR even without master tier', async () => {
      const user: TelegramUser = {
        id: 1,
        telegramID: 1,
        tg: 'user',
        birthday: '',
        firstName: '',
        lastName: '',
        bio: '',
        grade: '',
        company: '',
        avatarUrl: '',
        roles: ['MENTOR'],
      }
      setStoredUser(user)

      const { useUserLevel } = await freshImport()
      const { result } = withSetup(() => useUserLevel())
      expect(result.level.value).toBe('Хозяин')
      expect(result.levelIndex.value).toBe(2)
    })

    it('returns King for ADMIN', async () => {
      const user: TelegramUser = {
        id: 1,
        telegramID: 1,
        tg: 'user',
        birthday: '',
        firstName: '',
        lastName: '',
        bio: '',
        grade: '',
        company: '',
        avatarUrl: '',
        roles: ['ADMIN'],
      }
      setStoredUser(user)

      const { useUserLevel } = await freshImport()
      const { result } = withSetup(() => useUserLevel())
      expect(result.level.value).toBe('King')
      expect(result.levelIndex.value).toBe(3)
    })

    it('exposes maxLevel', async () => {
      const { useUserLevel } = await freshImport()
      const { result } = withSetup(() => useUserLevel())
      expect(result.maxLevel).toBe(3)
    })
  })
})
