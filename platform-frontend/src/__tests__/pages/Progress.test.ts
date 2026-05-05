import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'
import { describe, expect, it } from 'vitest'

const PROGRESS_VUE = resolve(__dirname, '../../pages/Progress.vue')
const progressSrc = readFileSync(PROGRESS_VUE, 'utf8')

describe('Progress.vue lazy-loading guard', () => {
  // Регрессионный тест на #331 review: без defineAsyncComponent панели
  // попадают в основной chunk Progress, и юзер на /progress?tab=today
  // скачивает все 4 «соседних» панели (~30kB gzip) впустую. Если кто-то
  // случайно вернёт статический импорт — этот тест упадёт.
  it.each([
    'LeaderboardPanel',
    'AchievementsPanel',
    'MyStatsPanel',
    'KudosPanel',
  ])('panel %s is lazy-loaded via defineAsyncComponent', (panel) => {
    const re = new RegExp(`defineAsyncComponent\\(\\(\\)\\s*=>\\s*import\\([^)]*${panel}\\.vue`)
    expect(progressSrc).toMatch(re)
  })

  it('does not statically import panels via @/components/progress barrel', () => {
    // Если барьер вернёт panel-ы в default-экспорт, они снова попадут в Progress chunk.
    // Здесь намеренно ловим именно статический import { LeaderboardPanel } from '@/components/progress'.
    const barrelLine = progressSrc
      .split('\n')
      .find(l => l.includes('@/components/progress\'') && l.startsWith('import'))
    expect(barrelLine).toBeDefined()
    for (const panel of ['LeaderboardPanel', 'AchievementsPanel', 'MyStatsPanel', 'KudosPanel'])
      expect(barrelLine).not.toContain(panel)
  })
})

describe('Progress.vue tab selection logic', () => {
  // Логика повторена из Progress.vue (строки 60-67), чтобы протестировать
  // поведение без тяжёлого mount-а. Если поменяется в src — поменяется тут.
  const VALID_TABS = ['today', 'period', 'history', 'sources', 'leaderboard', 'achievements', 'stats', 'kudos'] as const
  type TabKey = typeof VALID_TABS[number]

  function pickInitialTab(queryTab: string | undefined): TabKey {
    return (VALID_TABS as readonly string[]).includes(queryTab ?? '') ? (queryTab as TabKey) : 'today'
  }

  it.each(VALID_TABS)('accepts whitelisted tab %s', (tab) => {
    expect(pickInitialTab(tab)).toBe(tab)
  })

  it.each([
    ['unknown'],
    [''],
    [undefined],
    ['leaderboard '], // hidden trailing space
    ['LEADERBOARD'], // case-sensitive
    ['<script>alert(1)</script>'],
  ])('falls back to "today" for invalid value %s', (val) => {
    expect(pickInitialTab(val as string | undefined)).toBe('today')
  })
})

describe('Progress.vue VALID_TABS guard', () => {
  // Если кто-то добавит/удалит таб, тут увидит — и обновит whitelist
  // в Progress.vue одновременно с тестом.
  it('VALID_TABS array contains exactly the 8 expected keys', () => {
    const match = progressSrc.match(/const VALID_TABS:\s*TabKey\[\]\s*=\s*\[([^\]]+)\]/)
    expect(match).toBeTruthy()
    const arr = match![1].split(',').map(s => s.trim().replace(/['"]/g, '')).filter(Boolean)
    expect(arr).toEqual(['today', 'period', 'history', 'sources', 'leaderboard', 'achievements', 'stats', 'kudos'])
  })
})
