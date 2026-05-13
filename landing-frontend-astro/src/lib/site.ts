export const SITE = {
  name: 'IT-ХОЗЯЕВА',
  shortName: 'IT-ХОЗЯЕВА',
  baseDescription: 'ИТ-сообщество 250+ специалистов на стыке технологий и ИИ. Менторство, вайбкодинг, нетворкинг, собеседования. Подписка от 520 ₽/мес.',
  url: 'https://ithozyaeva.ru',
  ogImage: 'https://ithozyaeva.ru/og-image.png',
  locale: 'ru_RU',
  twitterHandle: '@ithozyaeva',
} as const

export const TARIFFS = {
  brigadir: {
    name: 'Бригадир',
    slug: 'brigadir',
    price: 520,
    boostyUrl: 'https://boosty.to/jointime/purchase/3150816',
    description: 'Доступ к материалам, ИТ-чаты, вебинары, вакансии и реферальные ссылки',
  },
  hozyain: {
    name: 'ХОЗЯИН',
    slug: 'hozyain',
    price: 1000,
    boostyUrl: 'https://boosty.to/jointime/purchase/3150814',
    description: 'Все возможности Бригадира + приоритетная поддержка, разбор резюме, доступ к базе менторов',
  },
  master: {
    name: 'МАСТЕР',
    slug: 'master',
    price: 5200,
    boostyUrl: 'https://boosty.to/jointime/purchase/967625',
    description: 'Все возможности ХОЗЯИНа + реклама ресурсов, верхняя позиция в таблице менторов, личная консультация',
  },
} as const
