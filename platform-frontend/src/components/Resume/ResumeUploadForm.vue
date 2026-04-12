<script setup lang="ts">
import type { WorkFormat } from '@/models/resume'
import { Loader2, UploadCloud } from 'lucide-vue-next'
import { reactive, ref } from 'vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Textarea } from '@/components/ui/textarea'
import { Typography } from '@/components/ui/typography'
import { handleError } from '@/services/errorService'
import { resumeService } from '@/services/resume'

const emit = defineEmits<{
  uploaded: []
}>()

const isUploading = ref(false)
const fileInput = ref<HTMLInputElement | null>(null)
const selectedFile = ref<File | null>(null)

const form = reactive({
  workExperience: '',
  desiredPosition: '',
  workFormat: '' as WorkFormat,
})

const workFormatOptions: { label: string, value: WorkFormat }[] = [
  { label: '— Не выбрано —', value: '' },
  { label: 'Удалёнка', value: 'REMOTE' },
  { label: 'Гибрид', value: 'HYBRID' },
  { label: 'Офис', value: 'OFFICE' },
]

function onFileChange(event: Event) {
  const target = event.target as HTMLInputElement
  selectedFile.value = target.files?.[0] ?? null
}

function resetForm() {
  form.workExperience = ''
  form.desiredPosition = ''
  form.workFormat = ''
  selectedFile.value = null
  if (fileInput.value)
    fileInput.value.value = ''
}

async function handleUpload() {
  if (!selectedFile.value)
    return

  isUploading.value = true
  try {
    await resumeService.upload({
      file: selectedFile.value,
      workExperience: form.workExperience || undefined,
      desiredPosition: form.desiredPosition || undefined,
      workFormat: form.workFormat || undefined,
    })
    resetForm()
    emit('uploaded')
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isUploading.value = false
  }
}
</script>

<template>
  <div class="rounded-sm border bg-card p-4 md:p-6 space-y-4 terminal-card">
    <div class="flex items-center space-x-3">
      <UploadCloud class="text-accent" />
      <div>
        <Typography
          variant="h3"
          as="h2"
        >
          Загрузить резюме
        </Typography>
        <p class="text-sm text-muted-foreground">
          Загружайте файлы в формате PDF/DOC/DOCX. Поля заполнятся автоматически и их можно поправить.
        </p>
      </div>
    </div>

    <div class="space-y-4">
      <label class="block text-sm font-medium text-muted-foreground">Файл резюме</label>
      <input
        ref="fileInput"
        type="file"
        accept=".pdf,.doc,.docx"
        class="block w-full border border-input rounded-sm bg-transparent px-4 py-3 cursor-pointer text-sm"
        @change="onFileChange"
      >

      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label class="block text-sm font-medium text-muted-foreground mb-1">Желаемая должность (необязательно)</label>
          <Input
            v-model="form.desiredPosition"
            placeholder="Product Manager"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-muted-foreground mb-1">Формат работы</label>
          <Select v-model="form.workFormat">
            <SelectTrigger>
              <SelectValue placeholder="— Не выбрано —" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem
                v-for="option in workFormatOptions"
                :key="option.value"
                :value="option.value"
              >
                {{ option.label }}
              </SelectItem>
            </SelectContent>
          </Select>
        </div>
      </div>

      <div>
        <label class="block text-sm font-medium text-muted-foreground mb-1">Опыт (необязательно)</label>
        <Textarea
          v-model="form.workExperience"
          placeholder="5 лет в разработке банковских приложений"
          rows="3"
        />
      </div>

      <Button
        class="w-full mt-2 gap-2"
        :disabled="!selectedFile || isUploading"
        @click="handleUpload"
      >
        <Loader2
          v-if="isUploading"
          class="h-4 w-4 animate-spin"
        />
        <span v-else>Загрузить резюме</span>
      </Button>
    </div>
  </div>
</template>
