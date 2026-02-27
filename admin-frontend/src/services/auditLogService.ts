import type { Registry } from '@/models/registry'
import { ref } from 'vue'
import api from '@/lib/api'
import { cleanParams } from '@/lib/utils'
import { handleError } from '@/services/errorService'

export interface AuditLog {
  id: number
  actorId: number
  actorName: string
  actorType: 'admin' | 'platform'
  action: 'create' | 'update' | 'delete' | 'approve'
  entityType: string
  entityId: number
  entityName: string
  createdAt: string
}

export interface AuditLogFilters {
  actorType?: string
  action?: string
  entityType?: string
}

class AuditLogService {
  public isLoading = ref(false)
  public items = ref<Registry<AuditLog>>({ items: [], total: 0 })
  public pagination = ref({ limit: 20, offset: 0 })
  public filters = ref<AuditLogFilters>({})

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
      const response = await api.get('audit-logs', { searchParams }).json<Registry<AuditLog>>()
      this.items.value = response
    }
    catch (error) {
      handleError(error)
    }
    finally {
      this.isLoading.value = false
    }
  }

  applyFilters = (filters: AuditLogFilters) => {
    this.filters.value = filters
    this.pagination.value.offset = 0
    this.search()
  }
}

export const auditLogService = new AuditLogService()
