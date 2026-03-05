import { ref } from 'vue'
import { useToast } from '@/components/ui/toast'
import api from '@/lib/api'
import { handleError } from '@/services/errorService'

export interface Season {
  id: number
  title: string
  startDate: string
  endDate: string
  status: 'ACTIVE' | 'FINISHED'
  createdAt: string
}

export interface SeasonCreateRequest {
  title: string
  startDate: string
  endDate: string
}

class SeasonService {
  public isLoading = ref(false)
  public items = ref<Season[]>([])

  private toast = useToast()

  getAll = async () => {
    try {
      this.isLoading.value = true
      const response = await api.get('seasons').json<Season[]>()
      this.items.value = response ?? []
    }
    catch (error) {
      handleError(error)
    }
    finally {
      this.isLoading.value = false
    }
  }

  create = async (data: SeasonCreateRequest): Promise<boolean> => {
    try {
      this.isLoading.value = true
      await api.post('seasons', { json: data }).json()

      this.toast.toast({
        title: 'Успешно',
        description: 'Сезон создан',
      })

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

  finish = async (id: number): Promise<boolean> => {
    try {
      this.isLoading.value = true
      await api.post(`seasons/${id}/finish`).json()

      this.toast.toast({
        title: 'Успешно',
        description: 'Сезон завершён',
      })

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

export const seasonService = new SeasonService()
