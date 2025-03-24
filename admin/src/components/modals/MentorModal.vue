<script setup lang="ts">
import MentorForm from '@/components/forms/MentorForm.vue'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'

defineProps<{
  isOpen: boolean
  mentorId?: number
  memberId: number
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
    <DialogContent class="sm:max-w-[1000px]">
      <DialogHeader>
        <DialogTitle>{{ mentorId ? 'Редактирование ментора' : 'Добавление ментора' }}</DialogTitle>
      </DialogHeader>
      <MentorForm
        :mentor-id="mentorId"
        :member-id="memberId"
        @saved="handleSaved"
        @cancel="handleCancel"
      />
    </DialogContent>
  </Dialog>
</template>
