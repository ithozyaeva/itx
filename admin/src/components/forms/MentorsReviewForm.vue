<script lang="ts" setup>
import type { Mentor, MentrosReviewFormData } from '@/models/mentors'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Textarea } from '@/components/ui/textarea'
import { requiredRule, useFormValidation } from '@/composables/useFormValidation'
import { formatDateToInput } from '@/lib/utils'
import { mentorService } from '@/services/mentorService'
import { mentorsReviewService } from '@/services/mentorsReviewService'
import { onMounted, ref } from 'vue'
import { toast } from '../ui/toast'

const props = defineProps<{
  reviewId?: number
}>()

const emit = defineEmits(['saved', 'cancel'])
const isLoading = ref(false)
const mentor = ref<Mentor>()
const allMentors = ref<Mentor[]>([])

const validationRules = {
  author: [requiredRule],
  text: [requiredRule],
  date: [requiredRule],
  serviceId: [requiredRule],
  id: [],
  authorId: [],
}

// Используем composable для валидации формы
const { values, errors, touched, validate, handleBlur, isValid } = useFormValidation<MentrosReviewFormData>({
  date: new Date().toISOString(),
  serviceId: 0,
  author: '',
  text: '',
}, validationRules)

// Загрузка данных участника при редактировании
onMounted(async () => {
  if (!props.reviewId) {
    allMentors.value = (await mentorService.getAll())!.items
  }
  if (props.reviewId) {
    const review = await mentorsReviewService.getById(props.reviewId)
    if (review) {
      values.value = {
        id: review?.id,
        author: review?.author,
        serviceId: review?.serviceId,
        text: review?.text ?? '',
        date: formatDateToInput(review?.date),
      }
    }
  }
})

// Обработчик отправки формы
async function handleSubmit(e: Event) {
  e.preventDefault()

  if (!validate()) {
    return
  }

  try {
    if (props.reviewId) {
      await mentorsReviewService.update(props.reviewId, {
        date: values.value.date,
        text: values.value.text,
        id: props.reviewId,
        serviceId: values.value.serviceId,
        author: values.value.author,
      })
    }
    else {
      await mentorsReviewService.create({
        date: values.value.date,
        text: values.value.text,
        serviceId: values.value.serviceId,
        author: values.value.author,
      })
    }

    emit('saved')
  }
  catch (err) {
    console.error(err)
    toast({
      variant: 'destructive',
      title: 'Не удалось сохранить отзыв.',
    })
  }
}

function handleCancel() {
  emit('cancel')
}
</script>

<template>
  <form class="space-y-6" @submit="handleSubmit">
    <div class="space-y-4">
      <div v-if="!props.reviewId" class="space-y-2">
        <Label for="mentor">Ментор</Label>
        <Select id="mentor" v-model="mentor">
          <SelectTrigger>
            <SelectValue placeholder="Выберите ментора" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem v-for="mentorItem in allMentors" :key="mentorItem.id" :value="mentorItem">
              {{ mentorItem.name }}
            </SelectItem>
          </SelectContent>
        </Select>
      </div>
      <div v-if="mentor" class="space-y-2">
        <Label for="service">Услуга ментора</Label>
        <Select id="service" v-model="values.serviceId">
          <SelectTrigger>
            <SelectValue placeholder="Выберите услугу" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem v-for="service in mentor.services" :key="service.id" :value="service.id">
              {{ service.name }}
            </SelectItem>
          </SelectContent>
        </Select>
      </div>
      <div class="space-y-2">
        <Label for="author">Автор</Label>
        <Input
          id="author"
          v-model="values.author"
          placeholder="Автор"
          :class="{ 'border-red-500': touched.author && errors.author }"
          @blur="handleBlur('author')"
        />
        <p v-if="touched.author && errors.author" class="text-sm text-red-500">
          {{ errors.author }}
        </p>
      </div>

      <div class="space-y-2">
        <Label for="date">Дата отзыва</Label>
        <Input
          id="date"
          v-model="values.date"
          type="date"
          placeholder="Отзыв"
          :class="{ 'border-red-500': touched.date && errors.date }"
          @blur="handleBlur('text')"
        />
        <p v-if="touched.date && errors.date" class="text-sm text-red-500">
          {{ errors.date }}
        </p>
      </div>
      <div class="space-y-2">
        <Label for="text">Текст отзыва</Label>
        <Textarea
          id="text"
          v-model="values.text"
          placeholder="Отзыв"
          :class="{ 'border-red-500': touched.text && errors.text }"
          @blur="handleBlur('text')"
        />
        <p v-if="touched.text && errors.text" class="text-sm text-red-500">
          {{ errors.text }}
        </p>
      </div>
    </div>

    <div class="flex justify-end space-x-2">
      <Button type="button" variant="outline" @click="handleCancel">
        Отмена
      </Button>
      <Button type="submit" :disabled="isLoading || !isValid">
        {{ isLoading ? 'Сохранение...' : props.reviewId ? 'Сохранить' : 'Добавить' }}
      </Button>
    </div>
  </form>
</template>
