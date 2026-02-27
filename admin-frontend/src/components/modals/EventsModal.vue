<script setup lang="ts">
import EventForm from '@/components/forms/EventForm.vue'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'

defineProps<{
  isOpen: boolean
  eventId?: number
}>()

const emit = defineEmits(['update:isOpen', 'saved'])

function handleSaved() {
  emit('saved')
  emit('update:isOpen', false)
}

function handleCancel() {
  emit('update:isOpen', false)
}
</script>

<template>
  <Dialog :open="isOpen" @update:open="$emit('update:isOpen', $event)">
    <DialogContent class="sm:max-w-[1000px] max-h-full overflow-auto">
      <DialogHeader>
        <DialogTitle>{{ eventId ? 'Редактирование события' : 'Добавление события' }}</DialogTitle>
      </DialogHeader>
      <EventForm
        :event-id="eventId"
        @saved="handleSaved"
        @cancel="handleCancel"
      />
    </DialogContent>
  </Dialog>
</template>
