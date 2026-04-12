import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '../components/table'

describe('Table', () => {
  it('renders a table element', () => {
    const wrapper = mount(Table, {
      slots: { default: '<tbody><tr><td>Cell</td></tr></tbody>' },
    })
    expect(wrapper.find('table').exists()).toBe(true)
  })

  it('wraps in scrollable container', () => {
    const wrapper = mount(Table, {
      slots: { default: '' },
    })
    expect(wrapper.classes().join(' ')).toContain('overflow-auto')
  })
})

describe('TableHead', () => {
  it('renders th with mono font for terminal style', () => {
    const wrapper = mount(TableHead, {
      slots: { default: 'Column' },
    })
    expect(wrapper.element.tagName).toBe('TH')
    const classes = wrapper.classes().join(' ')
    expect(classes).toContain('font-mono')
    expect(classes).toContain('uppercase')
    expect(classes).toContain('tracking-wider')
  })
})

describe('TableRow', () => {
  it('renders tr with hover effect', () => {
    const wrapper = mount(TableRow, {
      slots: { default: '<td>Data</td>' },
    })
    expect(wrapper.element.tagName).toBe('TR')
    expect(wrapper.classes().join(' ')).toContain('hover:bg-accent/10')
  })
})

describe('TableCell', () => {
  it('renders td', () => {
    const wrapper = mount(TableCell, {
      slots: { default: 'Value' },
    })
    expect(wrapper.element.tagName).toBe('TD')
    expect(wrapper.text()).toBe('Value')
  })
})

describe('TableHeader', () => {
  it('renders thead', () => {
    const wrapper = mount(TableHeader, {
      slots: { default: '<tr><th>H</th></tr>' },
    })
    expect(wrapper.element.tagName).toBe('THEAD')
  })
})

describe('TableBody', () => {
  it('renders tbody', () => {
    const wrapper = mount(TableBody, {
      slots: { default: '<tr><td>D</td></tr>' },
    })
    expect(wrapper.element.tagName).toBe('TBODY')
  })
})
