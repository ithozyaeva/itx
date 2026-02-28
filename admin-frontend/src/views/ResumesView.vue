<script setup lang="ts">
import type { ResumeFilter, WorkFormat } from '@/models/resume'
import { Typography } from 'itx-ui-kit'
import { onMounted, reactive, ref } from 'vue'
import Download from '~icons/lucide/download'
import RefreshCw from '~icons/lucide/refresh-cw'
import Search from '~icons/lucide/search'
import AdminLayout from '@/components/layout/AdminLayout.vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { resumeService } from '@/services/resumeService'

const filters = reactive<ResumeFilter>({
  workFormat: '',
  desiredPosition: '',
  workExperience: '',
})

const workFormatOptions: { label: string, value: WorkFormat }[] = [
  { label: 'Любой формат', value: '' },
  { label: 'Удалёнка', value: 'REMOTE' },
  { label: 'Гибрид', value: 'HYBRID' },
  { label: 'Офис', value: 'OFFICE' },
]

const isDownloading = ref(false)

async function applyFilters() {
  await resumeService.searchWithFilters(filters)
}

async function resetFilters() {
  filters.workFormat = ''
  filters.desiredPosition = ''
  filters.workExperience = ''
  resumeService.clearPagination()
  await resumeService.searchWithFilters()
}

async function downloadArchive() {
  isDownloading.value = true
  try {
    await resumeService.downloadArchive(filters)
  }
  finally {
    isDownloading.value = false
  }
}

function formatWorkFormat(value?: WorkFormat) {
  const match = workFormatOptions.find(option => option.value === value)
  return match?.label ?? 'Не указано'
}

onMounted(() => resumeService.searchWithFilters(filters))
</script>

<template>
  <AdminLayout>
    <div class="space-y-6">
      <div class="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div>
          <Typography variant="h2" as="h1">
            Резюме пользователей
          </Typography>
          <p class="text-muted-foreground">
            Просматривайте загруженные резюме, фильтруйте их и выгружайте архивом.
          </p>
        </div>
        <Button class="gap-2" :disabled="isDownloading" @click="downloadArchive">
          <RefreshCw v-if="isDownloading" class="h-4 w-4 animate-spin" />
          <template v-else>
            <Download class="h-4 w-4" />
            Скачать архив
          </template>
        </Button>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Фильтры</CardTitle>
          <CardDescription>Уточните выборку по формату работы, должности или опыту</CardDescription>
        </CardHeader>
        <CardContent>
          <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div class="space-y-2">
              <label class="text-sm font-medium text-muted-foreground">Формат работы</label>
              <select v-model="filters.workFormat" class="w-full border border-input rounded-xl bg-transparent px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring">
                <option v-for="option in workFormatOptions" :key="option.value" :value="option.value">
                  {{ option.label }}
                </option>
              </select>
            </div>
            <div class="space-y-2">
              <label class="text-sm font-medium text-muted-foreground">Желаемая должность</label>
              <Input v-model="filters.desiredPosition" placeholder="Product менеджер" />
            </div>
            <div class="space-y-2">
              <label class="text-sm font-medium text-muted-foreground">Ключевые слова опыта</label>
              <Input v-model="filters.workExperience" placeholder="финтех, маркетинг" />
            </div>
          </div>
          <div class="flex gap-3 mt-4">
            <Button class="gap-2" variant="secondary" @click="applyFilters">
              <Search class="h-4 w-4" />
              Применить
            </Button>
            <Button variant="ghost" @click="resetFilters">
              Сбросить
            </Button>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Пользователь</TableHead>
                <TableHead>Должность</TableHead>
                <TableHead>Формат</TableHead>
                <TableHead>Опыт</TableHead>
                <TableHead>Файл</TableHead>
                <TableHead>Дата</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow v-if="resumeService.items.value.total === 0">
                <TableCell colspan="6" class="text-center text-muted-foreground py-10">
                  Резюме не найдены
                </TableCell>
              </TableRow>
              <TableRow
                v-for="resume in resumeService.items.value.items"
                :key="resume.id"
                class="align-top"
              >
                <TableCell>
                  <div class="flex flex-col">
                    <span class="font-medium">{{ resume.member?.firstName }} {{ resume.member?.lastName }}</span>
                    <span class="text-muted-foreground text-sm">@{{ resume.member?.tg }}</span>
                  </div>
                </TableCell>
                <TableCell>{{ resume.desiredPosition || '—' }}</TableCell>
                <TableCell>{{ formatWorkFormat(resume.workFormat) }}</TableCell>
                <TableCell class="max-w-[200px] lg:max-w-[280px] whitespace-pre-line break-words text-sm text-muted-foreground">
                  {{ resume.workExperience || '—' }}
                </TableCell>
                <TableCell>{{ resume.fileName }}</TableCell>
                <TableCell>{{ new Date(resume.createdAt).toLocaleString() }}</TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  </AdminLayout>
</template>
