<script setup lang="ts">
import { Typography } from 'itx-ui-kit'
import { X } from 'lucide-vue-next'
import { computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import ChevronLeft from '~icons/lucide/chevron-left'
import ChevronRight from '~icons/lucide/chevron-right'
import { Button } from '@/components/ui/button'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'
import { useSidebar } from '@/composables/useSidebar'

const { isCollapsed, isMobileOpen, sidebarItems, toggleSidebar, closeMobileSidebar } = useSidebar()
const route = useRoute()
const router = useRouter()

const sidebarWidth = computed(() => {
  return isCollapsed.value ? 'w-16' : 'w-64'
})

function navigateTo(path: string) {
  router.push(path)
  closeMobileSidebar()
}

// Close mobile sidebar on route change
watch(() => route.path, () => {
  closeMobileSidebar()
})
</script>

<template>
  <!-- Mobile overlay backdrop -->
  <Transition name="fade">
    <div
      v-if="isMobileOpen"
      class="fixed inset-0 bg-black/50 z-40 lg:hidden"
      @click="closeMobileSidebar"
    />
  </Transition>

  <!-- Sidebar -->
  <div
    class="h-screen border-r border-border bg-primary text-primary-foreground transition-all duration-300 flex-shrink-0"
    :class="[
      sidebarWidth,
      isMobileOpen
        ? 'fixed inset-y-0 left-0 z-50 w-64 lg:relative lg:z-auto'
        : 'hidden lg:block',
    ]"
  >
    <div class="flex flex-col h-full">
      <div class="flex items-center justify-between p-4 border-b border-border/20">
        <Typography v-if="!isCollapsed || isMobileOpen" variant="h4" as="h1">
          Админ-панель
        </Typography>
        <!-- Close button on mobile -->
        <Button
          class="lg:hidden"
          variant="ghost"
          size="icon"
          @click="closeMobileSidebar"
        >
          <X class="h-5 w-5" />
        </Button>
        <!-- Collapse toggle on desktop -->
        <Button
          class="hidden lg:flex"
          variant="ghost"
          size="icon"
          @click="toggleSidebar"
        >
          <ChevronRight v-if="isCollapsed" class="h-5 w-5" />
          <ChevronLeft v-else class="h-5 w-5" />
        </Button>
      </div>

      <div class="flex-1 overflow-y-auto py-4">
        <TooltipProvider>
          <ul class="space-y-2 px-2">
            <li v-for="item in sidebarItems" :key="item.path">
              <!-- Collapsed desktop: icon with tooltip -->
              <Tooltip v-if="isCollapsed && !isMobileOpen">
                <TooltipTrigger as-child>
                  <Button
                    variant="ghost"
                    class="w-full justify-center py-2 text-primary-foreground hover:bg-accent hover:text-accent-foreground"
                    :class="[
                      route.path.startsWith(item.path) ? 'bg-accent text-accent-foreground' : '',
                    ]"
                    @click="navigateTo(item.path)"
                  >
                    <component :is="item.icon" class="h-5 w-5" />
                  </Button>
                </TooltipTrigger>
                <TooltipContent side="right">
                  {{ item.title }}
                </TooltipContent>
              </Tooltip>

              <!-- Expanded or mobile: icon + text -->
              <Button
                v-else
                variant="ghost"
                class="w-full justify-start py-2 text-primary-foreground hover:bg-accent hover:text-accent-foreground"
                :class="[
                  route.path.startsWith(item.path) ? 'bg-accent text-accent-foreground' : '',
                ]"
                @click="navigateTo(item.path)"
              >
                <component :is="item.icon" class="h-5 w-5 mr-2" />
                <span>{{ item.title }}</span>
              </Button>
            </li>
          </ul>
        </TooltipProvider>
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
