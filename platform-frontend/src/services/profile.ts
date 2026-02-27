import type { Contacts, Mentor, ProfTag, Service, TelegramUser } from '@/models/profile'
import { useUser } from '@/composables/useUser'
import ky from 'ky'
import { apiClient } from './api'

const localStorageUser = useUser<TelegramUser | Mentor>()

export const profileSiervice = {
  async getMe() {
    try {
      const response = await apiClient.get('members/me')

      if (response.ok) {
        const data = await response.json<TelegramUser | Mentor>()

        localStorageUser.value = data
      }
    }
    catch (err) {
      console.error('Произошла ошибка', err)
    }
  },
  async updateMe(newUser: { firstName: string, lastName: string, birthday: string }) {
    const response = await apiClient.patch('members/me', { json: newUser })
    const data = await response.json<TelegramUser>()
    localStorageUser.value = data

    return data
  },

  async updateTags(profTags: ProfTag[]) {
    const response = await apiClient.post('mentors/me/update-prof-tags', { json: { profTags } })
    const data = await response.json<Mentor>()

    localStorageUser.value = data
    return data
  },

  async updateServices(services: Service[]) {
    const response = await apiClient.post('mentors/me/update-services', { json: { services } })
    const data = await response.json<Mentor>()

    localStorageUser.value = data
    return data
  },

  async updateContacts(contacts: Contacts[]) {
    const response = await apiClient.post('mentors/me/update-contacts', { json: { contacts } })
    const data = await response.json<Mentor>()

    localStorageUser.value = data
    return data
  },

  async updateMentorInfo(info: { occupation: string, experience: string }) {
    const response = await apiClient.post('mentors/me/update-info', { json: info })
    const data = await response.json<Mentor>()

    localStorageUser.value = data
    return data
  },

  async getAllProfTags() {
    const response = await ky.get('/api/profTags')
    return (await response.json<{ items: ProfTag[] }>()).items
  },
}
