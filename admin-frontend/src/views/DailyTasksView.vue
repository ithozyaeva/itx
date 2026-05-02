<script setup lang="ts">
import type { DailyTask, DailyTaskCreateRequest, DailyTaskTier } from '@/services/dailyTaskService'
import { Pencil, Plus, Trash2 } from 'lucide-vue-next'
import { onMounted, ref } from 'vue'
import AdminLayout from '@/components/layout/AdminLayout.vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Typography } from '@/components/ui/typography'
import { dailyTaskService } from '@/services/dailyTaskService'

const showModal = ref(false)
const editingId = ref<number | null>(null)
const confirmDeleteId = ref<number | null>(null)

const tierOptions: { value: DailyTaskTier, label: string, points: number }[] = [
  { value: 'engagement', label: 'Engagement', points: 10 },
  { value: 'light', label: 'Light Action', points: 15 },
  { value: 'meaningful', label: 'Meaningful', points: 20 },
  { value: 'big', label: 'Big Action', points: 30 },
]

function emptyForm(): DailyTaskCreateRequest {
  return {
    code: '',
    title: '',
    description: '',
    icon: 'circle',
    tier: 'engagement',
    points: 10,
    target: 1,
    triggerKey: '',
    active: true,
  }
}

const form = ref<DailyTaskCreateRequest>(emptyForm())

function tierBadgeClass(tier: DailyTaskTier) {
  return {
    engagement: 'bg-blue-500/10 text-blue-500',
    light: 'bg-green-500/10 text-green-500',
    meaningful: 'bg-orange-500/10 text-orange-500',
    big: 'bg-purple-500/10 text-purple-500',
  }[tier]
}

function openCreate() {
  editingId.value = null
  form.value = emptyForm()
  showModal.value = true
}

function openEdit(task: DailyTask) {
  editingId.value = task.id
  form.value = {
    code: task.code,
    title: task.title,
    description: task.description,
    icon: task.icon,
    tier: task.tier,
    points: task.points,
    target: task.target,
    triggerKey: task.triggerKey,
    active: task.active,
  }
  showModal.value = true
}

async function handleSubmit() {
  const ok = editingId.value !== null
    ? await dailyTaskService.update(editingId.value, form.value)
    : await dailyTaskService.create(form.value)
  if (ok) {
    showModal.value = false
    editingId.value = null
    form.value = emptyForm()
  }
}

async function handleDelete() {
  if (confirmDeleteId.value === null)
    return
  await dailyTaskService.delete(confirmDeleteId.value)
  confirmDeleteId.value = null
}

onMounted(() => {
  dailyTaskService.getAll()
})
</script>

<template>
  <AdminLayout>
    <div class="space-y-6">
      <div class="flex items-center justify-between">
        <div>
          <Typography variant="h2" as="h1">
            Дейлики
          </Typography>
          <p class="text-sm text-muted-foreground mt-1">
            Пул задач, из которого крон выбирает 5 случайных на каждый МСК-день.
          </p>
        </div>
        <Button size="sm" @click="openCreate">
          <Plus class="h-4 w-4 mr-1" />
          Добавить задачу
        </Button>
      </div>

      <Card>
        <CardContent class="p-0">
          <div class="overflow-x-auto">
            <table class="w-full text-sm">
              <thead>
                <tr class="border-b">
                  <th class="text-left py-3 px-4 font-medium">
                    Code
                  </th>
                  <th class="text-left py-3 px-4 font-medium">
                    Название
                  </th>
                  <th class="text-center py-3 px-4 font-medium">
                    Tier
                  </th>
                  <th class="text-center py-3 px-4 font-medium">
                    Баллы
                  </th>
                  <th class="text-center py-3 px-4 font-medium">
                    Target
                  </th>
                  <th class="text-left py-3 px-4 font-medium">
                    trigger_key
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
                  v-for="task in dailyTaskService.items.value"
                  :key="task.id"
                  class="border-b last:border-0 hover:bg-muted/50"
                  :class="!task.active ? 'opacity-60' : ''"
                >
                  <td class="py-3 px-4 font-mono text-xs">
                    {{ task.code }}
                  </td>
                  <td class="py-3 px-4">
                    <div class="font-medium">
                      {{ task.title }}
                    </div>
                    <div v-if="task.description" class="text-xs text-muted-foreground mt-0.5">
                      {{ task.description }}
                    </div>
                  </td>
                  <td class="py-3 px-4 text-center">
                    <span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium" :class="tierBadgeClass(task.tier)">
                      {{ task.tier }}
                    </span>
                  </td>
                  <td class="py-3 px-4 text-center text-yellow-500 font-medium">
                    +{{ task.points }}
                  </td>
                  <td class="py-3 px-4 text-center">
                    {{ task.target }}
                  </td>
                  <td class="py-3 px-4 font-mono text-xs text-muted-foreground">
                    {{ task.triggerKey }}
                  </td>
                  <td class="py-3 px-4 text-center">
                    <span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium" :class="task.active ? 'bg-green-500/10 text-green-500' : 'bg-muted text-muted-foreground'">
                      {{ task.active ? 'Активна' : 'Выключена' }}
                    </span>
                  </td>
                  <td class="py-3 px-4 text-right space-x-1">
                    <Button variant="ghost" size="sm" @click="openEdit(task)">
                      <Pencil class="h-4 w-4" />
                    </Button>
                    <Button variant="ghost" size="sm" @click="confirmDeleteId = task.id">
                      <Trash2 class="h-4 w-4 text-destructive" />
                    </Button>
                  </td>
                </tr>
                <tr v-if="dailyTaskService.items.value.length === 0 && !dailyTaskService.isLoading.value">
                  <td colspan="8" class="py-8 text-center text-muted-foreground">
                    Задач пока нет
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </CardContent>
      </Card>

      <!-- Confirm delete -->
      <Teleport to="body">
        <div
          v-if="confirmDeleteId"
          class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
          @mousedown.self="confirmDeleteId = null"
        >
          <Card class="w-full max-w-sm">
            <CardContent class="p-6 space-y-4">
              <p class="text-sm">
                Удалить задачу из пула? Будущие наборы перестанут её получать. Прошлые сеты не затрагиваются.
              </p>
              <div class="flex justify-end gap-2">
                <Button variant="outline" size="sm" @click="confirmDeleteId = null">
                  Отмена
                </Button>
                <Button variant="destructive" size="sm" @click="handleDelete">
                  Удалить
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>
      </Teleport>

      <!-- Create/Edit modal -->
      <Teleport to="body">
        <div
          v-if="showModal"
          class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
          @mousedown.self="showModal = false"
        >
          <Card class="w-full max-w-lg max-h-[90vh] overflow-y-auto">
            <CardHeader>
              <CardTitle>{{ editingId !== null ? 'Редактирование задачи' : 'Новая задача' }}</CardTitle>
            </CardHeader>
            <CardContent>
              <form class="space-y-4" @submit.prevent="handleSubmit">
                <div class="grid grid-cols-2 gap-4">
                  <div>
                    <label class="text-sm font-medium mb-1 block">Code (machine key)</label>
                    <input
                      v-model="form.code"
                      type="text"
                      required
                      :disabled="editingId !== null"
                      class="w-full border rounded-md px-3 py-2 text-sm bg-background font-mono disabled:opacity-60"
                      placeholder="post_comment"
                    >
                  </div>
                  <div>
                    <label class="text-sm font-medium mb-1 block">Icon (Lucide name)</label>
                    <input
                      v-model="form.icon"
                      type="text"
                      class="w-full border rounded-md px-3 py-2 text-sm bg-background font-mono"
                      placeholder="message-square"
                    >
                  </div>
                </div>

                <div>
                  <label class="text-sm font-medium mb-1 block">Заголовок</label>
                  <input
                    v-model="form.title"
                    type="text"
                    required
                    class="w-full border rounded-md px-3 py-2 text-sm bg-background"
                    placeholder="Напиши комментарий"
                  >
                </div>

                <div>
                  <label class="text-sm font-medium mb-1 block">Описание</label>
                  <textarea
                    v-model="form.description"
                    rows="2"
                    class="w-full border rounded-md px-3 py-2 text-sm bg-background resize-none"
                    placeholder="Поделись мыслью под любым материалом"
                  />
                </div>

                <div class="grid grid-cols-3 gap-4">
                  <div>
                    <label class="text-sm font-medium mb-1 block">Tier</label>
                    <select
                      v-model="form.tier"
                      class="w-full border rounded-md px-3 py-2 text-sm bg-background"
                    >
                      <option v-for="opt in tierOptions" :key="opt.value" :value="opt.value">
                        {{ opt.label }} ({{ opt.points }})
                      </option>
                    </select>
                  </div>
                  <div>
                    <label class="text-sm font-medium mb-1 block">Баллы</label>
                    <input
                      v-model.number="form.points"
                      type="number"
                      min="0"
                      required
                      class="w-full border rounded-md px-3 py-2 text-sm bg-background"
                    >
                  </div>
                  <div>
                    <label class="text-sm font-medium mb-1 block">Target</label>
                    <input
                      v-model.number="form.target"
                      type="number"
                      min="1"
                      required
                      class="w-full border rounded-md px-3 py-2 text-sm bg-background"
                    >
                  </div>
                </div>

                <div>
                  <label class="text-sm font-medium mb-1 block">trigger_key</label>
                  <input
                    v-model="form.triggerKey"
                    type="text"
                    required
                    class="w-full border rounded-md px-3 py-2 text-sm bg-background font-mono"
                    placeholder="post_comment"
                  >
                  <p class="text-xs text-muted-foreground mt-1">
                    Имя триггера в коде. Должно совпадать со значением, передаваемым в TrackDailyTrigger().
                  </p>
                </div>

                <label class="flex items-center gap-2 text-sm">
                  <input
                    v-model="form.active"
                    type="checkbox"
                    class="h-4 w-4"
                  >
                  Активна (попадает в выбор cron'а)
                </label>

                <div class="flex justify-end gap-2 pt-2">
                  <Button type="button" variant="outline" @click="showModal = false">
                    Отмена
                  </Button>
                  <Button type="submit">
                    {{ editingId !== null ? 'Сохранить' : 'Создать' }}
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
