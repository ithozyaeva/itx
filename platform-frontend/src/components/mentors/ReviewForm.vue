<script setup lang="ts">
import type { Service } from '@/models/profile'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { mentorsService } from '@/services/mentors'
import { Loader2 } from 'lucide-vue-next'
import { ref } from 'vue'

const props = defineProps<{
  mentorId: number
  services: Service[]
}>()

const emit = defineEmits<{
  submitted: []
}>()

const selectedServiceId = ref<number | ''>('')
const text = ref('')
const isSubmitting = ref(false)

async function handleSubmit() {
  if (!selectedServiceId.value || !text.value.trim())
    return

  isSubmitting.value = true
  try {
    await mentorsService.addReview(props.mentorId, Number(selectedServiceId.value), text.value.trim())
    selectedServiceId.value = ''
    text.value = ''
    emit('submitted')
  }
  catch (error) {
    console.error('Ошибка при отправке отзыва:', error)
  }
  finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <div class="space-y-4">
    <h4 class="font-semibold">
      Оставить отзыв
    </h4>
    <select
      v-model="selectedServiceId"
      class="w-full border border-input rounded-xl bg-transparent px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
    >
      <option value="" disabled>
        Выберите услугу
      </option>
      <option v-for="service in services" :key="service.id" :value="service.id">
        {{ service.name }}
      </option>
    </select>
    <Textarea
      v-model="text"
      placeholder="Ваш отзыв..."
      rows="3"
    />
    <Button
      :disabled="!selectedServiceId || !text.trim() || isSubmitting"
      @click="handleSubmit"
    >
      <Loader2 v-if="isSubmitting" class="h-4 w-4 animate-spin mr-2" />
      Отправить отзыв
    </Button>
  </div>
</template>
