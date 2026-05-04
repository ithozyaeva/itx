import type { IconTone } from '@/components/progress/TintedIcon.vue'
import type { DailyTaskTier } from '@/models/dailies'

export interface TierBadgeMeta {
  label: string
  tone: IconTone
}

// Лейблы и тон бейджа для tier'а ежедневной задачи. Используется в TaskCard
// и в любом будущем месте, где надо отрисовать tier — словарь один на всё
// приложение.
export const tierBadge: Record<DailyTaskTier, TierBadgeMeta> = {
  engagement: { label: 'Просмотр', tone: 'blue' },
  light: { label: 'Действие', tone: 'green' },
  meaningful: { label: 'Контент', tone: 'orange' },
  big: { label: 'Серьёзное', tone: 'purple' },
}
