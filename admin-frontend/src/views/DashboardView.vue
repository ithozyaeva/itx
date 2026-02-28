<script setup lang="ts">
import type { ChartStats, DashboardStats } from '@/services/statsService'
import {
  BarElement,
  CategoryScale,
  Chart as ChartJS,
  Legend,
  LinearScale,
  LineElement,
  PointElement,
  Title,
  Tooltip,
} from 'chart.js'
import { Typography } from 'itx-ui-kit'
import { Calendar, FileText, Folder, MessageSquare, Users } from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'
import { Bar, Line } from 'vue-chartjs'
import AdminLayout from '@/components/layout/AdminLayout.vue'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { useCardReveal } from '@/composables/useCardReveal'
import { statsService } from '@/services/statsService'

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, BarElement, Title, Tooltip, Legend)

const containerRef = ref<HTMLElement | null>(null)
useCardReveal(containerRef)

const stats = ref<DashboardStats | null>(null)
const chartStats = ref<ChartStats | null>(null)

onMounted(async () => {
  try {
    const [s, c] = await Promise.all([
      statsService.getStats(),
      statsService.getChartStats(),
    ])
    stats.value = s
    chartStats.value = c
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

const memberGrowthData = computed(() => ({
  labels: chartStats.value?.memberGrowth.map(m => m.month) ?? [],
  datasets: [{
    label: 'Участники',
    data: chartStats.value?.memberGrowth.map(m => m.count) ?? [],
    borderColor: 'hsl(var(--primary))',
    backgroundColor: 'hsl(var(--primary) / 0.1)',
    fill: true,
    tension: 0.3,
  }],
}))

const eventAttendanceData = computed(() => ({
  labels: chartStats.value?.eventAttendance.map(m => m.month) ?? [],
  datasets: [{
    label: 'Посещаемость',
    data: chartStats.value?.eventAttendance.map(m => m.count) ?? [],
    backgroundColor: 'hsl(var(--primary) / 0.7)',
    borderRadius: 4,
  }],
}))

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { display: false },
  },
  scales: {
    y: { beginAtZero: true },
  },
}
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

      <div v-if="chartStats" class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <Card data-reveal>
          <CardHeader>
            <CardTitle>Рост участников</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="h-64">
              <Line :data="memberGrowthData" :options="chartOptions" />
            </div>
          </CardContent>
        </Card>

        <Card data-reveal>
          <CardHeader>
            <CardTitle>Посещаемость событий</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="h-64">
              <Bar :data="eventAttendanceData" :options="chartOptions" />
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  </AdminLayout>
</template>
