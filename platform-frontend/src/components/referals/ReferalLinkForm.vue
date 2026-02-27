<script setup lang="ts">
import type { Grade, ReferalLink } from '@/models/referals'
import ProfTagsInput from '@/components/common/ProfTagsInput.vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { useDictionary } from '@/composables/useDictionary'
import { Typography } from 'itx-ui-kit'
import { Loader2 } from 'lucide-vue-next'
import { ref } from 'vue'

const props = defineProps<{ link?: Partial<ReferalLink>, title?: string, isSaving?: boolean }>()

const emit = defineEmits(['save', 'cancel'])

const formData = ref<Partial<ReferalLink>>({
  company: props.link?.company || '',
  grade: props.link?.grade || 'junior',
  profTags: props.link?.profTags || [],
  vacationsCount: props.link?.vacationsCount || 1,
})

function handleSave() {
  emit('save', formData.value)
}

function handleCancel() {
  emit('cancel')
}

const { grades } = useDictionary<Grade>(['grades'])
</script>

<template>
  <div class="flex flex-col gap-4" style="min-height: 150px;">
    <Typography variant="h4" as="h3">
      {{ title }}
    </Typography>
    <div class="flex flex-col gap-2">
      <div>
        <Label for="company" class="block text-sm font-medium text-muted-foreground">Компания</Label>
        <Input id="company" v-model="formData.company" type="text" class="mt-1 block w-full rounded-md border-input shadow-sm" />
      </div>
      <div>
        <Label for="grade" class="block text-sm font-medium text-muted-foreground">Грейд</Label>
        <Select id="grade" v-model="formData.grade" class="mt-1 block w-full rounded-md border-input shadow-sm">
          <SelectTrigger>
            <SelectValue placeholder="Выберите грейд" />
          </SelectTrigger>
          <SelectContent class="w-full">
            <SelectItem v-for="grade in grades" :key="grade.value" :value="grade.value">
              {{ grade.label }}
            </SelectItem>
          </SelectContent>
        </Select>
      </div>
      <div>
        <Label class="block text-sm font-medium text-muted-foreground">Навыки</Label>
        <ProfTagsInput v-model="formData.profTags" />
      </div>
      <div>
        <Label for="vacations" class="block text-sm font-medium text-muted-foreground">Количество вакансий</Label>
        <Input id="vacations" v-model.number="formData.vacationsCount" type="number" min="1" class="mt-1 block w-full rounded-md border-input shadow-sm" />
      </div>
    </div>
    <div class="flex justify-end gap-2 mt-4">
      <Button variant="outline" class="px-4 py-2 text-sm font-medium" @click="handleCancel">
        Отменить
      </Button>
      <Button class="px-4 py-2 text-sm font-medium gap-2" :disabled="isSaving" @click="handleSave">
        <Loader2 v-if="isSaving" class="h-4 w-4 animate-spin" />
        Сохранить
      </Button>
    </div>
  </div>
</template>
