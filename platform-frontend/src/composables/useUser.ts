import type { RemovableRef } from '@vueuse/core'
import type { Mentor, TelegramUser } from '@/models/profile'
import { useLocalStorage } from '@vueuse/core'
import { computed } from 'vue'
import { getSubscriptionLevel, getSubscriptionLevelIndex, SUBSCRIPTION_LEVELS } from '@/models/profile'

// Версионируем ключ, чтобы при добавлении новых полей в TelegramUser/Mentor
// (например subscriptionTier) старый кэш не блокировал свежие данные.
const TG_USER_KEY = 'tg_user:v2'
// TTL делаем равным интервалу рефреша токена с запасом — после часа простоя
// клиент в любом случае дёргает /me и обновит кэш.
const TG_USER_TTL_MS = 60 * 60 * 1000

let userRef: RemovableRef<null | TelegramUser> | null = null

export function useUser<TUser extends TelegramUser | Mentor = TelegramUser>(): RemovableRef<null | TUser> {
  if (!userRef) {
    userRef = useLocalStorage<TUser>(TG_USER_KEY, null, {
      serializer: {
        read: (raw: string) => {
          try {
            const parsed = JSON.parse(raw)
            if (!parsed || typeof parsed !== 'object')
              return null
            const { data, savedAt } = parsed as { data: unknown, savedAt: number }
            if (typeof savedAt !== 'number' || Date.now() - savedAt > TG_USER_TTL_MS)
              return null
            return data as TUser
          }
          catch {
            return null
          }
        },
        write: value => JSON.stringify({ data: value, savedAt: Date.now() }),
      },
    })
  }
  return userRef as RemovableRef<null | TUser>
}

export function isUserSubscribed() {
  const user = useUser()

  return computed(() => !user.value?.roles.includes('UNSUBSCRIBER'))
}

export function isUserMentor() {
  return computed(() => isMentor(useUser().value))
}

export function isUserAdmin() {
  const user = useUser()
  return computed(() => user.value?.roles?.includes('ADMIN') ?? false)
}

export function canViewAdminPanel() {
  const user = useUser()
  return computed(() => {
    const roles = user.value?.roles
    if (!roles)
      return false
    return roles.includes('ADMIN') || roles.includes('EVENT_MAKER')
  })
}

export function useUserLevel() {
  const user = useUser()
  const level = computed(() => user.value ? getSubscriptionLevel(user.value.roles, user.value.subscriptionTier) : SUBSCRIPTION_LEVELS[0])
  const levelIndex = computed(() => user.value ? getSubscriptionLevelIndex(user.value.roles, user.value.subscriptionTier) : 0)
  const maxLevel = SUBSCRIPTION_LEVELS.length - 1

  return { level, levelIndex, maxLevel }
}

function isMentor(user: TelegramUser | Mentor): user is Mentor {
  return user?.roles?.includes('MENTOR')
}
