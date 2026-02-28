<script setup lang="ts">
import { reactive } from 'vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import Card from './ui/card/Card.vue'

export interface MentorSearchFilters {
  name: string
  tag: string
}

const emit = defineEmits<{
  apply: [filters: MentorSearchFilters]
}>()

const localFilters = reactive<MentorSearchFilters>({
  name: '',
  tag: '',
})

function applyFilters() {
  emit('apply', { ...localFilters })
}

function resetFilters() {
  localFilters.name = ''
  localFilters.tag = ''
  emit('apply', { ...localFilters })
}
</script>

<template>
  <Card class="p-4 rounded-lg mb-6">
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div>
        <Label for="name" class="text-sm font-medium">Имя / Username</Label>
        <Input
          id="name"
          v-model="localFilters.name"
          placeholder="Поиск по имени или TG"
          class="mt-1"
        />
      </div>

      <div>
        <Label for="tag" class="text-sm font-medium">Профессиональный тег</Label>
        <Input
          id="tag"
          v-model="localFilters.tag"
          placeholder="Поиск по тегу"
          class="mt-1"
        />
      </div>

      <div class="md:col-span-2 flex justify-end gap-2">
        <Button
          variant="outline"
          size="sm"
          @click="resetFilters"
        >
          Сбросить
        </Button>
        <Button
          size="sm"
          @click="applyFilters"
        >
          Применить
        </Button>
      </div>
    </div>
  </Card>
</template>
