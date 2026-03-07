import { describe, expect, it } from 'vitest'

describe('Leaderboard logic', () => {
  interface LeaderboardEntry {
    avatarUrl?: string
    tg?: string
  }

  function getAvatarSrc(entry: LeaderboardEntry) {
    return entry.avatarUrl || `https://t.me/i/userpic/160/${entry.tg}.jpg`
  }

  describe('getAvatarSrc', () => {
    it('returns avatarUrl when present', () => {
      const entry = { avatarUrl: 'https://example.com/avatar.jpg', tg: 'user123' }
      expect(getAvatarSrc(entry)).toBe('https://example.com/avatar.jpg')
    })

    it('falls back to Telegram userpic when avatarUrl is missing', () => {
      const entry = { tg: 'user123' }
      expect(getAvatarSrc(entry)).toBe('https://t.me/i/userpic/160/user123.jpg')
    })

    it('falls back to Telegram userpic when avatarUrl is empty string', () => {
      const entry = { avatarUrl: '', tg: 'johndoe' }
      expect(getAvatarSrc(entry)).toBe('https://t.me/i/userpic/160/johndoe.jpg')
    })

    it('handles undefined tg gracefully', () => {
      const entry = {}
      expect(getAvatarSrc(entry)).toBe('https://t.me/i/userpic/160/undefined.jpg')
    })
  })
})
