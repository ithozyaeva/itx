<script setup lang="ts">
import type { Service } from '@/models/profile'
import { Loader2 } from 'lucide-vue-next'
import { ref } from 'vue'
import { Button } from '@/components/ui/button'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Textarea } from '@/components/ui/textarea'
import { handleError } from '@/services/errorService'
import { mentorsService } from '@/services/mentors'

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
    handleError(error)
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
    <Select
      :model-value="String(selectedServiceId)"
      @update:model-value="selectedServiceId = $event ? Number($event) : ''"
    >
      <SelectTrigger>
        <SelectValue placeholder="Выберите услугу" />
      </SelectTrigger>
      <SelectContent>
        <SelectItem v-for="service in services" :key="service.id" :value="String(service.id)">
          {{ service.name }}
        </SelectItem>
      </SelectContent>
    </Select>
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
