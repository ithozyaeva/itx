// Build-time fetch менторов. Поверх — статичная карточная сетка с фильтрацией
// по тегам через классы и data-атрибуты (без рантайм-фреймворков).

import process from 'node:process'

export interface ProfTag {
  id: number
  title: string
}

export interface Contact {
  id: number
  type: number
  link: string
  ownerId: number
}

export interface Service {
  id: number
  name: string
  price: number
  ownerId: number
}

export interface MentorRaw {
  id: number
  telegramID: number
  tg: string
  firstName: string
  lastName: string
  bio: string
  avatarUrl?: string
  occupation: string
  experience: string
  profTags: ProfTag[]
  contacts: Contact[]
  services?: Service[]
  order?: number
}

export interface MentorView {
  id: number
  avatar: string
  name: string
  position: string
  description: string
  tags: string[]
  tagsLower: string[]
  link: string
  initials: string
}

const API_BASE = (import.meta.env.PUBLIC_API_BASE as string | undefined) ?? 'https://ithozyaeva.ru'

// На прод-деплое (CI=true) пустой ответ от /api/mentors означает, что страница
// уйдёт без секции менторов — это catastrophic для главной. Падаем сборкой,
// чтобы алертнуть деплой-workflow. В dev — graceful fallback (можно билдить
// без сети).
function failBuildIfEmptyInCi(name: string, count: number): void {
  if (count === 0 && process.env.CI === 'true') {
    throw new Error(`[${name}] API returned 0 items during CI build — refusing to deploy with empty section. Check API_BASE=${API_BASE} reachability.`)
  }
}

// Приоритет контактов: Telegram(1) > LinkedIn(5) > GitHub(6) > Email(2) > VK(7) > Сайт(8) > Телефон(3) > Другое(4)
const CONTACT_PRIORITY = [1, 5, 6, 2, 7, 8, 3, 4]

function getBestContactLink(contacts: Contact[] | undefined): string {
  if (!contacts || contacts.length === 0)
    return '#'
  for (const type of CONTACT_PRIORITY) {
    const c = contacts.find(x => x.type === type)
    if (c?.link)
      return c.link
  }
  return contacts[0]?.link ?? '#'
}

function deterministicShuffle<T>(arr: T[], seed: number): T[] {
  // Mulberry32 seeded RNG для воспроизводимой перетасовки при каждом билде.
  let s = seed >>> 0
  function rand() {
    s = (s + 0x6D2B79F5) >>> 0
    let r = s
    r = Math.imul(r ^ (r >>> 15), r | 1)
    r ^= r + Math.imul(r ^ (r >>> 7), r | 61)
    return ((r ^ (r >>> 14)) >>> 0) / 4294967296
  }
  const out = [...arr]
  for (let i = out.length - 1; i > 0; i--) {
    const j = Math.floor(rand() * (i + 1))
    ;[out[i], out[j]] = [out[j], out[i]]
  }
  return out
}

function toView(m: MentorRaw): MentorView {
  const tags = m.profTags?.map(t => t.title.trim()) ?? []
  const initials = `${m.firstName?.[0] ?? ''}${m.lastName?.[0] ?? ''}`.toUpperCase() || '?'
  return {
    id: m.id,
    avatar: m.avatarUrl?.trim() || `https://t.me/i/userpic/160/${m.tg}.jpg`,
    name: `${m.firstName} ${m.lastName}`.trim(),
    position: m.occupation?.trim() || 'Не указано',
    description: m.experience?.trim() ?? '',
    tags,
    tagsLower: tags.map(t => t.toLowerCase()),
    link: getBestContactLink(m.contacts),
    initials,
  }
}

export async function getMentors(): Promise<MentorView[]> {
  let items: MentorRaw[] = []
  try {
    const res = await fetch(`${API_BASE}/api/mentors`, {
      headers: { Accept: 'application/json' },
    })
    if (!res.ok) {
      console.warn(`[mentors] api returned ${res.status}, skipping`)
    }
    else {
      const data = await res.json() as { items?: MentorRaw[] }
      items = data.items ?? []
    }
  }
  catch (e) {
    console.warn('[mentors] fetch failed:', e)
  }
  failBuildIfEmptyInCi('mentors', items.length)
  return deterministicShuffle(items.map(toView), Date.now() & 0xFFFF)
}

export function collectTags(mentors: MentorView[]): string[] {
  const set = new Set<string>()
  for (const m of mentors) {
    for (const t of m.tags)
      set.add(t)
  }
  return Array.from(set).sort((a, b) => a.localeCompare(b, 'ru'))
}
