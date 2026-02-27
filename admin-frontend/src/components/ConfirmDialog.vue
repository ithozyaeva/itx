<script setup lang="ts">
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog'
import { Button } from '@/components/ui/button'

withDefaults(defineProps<{
  title?: string
  description?: string
  confirmLabel?: string
  variant?: 'default' | 'destructive'
}>(), {
  title: 'Вы уверены?',
  description: 'Это действие нельзя отменить.',
  confirmLabel: 'Подтвердить',
  variant: 'destructive',
})

const emit = defineEmits<{
  confirm: []
}>()
</script>

<template>
  <AlertDialog>
    <AlertDialogTrigger>
      <slot name="trigger" />
    </AlertDialogTrigger>
    <AlertDialogContent>
      <AlertDialogHeader>
        <AlertDialogTitle>{{ title }}</AlertDialogTitle>
        <AlertDialogDescription>{{ description }}</AlertDialogDescription>
      </AlertDialogHeader>
      <AlertDialogFooter>
        <AlertDialogCancel>
          <Button variant="outline">
            Отмена
          </Button>
        </AlertDialogCancel>
        <AlertDialogAction>
          <Button :variant="variant" @click="emit('confirm')">
            {{ confirmLabel }}
          </Button>
        </AlertDialogAction>
      </AlertDialogFooter>
    </AlertDialogContent>
  </AlertDialog>
</template>
