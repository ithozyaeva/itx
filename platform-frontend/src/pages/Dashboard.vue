<script setup lang="ts">
import type { CommunityEvent } from '@/models/event'
import { Loader2, Lock } from 'lucide-vue-next'
import { onMounted, ref } from 'vue'
import EventsSection from '@/components/dashboard/EventsSection.vue'
import GreetingCard from '@/components/dashboard/GreetingCard.vue'
import NearestEvent from '@/components/dashboard/NearestEvent.vue'
import ReferralsSidebar from '@/components/dashboard/ReferralsSidebar.vue'
import { useUserLevel } from '@/composables/useUser'
import { handleError } from '@/services/errorService'
import { eventsService } from '@/services/events'

const nearestEvent = ref<CommunityEvent | null>(null)
const isLoading = ref(false)
const { levelIndex } = useUserLevel()

async function loadNearestEvent() {
  isLoading.value = true
  try {
    const result = await eventsService.searchNext(1, 0)
    nearestEvent.value = result.items[0] ?? null
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isLoading.value = false
  }
}

onMounted(() => loadNearestEvent())
</script>

<template>
  <div class="flex gap-6 px-4 py-6 md:py-8">
    <!-- Main content -->
    <div class="flex-1 min-w-0">
      <GreetingCard />

      <div
        v-if="isLoading"
        class="flex justify-center py-12"
      >
        <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
      </div>

      <NearestEvent
        v-else-if="nearestEvent"
        :event="nearestEvent"
        class="mt-6"
      />

      <EventsSection />

      <!-- Locked content teaser -->
      <div
        v-if="levelIndex < 2"
        class="mt-6 rounded-3xl border border-dashed border-muted-foreground/30 bg-muted/30 p-6 text-center"
      >
        <Lock class="h-8 w-8 mx-auto text-muted-foreground/50 mb-2" />
        <p class="text-muted-foreground font-medium">
          Больше контента доступно на уровне «Хозяин»
        </p>
        <p class="text-sm text-muted-foreground/60 mt-1">
          Повышайте уровень, чтобы разблокировать эксклюзивные материалы
        </p>
      </div>
    </div>

    <!-- Right sidebar (desktop only) -->
    <div class="hidden lg:block w-80 shrink-0">
      <ReferralsSidebar />
    </div>
  </div>
</template>
