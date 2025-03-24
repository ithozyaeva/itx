<script setup lang="ts">
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { computed } from 'vue'

const props = defineProps<{
  occupation: string
  experience: string
  order: number
  errors: Record<string, string | null>
  touched: Record<string, boolean>
}>()

const emit = defineEmits(['update:occupation', 'update:experience', 'update:order'])

const localOccupation = computed({
  get: () => props.occupation,
  set: value => emit('update:occupation', value),
})

const localExperience = computed({
  get: () => props.experience,
  set: value => emit('update:experience', value),
})

const localOrder = computed({
  get: () => props.order,
  set: value => emit('update:order', Number(value)),
})
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle>Основная информация</CardTitle>
    </CardHeader>
    <CardContent class="space-y-4">
      <!-- Специализация -->
      <div class="space-y-2">
        <Label for="occupation">Специализация</Label>
        <Input
          id="occupation"
          v-model="localOccupation"
          placeholder="Введите специализацию"
          :class="{ 'border-destructive': errors.occupation && touched.occupation }"
        />
        <p v-if="errors.occupation && touched.occupation" class="text-sm text-destructive">
          {{ errors.occupation }}
        </p>
      </div>

      <!-- Опыт работы -->
      <div class="space-y-2">
        <Label for="experience">Опыт работы</Label>
        <Textarea
          id="experience"
          v-model="localExperience"
          placeholder="Опишите опыт работы"
          rows="4"
          :class="{ 'border-destructive': errors.experience && touched.experience }"
        />
        <p v-if="errors.experience && touched.experience" class="text-sm text-destructive">
          {{ errors.experience }}
        </p>
      </div>

      <!-- Порядок отображения -->
      <div class="space-y-2">
        <Label for="order">Порядок отображения</Label>
        <Input
          id="order"
          v-model="localOrder"
          type="number"
          placeholder="Введите порядок отображения"
        />
      </div>
    </CardContent>
  </Card>
</template>
