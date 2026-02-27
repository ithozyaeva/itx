<script setup lang="ts">
import type { DashboardStats } from '@/services/statsService'
import { Typography } from 'itx-ui-kit'
import { Calendar, FileText, Folder, MessageSquare, Users } from 'lucide-vue-next'
import { onMounted, ref } from 'vue'
import AdminLayout from '@/components/layout/AdminLayout.vue'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { useCardReveal } from '@/composables/useCardReveal'
import { statsService } from '@/services/statsService'

const containerRef = ref<HTMLElement | null>(null)
useCardReveal(containerRef)

const stats = ref<DashboardStats | null>(null)

onMounted(async () => {
  try {
    stats.value = await statsService.getStats()
  }
  catch (error) {
    console.error('Ошибка загрузки статистики:', error)
  }
})

const statCards = [
  { key: 'totalMembers', label: 'Участники', icon: Users },
  { key: 'totalMentors', label: 'Менторы', icon: Users },
  { key: 'upcomingEvents', label: 'Предстоящие события', icon: Calendar },
  { key: 'pastEvents', label: 'Прошедшие события', icon: Calendar },
  { key: 'pendingReviews', label: 'Ожидают публикации', icon: MessageSquare },
  { key: 'approvedReviews', label: 'Опубликованные отзывы', icon: MessageSquare },
  { key: 'referralLinks', label: 'Реферальные ссылки', icon: Folder },
  { key: 'resumes', label: 'Резюме', icon: FileText },
] as const
</script>

<template>
  <AdminLayout>
    <div ref="containerRef" class="space-y-6">
      <Typography variant="h2" as="h1">
        Дашборд
      </Typography>

      <div v-if="stats" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <Card v-for="card in statCards" :key="card.key" data-reveal>
          <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle class="text-sm font-medium">
              {{ card.label }}
            </CardTitle>
            <component :is="card.icon" class="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <p class="text-3xl font-bold">
              {{ stats[card.key] }}
            </p>
          </CardContent>
        </Card>
      </div>
    </div>
  </AdminLayout>
</template>
