import type { GlobalBan, ModerationAction, VotebanRow } from '@/models/moderation'
import { ref } from 'vue'
import { useToast } from '@/components/ui/toast'
import api from '@/lib/api'
import { handleError } from '@/services/errorService'

interface List<T> { items: T[], total: number }

class ModerationService {
  public isLoading = ref(false)
  public sanctions = ref<ModerationAction[]>([])
  public actions = ref<ModerationAction[]>([])
  public globalBans = ref<GlobalBan[]>([])
  public votebans = ref<VotebanRow[]>([])

  private toast = useToast()

  fetchSanctions = async () => {
    try {
      this.isLoading.value = true
      const res = await api.get('moderation/sanctions').json<List<ModerationAction>>()
      this.sanctions.value = res.items
    }
    catch (e) {
      handleError(e)
    }
    finally {
      this.isLoading.value = false
    }
  }

  fetchActions = async () => {
    try {
      this.isLoading.value = true
      const res = await api.get('moderation/actions').json<List<ModerationAction>>()
      this.actions.value = res.items
    }
    catch (e) {
      handleError(e)
    }
    finally {
      this.isLoading.value = false
    }
  }

  fetchGlobalBans = async () => {
    try {
      this.isLoading.value = true
      const res = await api.get('moderation/global-bans').json<List<GlobalBan>>()
      this.globalBans.value = res.items
    }
    catch (e) {
      handleError(e)
    }
    finally {
      this.isLoading.value = false
    }
  }

  fetchVotebans = async () => {
    try {
      this.isLoading.value = true
      const res = await api.get('moderation/votebans').json<List<VotebanRow>>()
      this.votebans.value = res.items
    }
    catch (e) {
      handleError(e)
    }
    finally {
      this.isLoading.value = false
    }
  }

  revokeSanction = async (actionId: number) => {
    try {
      await api.post(`moderation/sanctions/${actionId}/revoke`)
      this.toast.toast({
        title: 'Команда отправлена',
        description: 'Бот снимет санкцию в Telegram в течение нескольких секунд.',
      })
      await this.fetchSanctions()
      return true
    }
    catch (e) {
      handleError(e)
      return false
    }
  }

  revokeGlobalBan = async (userId: number) => {
    try {
      await api.delete(`moderation/global-bans/${userId}`)
      this.toast.toast({
        title: 'Команда отправлена',
        description: 'Бот снимает глобальный бан и разблокирует во всех чатах.',
      })
      await this.fetchGlobalBans()
      return true
    }
    catch (e) {
      handleError(e)
      return false
    }
  }

  cancelVoteban = async (votebanId: number) => {
    try {
      await api.post(`moderation/votebans/${votebanId}/cancel`)
      this.toast.toast({
        title: 'Голосование отменено',
        description: 'Бот обновит сообщение в чате.',
      })
      await this.fetchVotebans()
      return true
    }
    catch (e) {
      handleError(e)
      return false
    }
  }
}

export const moderationService = new ModerationService()
