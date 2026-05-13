import type { ChallengesResponse } from '@/models/challenges'
import { ref } from 'vue'
import { challengesService } from '@/services/challenges'
import { useSSE } from './useSSE'

const data = ref<ChallengesResponse | null>(null)
const loading = ref(false)
const error = ref<string | null>(null)

async function fetchAll() {
  loading.value = true
  // Сбрасываем прошлую ошибку в начале попытки: иначе после первого
  // network glitch error.value остаётся выставленной навсегда — даже
  // когда SSE-триггер перетянул свежие данные, UI продолжает показывать
  // «не удалось загрузить» рядом с актуальным data.
  error.value = null
  try {
    data.value = await challengesService.getMine()
  }
  catch (e) {
    data.value = null
    error.value = (e as Error).message
  }
  finally {
    loading.value = false
  }
}

export function useChallenges() {
  useSSE('challenges', fetchAll)
  useSSE('points', fetchAll)
  return {
    data,
    loading,
    error,
    fetchAll,
  }
}
