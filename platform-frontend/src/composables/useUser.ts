import type { RemovableRef } from '@vueuse/core'
import type { Mentor, TelegramUser } from '@/models/profile'
import { useLocalStorage } from '@vueuse/core'
import { computed } from 'vue'
import { getSubscriptionLevel, getSubscriptionLevelIndex, SUBSCRIPTION_LEVELS } from '@/models/profile'

let userRef: RemovableRef<null | TelegramUser> | null = null

export function useUser<TUser extends TelegramUser | Mentor = TelegramUser>(): RemovableRef<null | TUser> {
  if (!userRef) {
    userRef = useLocalStorage<TUser>('tg_user', null, {
      serializer: { read: JSON.parse, write: JSON.stringify },
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

export function useUserLevel() {
  const user = useUser()
  const level = computed(() => user.value ? getSubscriptionLevel(user.value.roles) : SUBSCRIPTION_LEVELS[0])
  const levelIndex = computed(() => user.value ? getSubscriptionLevelIndex(user.value.roles) : 0)
  const maxLevel = SUBSCRIPTION_LEVELS.length - 1

  return { level, levelIndex, maxLevel }
}

function isMentor(user: TelegramUser | Mentor): user is Mentor {
  return user?.roles?.includes('MENTOR')
}
