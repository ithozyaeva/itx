<script setup lang="ts">
import type { ChallengeInstance, ChallengeKind, ChallengeTemplate, ChallengeTemplateRequest } from '@/services/challengeAdminService'
import { Calendar, Pencil, Plus, Trash2 } from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'
import AdminLayout from '@/components/layout/AdminLayout.vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Typography } from '@/components/ui/typography'
import { challengeAdminService } from '@/services/challengeAdminService'

const showModal = ref(false)
const editingId = ref<number | null>(null)
const confirmDeleteId = ref<number | null>(null)
const recentInstances = ref<ChallengeInstance[]>([])
const showInstances = ref(false)

const templateById = computed(() => {
  const map = new Map<number, ChallengeTemplate>()
  for (const t of challengeAdminService.items.value)
    map.set(t.id, t)
  return map
})

function formatInstanceDate(d: string) {
  return new Date(d).toLocaleDateString('ru-RU', { day: '2-digit', month: 'short' })
}

async function toggleInstances() {
  showInstances.value = !showInstances.value
  if (showInstances.value && recentInstances.value.length === 0)
    recentInstances.value = await challengeAdminService.recentInstances(30)
}

function emptyForm(): ChallengeTemplateRequest {
  return {
    code: '',
    title: '',
    description: '',
    icon: 'trophy',
    kind: 'weekly',
    metricKey: '',
    target: 5,
    rewardPoints: 200,
    achievementCode: null,
    active: true,
  }
}

const form = ref<ChallengeTemplateRequest>(emptyForm())

function openCreate() {
  editingId.value = null
  form.value = emptyForm()
  showModal.value = true
}

function openEdit(t: ChallengeTemplate) {
  editingId.value = t.id
  form.value = {
    code: t.code,
    title: t.title,
    description: t.description,
    icon: t.icon,
    kind: t.kind,
    metricKey: t.metricKey,
    target: t.target,
    rewardPoints: t.rewardPoints,
    achievementCode: t.achievementCode,
    active: t.active,
  }
  showModal.value = true
}

async function handleSubmit() {
  const ok = editingId.value !== null
    ? await challengeAdminService.update(editingId.value, form.value)
    : await challengeAdminService.create(form.value)
  if (ok) {
    showModal.value = false
    editingId.value = null
    form.value = emptyForm()
  }
}

async function handleDelete() {
  if (confirmDeleteId.value === null)
    return
  await challengeAdminService.delete(confirmDeleteId.value)
  confirmDeleteId.value = null
}

function kindBadgeClass(kind: ChallengeKind) {
  return kind === 'weekly'
    ? 'bg-blue-500/10 text-blue-500'
    : 'bg-purple-500/10 text-purple-500'
}

onMounted(() => {
  challengeAdminService.getAll()
})
</script>

<template>
  <AdminLayout>
    <div class="space-y-6">
      <div class="flex items-center justify-between">
        <div>
          <Typography variant="h2" as="h1">
            Челленджи
          </Typography>
          <p class="text-sm text-muted-foreground mt-1">
            Шаблоны еженедельных и ежемесячных челленджей. Cron каждый понедельник МСК выбирает 3 случайных weekly, 1-го числа — 1 monthly.
          </p>
        </div>
        <Button size="sm" @click="openCreate">
          <Plus class="h-4 w-4 mr-1" />
          Добавить шаблон
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
                    Тип
                  </th>
                  <th class="text-left py-3 px-4 font-medium">
                    metric_key
                  </th>
                  <th class="text-center py-3 px-4 font-medium">
                    Target
                  </th>
                  <th class="text-center py-3 px-4 font-medium">
                    Награда
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
                  v-for="t in challengeAdminService.items.value"
                  :key="t.id"
                  class="border-b last:border-0 hover:bg-muted/50"
                  :class="!t.active ? 'opacity-60' : ''"
                >
                  <td class="py-3 px-4 font-mono text-xs">
                    {{ t.code }}
                  </td>
                  <td class="py-3 px-4">
                    <div class="font-medium">
                      {{ t.title }}
                    </div>
                    <div v-if="t.description" class="text-xs text-muted-foreground mt-0.5">
                      {{ t.description }}
                    </div>
                  </td>
                  <td class="py-3 px-4 text-center">
                    <span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium" :class="kindBadgeClass(t.kind)">
                      {{ t.kind }}
                    </span>
                  </td>
                  <td class="py-3 px-4 font-mono text-xs text-muted-foreground">
                    {{ t.metricKey }}
                  </td>
                  <td class="py-3 px-4 text-center">
                    {{ t.target }}
                  </td>
                  <td class="py-3 px-4 text-center text-yellow-500 font-medium">
                    +{{ t.rewardPoints }}
                  </td>
                  <td class="py-3 px-4 text-center">
                    <span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium" :class="t.active ? 'bg-green-500/10 text-green-500' : 'bg-muted text-muted-foreground'">
                      {{ t.active ? 'Активен' : 'Выключен' }}
                    </span>
                  </td>
                  <td class="py-3 px-4 text-right space-x-1">
                    <Button variant="ghost" size="sm" @click="openEdit(t)">
                      <Pencil class="h-4 w-4" />
                    </Button>
                    <Button variant="ghost" size="sm" @click="confirmDeleteId = t.id">
                      <Trash2 class="h-4 w-4 text-destructive" />
                    </Button>
                  </td>
                </tr>
                <tr v-if="challengeAdminService.items.value.length === 0 && !challengeAdminService.isLoading.value">
                  <td colspan="8" class="py-8 text-center text-muted-foreground">
                    Шаблонов пока нет
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </CardContent>
      </Card>

      <!-- Recent instances audit -->
      <div>
        <button
          type="button"
          class="text-sm text-muted-foreground hover:text-foreground transition-colors flex items-center gap-1.5"
          @click="toggleInstances"
        >
          <Calendar class="h-4 w-4" />
          {{ showInstances ? 'Скрыть историю инстансов' : 'История запущенных челленджей (последние 30)' }}
        </button>
        <Card v-if="showInstances" class="mt-3">
          <CardContent class="p-0">
            <div class="overflow-x-auto">
              <table class="w-full text-sm">
                <thead>
                  <tr class="border-b">
                    <th class="text-left py-3 px-4 font-medium">
                      Period
                    </th>
                    <th class="text-left py-3 px-4 font-medium">
                      Шаблон
                    </th>
                    <th class="text-center py-3 px-4 font-medium">
                      Тип
                    </th>
                    <th class="text-left py-3 px-4 font-medium">
                      Период действия
                    </th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="inst in recentInstances" :key="inst.id" class="border-b last:border-0 hover:bg-muted/50">
                    <td class="py-3 px-4 font-mono text-xs">
                      {{ inst.periodKey }}
                    </td>
                    <td class="py-3 px-4">
                      {{ templateById.get(inst.templateId)?.title ?? `#${inst.templateId}` }}
                      <span class="text-xs text-muted-foreground ml-2 font-mono">
                        {{ templateById.get(inst.templateId)?.code }}
                      </span>
                    </td>
                    <td class="py-3 px-4 text-center text-xs text-muted-foreground">
                      {{ inst.kind }}
                    </td>
                    <td class="py-3 px-4 text-xs text-muted-foreground">
                      {{ formatInstanceDate(inst.startsAt) }} — {{ formatInstanceDate(inst.endsAt) }}
                    </td>
                  </tr>
                  <tr v-if="recentInstances.length === 0">
                    <td colspan="4" class="py-6 text-center text-muted-foreground">
                      Истории пока нет
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </CardContent>
        </Card>
      </div>

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
                Удалить шаблон? Запущенные инстансы (challenge_instances) не затрагиваются.
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
              <CardTitle>{{ editingId !== null ? 'Редактирование шаблона' : 'Новый шаблон' }}</CardTitle>
            </CardHeader>
            <CardContent>
              <form class="space-y-4" @submit.prevent="handleSubmit">
                <div class="grid grid-cols-2 gap-4">
                  <div>
                    <label class="text-sm font-medium mb-1 block">Code</label>
                    <input
                      v-model="form.code"
                      type="text"
                      required
                      :disabled="editingId !== null"
                      class="w-full border rounded-md px-3 py-2 text-sm bg-background font-mono disabled:opacity-60"
                      placeholder="w_event_guest"
                    >
                  </div>
                  <div>
                    <label class="text-sm font-medium mb-1 block">Icon</label>
                    <input
                      v-model="form.icon"
                      type="text"
                      class="w-full border rounded-md px-3 py-2 text-sm bg-background font-mono"
                      placeholder="trophy"
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
                    placeholder="Гость недели"
                  >
                </div>

                <div>
                  <label class="text-sm font-medium mb-1 block">Описание</label>
                  <textarea
                    v-model="form.description"
                    rows="2"
                    class="w-full border rounded-md px-3 py-2 text-sm bg-background resize-none"
                    placeholder="Посети 2 события за неделю"
                  />
                </div>

                <div class="grid grid-cols-3 gap-4">
                  <div>
                    <label class="text-sm font-medium mb-1 block">Тип</label>
                    <select
                      v-model="form.kind"
                      class="w-full border rounded-md px-3 py-2 text-sm bg-background"
                    >
                      <option value="weekly">
                        Weekly
                      </option>
                      <option value="monthly">
                        Monthly
                      </option>
                    </select>
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
                  <div>
                    <label class="text-sm font-medium mb-1 block">Награда</label>
                    <input
                      v-model.number="form.rewardPoints"
                      type="number"
                      min="0"
                      required
                      class="w-full border rounded-md px-3 py-2 text-sm bg-background"
                    >
                  </div>
                </div>

                <div>
                  <label class="text-sm font-medium mb-1 block">metric_key</label>
                  <input
                    v-model="form.metricKey"
                    type="text"
                    required
                    class="w-full border rounded-md px-3 py-2 text-sm bg-background font-mono"
                    placeholder="events_attended"
                  >
                  <p class="text-xs text-muted-foreground mt-1">
                    Должно совпадать со значением, передаваемым в TrackChallengeMetric() в коде хендлеров.
                  </p>
                </div>

                <div>
                  <label class="text-sm font-medium mb-1 block">Achievement code (опционально)</label>
                  <input
                    v-model="form.achievementCode"
                    type="text"
                    class="w-full border rounded-md px-3 py-2 text-sm bg-background font-mono"
                    placeholder="achievement_owner_month"
                  >
                </div>

                <label class="flex items-center gap-2 text-sm">
                  <input
                    v-model="form.active"
                    type="checkbox"
                    class="h-4 w-4"
                  >
                  Активен (попадает в выбор cron'а)
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
