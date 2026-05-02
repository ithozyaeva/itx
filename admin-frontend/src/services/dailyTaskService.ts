import { ref } from 'vue'
import { useToast } from '@/components/ui/toast'
import api from '@/lib/api'
import { handleError } from '@/services/errorService'

export type DailyTaskTier = 'engagement' | 'light' | 'meaningful' | 'big'

export interface DailyTask {
  id: number
  code: string
  title: string
  description: string
  icon: string
  tier: DailyTaskTier
  points: number
  target: number
  triggerKey: string
  active: boolean
  createdAt: string
}

export interface DailyTaskCreateRequest {
  code: string
  title: string
  description: string
  icon: string
  tier: DailyTaskTier
  points: number
  target: number
  triggerKey: string
  active: boolean
}

class DailyTaskService {
  public isLoading = ref(false)
  public items = ref<DailyTask[]>([])

  private toast = useToast()

  getAll = async () => {
    try {
      this.isLoading.value = true
      const resp = await api.get('daily-tasks').json<{ items: DailyTask[] }>()
      this.items.value = resp?.items ?? []
    }
    catch (error) {
      handleError(error)
    }
    finally {
      this.isLoading.value = false
    }
  }

  create = async (data: DailyTaskCreateRequest): Promise<boolean> => {
    try {
      this.isLoading.value = true
      await api.post('daily-tasks', { json: data }).json()
      this.toast.toast({ title: 'Успешно', description: 'Задание создано' })
      await this.getAll()
      return true
    }
    catch (error) {
      handleError(error)
      return false
    }
    finally {
      this.isLoading.value = false
    }
  }

  update = async (id: number, data: DailyTaskCreateRequest): Promise<boolean> => {
    try {
      this.isLoading.value = true
      await api.put(`daily-tasks/${id}`, { json: data }).json()
      this.toast.toast({ title: 'Успешно', description: 'Задание обновлено' })
      await this.getAll()
      return true
    }
    catch (error) {
      handleError(error)
      return false
    }
    finally {
      this.isLoading.value = false
    }
  }

  delete = async (id: number): Promise<boolean> => {
    try {
      this.isLoading.value = true
      await api.delete(`daily-tasks/${id}`)
      this.toast.toast({ title: 'Успешно', description: 'Задание удалено' })
      await this.getAll()
      return true
    }
    catch (error) {
      handleError(error)
      return false
    }
    finally {
      this.isLoading.value = false
    }
  }
}

export const dailyTaskService = new DailyTaskService()
