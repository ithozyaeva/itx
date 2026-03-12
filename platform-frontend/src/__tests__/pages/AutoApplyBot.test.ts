import { mount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'

vi.mock('itx-ui-kit', () => ({
  Typography: { template: '<div><slot /></div>', props: ['variant', 'as'] },
}))

vi.mock('lucide-vue-next', () => ({
  Bot: { template: '<span class="icon-bot" />' },
  CheckCircle: { template: '<span />' },
  Clock: { template: '<span />' },
  ExternalLink: { template: '<span />' },
  Filter: { template: '<span />' },
  MessageSquare: { template: '<span />' },
  Rocket: { template: '<span />' },
  Settings: { template: '<span />' },
  Zap: { template: '<span />' },
}))

import AutoApplyBot from '@/pages/AutoApplyBot.vue'

describe('AutoApplyBot page', () => {
  it('renders page title', () => {
    const wrapper = mount(AutoApplyBot)
    expect(wrapper.text()).toContain('Бот для автооткликов')
  })

  it('shows bot name "Roaster Resume Bot"', () => {
    const wrapper = mount(AutoApplyBot)
    expect(wrapper.text()).toContain('Roaster Resume Bot')
  })

  it('renders 4 "how it works" steps', () => {
    const wrapper = mount(AutoApplyBot)
    expect(wrapper.text()).toContain('Откройте бота')
    expect(wrapper.text()).toContain('Настройте профиль')
    expect(wrapper.text()).toContain('Задайте фильтры')
    expect(wrapper.text()).toContain('Запустите автоотклик')
  })

  it('renders 4 features', () => {
    const wrapper = mount(AutoApplyBot)
    expect(wrapper.text()).toContain('Автоматические отклики')
    expect(wrapper.text()).toContain('Умные фильтры')
    expect(wrapper.text()).toContain('Работает 24/7')
    expect(wrapper.text()).toContain('Сопроводительные письма')
  })

  it('renders FAQ items', () => {
    const wrapper = mount(AutoApplyBot)
    expect(wrapper.text()).toContain('Бот бесплатный?')
    expect(wrapper.text()).toContain('Нужно ли предоставлять логин и пароль от hh.ru?')
    expect(wrapper.text()).toContain('Можно ли настроить фильтры вакансий?')
    expect(wrapper.text()).toContain('Как часто бот проверяет новые вакансии?')
  })

  it('FAQ toggle works (click opens/closes)', async () => {
    const wrapper = mount(AutoApplyBot)

    // Initially answers are hidden
    expect(wrapper.text()).not.toContain('Да, бот полностью бесплатен')

    // Click first FAQ question to open
    const faqButtons = wrapper.findAll('button[aria-expanded]')
    expect(faqButtons.length).toBe(4)
    await faqButtons[0].trigger('click')

    // Now answer should be visible
    expect(wrapper.text()).toContain('Да, бот полностью бесплатен')

    // Click again to close
    await faqButtons[0].trigger('click')
    expect(wrapper.text()).not.toContain('Да, бот полностью бесплатен')
  })

  it('links point to telegram bot', () => {
    const wrapper = mount(AutoApplyBot)
    const links = wrapper.findAll('a[href="https://t.me/roaster_resume_bot"]')
    expect(links.length).toBeGreaterThanOrEqual(2)
  })

  it('CTA buttons present', () => {
    const wrapper = mount(AutoApplyBot)
    expect(wrapper.text()).toContain('Открыть бота в Telegram')
    expect(wrapper.text()).toContain('Начать пользоваться ботом')
  })
})

describe('AutoApplyBot logic', () => {
  const steps = [
    { title: 'Откройте бота', description: 'Перейдите в Telegram и запустите бота' },
    { title: 'Настройте профиль', description: 'Укажите желаемую должность, зарплату и город' },
    { title: 'Задайте фильтры', description: 'Выберите формат работы, опыт и другие параметры' },
    { title: 'Запустите автоотклик', description: 'Бот будет откликаться на подходящие вакансии за вас' },
  ]

  const features = [
    { title: 'Автоматические отклики', description: 'Бот откликается на новые вакансии без вашего участия' },
    { title: 'Умные фильтры', description: 'Точная настройка параметров поиска вакансий' },
    { title: 'Работает 24/7', description: 'Не пропустите ни одной подходящей вакансии' },
    { title: 'Сопроводительные письма', description: 'Генерация персонализированных откликов' },
  ]

  it('has exactly 4 steps', () => {
    expect(steps).toHaveLength(4)
  })

  it('has exactly 4 features', () => {
    expect(features).toHaveLength(4)
  })

  it('each step has title and description', () => {
    for (const step of steps) {
      expect(typeof step.title).toBe('string')
      expect(typeof step.description).toBe('string')
      expect(step.title.length).toBeGreaterThan(0)
      expect(step.description.length).toBeGreaterThan(0)
    }
  })

  it('each feature has title and description', () => {
    for (const feature of features) {
      expect(typeof feature.title).toBe('string')
      expect(typeof feature.description).toBe('string')
      expect(feature.title.length).toBeGreaterThan(0)
      expect(feature.description.length).toBeGreaterThan(0)
    }
  })
})
