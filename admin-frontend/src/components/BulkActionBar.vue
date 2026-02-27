<script setup lang="ts">
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import { Button } from '@/components/ui/button'

defineProps<{
  count: number
  actions: { label: string, variant?: 'default' | 'destructive', handler: () => void }[]
}>()

const emit = defineEmits<{
  clear: []
}>()
</script>

<template>
  <div
    v-if="count > 0"
    class="fixed bottom-4 left-1/2 -translate-x-1/2 z-50 flex items-center gap-3 bg-background border rounded-2xl shadow-lg px-6 py-3"
  >
    <span class="text-sm font-medium">Выбрано: {{ count }}</span>
    <ConfirmDialog
      v-for="(action, idx) in actions"
      :key="idx"
      :title="`${action.label}?`"
      :description="`Будет применено к ${count} элементам.`"
      :confirm-label="action.label"
      :variant="action.variant ?? 'destructive'"
      @confirm="action.handler"
    >
      <template #trigger>
        <Button size="sm" :variant="action.variant ?? 'destructive'">
          {{ action.label }}
        </Button>
      </template>
    </ConfirmDialog>
    <Button size="sm" variant="ghost" @click="emit('clear')">
      Снять выбор
    </Button>
  </div>
</template>
