<script setup lang="ts">
import { AlertTriangle, RefreshCw, WifiOff } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'

const props = withDefaults(defineProps<{
  message?: string
  type?: 'network' | 'server' | 'unknown'
}>(), {
  message: 'Произошла ошибка при загрузке данных',
  type: 'unknown',
})

const emit = defineEmits<{
  retry: []
}>()
</script>

<template>
  <div class="flex flex-col items-center justify-center py-12 px-4 text-center">
    <div class="flex items-center justify-center w-16 h-16 rounded-sm bg-destructive/10 mb-4">
      <WifiOff
        v-if="props.type === 'network'"
        class="h-8 w-8 text-destructive"
      />
      <AlertTriangle
        v-else
        class="h-8 w-8 text-destructive"
      />
    </div>
    <h3 class="text-base font-semibold mb-1">
      {{ props.type === 'network' ? 'Нет соединения' : 'Ошибка загрузки' }}
    </h3>
    <p class="text-sm text-muted-foreground max-w-sm mb-4">
      {{ message }}
    </p>
    <Button
      variant="outline"
      size="sm"
      @click="emit('retry')"
    >
      <RefreshCw class="h-4 w-4 mr-2" />
      Попробовать снова
    </Button>
  </div>
</template>
