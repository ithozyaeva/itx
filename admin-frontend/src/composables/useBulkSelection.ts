import { computed, ref } from 'vue'

export function useBulkSelection() {
  const selectedIds = ref<Set<number>>(new Set())

  const count = computed(() => selectedIds.value.size)
  const ids = computed(() => Array.from(selectedIds.value))

  function isSelected(id: number) {
    return selectedIds.value.has(id)
  }

  function toggleItem(id: number) {
    const next = new Set(selectedIds.value)
    if (next.has(id)) {
      next.delete(id)
    }
    else {
      next.add(id)
    }
    selectedIds.value = next
  }

  function toggleAll(allIds: number[]) {
    if (selectedIds.value.size === allIds.length) {
      selectedIds.value = new Set()
    }
    else {
      selectedIds.value = new Set(allIds)
    }
  }

  function clearSelection() {
    selectedIds.value = new Set()
  }

  return {
    selectedIds,
    count,
    ids,
    isSelected,
    toggleItem,
    toggleAll,
    clearSelection,
  }
}
