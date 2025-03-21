<script setup lang="ts">
import type { Contacts, MentorFormData, ProfTag, Service } from '@/models/mentors'
import { Button } from '@/components/ui/button'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { requiredRule, useFormValidation } from '@/composables/useFormValidation'
import { handleError } from '@/services/errorService'
import { memberService } from '@/services/memberService'
import { mentorService } from '@/services/mentorService'
import { profTagService } from '@/services/profTagService'
import { onMounted, ref } from 'vue'
import MentorContacts from './mentor/MentorContacts.vue'
import MentorGeneralInfo from './mentor/MentorGeneralInfo.vue'
import MentorProfTags from './mentor/MentorProfTags.vue'
import MentorServices from './mentor/MentorServices.vue'

const props = defineProps<{
  mentorId?: number
  memberId: number
}>()

const emit = defineEmits(['saved', 'cancel'])
const isLoading = ref(false)
const profTags = ref<ProfTag[]>([])
const activeTab = ref('general')

// Используем composable для валидации формы
const { values, errors, touched, validate, isValid, resetForm } = useFormValidation<MentorFormData>({
  memberId: props.memberId,
  occupation: '',
  experience: '',
  order: 0,
  profTags: [],
  contacts: [] as Contacts[],
  services: [] as Service[],
  id: props.mentorId,
}, {
  memberId: [requiredRule],
  occupation: [requiredRule],
  order: [],
  contacts: [],
  services: [],
  id: [],
  experience: [],
  profTags: [],
})

// Загрузка данных для формы
onMounted(async () => {
  isLoading.value = true

  // Загружаем список профессиональных тегов
  const tagsData = await profTagService.getAll()
  profTags.value = tagsData.items

  // Если редактируем существующего ментора, загружаем его данные
  if (props.mentorId) {
    const mentor = await mentorService.getById(props.mentorId)
    if (mentor) {
      values.value.occupation = mentor.occupation
      values.value.experience = mentor.experience
      values.value.order = mentor.order
      values.value.profTags = mentor.profTags || []
      values.value.contacts = mentor.contacts || []
      values.value.services = mentor.services || []
    }
  }

  isLoading.value = false
})

// Обработчик отправки формы
async function handleSubmit(e: Event) {
  e.preventDefault()

  if (!validate()) {
    return
  }
  try {
    if (props.mentorId) {
    // Обновление существующего ментора

      await mentorService.update(props.mentorId, {
        id: props.mentorId,
        memberId: Number(values.value.memberId),
        occupation: values.value.occupation,
        experience: values.value.experience,
        order: Number(values.value.order),
        profTags: values.value.profTags,
        contacts: values.value.contacts,
        services: values.value.services,
      })
    }
    else {
    // Создание нового ментора
      await mentorService.create({
        memberId: Number(values.value.memberId),
        occupation: values.value.occupation,
        experience: values.value.experience,
        order: Number(values.value.order),
        profTags: values.value.profTags,
        contacts: values.value.contacts,
        services: values.value.services,
      })
    }

    resetForm()
    memberService.search()
    emit('saved')
  }
  catch (err) {
    handleError(err)
  }
}

function handleCancel() {
  emit('cancel')
}
</script>

<template>
  <form class="space-y-6" @submit="handleSubmit" @keydown.enter.prevent>
    <Tabs v-model="activeTab" class="w-full">
      <TabsList class="grid grid-cols-4 mb-4">
        <TabsTrigger value="general">
          Основная информация
        </TabsTrigger>
        <TabsTrigger value="tags">
          Профессиональные теги
        </TabsTrigger>
        <TabsTrigger value="contacts">
          Контакты
        </TabsTrigger>
        <TabsTrigger value="services">
          Услуги
        </TabsTrigger>
      </TabsList>

      <!-- Вкладка с общей информацией -->
      <TabsContent value="general" class="space-y-4">
        <MentorGeneralInfo
          v-model:occupation="values.occupation"
          v-model:experience="values.experience"
          v-model:order="values.order"
          :errors="errors"
          :touched="touched"
        />
      </TabsContent>

      <!-- Вкладка с профессиональными тегами -->
      <TabsContent value="tags" class="space-y-4">
        <MentorProfTags
          v-model:prof-tags="values.profTags"
        />
      </TabsContent>

      <TabsContent value="contacts" class="space-y-6">
        <MentorContacts
          v-model:contacts="values.contacts"
        />
      </TabsContent>

      <TabsContent value="services" class="space-y-6">
        <MentorServices
          v-model:services="values.services"
        />
      </TabsContent>
    </Tabs>

    <div class="flex justify-end space-x-2">
      <Button type="button" variant="outline" @click="handleCancel">
        Отмена
      </Button>
      <Button type="submit" :disabled="isLoading || !isValid">
        {{ isLoading ? 'Сохранение...' : props.mentorId ? 'Сохранить' : 'Добавить' }}
      </Button>
    </div>
  </form>
</template>
