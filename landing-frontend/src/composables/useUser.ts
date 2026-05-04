import type { RemovableRef } from '@vueuse/core'
import type { TelegramUser } from '@/services/auth'
import { useLocalStorage } from '@vueuse/core'

// Ключ и сериализатор должны быть синхронны с platform-frontend/src/composables/useUser.ts:
// иначе состояние логина не разделяется между лендингом и /platform, и юзер,
// зашедший с лендинга, бьётся об пустой `tg_user` на платформе и улетает обратно.
const TG_USER_KEY = 'tg_user:v2'
const TG_USER_TTL_MS = 60 * 60 * 1000

export function useUser(): RemovableRef<null | TelegramUser> {
  return useLocalStorage<TelegramUser>(TG_USER_KEY, null, {
    serializer: {
      read: (raw: string) => {
        try {
          const parsed = JSON.parse(raw)
          if (!parsed || typeof parsed !== 'object')
            return null as unknown as TelegramUser
          const { data, savedAt } = parsed as { data: unknown, savedAt: number }
          if (typeof savedAt !== 'number' || Date.now() - savedAt > TG_USER_TTL_MS)
            return null as unknown as TelegramUser
          return data as TelegramUser
        }
        catch {
          return null as unknown as TelegramUser
        }
      },
      write: value => JSON.stringify({ data: value, savedAt: Date.now() }),
    },
  })
}

export function useConfirmedPrivacy(): RemovableRef<boolean> {
  return useLocalStorage<boolean>('confirmed_privacy', false)
}
