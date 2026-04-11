<script setup lang="ts">
import { ChevronRight, Home } from 'lucide-vue-next'
import { computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useBreadcrumb } from '@/composables/useBreadcrumb'

const route = useRoute()
const router = useRouter()
const { dynamicLabel, clearLabel } = useBreadcrumb()

watch(() => route.path, () => {
  clearLabel()
})

interface BreadcrumbItem {
  label: string
  to?: string
}

const breadcrumbs = computed<BreadcrumbItem[]>(() => {
  const meta = route.meta.breadcrumb as BreadcrumbItem[] | undefined
  if (!meta)
    return []

  return meta.map((item, index) => {
    if (index === meta.length - 1 && dynamicLabel.value) {
      return { ...item, label: dynamicLabel.value }
    }
    return item
  })
})
</script>

<template>
  <nav
    v-if="breadcrumbs.length > 0"
    aria-label="Навигация"
    class="flex items-center gap-1.5 text-xs text-muted-foreground px-4 pt-3"
  >
    <button
      class="flex items-center hover:text-foreground transition-colors"
      aria-label="Главная"
      @click="router.push('/')"
    >
      <Home class="h-3 w-3" />
    </button>
    <template
      v-for="(crumb, index) in breadcrumbs"
      :key="index"
    >
      <ChevronRight class="h-3 w-3 opacity-50" />
      <button
        v-if="crumb.to && index < breadcrumbs.length - 1"
        class="hover:text-foreground transition-colors"
        @click="router.push(crumb.to)"
      >
        {{ crumb.label }}
      </button>
      <span
        v-else
        class="text-foreground/70"
      >
        {{ crumb.label }}
      </span>
    </template>
  </nav>
</template>
