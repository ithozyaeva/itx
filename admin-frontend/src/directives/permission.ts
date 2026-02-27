import type { Directive, DirectiveBinding } from 'vue'
import type { Permission } from '@/types/permissions'
import { watch } from 'vue'
import { usePermissionsStore } from '@/stores/permissions'

export const permissionDirective: Directive<HTMLElement, Permission | Permission[]> = {
  mounted(el, binding) {
    const store = usePermissionsStore()
    watch(
      () => store.userPermissions,
      () => updateElementVisibility(el, binding, store),
      { immediate: true, deep: true },
    )
  },
}

function updateElementVisibility(
  el: HTMLElement,
  binding: DirectiveBinding<Permission | Permission[]>,
  store: ReturnType<typeof usePermissionsStore>,
) {
  const { value } = binding

  if (!value) {
    console.warn('v-permission directive requires a permission value')
    return
  }

  let hasAccess = false

  if (Array.isArray(value)) {
    hasAccess = store.hasAnyPermission(value)
  }
  else {
    hasAccess = store.hasPermission(value)
  }

  el.style.display = hasAccess ? '' : 'none'
}

export default permissionDirective
