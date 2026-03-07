import type { Mentor } from '@/models/profile'
import { mount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'

vi.mock('itx-ui-kit', () => ({
  Typography: { template: '<div><slot /></div>' },
  Tag: { template: '<span><slot /></span>' },
}))

import MentorCard from '@/components/mentors/MentorCard.vue'

function createMentor(overrides: Partial<Mentor> = {}): Mentor {
  return {
    id: 1,
    telegramID: 123,
    tg: 'johndoe',
    birthday: '1990-01-01',
    firstName: 'John',
    lastName: 'Doe',
    bio: 'bio',
    grade: 'Senior',
    company: 'ACME',
    avatarUrl: '',
    roles: ['MENTOR'],
    occupation: 'Engineer',
    experience: '10 years',
    profTags: [{ id: 1, title: 'Go' }, { id: 2, title: 'Vue' }],
    contacts: [],
    services: [],
    ...overrides,
  }
}

describe('MentorCard', () => {
  const globalConfig = {
    stubs: {
      RouterLink: {
        template: '<a :href="to"><slot /></a>',
        props: ['to'],
      },
    },
  }

  it('renders mentor name', () => {
    const mentor = createMentor()
    const wrapper = mount(MentorCard, {
      props: { mentor },
      global: globalConfig,
    })
    expect(wrapper.text()).toContain('John')
    expect(wrapper.text()).toContain('Doe')
  })

  it('renders occupation when present', () => {
    const wrapper = mount(MentorCard, {
      props: { mentor: createMentor({ occupation: 'Software Engineer' }) },
      global: globalConfig,
    })
    expect(wrapper.text()).toContain('Software Engineer')
  })

  it('does not render occupation when absent', () => {
    const wrapper = mount(MentorCard, {
      props: { mentor: createMentor({ occupation: '' }) },
      global: globalConfig,
    })
    // The paragraph with occupation should not be visible
    const occupationP = wrapper.findAll('p').find(p => p.classes().includes('text-muted-foreground'))
    // Either doesn't exist or doesn't have text
    expect(occupationP?.exists() ?? false).toBe(false)
  })

  it('renders experience when present', () => {
    const wrapper = mount(MentorCard, {
      props: { mentor: createMentor({ experience: '5 years of Go' }) },
      global: globalConfig,
    })
    expect(wrapper.text()).toContain('5 years of Go')
  })

  it('renders prof tags', () => {
    const wrapper = mount(MentorCard, {
      props: { mentor: createMentor() },
      global: globalConfig,
    })
    expect(wrapper.text()).toContain('Go')
    expect(wrapper.text()).toContain('Vue')
  })

  it('does not render tags section when profTags is empty', () => {
    const wrapper = mount(MentorCard, {
      props: { mentor: createMentor({ profTags: [] }) },
      global: globalConfig,
    })
    expect(wrapper.find('.flex-wrap').exists()).toBe(false)
  })

  it('renders telegram link when tg is present', () => {
    const wrapper = mount(MentorCard, {
      props: { mentor: createMentor({ tg: 'myhandle' }) },
      global: globalConfig,
    })
    const link = wrapper.find('a[href="https://t.me/myhandle"]')
    expect(link.exists()).toBe(true)
    expect(link.text()).toContain('@myhandle')
  })

  it('links to correct mentor page', () => {
    const wrapper = mount(MentorCard, {
      props: { mentor: createMentor({ id: 42 }) },
      global: globalConfig,
    })
    const link = wrapper.find('a[href="/mentors/42"]')
    expect(link.exists()).toBe(true)
  })
})
