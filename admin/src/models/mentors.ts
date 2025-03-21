/**
 * Профессиональные теги.
 *
 * @prop id - идентификатор тега.
 * @prop title - название тега.
 */

export interface ProfTag {
  id: number
  title: string
}

/**
 * Ментор.
 *
 * @prop id - идентификатор ментора.
 * @prop name - имя ментора.
 * @prop tg - ссылка на Telegram.
 * @prop order - порядковый номер ментора.
 * @prop occupation - профессия ментора.
 * @prop experience - опыт работы ментора.
 * @prop profTags - профессиональные теги ментора.
 * @prop services - услуги ментора.
 * @prop contacts - контакты ментора.
 */
export interface Mentor {
  id: number
  name: string
  tg: string
  order: number
  memberId: number
  occupation: string
  experience: string
  profTags: ProfTag[]
  services: Service[]
  contacts: Contacts[]
}

/**
 * Услуга ментора.
 *
 * @prop id - идентификатор услуги.
 * @prop name - название услуги.
 * @prop price - стоимость услуги.
 * @prop reviews - отзывы о услуге.
 */
export interface Service {
  id: number
  name: string
  price: number
  reviews: ReviewOnService[]
}

/**
 * Контакты ментора.
 *
 * @prop id - идентификатор контакта.
 * @prop type - тип контакта.
 * @prop link - ссылка на контакт.
 */
export interface Contacts {
  id: number
  type: number
  link: string
}

/**
 * Отзыв о услуге.
 *
 * @prop id - идентификатор отзыва.
 * @prop author - автор отзыва.
 * @prop Text - текст отзыва.
 * @prop Date - дата отзыва.
 */
export interface ReviewOnService {
  id: number
  author: string
  serviceId: number
  text: string
  date: string
  serviceName: string
  mentorName: string
}

export interface MentorFormData {
  id?: number
  memberId: number
  occupation: string
  experience: string
  order: number
  profTags: ProfTagFormData[]
  contacts: ContactFormData[]
  services: ServiceFormData[]
}

export interface ContactFormData {
  id?: number
  type: number
  link: string
}

export interface ServiceFormData {
  id?: number
  name: string
  price: number
}

export interface ProfTagFormData {
  id?: number
  title: string
}

export interface MentrosReviewFormData {
  id?: number
  serviceId: number
  author: string
  text: string
  date: string
}
