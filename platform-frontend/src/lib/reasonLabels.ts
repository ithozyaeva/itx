export const reasonLabels: Record<string, string> = {
  event_attend: 'Посещение событий',
  event_host: 'Проведение событий',
  review_community: 'Отзывы (сообщество)',
  review_service: 'Отзывы (услуги)',
  resume_upload: 'Загрузка резюме',
  referal_create: 'Создание рефералов',
  referal_conversion: 'Конверсия рефералов',
  profile_complete: 'Заполнение профиля',
  weekly_activity: 'Еженедельная активность',
  monthly_active: 'Ежемесячная активность',
  streak_4weeks: 'Серия 4 недели',
  admin_manual: 'Начисление вручную',
  task_create: 'Создание заданий',
  task_execute: 'Выполнение заданий',
  marketplace_create: 'Публикация объявлений',
  marketplace_buy: 'Покупки',
  chat_quest: 'Квесты в чатах',
  chatter_of_week: 'Чаттер недели',
  kudos_received: 'Благодарности',
  raffle_spend: 'Розыгрыши',
  casino_bet: 'Ставки мини-игр',
  casino_win: 'Выигрыши мини-игр',
}

// «Способы получить баллы» — крупные действия платформы, у которых есть
// постоянная сумма награды. Используются:
//   - на табе «Способы» страницы /progress как карточки-ссылки;
//   - в toast «Задание выполнено! +N» после первого срабатывания.
// До бэкенд-эндпоинта /points/sources (todo) — единый словарь здесь.
export interface PointSource {
  reason: string
  label: string
  shortLabel?: string // короткое название для toast (если нужно)
  to: string
  points: number
}

export const pointSources: PointSource[] = [
  { reason: 'event_attend', label: 'Запишись на событие', to: '/events', points: 10 },
  { reason: 'event_host', label: 'Проведи событие', to: '/events', points: 25 },
  { reason: 'review_community', label: 'Оставь отзыв на сообщество', shortLabel: 'Оставь отзыв', to: '/my-reviews', points: 15 },
  { reason: 'resume_upload', label: 'Загрузи резюме', to: '/resumes', points: 10 },
  { reason: 'profile_complete', label: 'Заполни профиль', to: '/me', points: 20 },
  { reason: 'referal_create', label: 'Создай реферальную ссылку', to: '/referals', points: 5 },
  { reason: 'referal_conversion', label: 'Получи конверсию по рефералу', shortLabel: 'Конверсия реферала', to: '/referals', points: 30 },
]
