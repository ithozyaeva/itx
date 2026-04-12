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
import { Calendar, ClipboardList, FileText, Folder, MessageSquare, Users } from 'lucide-vue-next'
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
  { key: 'totalMembers', label: 'Участники', icon: Users, color: 'accent' },
  { key: 'totalMentors', label: 'Менторы', icon: Users, color: 'term-cyan' },
  { key: 'upcomingEvents', label: 'Предстоящие события', icon: Calendar, color: 'term-amber' },
  { key: 'pastEvents', label: 'Прошедшие события', icon: Calendar, color: 'muted-foreground' },
  { key: 'pendingReviews', label: 'Ожидают публикации', icon: MessageSquare, color: 'term-magenta' },
  { key: 'approvedReviews', label: 'Опубликованные', icon: MessageSquare, color: 'accent' },
  { key: 'referralLinks', label: 'Реф. ссылки', icon: Folder, color: 'term-cyan' },
  { key: 'resumes', label: 'Резюме', icon: FileText, color: 'term-amber' },
  { key: 'openTasks', label: 'Открытые', icon: ClipboardList, color: 'accent' },
  { key: 'inProgressTasks', label: 'В работе', icon: ClipboardList, color: 'term-amber' },
  { key: 'doneTasks', label: 'Выполнены', icon: ClipboardList, color: 'term-cyan' },
  { key: 'approvedTasks', label: 'Приняты', icon: ClipboardList, color: 'accent' },
] as const

const memberGrowthData = computed(() => ({
  labels: chartStats.value?.memberGrowth.map(m => m.month) ?? [],
  datasets: [{
    label: 'Участники',
    data: chartStats.value?.memberGrowth.map(m => m.count) ?? [],
    borderColor: '#4ade80',
    backgroundColor: 'rgba(74, 222, 128, 0.08)',
    fill: true,
    tension: 0.4,
    pointBackgroundColor: '#4ade80',
    pointBorderColor: '#4ade80',
    pointRadius: 3,
    pointHoverRadius: 5,
    borderWidth: 2,
  }],
}))

const eventAttendanceData = computed(() => ({
  labels: chartStats.value?.eventAttendance.map(m => m.month) ?? [],
  datasets: [{
    label: 'Посещаемость',
    data: chartStats.value?.eventAttendance.map(m => m.count) ?? [],
    backgroundColor: 'rgba(94, 234, 212, 0.6)',
    hoverBackgroundColor: 'rgba(94, 234, 212, 0.8)',
    borderRadius: 2,
    borderWidth: 0,
  }],
}))

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { display: false },
  },
  scales: {
    y: {
      beginAtZero: true,
      grid: {
        color: 'rgba(74, 222, 128, 0.06)',
      },
      ticks: {
        font: { family: 'JetBrains Mono', size: 10 },
      },
    },
    x: {
      grid: { display: false },
      ticks: {
        font: { family: 'JetBrains Mono', size: 10 },
      },
    },
  },
}
</script>

<template>
  <AdminLayout>
    <div ref="containerRef" class="space-y-6">
      <div class="flex items-center gap-3">
        <h1 class="text-xl font-semibold">
          Дашборд
        </h1>
        <span class="font-mono text-[10px] text-muted-foreground uppercase tracking-wider">// overview</span>
      </div>

      <div v-if="stats" class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-3">
        <div
          v-for="card in statCards"
          :key="card.key"
          class="terminal-stat bg-card text-card-foreground p-3 lg:p-4"
          data-reveal
        >
          <div class="flex items-center justify-between mb-2">
            <span class="font-mono text-[10px] uppercase tracking-wider text-muted-foreground">{{ card.label }}</span>
            <component :is="card.icon" class="h-3.5 w-3.5 text-muted-foreground/60" />
          </div>
          <p class="text-2xl lg:text-3xl font-bold font-mono tabular-nums">
            {{ stats[card.key] }}
          </p>
        </div>
      </div>

      <div v-if="chartStats" class="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <Card data-reveal>
          <CardHeader class="p-4 pb-2">
            <CardTitle class="text-sm font-mono flex items-center gap-2">
              <span class="w-1.5 h-1.5 rounded-full bg-accent" />
              Рост участников
            </CardTitle>
          </CardHeader>
          <CardContent class="p-4 pt-0">
            <div class="h-52 lg:h-64">
              <Line :data="memberGrowthData" :options="chartOptions" />
            </div>
          </CardContent>
        </Card>

        <Card data-reveal>
          <CardHeader class="p-4 pb-2">
            <CardTitle class="text-sm font-mono flex items-center gap-2">
              <span class="w-1.5 h-1.5 rounded-full bg-term-cyan" />
              Посещаемость событий
            </CardTitle>
          </CardHeader>
          <CardContent class="p-4 pt-0">
            <div class="h-52 lg:h-64">
              <Bar :data="eventAttendanceData" :options="chartOptions" />
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  </AdminLayout>
</template>
