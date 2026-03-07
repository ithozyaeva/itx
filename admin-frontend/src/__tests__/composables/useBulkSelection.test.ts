import { describe, expect, it } from 'vitest'
import { useBulkSelection } from '@/composables/useBulkSelection'

describe('useBulkSelection', () => {
  it('starts with empty selection and count 0', () => {
    const { selectedIds, count } = useBulkSelection()

    expect(selectedIds.value.size).toBe(0)
    expect(count.value).toBe(0)
  })

  it('toggleItem adds an item', () => {
    const { toggleItem, isSelected, count } = useBulkSelection()

    toggleItem(1)

    expect(isSelected(1)).toBe(true)
    expect(count.value).toBe(1)
  })

  it('toggleItem removes an already selected item', () => {
    const { toggleItem, isSelected, count } = useBulkSelection()

    toggleItem(1)
    toggleItem(1)

    expect(isSelected(1)).toBe(false)
    expect(count.value).toBe(0)
  })

  it('isSelected returns correct boolean', () => {
    const { toggleItem, isSelected } = useBulkSelection()

    toggleItem(1)

    expect(isSelected(1)).toBe(true)
    expect(isSelected(2)).toBe(false)
  })

  it('toggleAll selects all when not all are selected', () => {
    const { toggleAll, isSelected, count } = useBulkSelection()

    toggleAll([1, 2, 3])

    expect(isSelected(1)).toBe(true)
    expect(isSelected(2)).toBe(true)
    expect(isSelected(3)).toBe(true)
    expect(count.value).toBe(3)
  })

  it('toggleAll clears when all are already selected', () => {
    const { toggleAll, count } = useBulkSelection()

    toggleAll([1, 2, 3])
    toggleAll([1, 2, 3])

    expect(count.value).toBe(0)
  })

  it('clearSelection empties the selection', () => {
    const { toggleItem, clearSelection, count } = useBulkSelection()

    toggleItem(1)
    toggleItem(2)
    clearSelection()

    expect(count.value).toBe(0)
  })

  it('ids computed reflects current state', () => {
    const { toggleItem, ids } = useBulkSelection()

    toggleItem(3)
    toggleItem(7)

    expect(ids.value).toEqual(expect.arrayContaining([3, 7]))
    expect(ids.value).toHaveLength(2)
  })
})
