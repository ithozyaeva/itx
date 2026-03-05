<script setup lang="ts">
import type { SeasonCreateRequest } from '@/services/seasonService'
import { Typography } from 'itx-ui-kit'
import { Plus } from 'lucide-vue-next'
import { onMounted, ref } from 'vue'
import AdminLayout from '@/components/layout/AdminLayout.vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { seasonService } from '@/services/seasonService'

const showModal = ref(false)
const confirmFinishId = ref<number | null>(null)

const form = ref<SeasonCreateRequest>({
  title: '',
  startDate: '',
  endDate: '',
})

function resetForm() {
  form.value = { title: '', startDate: '', endDate: '' }
}

function openCreate() {
  resetForm()
  showModal.value = true
}

function formatDateForAPI(dateStr: string): string {
  if (!dateStr)
    return dateStr
  const d = new Date(dateStr)
  return d.toISOString()
}

async function handleSubmit() {
  const ok = await seasonService.create({
    title: form.value.title,
    startDate: formatDateForAPI(form.value.startDate),
    endDate: formatDateForAPI(form.value.endDate),
  })
  if (ok) {
    showModal.value = false
    resetForm()
  }
}

async function handleFinish() {
  if (!confirmFinishId.value)
    return
  await seasonService.finish(confirmFinishId.value)
  confirmFinishId.value = null
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString('ru-RU', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
  })
}

onMounted(() => {
  seasonService.getAll()
})
</script>

<template>
  <AdminLayout>
    <div class="space-y-6">
      <div class="flex items-center justify-between">
        <Typography
          variant="h2"
          as="h1"
        >
          Сезоны
        </Typography>
        <Button
          size="sm"
          @click="openCreate"
        >
          <Plus class="h-4 w-4 mr-1" />
          Создать сезон
        </Button>
      </div>

      <Card>
        <CardContent class="p-0">
          <div class="overflow-x-auto">
            <table class="w-full text-sm">
              <thead>
                <tr class="border-b">
                  <th class="text-left py-3 px-4 font-medium">
                    Название
                  </th>
                  <th class="text-left py-3 px-4 font-medium">
                    Период
                  </th>
                  <th class="text-center py-3 px-4 font-medium">
                    Статус
                  </th>
                  <th class="text-right py-3 px-4 font-medium">
                    Действия
                  </th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="season in seasonService.items.value"
                  :key="season.id"
                  class="border-b last:border-0 hover:bg-muted/50"
                >
                  <td class="py-3 px-4 font-medium">
                    {{ season.title }}
                  </td>
                  <td class="py-3 px-4 text-muted-foreground">
                    {{ formatDate(season.startDate) }} — {{ formatDate(season.endDate) }}
                  </td>
                  <td class="py-3 px-4 text-center">
                    <span
                      class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium"
                      :class="season.status === 'ACTIVE'
                        ? 'bg-green-500/10 text-green-500'
                        : 'bg-muted text-muted-foreground'"
                    >
                      {{ season.status === 'ACTIVE' ? 'Активный' : 'Завершён' }}
                    </span>
                  </td>
                  <td class="py-3 px-4 text-right">
                    <Button
                      v-if="season.status === 'ACTIVE'"
                      variant="outline"
                      size="sm"
                      @click="confirmFinishId = season.id"
                    >
                      Завершить
                    </Button>
                  </td>
                </tr>
                <tr v-if="seasonService.items.value.length === 0 && !seasonService.isLoading.value">
                  <td
                    colspan="4"
                    class="py-8 text-center text-muted-foreground"
                  >
                    Сезонов пока нет
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </CardContent>
      </Card>

      <!-- Confirm finish dialog -->
      <Teleport to="body">
        <div
          v-if="confirmFinishId"
          class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
          @click.self="confirmFinishId = null"
        >
          <Card class="w-full max-w-sm">
            <CardContent class="p-6 space-y-4">
              <p class="text-sm">
                Завершить сезон? Это действие нельзя отменить.
              </p>
              <div class="flex justify-end gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  @click="confirmFinishId = null"
                >
                  Отмена
                </Button>
                <Button
                  variant="destructive"
                  size="sm"
                  @click="handleFinish"
                >
                  Завершить
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>
      </Teleport>

      <!-- Create modal -->
      <Teleport to="body">
        <div
          v-if="showModal"
          class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
          @click.self="showModal = false"
        >
          <Card class="w-full max-w-lg">
            <CardHeader>
              <CardTitle>Новый сезон</CardTitle>
            </CardHeader>
            <CardContent>
              <form
                class="space-y-4"
                @submit.prevent="handleSubmit"
              >
                <div>
                  <label class="text-sm font-medium mb-1 block">Название</label>
                  <input
                    v-model="form.title"
                    type="text"
                    required
                    class="w-full border rounded-md px-3 py-2 text-sm bg-background"
                    placeholder="Сезон 1"
                  >
                </div>

                <div class="grid grid-cols-2 gap-4">
                  <div>
                    <label class="text-sm font-medium mb-1 block">Начало</label>
                    <input
                      v-model="form.startDate"
                      type="datetime-local"
                      required
                      class="w-full border rounded-md px-3 py-2 text-sm bg-background"
                    >
                  </div>
                  <div>
                    <label class="text-sm font-medium mb-1 block">Окончание</label>
                    <input
                      v-model="form.endDate"
                      type="datetime-local"
                      required
                      class="w-full border rounded-md px-3 py-2 text-sm bg-background"
                    >
                  </div>
                </div>

                <div class="flex justify-end gap-2 pt-2">
                  <Button
                    type="button"
                    variant="outline"
                    @click="showModal = false"
                  >
                    Отмена
                  </Button>
                  <Button type="submit">
                    Создать
                  </Button>
                </div>
              </form>
            </CardContent>
          </Card>
        </div>
      </Teleport>
    </div>
  </AdminLayout>
</template>
