import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import UiAccordion from '@/components/ui/UiAccordion.vue'
import UiButton from '@/components/ui/UiButton.vue'
import UiTag from '@/components/ui/UiTag.vue'
import UiTypography from '@/components/ui/UiTypography.vue'

describe('uiTypography', () => {
  it('renders as <p> by default', () => {
    const wrapper = mount(UiTypography, { slots: { default: 'Hello' } })
    expect(wrapper.element.tagName).toBe('P')
  })

  it('renders correct tag with `as` prop', () => {
    const wrapper = mount(UiTypography, {
      props: { as: 'span' },
      slots: { default: 'Text' },
    })
    expect(wrapper.element.tagName).toBe('SPAN')
  })

  it('applies variant class h1', () => {
    const wrapper = mount(UiTypography, {
      props: { variant: 'h1' },
      slots: { default: 'Heading' },
    })
    expect(wrapper.classes()).toContain('h1')
  })

  it('applies variant class h2', () => {
    const wrapper = mount(UiTypography, {
      props: { variant: 'h2' },
      slots: { default: 'Heading' },
    })
    expect(wrapper.classes()).toContain('h2')
  })

  it('applies variant class body-m', () => {
    const wrapper = mount(UiTypography, {
      props: { variant: 'body-m' },
      slots: { default: 'Body' },
    })
    expect(wrapper.classes()).toContain('body-m')
  })

  it('defaults to body-l variant class', () => {
    const wrapper = mount(UiTypography, { slots: { default: 'Text' } })
    expect(wrapper.classes()).toContain('body-l')
  })

  it('renders slot content', () => {
    const wrapper = mount(UiTypography, {
      slots: { default: '<em>rich</em> text' },
    })
    expect(wrapper.text()).toContain('rich text')
    expect(wrapper.find('em').exists()).toBe(true)
  })
})

describe('uiButton', () => {
  it('renders as <button> by default with filled variant', () => {
    const wrapper = mount(UiButton, { slots: { default: 'Click' } })
    expect(wrapper.element.tagName).toBe('BUTTON')
    expect(wrapper.classes()).toContain('filled')
  })

  it('renders disabled state', () => {
    const wrapper = mount(UiButton, {
      props: { disabled: true },
      slots: { default: 'Disabled' },
    })
    expect(wrapper.classes()).toContain('disabled')
    expect(wrapper.attributes('disabled')).toBeDefined()
  })

  it('renders stroke variant', () => {
    const wrapper = mount(UiButton, {
      props: { variant: 'stroke' },
      slots: { default: 'Stroke' },
    })
    expect(wrapper.classes()).toContain('stroke')
    expect(wrapper.classes()).not.toContain('filled')
  })

  it('renders dark-filled variant', () => {
    const wrapper = mount(UiButton, {
      props: { variant: 'dark-filled' },
      slots: { default: 'Dark' },
    })
    expect(wrapper.classes()).toContain('dark-filled')
  })

  it('renders as anchor when as="a"', () => {
    const wrapper = mount(UiButton, {
      props: { as: 'a' },
      slots: { default: 'Link' },
    })
    expect(wrapper.element.tagName).toBe('A')
  })

  it('renders slot content', () => {
    const wrapper = mount(UiButton, {
      slots: { default: 'Submit' },
    })
    expect(wrapper.text()).toBe('Submit')
  })
})

describe('uiTag', () => {
  it('renders with default variant', () => {
    const wrapper = mount(UiTag, { slots: { default: 'Vue' } })
    expect(wrapper.classes()).toContain('default')
    expect(wrapper.element.tagName).toBe('BUTTON')
  })

  it('renders with active variant', () => {
    const wrapper = mount(UiTag, {
      props: { variant: 'active' },
      slots: { default: 'Active' },
    })
    expect(wrapper.classes()).toContain('active')
    expect(wrapper.classes()).not.toContain('default')
  })

  it('renders disabled state', () => {
    const wrapper = mount(UiTag, {
      props: { disabled: true },
      slots: { default: 'Disabled' },
    })
    expect(wrapper.classes()).toContain('disabled')
    expect(wrapper.attributes('disabled')).toBeDefined()
    expect(wrapper.attributes('aria-disabled')).toBe('true')
  })

  it('renders slot content', () => {
    const wrapper = mount(UiTag, {
      slots: { default: 'TypeScript' },
    })
    expect(wrapper.text()).toBe('TypeScript')
  })
})

describe('uiAccordion', () => {
  it('renders with title', () => {
    const wrapper = mount(UiAccordion, {
      props: { title: 'FAQ Question' },
    })
    expect(wrapper.text()).toContain('FAQ Question')
  })

  it('starts closed (content wrapper collapsed)', () => {
    const wrapper = mount(UiAccordion, {
      props: { title: 'Question' },
      slots: { default: 'Answer' },
    })
    const contentWrapper = wrapper.find('.accordion-wrapper')
    expect(contentWrapper.classes()).not.toContain('accordion-wrapper--open')
    expect(wrapper.classes()).not.toContain('open')
  })

  it('toggles open on click', async () => {
    const wrapper = mount(UiAccordion, {
      props: { title: 'Question' },
      slots: { default: 'Answer text' },
    })

    await wrapper.find('.accordion-header').trigger('click')

    expect(wrapper.find('.accordion-wrapper').classes()).toContain('accordion-wrapper--open')
    expect(wrapper.classes()).toContain('open')

    // Click again to close
    await wrapper.find('.accordion-header').trigger('click')

    expect(wrapper.find('.accordion-wrapper').classes()).not.toContain('accordion-wrapper--open')
    expect(wrapper.classes()).not.toContain('open')
  })

  it('shows plus icon when closed and cross icon when open', async () => {
    const wrapper = mount(UiAccordion, {
      props: { title: 'Question' },
    })

    // Closed: plus icon (path with "M20 5v30M35 20H5")
    let paths = wrapper.findAll('path')
    expect(paths.some(p => p.attributes('d')?.includes('M20 5v30'))).toBe(true)

    await wrapper.find('.accordion-header').trigger('click')

    // Open: cross icon (path with "31 9 9 31")
    paths = wrapper.findAll('path')
    expect(paths.some(p => p.attributes('d')?.includes('31 9 9 31'))).toBe(true)
  })

  it('renders content prop as fallback', () => {
    const wrapper = mount(UiAccordion, {
      props: { title: 'Q', content: 'Fallback answer' },
    })
    expect(wrapper.text()).toContain('Fallback answer')
  })
})
