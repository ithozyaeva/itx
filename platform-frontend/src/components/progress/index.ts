export { default as DailyCheckInWidget } from './DailyCheckInWidget.vue'
export { default as Deadline } from './Deadline.vue'
export { default as PointsBadge } from './PointsBadge.vue'
export { default as ProgressBar } from './ProgressBar.vue'
export { default as TaskCard } from './TaskCard.vue'
export { default as TaskCardSkeleton } from './TaskCardSkeleton.vue'
export type { IconTone } from './TintedIcon.vue'
export { default as TintedIcon } from './TintedIcon.vue'

// Panel-компоненты (Leaderboard/Achievements/MyStats/Kudos) намеренно
// НЕ экспортируются здесь: их следует импортировать напрямую через
// defineAsyncComponent(() => import('@/components/progress/XxxPanel.vue')),
// чтобы тяжёлые табы /progress грузились лениво. Импорт через barrel
// сделает их статической зависимостью и попадёт в основной chunk.
