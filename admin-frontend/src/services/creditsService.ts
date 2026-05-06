import type { AdminAwardCreditsRequest, AdminCreditTransaction } from '@/models/credits'
import type { Registry } from '@/models/registry'
import { ref } from 'vue'
import { useToast } from '@/components/ui/toast'
import api from '@/lib/api'
import { cleanParams } from '@/lib/utils'
import { handleError } from '@/services/errorService'

export interface CreditsFilters {
  username?: string
}

class CreditsService {
  public isLoading = ref(false)
  public items = ref<Registry<AdminCreditTransaction>>({ items: [], total: 0 })
  public pagination = ref({ limit: 20, offset: 0 })
  public filters = ref<CreditsFilters>({})

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
      const response = await api.get('credits', { searchParams }).json<Registry<AdminCreditTransaction>>()
      this.items.value = response
    }
    catch (error) {
      handleError(error)
    }
    finally {
      this.isLoading.value = false
    }
  }

  applyFilters = (filters: CreditsFilters) => {
    this.filters.value = filters
    this.pagination.value.offset = 0
    this.search()
  }

  award = async (data: AdminAwardCreditsRequest): Promise<boolean> => {
    try {
      this.isLoading.value = true
      await api.post('credits', { json: data }).json()
      this.toast.toast({
        title: 'Успешно',
        description: 'Кредиты начислены',
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

export const creditsService = new CreditsService()
