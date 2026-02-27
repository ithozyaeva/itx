import type { Mentor, TelegramUser } from '@/models/profile'
import type { RemovableRef } from '@vueuse/core'
import { useLocalStorage } from '@vueuse/core'
import { computed } from 'vue'

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

function isMentor(user: TelegramUser | Mentor): user is Mentor {
  return user?.roles?.includes('MENTOR')
}
