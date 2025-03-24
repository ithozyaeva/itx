<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'
import { useSidebar } from '@/composables/useSidebar'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import ChevronLeft from '~icons/lucide/chevron-left'
import ChevronRight from '~icons/lucide/chevron-right'

const { isCollapsed, sidebarItems, toggleSidebar } = useSidebar()
const route = useRoute()
const router = useRouter()

const sidebarWidth = computed(() => {
  return isCollapsed.value ? 'w-16' : 'w-64'
})

function navigateTo(path: string) {
  router.push(path)
}
</script>

<template>
  <div
    class="h-screen border-r border-gray-200 transition-all duration-300"
    :class="sidebarWidth"
  >
    <div class="flex flex-col h-full">
      <div class="flex items-center justify-between p-4 border-b border-gray-200">
        <h1 v-if="!isCollapsed" class="text-xl font-bold">
          Админ-панель
        </h1>
        <Button variant="ghost" size="icon" @click="toggleSidebar">
          <ChevronRight v-if="isCollapsed" class="h-5 w-5" />
          <ChevronLeft v-else class="h-5 w-5" />
        </Button>
      </div>

      <div class="flex-1 overflow-y-auto py-4">
        <TooltipProvider>
          <ul class="space-y-2 px-2">
            <li v-for="item in sidebarItems" :key="item.path">
              <Tooltip v-if="isCollapsed">
                <TooltipTrigger as-child>
                  <Button
                    variant="ghost"
                    class="w-full justify-center py-2"
                    :class="[
                      route.path.startsWith(item.path) ? 'bg-accent' : '',
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

              <Button
                v-else
                variant="ghost"
                class="w-full justify-start py-2"
                :class="[
                  route.path.startsWith(item.path) ? 'bg-accent' : '',
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
