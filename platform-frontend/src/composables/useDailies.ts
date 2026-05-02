import type { CheckInResponse, DailyTodayResponse, StreakResponse } from '@/models/dailies'
import { ref } from 'vue'
import { dailiesService } from '@/services/dailies'
import { useSSE } from './useSSE'

const today = ref<DailyTodayResponse | null>(null)
const streak = ref<StreakResponse | null>(null)
const loading = ref(false)
const checkingIn = ref(false)
const error = ref<string | null>(null)

async function fetchToday() {
  try {
    today.value = await dailiesService.getToday()
  }
  catch (e) {
    today.value = null
    error.value = (e as Error).message
  }
}

async function fetchStreak() {
  try {
    streak.value = await dailiesService.getStreak()
  }
  catch (e) {
    streak.value = null
    error.value = (e as Error).message
  }
}

async function refresh() {
  loading.value = true
  await Promise.allSettled([fetchToday(), fetchStreak()])
  loading.value = false
}

async function checkIn(): Promise<CheckInResponse | null> {
  if (checkingIn.value)
    return null
  checkingIn.value = true
  try {
    const resp = await dailiesService.checkIn()
    streak.value = resp.streak
    fetchToday()
    return resp
  }
  catch (e) {
    error.value = (e as Error).message
    return null
  }
  finally {
    checkingIn.value = false
  }
}

export function useDailies() {
  // SSE-подписки на live-обновления баллов/стрика/дейликов.
  // Подписки внутри useSSE сами отвяжутся при unmount компонента.
  useSSE('dailies', fetchToday)
  useSSE('streak', fetchStreak)
  useSSE('points', fetchToday)

  return {
    today,
    streak,
    loading,
    checkingIn,
    error,
    refresh,
    fetchToday,
    fetchStreak,
    checkIn,
  }
}
