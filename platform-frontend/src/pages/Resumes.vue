<script setup lang="ts">
import type { Resume, WorkFormat } from '@/models/resume'
import { Download, FileText, Loader2, Pencil, Trash2 } from 'lucide-vue-next'
import { onMounted, ref } from 'vue'
import EmptyState from '@/components/common/EmptyState.vue'
import ErrorState from '@/components/common/ErrorState.vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import ResumeUploadForm from '@/components/Resume/ResumeUploadForm.vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Textarea } from '@/components/ui/textarea'
import { Typography } from '@/components/ui/typography'
import { formatShortDate } from '@/lib/utils'
import { handleError } from '@/services/errorService'
import { resumeService } from '@/services/resume'

const resumes = ref<Resume[]>([])
const isLoading = ref(false)
const loadError = ref<string | null>(null)
const isSavingEdit = ref(false)
const deletingId = ref<number | null>(null)

const editingId = ref<number | null>(null)
const editForm = ref({
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

const showUploadForm = ref(false)

async function loadResumes() {
  isLoading.value = true
  loadError.value = null
  try {
    resumes.value = await resumeService.listMine()
  }
  catch (error) {
    loadError.value = (await handleError(error)).message
  }
  finally {
    isLoading.value = false
  }
}

function handleUploaded() {
  showUploadForm.value = false
  loadResumes()
}

function startEdit(resume: Resume) {
  editingId.value = resume.id
  editForm.value = {
    workExperience: resume.workExperience || '',
    desiredPosition: resume.desiredPosition || '',
    workFormat: (resume.workFormat || '') as WorkFormat,
  }
}

async function saveEdit(resume: Resume) {
  if (editingId.value !== resume.id)
    return
  isSavingEdit.value = true
  try {
    const updated = await resumeService.update(resume.id, {
      workExperience: editForm.value.workExperience,
      desiredPosition: editForm.value.desiredPosition,
      workFormat: editForm.value.workFormat || undefined,
    })
    const index = resumes.value.findIndex(item => item.id === resume.id)
    if (index !== -1)
      resumes.value[index] = updated
    editingId.value = null
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isSavingEdit.value = false
  }
}

async function deleteResume(resume: Resume) {
  if (deletingId.value)
    return
  deletingId.value = resume.id
  try {
    await resumeService.delete(resume.id)
    resumes.value = resumes.value.filter(item => item.id !== resume.id)
  }
  catch (error) {
    handleError(error)
  }
  finally {
    deletingId.value = null
  }
}

async function downloadResume(resume: Resume) {
  try {
    const data = await resumeService.download(resume.id)
    if (!data?.url) {
      handleError(new Error('Ссылка для скачивания недоступна'))
      return
    }
    window.open(data.url, '_blank')
  }
  catch (error) {
    handleError(error)
  }
}

function formatWorkFormat(format?: WorkFormat) {
  if (!format)
    return 'Не указано'
  const match = workFormatOptions.find(option => option.value === format)
  return match?.label ?? 'Не указано'
}

onMounted(loadResumes)
</script>

<template>
  <div class="min-h-screen pt-20 pb-10">
    <div class="container mx-auto px-4 max-w-4xl">
      <div class="flex items-center justify-between mb-6">
        <Typography
          variant="h2"
          as="h1"
        >
          Мои резюме
        </Typography>
        <Button
          v-if="!showUploadForm"
          @click="showUploadForm = true"
        >
          Загрузить резюме
        </Button>
      </div>

      <div
        v-if="showUploadForm"
        class="mb-6"
      >
        <ResumeUploadForm @uploaded="handleUploaded" />
      </div>

      <div v-if="isLoading" class="flex justify-center py-12">
        <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
      </div>

      <ErrorState
        v-else-if="loadError"
        :message="loadError"
        @retry="loadResumes"
      />

      <EmptyState
        v-else-if="resumes.length === 0"
        :icon="FileText"
        title="Резюме пока нет"
        description="Загрузите своё первое резюме"
        action-label="Загрузить резюме"
        @action="showUploadForm = true"
      />

      <div v-else class="space-y-4">
        <div
          v-for="resume in resumes"
          :key="resume.id"
          class="rounded-3xl border bg-card p-4"
        >
          <div class="flex flex-col md:flex-row md:items-start md:justify-between gap-4">
            <div class="min-w-0 flex-1">
              <p class="text-base font-medium break-all">
                {{ resume.fileName }}
              </p>
              <p class="text-sm text-muted-foreground">
                Загружено {{ formatShortDate(resume.createdAt) }}
              </p>
            </div>
            <div class="flex gap-2 shrink-0">
              <Button
                variant="secondary"
                size="sm"
                class="gap-2"
                @click="downloadResume(resume)"
              >
                <Download class="h-4 w-4" />
              </Button>
              <Button
                variant="secondary"
                size="sm"
                class="gap-2"
                @click="startEdit(resume)"
              >
                <Pencil class="h-4 w-4" />
              </Button>
              <ConfirmDialog
                title="Удалить резюме?"
                description="Резюме будет удалено без возможности восстановления."
                confirm-label="Удалить"
                @confirm="deleteResume(resume)"
              >
                <template #trigger>
                  <Button
                    variant="destructive"
                    size="sm"
                    class="gap-2"
                  >
                    <Trash2 class="h-4 w-4" />
                  </Button>
                </template>
              </ConfirmDialog>
            </div>
          </div>

          <dl class="mt-4 grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
            <div>
              <dt class="text-muted-foreground">
                Должность
              </dt>
              <dd class="font-medium text-foreground">
                {{ resume.desiredPosition || 'Не указано' }}
              </dd>
            </div>
            <div>
              <dt class="text-muted-foreground">
                Формат работы
              </dt>
              <dd class="font-medium text-foreground">
                {{ formatWorkFormat(resume.workFormat as WorkFormat) }}
              </dd>
            </div>
            <div>
              <dt class="text-muted-foreground">
                Дата загрузки
              </dt>
              <dd class="font-medium text-foreground">
                {{ formatShortDate(resume.createdAt) }}
              </dd>
            </div>
          </dl>

          <div class="mt-4 text-sm">
            <p class="text-muted-foreground">
              Опыт
            </p>
            <p class="text-foreground whitespace-pre-line break-words">
              {{ resume.workExperience || 'Не указано' }}
            </p>
          </div>

          <div
            v-if="editingId === resume.id"
            class="mt-4 border-t pt-4 space-y-3"
          >
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <Input
                v-model="editForm.desiredPosition"
                placeholder="Желаемая должность"
              />
              <Select v-model="editForm.workFormat">
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
            <Textarea
              v-model="editForm.workExperience"
              rows="3"
              placeholder="Опыт работы"
            />
            <div class="flex gap-2 justify-end">
              <Button
                variant="ghost"
                @click="editingId = null"
              >
                Отмена
              </Button>
              <Button
                class="gap-2"
                :disabled="isSavingEdit"
                @click="saveEdit(resume)"
              >
                <Loader2
                  v-if="isSavingEdit"
                  class="h-4 w-4 animate-spin"
                />
                <span v-else>Сохранить</span>
              </Button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
