import type { RemovableRef } from '@vueuse/core'
import type { Mentor, SubscriptionTierSlug, TelegramUser } from '@/models/profile'
import { useLocalStorage } from '@vueuse/core'
import { computed } from 'vue'
import { getSubscriptionLevel, getSubscriptionLevelIndex, SUBSCRIPTION_LEVELS } from '@/models/profile'

// Минимальный level для каждого slug — должен совпадать с миграцией
// 20260319000000_create_subscription_system + 20260424150000_add_king_tier:
// beginner=1, foreman=2, master=3, king=4. Используется для фронтенд-гейтинга
// разделов с meta.requiresMinTier (источник правды — backend RequireMinTier).
const TIER_SLUG_LEVELS: Record<SubscriptionTierSlug, number> = {
  beginner: 1,
  foreman: 2,
  master: 3,
  king: 4,
}

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

  // Источник правды — наличие effective tier (manual override или resolved
  // через anchor-чат). Старый признак roles.includes('UNSUBSCRIBER') не
  // различал «в main chat без оплаченного тира», поэтому переключились
  // на tier.id напрямую — синхронно с backend RequireSubscription.
  return computed(() => user.value?.subscriptionTier?.id != null)
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

// hasMinTier — гейт по минимальному уровню подписки (mirror к backend
// AuthMiddleware.RequireMinTier). ADMIN всегда проходит — синхронно с
// бэкенд-логикой ролей. Возвращает Ref<boolean>, чтобы использовать
// в реактивных guard'ах (роутер, сайдбар).
export function hasMinTier(slug: SubscriptionTierSlug) {
  const user = useUser()
  return computed(() => {
    if (user.value?.roles?.includes('ADMIN'))
      return true
    const tierLevel = user.value?.subscriptionTier?.level
    if (tierLevel == null)
      return false
    return tierLevel >= TIER_SLUG_LEVELS[slug]
  })
}

function isMentor(user: TelegramUser | Mentor): user is Mentor {
  return user?.roles?.includes('MENTOR')
}
