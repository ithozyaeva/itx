<script lang="ts" setup>
import type { ReviewOnCommunityFormData } from '@/models/reviewOnCommunity'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { requiredRule, useFormValidation } from '@/composables/useFormValidation'
import { formatDateToInput } from '@/lib/utils'
import { reviewOnCommunityService } from '@/services/reviewOnCommunityService'
import { onMounted, ref } from 'vue'
import { toast } from '../ui/toast'

const props = defineProps<{
  reviewId?: number
}>()

const emit = defineEmits(['saved', 'cancel'])
const isLoading = ref(false)

const validationRules = {
  authorTg: [requiredRule],
  text: [requiredRule],
  date: [requiredRule],
  id: [],
  authorId: [],
}

// Используем composable для валидации формы
const { values, errors, touched, validate, handleBlur, isValid } = useFormValidation<ReviewOnCommunityFormData>({
  date: new Date().toISOString(),
  authorTg: '',
  text: '',
}, validationRules)

// Загрузка данных участника при редактировании
onMounted(async () => {
  if (props.reviewId) {
    const review = await reviewOnCommunityService.getById(props.reviewId)
    values.value = {
      id: review?.id,
      authorId: review?.authorId,
      text: review?.text ?? '',
      date: formatDateToInput(review?.date),
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
      await reviewOnCommunityService.update(props.reviewId, values.value)
    }
    else {
      await reviewOnCommunityService.create(values.value)
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
      <div v-if="!values.authorId" class="space-y-2">
        <Label for="authorTg">ТГ автора</Label>
        <Input
          id="authorTg"
          v-model="values.authorTg"
          placeholder="Введите tg участника"
          :class="{ 'border-red-500': touched.authorTg && errors.authorTg }"
          @blur="handleBlur('authorTg')"
        />
        <p v-if="touched.authorTg && errors.authorTg" class="text-sm text-red-500">
          {{ errors.authorTg }}
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
