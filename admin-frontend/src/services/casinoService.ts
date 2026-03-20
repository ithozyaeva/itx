import type { CasinoAdminStats, CasinoBet } from '@/models/casino'
import type { Registry } from '@/models/registry'
import { ref } from 'vue'
import api from '@/lib/api'
import { cleanParams } from '@/lib/utils'
import { handleError } from '@/services/errorService'

export interface CasinoFilters {
  username?: string
  game?: string
}

class CasinoService {
  public isLoading = ref(false)
  public items = ref<Registry<CasinoBet>>({ items: [], total: 0 })
  public stats = ref<CasinoAdminStats | null>(null)
  public pagination = ref({ limit: 20, offset: 0 })
  public filters = ref<CasinoFilters>({})

  changePagination = (page: number) => {
    this.pagination.value.offset = (page - 1) * this.pagination.value.limit
    this.searchBets()
  }

  clearPagination = () => {
    this.pagination.value.offset = 0
  }

  getStats = async () => {
    try {
      const response = await api.get('minigames/stats').json<CasinoAdminStats>()
      this.stats.value = response
    }
    catch (error) {
      handleError(error)
    }
  }

  searchBets = async () => {
    try {
      this.isLoading.value = true
      const searchParams = cleanParams({
        limit: this.pagination.value.limit,
        offset: this.pagination.value.offset,
        ...this.filters.value,
      })
      const response = await api.get('minigames/bets', { searchParams }).json<Registry<CasinoBet>>()
      this.items.value = response
    }
    catch (error) {
      handleError(error)
    }
    finally {
      this.isLoading.value = false
    }
  }

  applyFilters = (filters: CasinoFilters) => {
    this.filters.value = filters
    this.pagination.value.offset = 0
    this.searchBets()
  }
}

export const casinoService = new CasinoService()
