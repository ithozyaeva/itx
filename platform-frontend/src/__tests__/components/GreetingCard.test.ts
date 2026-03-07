import { mount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'
import { ref } from 'vue'

const mockUserData = ref({
  id: 1,
  telegramID: 123,
  tg: 'testuser',
  birthday: '',
  firstName: 'Алексей',
  lastName: 'Иванов',
  bio: '',
  grade: '',
  company: '',
  avatarUrl: '',
  roles: ['SUBSCRIBER'] as string[],
  createdAt: new Date().toISOString(),
})

vi.mock('@/composables/useUser', () => ({
  useUser: () => mockUserData,
  useUserLevel: () => ({
    level: ref('Бригадир'),
    levelIndex: ref(1),
  }),
}))

vi.mock('@/models/profile', async (importOriginal) => {
  const original = await importOriginal() as any
  return { ...original }
})

import GreetingCard from '@/components/dashboard/GreetingCard.vue'

// Test the pure logic functions extracted from GreetingCard

describe('GreetingCard logic', () => {
  describe('pluralizeDays', () => {
    function pluralizeDays(n: number): string {
      if (n % 10 === 1 && n % 100 !== 11)
        return 'день'
      if (n % 10 >= 2 && n % 10 <= 4 && (n % 100 < 10 || n % 100 >= 20))
        return 'дня'
      return 'дней'
    }

    it('returns "день" for 1', () => {
      expect(pluralizeDays(1)).toBe('день')
    })

    it('returns "дня" for 2, 3, 4', () => {
      expect(pluralizeDays(2)).toBe('дня')
      expect(pluralizeDays(3)).toBe('дня')
      expect(pluralizeDays(4)).toBe('дня')
    })

    it('returns "дней" for 5-20', () => {
      for (let i = 5; i <= 20; i++) {
        expect(pluralizeDays(i)).toBe('дней')
      }
    })

    it('returns "день" for 21, 31, 101', () => {
      expect(pluralizeDays(21)).toBe('день')
      expect(pluralizeDays(31)).toBe('день')
      expect(pluralizeDays(101)).toBe('день')
    })

    it('returns "дня" for 22, 33, 44', () => {
      expect(pluralizeDays(22)).toBe('дня')
      expect(pluralizeDays(33)).toBe('дня')
      expect(pluralizeDays(44)).toBe('дня')
    })

    it('returns "дней" for 11, 12, 111, 112', () => {
      expect(pluralizeDays(11)).toBe('дней')
      expect(pluralizeDays(12)).toBe('дней')
      expect(pluralizeDays(111)).toBe('дней')
      expect(pluralizeDays(112)).toBe('дней')
    })

    it('returns "дней" for 0', () => {
      expect(pluralizeDays(0)).toBe('дней')
    })
  })

  describe('daysSinceJoined calculation', () => {
    function calculateDaysSinceJoined(createdAt: string | undefined): number {
      if (!createdAt)
        return 1
      const created = new Date(createdAt)
      const now = new Date()
      created.setHours(0, 0, 0, 0)
      now.setHours(0, 0, 0, 0)
      const diff = Math.floor((now.getTime() - created.getTime()) / (1000 * 60 * 60 * 24))
      return diff + 1
    }

    it('returns 1 when createdAt is undefined', () => {
      expect(calculateDaysSinceJoined(undefined)).toBe(1)
    })

    it('returns 1 for today', () => {
      const today = new Date().toISOString()
      expect(calculateDaysSinceJoined(today)).toBe(1)
    })

    it('returns 2 for yesterday', () => {
      const yesterday = new Date()
      yesterday.setDate(yesterday.getDate() - 1)
      expect(calculateDaysSinceJoined(yesterday.toISOString())).toBe(2)
    })

    it('returns correct days for a week ago', () => {
      const weekAgo = new Date()
      weekAgo.setDate(weekAgo.getDate() - 7)
      expect(calculateDaysSinceJoined(weekAgo.toISOString())).toBe(8)
    })
  })

  describe('GreetingCard component rendering', () => {
    it('renders without errors', () => {
      const wrapper = mount(GreetingCard)
      expect(wrapper.exists()).toBe(true)
    })

    it('displays user first name', () => {
      const wrapper = mount(GreetingCard)
      expect(wrapper.text()).toContain('Алексей')
    })

    it('displays subscription level', () => {
      const wrapper = mount(GreetingCard)
      expect(wrapper.text()).toContain('Бригадир')
    })

    it('displays days since joined text', () => {
      const wrapper = mount(GreetingCard)
      expect(wrapper.text()).toContain('Ты в IT-Хозяевах уже')
    })
  })
})
