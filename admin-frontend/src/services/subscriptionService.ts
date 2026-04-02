import type { Registry } from '@/models/registry'
import { ref } from 'vue'
import { useToast } from '@/components/ui/toast'
import api from '@/lib/api'
import { cleanParams } from '@/lib/utils'
import { handleError } from '@/services/errorService'

export interface SubscriptionStats {
  totalUsers: number
  totalChats: number
  anchorChats: number
  contentChats: number
  tiers: SubscriptionTierStat[]
}

export interface SubscriptionTierStat {
  id: number
  slug: string
  name: string
  level: number
  users: number
}

export interface SubscriptionChat {
  id: number
  title: string
  chatType: string
  anchorForTierID?: number
  anchorTierName?: string
  tierIDs?: number[]
  tierNames?: string[]
  activeUsers: number
}

export interface SubscriptionChatDetail {
  id: number
  title: string
  chatType: string
  anchorForTierID?: number
  tierIDs: number[]
}

export interface SubscriptionUser {
  id: number
  username: string | null
  fullName: string
  isActive: boolean
  tierName?: string
  tierSlug?: string
  manualTierID?: number
  resolvedTierID?: number
  activeChats: number
  lastCheckAt?: string
  createdAt: string
}

export interface SubscriptionUserDetail extends SubscriptionUser {
  resolvedTierName?: string
  manualTierName?: string
  effectiveTierName?: string
  access: {
    chatID: number
    chatTitle?: string
    grantedAt: string
  }[]
}

class SubscriptionService {
  public isLoading = ref(false)
  public stats = ref<SubscriptionStats | null>(null)
  public users = ref<Registry<SubscriptionUser>>({ items: [], total: 0 })
  public chats = ref<SubscriptionChat[]>([])
  public tiers = ref<SubscriptionTierStat[]>([])
  public pagination = ref({ limit: 20, offset: 0 })

  private toast = useToast()

  changePagination = (page: number) => {
    this.pagination.value.offset = (page - 1) * this.pagination.value.limit
    this.searchUsers()
  }

  clearPagination = () => {
    this.pagination.value.offset = 0
  }

  fetchStats = async () => {
    try {
      this.isLoading.value = true
      this.stats.value = await api.get('subscriptions/stats').json<SubscriptionStats>()
    }
    catch (error) {
      handleError(error)
    }
    finally {
      this.isLoading.value = false
    }
  }

  fetchTiers = async () => {
    try {
      const response = await api.get('subscriptions/tiers').json<Registry<SubscriptionTierStat>>()
      this.tiers.value = response.items
    }
    catch (error) {
      handleError(error)
    }
  }

  fetchChats = async () => {
    try {
      const response = await api.get('subscriptions/chats').json<Registry<SubscriptionChat>>()
      this.chats.value = response.items
    }
    catch (error) {
      handleError(error)
    }
  }

  searchUsers = async () => {
    try {
      this.isLoading.value = true
      const searchParams = cleanParams({
        limit: this.pagination.value.limit,
        offset: this.pagination.value.offset,
      })
      this.users.value = await api.get('subscriptions/users', { searchParams }).json<Registry<SubscriptionUser>>()
    }
    catch (error) {
      handleError(error)
    }
    finally {
      this.isLoading.value = false
    }
  }

  getUser = async (id: number): Promise<SubscriptionUserDetail | null> => {
    try {
      return await api.get(`subscriptions/users/${id}`).json<SubscriptionUserDetail>()
    }
    catch (error) {
      handleError(error)
      return null
    }
  }

  setOverride = async (userId: number, tierSlug: string): Promise<boolean> => {
    try {
      await api.put(`subscriptions/users/${userId}/override`, { json: { tierSlug } }).json()
      this.toast.toast({ title: 'Успешно', description: 'Тир установлен' })
      return true
    }
    catch (error) {
      handleError(error)
      return false
    }
  }

  clearOverride = async (userId: number): Promise<boolean> => {
    try {
      await api.delete(`subscriptions/users/${userId}/override`).json()
      this.toast.toast({ title: 'Успешно', description: 'Ручной тир снят' })
      return true
    }
    catch (error) {
      handleError(error)
      return false
    }
  }

  getChatDetail = async (chatId: number): Promise<SubscriptionChatDetail | null> => {
    try {
      return await api.get(`subscriptions/chats/${chatId}`).json<SubscriptionChatDetail>()
    }
    catch (error) {
      handleError(error)
      return null
    }
  }

  createChat = async (data: { id: number, title: string, chatType: string, anchorForTierID?: number, tierIDs?: number[] }): Promise<boolean> => {
    try {
      await api.post('subscriptions/chats', { json: data }).json()
      this.toast.toast({ title: 'Успешно', description: 'Чат добавлен' })
      return true
    }
    catch (error) {
      handleError(error)
      return false
    }
  }

  updateChat = async (chatId: number, data: { title?: string, anchorForTierID?: number, clearAnchor?: boolean, tierIDs?: number[] }): Promise<boolean> => {
    try {
      await api.put(`subscriptions/chats/${chatId}`, { json: data }).json()
      this.toast.toast({ title: 'Успешно', description: 'Чат обновлён' })
      return true
    }
    catch (error) {
      handleError(error)
      return false
    }
  }

  deleteChat = async (chatId: number): Promise<boolean> => {
    try {
      await api.delete(`subscriptions/chats/${chatId}`).json()
      this.toast.toast({ title: 'Успешно', description: 'Чат удалён' })
      return true
    }
    catch (error) {
      handleError(error)
      return false
    }
  }

  revokeAccess = async (userId: number, chatId: number): Promise<boolean> => {
    try {
      await api.delete(`subscriptions/users/${userId}/access/${chatId}`).json()
      this.toast.toast({ title: 'Успешно', description: 'Доступ отозван' })
      return true
    }
    catch (error) {
      handleError(error)
      return false
    }
  }
}

export const subscriptionService = new SubscriptionService()
