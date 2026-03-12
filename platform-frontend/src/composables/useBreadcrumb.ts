import { ref } from 'vue'

const dynamicLabel = ref<string | null>(null)

export function useBreadcrumb() {
  function setLabel(label: string) {
    dynamicLabel.value = label
  }

  function clearLabel() {
    dynamicLabel.value = null
  }

  return { dynamicLabel, setLabel, clearLabel }
}
