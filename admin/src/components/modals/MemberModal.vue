<script setup lang="ts">
import MemberForm from '@/components/forms/MemberForm.vue'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'

defineProps<{
  isOpen: boolean
  memberId: number | null
}>()

const emit = defineEmits(['update:isOpen', 'saved'])

function handleClose() {
  emit('update:isOpen', false)
}

function handleSaved() {
  emit('saved')
  handleClose()
}
</script>

<template>
  <Dialog :open="isOpen" @update:open="handleClose">
    <DialogContent class="sm:max-w-[500px]">
      <DialogHeader>
        <DialogTitle>{{ Boolean(memberId) ? 'Редактирование участника' : 'Добавление участника' }}</DialogTitle>
      </DialogHeader>

      <MemberForm
        :member-id="memberId"
        @saved="handleSaved"
        @cancel="handleClose"
      />
    </DialogContent>
  </Dialog>
</template>
