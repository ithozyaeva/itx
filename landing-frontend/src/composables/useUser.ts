import type { RemovableRef } from '@vueuse/core'
import type { TelegramUser } from '@/services/auth'
import { useLocalStorage } from '@vueuse/core'

export function useUser(): RemovableRef<null | TelegramUser> {
  return useLocalStorage<TelegramUser>('tg_user', null, { serializer: { read: JSON.parse, write: JSON.stringify } })
}

export function useConfirmedPrivacy(): RemovableRef<boolean> {
  return useLocalStorage<boolean>('confirmed_privacy', false)
}
