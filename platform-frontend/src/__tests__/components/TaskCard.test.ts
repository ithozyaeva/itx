import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import TaskCard from '@/components/progress/TaskCard.vue'

const stubs = {
  TintedIcon: { props: ['icon', 'tone', 'done', 'size'], template: '<div class="ti" />' },
  ProgressBar: { props: ['progress', 'target', 'state', 'size'], template: '<div class="pb" :data-state="state" />' },
  PointsBadge: { props: ['amount', 'earned'], template: '<span class="pb-badge" :data-earned="String(earned)">{{ amount }}</span>' },
  Deadline: { props: ['endsAt'], template: '<span class="deadline" />' },
  RouterLink: { props: ['to'], template: '<a :href="to"><slot /></a>' },
}

const MockIcon = { template: '<span class="mock-icon" />' }

describe('TaskCard', () => {
  it('renders title, description and points', () => {
    const w = mount(TaskCard, {
      props: { title: 'Сделай X', description: 'Подробности', points: 10 },
      global: { stubs },
    })
    expect(w.text()).toContain('Сделай X')
    expect(w.text()).toContain('Подробности')
    expect(w.find('.pb-badge').text()).toContain('10')
  })

  it('shows progress counter alongside pillLabel (regression: tier+counter coexist)', () => {
    const w = mount(TaskCard, {
      props: {
        title: 'Дейлик',
        points: 5,
        progress: 3,
        target: 7,
        pillLabel: 'Контент',
      },
      global: { stubs },
    })
    // Бейдж типа
    expect(w.text()).toContain('Контент')
    // Счётчик прогресса не должен исчезать
    expect(w.text()).toContain('3 / 7')
  })

  it('uses custom progressLabel when provided', () => {
    const w = mount(TaskCard, {
      props: {
        title: 'Чат-квест',
        points: 20,
        progress: 4,
        target: 10,
        progressLabel: '4 / 10 сообщений',
      },
      global: { stubs },
    })
    expect(w.text()).toContain('4 / 10 сообщений')
  })

  it('shows "Выполнено" label when done=true and not awarded', () => {
    const w = mount(TaskCard, {
      props: {
        title: 'X',
        points: 10,
        progress: 5,
        target: 5,
        done: true,
        awarded: false,
      },
      global: { stubs },
    })
    expect(w.text()).toContain('Выполнено')
  })

  it('shows "Награда зачислена" when awarded=true', () => {
    const w = mount(TaskCard, {
      props: {
        title: 'X',
        points: 10,
        progress: 5,
        target: 5,
        done: true,
        awarded: true,
      },
      global: { stubs },
    })
    expect(w.text()).toContain('Награда зачислена')
  })

  it('progress bar gets state=done when done && !awarded', () => {
    const w = mount(TaskCard, {
      props: {
        title: 'X', points: 10, progress: 5, target: 5, done: true,
      },
      global: { stubs },
    })
    expect(w.find('.pb').attributes('data-state')).toBe('done')
  })

  it('progress bar gets state=awarded when awarded', () => {
    const w = mount(TaskCard, {
      props: {
        title: 'X', points: 10, progress: 5, target: 5, awarded: true,
      },
      global: { stubs },
    })
    expect(w.find('.pb').attributes('data-state')).toBe('awarded')
  })

  it('hides Deadline when done', () => {
    const w = mount(TaskCard, {
      props: {
        title: 'X', points: 10, progress: 5, target: 5,
        done: true, endsAt: new Date(Date.now() + 86400000).toISOString(),
      },
      global: { stubs },
    })
    expect(w.find('.deadline').exists()).toBe(false)
  })

  it('shows Deadline + achievement marker when active', () => {
    const w = mount(TaskCard, {
      props: {
        title: 'X', points: 10, progress: 1, target: 5,
        endsAt: new Date(Date.now() + 86400000).toISOString(),
        hasAchievement: true,
      },
      global: { stubs },
    })
    expect(w.find('.deadline').exists()).toBe(true)
    expect(w.text()).toContain('+ ачивка')
  })

  it('renders <RouterLink> when "to" prop is set', () => {
    const w = mount(TaskCard, {
      props: { title: 'X', points: 10, to: '/foo' },
      global: { stubs },
    })
    expect(w.find('a').exists()).toBe(true)
    expect(w.find('a').attributes('href')).toBe('/foo')
  })

  it('marks PointsBadge as earned when awarded', () => {
    const w = mount(TaskCard, {
      props: {
        title: 'X', points: 10, progress: 5, target: 5, awarded: true,
      },
      global: { stubs },
    })
    expect(w.find('.pb-badge').attributes('data-earned')).toBe('true')
  })

  it('does NOT show progress block when target is 0 or absent', () => {
    const w = mount(TaskCard, {
      props: { title: 'X', points: 10 },
      global: { stubs },
    })
    expect(w.find('.pb').exists()).toBe(false)
  })
})
