import { ref } from 'vue'
import { useToast } from '@/components/ui/toast'
import api from '@/lib/api'
import { handleError } from '@/services/errorService'

export interface Raffle {
  id: number
  title: string
  description: string
  prize: string
  ticketCost: number
  maxTickets: number
  endsAt: string
  status: 'ACTIVE' | 'FINISHED'
  winnerId?: number | null
  createdAt: string
}

export interface RaffleCreateRequest {
  title: string
  description: string
  prize: string
  ticketCost: number
  maxTickets: number
  endsAt: string
}

class RaffleService {
  public isLoading = ref(false)
  public items = ref<Raffle[]>([])

  private toast = useToast()

  getAll = async () => {
    try {
      this.isLoading.value = true
      const response = await api.get('raffles').json<Raffle[]>()
      this.items.value = response ?? []
    }
    catch (error) {
      handleError(error)
    }
    finally {
      this.isLoading.value = false
    }
  }

  create = async (data: RaffleCreateRequest): Promise<boolean> => {
    try {
      this.isLoading.value = true
      await api.post('raffles', { json: data }).json()

      this.toast.toast({
        title: 'Успешно',
        description: 'Розыгрыш создан',
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

  delete = async (id: number): Promise<boolean> => {
    try {
      this.isLoading.value = true
      await api.delete(`raffles/${id}`)

      this.toast.toast({
        title: 'Успешно',
        description: 'Розыгрыш удалён',
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

export const raffleService = new RaffleService()
