import { ref } from 'vue'
import { useToast } from '@/components/ui/toast'
import api from '@/lib/api'
import { handleError } from '@/services/errorService'

export type ChallengeKind = 'weekly' | 'monthly'

export interface ChallengeTemplate {
  id: number
  code: string
  title: string
  description: string
  icon: string
  kind: ChallengeKind
  metricKey: string
  target: number
  rewardPoints: number
  achievementCode: string | null
  active: boolean
  createdAt: string
}

export interface ChallengeTemplateRequest {
  code: string
  title: string
  description: string
  icon: string
  kind: ChallengeKind
  metricKey: string
  target: number
  rewardPoints: number
  achievementCode: string | null
  active: boolean
}

class ChallengeAdminService {
  public isLoading = ref(false)
  public items = ref<ChallengeTemplate[]>([])

  private toast = useToast()

  getAll = async () => {
    try {
      this.isLoading.value = true
      const resp = await api.get('challenges').json<{ items: ChallengeTemplate[] }>()
      this.items.value = resp?.items ?? []
    }
    catch (error) {
      handleError(error)
    }
    finally {
      this.isLoading.value = false
    }
  }

  create = async (data: ChallengeTemplateRequest): Promise<boolean> => {
    try {
      this.isLoading.value = true
      await api.post('challenges', { json: data }).json()
      this.toast.toast({ title: 'Успешно', description: 'Шаблон создан' })
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

  update = async (id: number, data: ChallengeTemplateRequest): Promise<boolean> => {
    try {
      this.isLoading.value = true
      await api.put(`challenges/${id}`, { json: data }).json()
      this.toast.toast({ title: 'Успешно', description: 'Шаблон обновлён' })
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
      await api.delete(`challenges/${id}`)
      this.toast.toast({ title: 'Успешно', description: 'Шаблон удалён' })
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

  recentInstances = async (limit = 30): Promise<ChallengeInstance[]> => {
    try {
      const resp = await api.get('challenges/instances', { searchParams: { limit } }).json<{ items: ChallengeInstance[] }>()
      return resp?.items ?? []
    }
    catch (error) {
      handleError(error)
      return []
    }
  }
}

export interface ChallengeInstance {
  id: number
  templateId: number
  kind: ChallengeKind
  startsAt: string
  endsAt: string
  periodKey: string
  createdAt: string
}

export const challengeAdminService = new ChallengeAdminService()
