<script setup lang="ts">
import { X } from 'lucide-vue-next'
import { watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import ExternalLink from '~icons/lucide/external-link'
import { Button } from '@/components/ui/button'
import { useSidebar } from '@/composables/useSidebar'

const { isMobileOpen, sidebarGroups, closeMobileSidebar } = useSidebar()
const route = useRoute()
const router = useRouter()

function isActive(path: string) {
  if (path === '/dashboard')
    return route.path === '/dashboard'
  return route.path === path || route.path.startsWith(`${path}/`)
}

function navigateTo(path: string) {
  router.push(path)
  closeMobileSidebar()
}

watch(() => route.path, () => {
  closeMobileSidebar()
})
</script>

<template>
  <Transition name="fade">
    <div
      v-if="isMobileOpen"
      class="fixed inset-0 bg-black/60 backdrop-blur-sm z-40 lg:hidden"
      @click="closeMobileSidebar"
    />
  </Transition>

  <div
    class="h-screen border-r border-sidebar-border bg-sidebar text-sidebar-foreground transition-all duration-300 flex-shrink-0 w-56"
    :class="[
      isMobileOpen
        ? 'fixed inset-y-0 left-0 z-50 w-56 lg:relative lg:z-auto'
        : 'hidden lg:block',
    ]"
  >
    <div class="flex flex-col h-full">
      <div class="flex items-center justify-between px-4 py-3 border-b border-sidebar-border">
        <div class="flex items-center gap-2">
          <span class="w-2 h-2 rounded-full bg-term-amber shadow-[0_0_6px_#ffb547]" />
          <span class="font-mono text-xs font-medium tracking-wider text-sidebar-foreground/80">[ADM] IT-X</span>
        </div>
        <button
          class="p-1.5 lg:hidden text-sidebar-foreground/60 hover:text-sidebar-foreground transition-colors"
          aria-label="Закрыть меню"
          @click="closeMobileSidebar"
        >
          <X class="w-5 h-5" />
        </button>
      </div>

      <div class="flex-1 overflow-y-auto py-3">
        <div
          v-for="(group, groupIndex) in sidebarGroups"
          :key="groupIndex"
          :class="groupIndex > 0 ? 'mt-3' : ''"
        >
          <p
            v-if="group.label"
            class="px-4 mb-1.5 text-[11px] font-mono font-medium uppercase tracking-wider text-sidebar-foreground/40"
          >
            // {{ group.label }}
          </p>
          <ul class="space-y-0.5 px-2">
            <li
              v-for="(item, itemIndex) in group.items"
              :key="item.path"
            >
              <Button
                variant="ghost"
                class="w-full justify-start py-2 cursor-pointer rounded-sm text-sidebar-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground"
                :class="[
                  isActive(item.path) ? 'bg-accent text-accent-foreground' : '',
                ]"
                @click="navigateTo(item.path)"
              >
                <span class="font-mono text-[10px] text-sidebar-foreground/40 mr-1.5 w-4 shrink-0">{{ String(itemIndex + 1).padStart(2, '0') }}</span>
                <component
                  :is="item.icon"
                  class="h-4 w-4 mr-2"
                />
                <span class="text-sm">{{ item.title }}</span>
              </Button>
            </li>
          </ul>
        </div>
      </div>

      <div class="border-t border-sidebar-border p-2">
        <a
          href="/platform"
          class="flex items-center gap-2 rounded-sm px-3 py-2 text-sm text-sidebar-foreground/70 hover:bg-sidebar-accent hover:text-sidebar-accent-foreground transition cursor-pointer"
        >
          <ExternalLink class="h-4 w-4" />
          <span>Платформа</span>
        </a>
      </div>
    </div>
  </div>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
