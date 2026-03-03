<script setup lang="ts">
import type { ChatActivityStats, DailyActivity, TopUser, TrackedChat, UserStats } from '@/services/chatActivityService'
import {
  CategoryScale,
  Chart as ChartJS,
  Filler,
  Legend,
  LinearScale,
  LineElement,
  PointElement,
  Title,
  Tooltip,
} from 'chart.js'
import { Typography } from 'itx-ui-kit'
import { ArrowDown, ArrowUp, BarChart3, Download, MessageSquare, Users, X } from 'lucide-vue-next'
import { computed, onMounted, ref, watch } from 'vue'
import { Line } from 'vue-chartjs'
import AdminLayout from '@/components/layout/AdminLayout.vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { useCardReveal } from '@/composables/useCardReveal'
import { chatActivityService } from '@/services/chatActivityService'

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend, Filler)

const containerRef = ref<HTMLElement | null>(null)
useCardReveal(containerRef)

const stats = ref<ChatActivityStats | null>(null)
const chartData = ref<DailyActivity[]>([])
const topUsers = ref<TopUser[]>([])
const chats = ref<TrackedChat[]>([])
const selectedChatId = ref<number | undefined>(undefined)
const selectedDays = ref(30)
const isLoading = ref(true)

// Фильтр по пользователю
const selectedUser = ref<TopUser | null>(null)
const userStats = ref<UserStats | null>(null)
const userChartData = ref<DailyActivity[]>([])

onMounted(async () => {
  try {
    const [s, c, t, ch] = await Promise.all([
      chatActivityService.getStats(),
      chatActivityService.getChart(undefined, 30),
      chatActivityService.getTopUsers(7, 5),
      chatActivityService.getChats(),
    ])
    stats.value = s
    chartData.value = c
    topUsers.value = t
    chats.value = ch
  }
  catch (error) {
    console.error('Ошибка загрузки статистики активности:', error)
  }
  finally {
    isLoading.value = false
  }
})

async function loadChart() {
  try {
    chartData.value = await chatActivityService.getChart(selectedChatId.value, selectedDays.value)
  }
  catch (error) {
    console.error('Ошибка загрузки графика:', error)
  }
}

watch([selectedChatId, selectedDays], () => {
  loadChart()
})

function calcChange(current: number, previous: number): number | null {
  if (previous === 0)
    return current > 0 ? 100 : null
  return Math.round(((current - previous) / previous) * 100)
}

const messagesChange = computed(() =>
  stats.value ? calcChange(stats.value.totalMessagesWeek, stats.value.totalMessagesLastWeek) : null,
)
const usersChange = computed(() =>
  stats.value ? calcChange(stats.value.uniqueUsersWeek, stats.value.uniqueUsersLastWeek) : null,
)

const summaryCards = computed(() => [
  { label: 'Сообщений сегодня', value: stats.value?.totalMessagesToday ?? 0, icon: MessageSquare, change: null },
  { label: 'Сообщений за неделю', value: stats.value?.totalMessagesWeek ?? 0, icon: BarChart3, change: messagesChange.value },
  { label: 'Активных сегодня', value: stats.value?.uniqueUsersToday ?? 0, icon: Users, change: null },
  { label: 'Активных за неделю', value: stats.value?.uniqueUsersWeek ?? 0, icon: Users, change: usersChange.value },
])

const lineChartData = computed(() => ({
  labels: chartData.value.map(d => d.date.slice(5)),
  datasets: [{
    label: 'Сообщения',
    data: chartData.value.map(d => d.count),
    borderColor: 'hsl(var(--primary))',
    backgroundColor: 'hsl(var(--primary) / 0.1)',
    fill: true,
    tension: 0.3,
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

const userLineChartData = computed(() => ({
  labels: userChartData.value.map(d => d.date.slice(5)),
  datasets: [{
    label: 'Сообщения',
    data: userChartData.value.map(d => d.count),
    borderColor: 'hsl(var(--primary))',
    backgroundColor: 'hsl(var(--primary) / 0.1)',
    fill: true,
    tension: 0.3,
  }],
}))

async function selectUser(user: TopUser) {
  selectedUser.value = user
  try {
    const [us, uc] = await Promise.all([
      chatActivityService.getUserStats(user.telegramUserId, 30),
      chatActivityService.getChart(undefined, 30, user.telegramUserId),
    ])
    userStats.value = us
    userChartData.value = uc
  }
  catch (error) {
    console.error('Ошибка загрузки статистики пользователя:', error)
  }
}

function clearUserSelection() {
  selectedUser.value = null
  userStats.value = null
  userChartData.value = []
}

async function handleExport() {
  try {
    await chatActivityService.exportCSV(selectedDays.value, selectedChatId.value)
  }
  catch (error) {
    console.error('Ошибка экспорта:', error)
  }
}
</script>

<template>
  <AdminLayout>
    <div ref="containerRef" class="space-y-6">
      <Typography variant="h2" as="h1">
        Активность чатов
      </Typography>

      <!-- Счётчики -->
      <div v-if="stats" class="grid grid-cols-2 lg:grid-cols-4 gap-3 lg:gap-4">
        <Card v-for="card in summaryCards" :key="card.label" data-reveal>
          <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2 p-3 lg:p-6 lg:pb-2">
            <CardTitle class="text-xs lg:text-sm font-medium">
              {{ card.label }}
            </CardTitle>
            <component :is="card.icon" class="h-4 w-4 text-muted-foreground hidden sm:block" />
          </CardHeader>
          <CardContent class="p-3 pt-0 lg:p-6 lg:pt-0">
            <p class="text-2xl lg:text-3xl font-bold">
              {{ card.value }}
            </p>
            <p
              v-if="card.change !== null"
              class="text-xs mt-1 flex items-center gap-0.5"
              :class="card.change >= 0 ? 'text-green-500' : 'text-red-500'"
            >
              <ArrowUp v-if="card.change >= 0" class="h-3 w-3" />
              <ArrowDown v-else class="h-3 w-3" />
              {{ Math.abs(card.change) }}% vs прошлая неделя
            </p>
          </CardContent>
        </Card>
      </div>

      <!-- Активность по чатам -->
      <div v-if="stats && stats.chatStats.length > 0">
        <Typography variant="h4" as="h2" class="mb-3">
          По чатам (за неделю)
        </Typography>
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
          <Card v-for="chat in stats.chatStats" :key="chat.chatId" data-reveal>
            <CardHeader class="p-3 lg:p-4 pb-1">
              <CardTitle class="text-xs lg:text-sm font-medium truncate">
                {{ chat.title }}
              </CardTitle>
            </CardHeader>
            <CardContent class="p-3 pt-0 lg:p-4 lg:pt-0">
              <p class="text-xl lg:text-2xl font-bold">
                {{ chat.count }}
              </p>
              <p class="text-xs text-muted-foreground">
                сообщений
              </p>
            </CardContent>
          </Card>
        </div>
      </div>

      <!-- График активности -->
      <Card data-reveal>
        <CardHeader class="p-3 lg:p-6">
          <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
            <CardTitle class="text-sm lg:text-base">
              График активности
            </CardTitle>
            <div class="flex gap-2 flex-wrap">
              <select
                v-model="selectedChatId"
                class="text-sm border rounded-md px-2 py-1 bg-background"
              >
                <option :value="undefined">
                  Все чаты
                </option>
                <option v-for="chat in chats" :key="chat.chatId" :value="chat.chatId">
                  {{ chat.title }}
                </option>
              </select>
              <select
                v-model="selectedDays"
                class="text-sm border rounded-md px-2 py-1 bg-background"
              >
                <option :value="7">
                  7 дней
                </option>
                <option :value="14">
                  14 дней
                </option>
                <option :value="30">
                  30 дней
                </option>
              </select>
              <Button size="sm" variant="outline" @click="handleExport">
                <Download class="h-4 w-4 mr-1" />
                CSV
              </Button>
            </div>
          </div>
        </CardHeader>
        <CardContent class="p-3 pt-0 lg:p-6 lg:pt-0">
          <div class="h-48 lg:h-64">
            <Line :data="lineChartData" :options="chartOptions" />
          </div>
        </CardContent>
      </Card>

      <!-- Топ пользователей -->
      <Card v-if="topUsers.length > 0" data-reveal>
        <CardHeader class="p-3 lg:p-6">
          <CardTitle class="text-sm lg:text-base">
            Топ-5 активных пользователей (за неделю)
          </CardTitle>
        </CardHeader>
        <CardContent class="p-3 pt-0 lg:p-6 lg:pt-0">
          <div class="overflow-x-auto">
            <table class="w-full text-sm">
              <thead>
                <tr class="border-b">
                  <th class="text-left py-2 font-medium">
                    #
                  </th>
                  <th class="text-left py-2 font-medium">
                    Пользователь
                  </th>
                  <th class="text-left py-2 font-medium hidden sm:table-cell">
                    Любимый чат
                  </th>
                  <th class="text-right py-2 font-medium">
                    Сообщений
                  </th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="(user, index) in topUsers"
                  :key="user.telegramUserId"
                  class="border-b last:border-0 cursor-pointer hover:bg-muted/50"
                  @click="selectUser(user)"
                >
                  <td class="py-2">
                    {{ index + 1 }}
                  </td>
                  <td class="py-2">
                    <span v-if="user.telegramUsername" class="font-medium">@{{ user.telegramUsername }}</span>
                    <span v-else class="font-medium">{{ user.telegramFirstName }}</span>
                  </td>
                  <td class="py-2 text-muted-foreground hidden sm:table-cell truncate max-w-48">
                    {{ user.topChat }}
                  </td>
                  <td class="py-2 text-right font-bold">
                    {{ user.count }}
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </CardContent>
      </Card>

      <!-- Статистика пользователя -->
      <Card v-if="selectedUser" data-reveal>
        <CardHeader class="p-3 lg:p-6">
          <div class="flex items-center justify-between">
            <CardTitle class="text-sm lg:text-base">
              Статистика: {{ selectedUser.telegramUsername ? `@${selectedUser.telegramUsername}` : selectedUser.telegramFirstName }}
            </CardTitle>
            <button class="p-1 rounded-md hover:bg-muted" @click="clearUserSelection">
              <X class="h-4 w-4" />
            </button>
          </div>
        </CardHeader>
        <CardContent class="p-3 pt-0 lg:p-6 lg:pt-0 space-y-4">
          <div v-if="userStats" class="grid grid-cols-3 gap-3">
            <div class="text-center p-3 bg-muted/50 rounded-lg">
              <p class="text-2xl font-bold">
                {{ userStats.totalMessages }}
              </p>
              <p class="text-xs text-muted-foreground">
                Сообщений
              </p>
            </div>
            <div class="text-center p-3 bg-muted/50 rounded-lg">
              <p class="text-2xl font-bold">
                {{ userStats.activeChats }}
              </p>
              <p class="text-xs text-muted-foreground">
                Чатов
              </p>
            </div>
            <div class="text-center p-3 bg-muted/50 rounded-lg">
              <p class="text-2xl font-bold">
                {{ userStats.avgPerDay }}
              </p>
              <p class="text-xs text-muted-foreground">
                Среднее/день
              </p>
            </div>
          </div>
          <div v-if="userChartData.length > 0" class="h-48">
            <Line :data="userLineChartData" :options="chartOptions" />
          </div>
        </CardContent>
      </Card>
    </div>
  </AdminLayout>
</template>
