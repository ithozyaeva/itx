import type { AdminAwardRequest, AdminPointTransaction } from '@/models/points'
import type { Registry } from '@/models/registry'
import { ref } from 'vue'
import { useToast } from '@/components/ui/toast'
import api from '@/lib/api'
import { cleanParams } from '@/lib/utils'
import { handleError } from '@/services/errorService'

export interface PointsFilters {
  memberId?: number
}

class PointsService {
  public isLoading = ref(false)
  public items = ref<Registry<AdminPointTransaction>>({ items: [], total: 0 })
  public pagination = ref({ limit: 20, offset: 0 })
  public filters = ref<PointsFilters>({})

  private toast = useToast()

  changePagination = (page: number) => {
    this.pagination.value.offset = (page - 1) * this.pagination.value.limit
    this.search()
  }

  clearPagination = () => {
    this.pagination.value.offset = 0
  }

  search = async () => {
    try {
      this.isLoading.value = true
      const searchParams = cleanParams({
        limit: this.pagination.value.limit,
        offset: this.pagination.value.offset,
        ...this.filters.value,
      })
      const response = await api.get('points', { searchParams }).json<Registry<AdminPointTransaction>>()
      this.items.value = response
    }
    catch (error) {
      handleError(error)
    }
    finally {
      this.isLoading.value = false
    }
  }

  applyFilters = (filters: PointsFilters) => {
    this.filters.value = filters
    this.pagination.value.offset = 0
    this.search()
  }

  award = async (data: AdminAwardRequest): Promise<boolean> => {
    try {
      this.isLoading.value = true
      await api.post('points', { json: data }).json()

      this.toast.toast({
        title: 'Успешно',
        description: 'Баллы успешно начислены',
      })

      await this.search()
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

  deleteTransaction = async (id: number): Promise<boolean> => {
    try {
      this.isLoading.value = true
      await api.delete(`points/${id}`).json()

      this.toast.toast({
        title: 'Успешно',
        description: 'Транзакция удалена',
      })

      await this.search()
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

export const pointsService = new PointsService()
