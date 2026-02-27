<script lang="ts" setup>
import { Briefcase, MessageSquare } from 'lucide-vue-next'
import { computed, defineProps } from 'vue'
import MailIcon from '~icons/mdi/email-variant'
import PhoneIcon from '~icons/mdi/phone'
import TgIcon from '~icons/mdi/telegram'
import TgImage from '@/components/ui/TgImage.vue'

const props = defineProps<{
  mentor: Mentor
}>()

// TODO: Первести в базе на строку
enum ContactType {
  Telegram = 1,
  Email = 2,
  Phone = 3,
  Other = 4,
}

interface Contact {
  id: number
  type: ContactType
  link: string
}

export interface Mentor {
  id: number
  tg: string
  firstName: string
  lastName: string
  occupation: string
  experience: string
  profTags: { id: number, title: string }[]
  contacts: Contact[]
  services: { id: number, name: string, price: number }[]
}

const contacts = computed(() => {
  const hasTypes: Record<ContactType, boolean> = {
    [ContactType.Telegram]: false,
    [ContactType.Email]: false,
    [ContactType.Phone]: false,
    [ContactType.Other]: false,
  }
  return props.mentor.contacts.map((item) => {
    const newItem = { ...item, colorType: hasTypes[item.type] ? 'secondary' : 'primary' }

    if (!hasTypes[item.type]) {
      hasTypes[item.type] = true
    }

    return newItem
  })
})

function getLink(contact: Contact) {
  switch (contact.type) {
    case ContactType.Telegram:
      return contact.link
    case ContactType.Email:
      return `mailto:${contact.link}`
    case ContactType.Phone:
      return `tel:${contact.link}`
    case ContactType.Other:
      return contact.link
  }
}
</script>

<template>
  <div class="bg-white rounded-xl shadow-md overflow-hidden hover:shadow-lg transition-shadow duration-300 p-4 md:p-6 h-fit">
    <div class="flex items-start md:items-center gap-10 flex-col md:flex-row">
      <TgImage
        :username="mentor.tg"
        class="rounded-full self-center"
        width="160"
        height="160"
      />
      <div>
        <h3 class="text-xl font-semibold text-gray-900">
          {{ mentor.firstName }} {{ mentor.lastName }}
        </h3>
        <p
          v-if="mentor.occupation"
          class="text-gray-600 mt-1"
        >
          {{ mentor.occupation }}
        </p>
        <div
          v-if="mentor.experience"
          class="mt-4"
        >
          <div class="flex items-start text-sm text-gray-500">
            <Briefcase class="h-4 w-4 mr-2 mt-1 shrink-0" />
            <span>{{ mentor.experience }}</span>
          </div>
        </div>
      </div>
    </div>

    <div class="mt-4">
      <div
        v-if="mentor.profTags?.length > 0"
        class="flex flex-wrap gap-2"
      >
        <span
          v-for="tag in mentor.profTags"
          :key="tag.id"
          class="px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded-full"
        >
          {{ tag.title }}
        </span>
      </div>

      <div
        v-if="mentor.services?.length > 0"
        class="mt-4"
      >
        <h4 class="text-sm font-medium text-gray-900 mb-2">
          Услуги:
        </h4>
        <ul class="space-y-2">
          <li
            v-for="service in mentor.services"
            :key="service.id"
            class="flex justify-between items-start text-sm gap-5"
          >
            <span>{{ service.name }}</span>
            <span class="font-medium whitespace-nowrap">{{
              service.price.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ' ') }} ₽</span>
          </li>
        </ul>
      </div>

      <div
        v-if="mentor.contacts?.length > 0"
        class="mt-5"
      >
        <hr class="mb-4">
        <div class="flex items-center gap-5 flex-wrap">
          <a
            v-for="contact in contacts"
            :key="contact.id"
            class="transition-colors duration-200 cursor-pointer"
            :href="getLink(contact)"
            target="_blank"
            :class="contact.colorType === 'primary' ? 'text-blue-500 hover:text-blue-700' : 'text-gray-500 hover:text-gray-700'"
          >
            <TgIcon
              v-if="contact.type === ContactType.Telegram"
              class="h-7 w-7"
            />
            <MailIcon
              v-else-if="contact.type === ContactType.Email"
              class="h-7 w-7"
            />
            <PhoneIcon
              v-else-if="contact.type === ContactType.Phone"
              class="h-7 w-7"
            />
            <MessageSquare
              v-else-if="contact.type === ContactType.Other"
              class="h-7 w-7"
            />
          </a>
        </div>
      </div>
    </div>
  </div>
</template>
