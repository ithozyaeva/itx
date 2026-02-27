<script setup lang="ts">
import type { Resume, WorkFormat } from '@/models/resume'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { resumeService } from '@/services/resume'
import { Typography } from 'itx-ui-kit'
import { FileText, Loader2, Pencil, Trash2, UploadCloud } from 'lucide-vue-next'
import { onMounted, reactive, ref } from 'vue'

const resumes = ref<Resume[]>([])
const isUploading = ref(false)
const isSavingEdit = ref(false)
const fileInput = ref<HTMLInputElement | null>(null)
const selectedFile = ref<File | null>(null)

const form = reactive({
  workExperience: '',
  desiredPosition: '',
  workFormat: '' as WorkFormat,
})

const editingId = ref<number | null>(null)
const editForm = reactive({
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

onMounted(async () => {
  await loadResumes()
})

async function loadResumes() {
  try {
    resumes.value = await resumeService.listMine()
  }
  catch (error) {
    console.error('Не удалось загрузить резюме', error)
  }
}

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
  if (!selectedFile.value) {
    return
  }
  isUploading.value = true
  try {
    const response = await resumeService.upload({
      file: selectedFile.value,
      workExperience: form.workExperience || undefined,
      desiredPosition: form.desiredPosition || undefined,
      workFormat: form.workFormat || undefined,
    })
    resumes.value = [response.resume, ...resumes.value]
    resetForm()
  }
  catch (error) {
    console.error(error)
  }
  finally {
    isUploading.value = false
  }
}

function startEdit(resume: Resume) {
  editingId.value = resume.id
  editForm.workExperience = resume.workExperience || ''
  editForm.desiredPosition = resume.desiredPosition || ''
  editForm.workFormat = (resume.workFormat || '') as WorkFormat
}

async function saveEdit(resume: Resume) {
  if (editingId.value !== resume.id)
    return
  isSavingEdit.value = true
  try {
    const updated = await resumeService.update(resume.id, {
      workExperience: editForm.workExperience,
      desiredPosition: editForm.desiredPosition,
      workFormat: editForm.workFormat || undefined,
    })
    const index = resumes.value.findIndex(item => item.id === resume.id)
    if (index !== -1)
      resumes.value[index] = updated
    editingId.value = null
  }
  catch (error) {
    console.error(error)
  }
  finally {
    isSavingEdit.value = false
  }
}

async function deleteResume(resume: Resume) {
  try {
    await resumeService.delete(resume.id)
    resumes.value = resumes.value.filter(item => item.id !== resume.id)
  }
  catch (error) {
    console.error(error)
  }
}

function formatWorkFormat(format?: WorkFormat) {
  if (!format)
    return 'Не указано'
  const match = workFormatOptions.find(option => option.value === format)
  return match?.label ?? 'Не указано'
}
</script>

<template>
  <div class="p-6 md:p-8 bg-card backdrop-blur border border-border shadow-lg rounded-3xl space-y-6">
    <div class="flex items-center space-x-3">
      <UploadCloud class="text-accent" />
      <div>
        <Typography variant="h3" as="h2">
          Резюме
        </Typography>
        <p class="text-sm text-muted-foreground">
          Загружайте файлы в формате PDF/DOC/DOCX. Поля заполнятся автоматически и их можно поправить.
        </p>
      </div>
    </div>

    <div class="space-y-4">
      <label class="block text-sm font-medium text-muted-foreground">Файл резюме</label>
      <input ref="fileInput" type="file" accept=".pdf,.doc,.docx" class="block w-full border border-input rounded-xl bg-transparent px-4 py-3 cursor-pointer text-sm" @change="onFileChange">

      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label class="block text-sm font-medium text-muted-foreground mb-1">Желаемая должность (необязательно)</label>
          <Input v-model="form.desiredPosition" placeholder="Product Manager" />
        </div>
        <div>
          <label class="block text-sm font-medium text-muted-foreground mb-1">Формат работы</label>
          <select v-model="form.workFormat" class="w-full border border-input rounded-xl bg-transparent px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring">
            <option v-for="option in workFormatOptions" :key="option.value" :value="option.value">
              {{ option.label }}
            </option>
          </select>
        </div>
      </div>

      <div>
        <label class="block text-sm font-medium text-muted-foreground mb-1">Опыт (необязательно)</label>
        <Textarea v-model="form.workExperience" placeholder="5 лет в разработке банковских приложений" rows="3" />
      </div>

      <Button class="w-full mt-2 gap-2" :disabled="isUploading" @click="handleUpload">
        <Loader2 v-if="isUploading" class="h-4 w-4 animate-spin" />
        <span v-else>Загрузить резюме</span>
      </Button>
    </div>

    <div class="space-y-4">
      <Typography variant="h4" as="h3" class="flex items-center gap-2">
        <FileText class="text-accent" />
        Мои резюме
      </Typography>

      <div v-if="!resumes.length" class="text-muted-foreground text-sm">
        Резюме пока не загружены.
      </div>

      <div v-else class="space-y-4">
        <div v-for="resume in resumes" :key="resume.id" class="border border-border rounded-xl p-4 bg-card shadow-sm">
          <div class="flex flex-col md:flex-row md:items-start md:justify-between gap-4">
            <div>
              <p class="text-sm text-muted-foreground">
                Файл
              </p>
              <p class="text-base font-medium">
                {{ resume.fileName }}
              </p>
            </div>
            <div class="flex gap-2">
              <Button variant="secondary" size="sm" class="gap-2" @click="startEdit(resume)">
                <Pencil class="h-4 w-4" />
              </Button>
              <Button variant="destructive" size="sm" class="gap-2" @click="deleteResume(resume)">
                <Trash2 class="h-4 w-4" />
              </Button>
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
                {{ new Date(resume.createdAt).toLocaleDateString() }}
              </dd>
            </div>
          </dl>
          <div class="mt-4 text-sm">
            <p class="text-muted-foreground">
              Опыт
            </p>
            <p class="text-foreground whitespace-pre-line">
              {{ resume.workExperience || 'Не указано' }}
            </p>
          </div>

          <div v-if="editingId === resume.id" class="mt-4 border-t pt-4 space-y-3">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <Input v-model="editForm.desiredPosition" placeholder="Желаемая должность" />
              <select v-model="editForm.workFormat" class="w-full border border-input rounded-xl bg-transparent px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring">
                <option v-for="option in workFormatOptions" :key="option.value" :value="option.value">
                  {{ option.label }}
                </option>
              </select>
            </div>
            <Textarea v-model="editForm.workExperience" rows="3" placeholder="Опыт работы" />
            <div class="flex gap-2 justify-end">
              <Button variant="ghost" @click="editingId = null">
                Отмена
              </Button>
              <Button class="gap-2" :disabled="isSavingEdit" @click="saveEdit(resume)">
                <Loader2 v-if="isSavingEdit" class="h-4 w-4 animate-spin" />
                <span v-else>Сохранить</span>
              </Button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
