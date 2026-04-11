<script setup lang="ts">
import type { TrackedChat } from '@/services/chatActivityService'
import type { ChatQuest, ChatQuestCreateRequest } from '@/services/chatQuestService'
import { Pencil, Plus, Trash2 } from 'lucide-vue-next'
import { onMounted, ref } from 'vue'
import AdminLayout from '@/components/layout/AdminLayout.vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Typography } from '@/components/ui/typography'
import { useCardReveal } from '@/composables/useCardReveal'
import { chatActivityService } from '@/services/chatActivityService'
import { chatQuestService } from '@/services/chatQuestService'

const containerRef = ref<HTMLElement | null>(null)
useCardReveal(containerRef)

const quests = ref<ChatQuest[]>([])
const total = ref(0)
const chats = ref<TrackedChat[]>([])
const isLoading = ref(true)
const showModal = ref(false)
const editingQuest = ref<ChatQuest | null>(null)
const deleteId = ref<number | null>(null)

const form = ref<ChatQuestCreateRequest>({
  title: '',
  description: '',
  questType: 'message_count',
  chatId: null,
  targetCount: 10,
  pointsReward: 10,
  startsAt: '',
  endsAt: '',
  isActive: true,
})

function resetForm() {
  form.value = {
    title: '',
    description: '',
    questType: 'message_count',
    chatId: null,
    targetCount: 10,
    pointsReward: 10,
    startsAt: '',
    endsAt: '',
    isActive: true,
  }
  editingQuest.value = null
}

function openCreate() {
  resetForm()
  showModal.value = true
}

function openEdit(quest: ChatQuest) {
  editingQuest.value = quest
  form.value = {
    title: quest.title,
    description: quest.description,
    questType: quest.questType,
    chatId: quest.chatId,
    targetCount: quest.targetCount,
    pointsReward: quest.pointsReward,
    startsAt: quest.startsAt.slice(0, 16),
    endsAt: quest.endsAt.slice(0, 16),
    isActive: quest.isActive,
  }
  showModal.value = true
}

async function loadQuests() {
  try {
    const res = await chatQuestService.getAll(50, 0)
    quests.value = res.items
    total.value = res.total
  }
  catch (error) {
    console.error('Ошибка загрузки заданий:', error)
  }
}

function formatDateForAPI(dateStr: string): string {
  if (!dateStr)
    return dateStr
  // datetime-local даёт "2026-03-04T17:00", Go ожидает RFC3339
  const d = new Date(dateStr)
  return d.toISOString()
}

async function handleSubmit() {
  try {
    const payload = {
      ...form.value,
      startsAt: formatDateForAPI(form.value.startsAt),
      endsAt: formatDateForAPI(form.value.endsAt),
    }
    if (editingQuest.value) {
      await chatQuestService.update(editingQuest.value.id, payload)
    }
    else {
      await chatQuestService.create(payload)
    }
    showModal.value = false
    resetForm()
    await loadQuests()
  }
  catch (error) {
    console.error('Ошибка сохранения задания:', error)
  }
}

function confirmDelete(id: number) {
  deleteId.value = id
}

async function handleDelete() {
  if (!deleteId.value)
    return

  try {
    await chatQuestService.remove(deleteId.value)
    deleteId.value = null
    await loadQuests()
  }
  catch (error) {
    console.error('Ошибка удаления задания:', error)
  }
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString('ru-RU', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

function isQuestActive(quest: ChatQuest) {
  const now = new Date()
  return quest.isActive && new Date(quest.startsAt) <= now && new Date(quest.endsAt) >= now
}

onMounted(async () => {
  try {
    const [, ch] = await Promise.all([
      loadQuests(),
      chatActivityService.getChats(),
    ])
    chats.value = ch
  }
  catch (error) {
    console.error('Ошибка загрузки:', error)
  }
  finally {
    isLoading.value = false
  }
})
</script>

<template>
  <AdminLayout>
    <div ref="containerRef" class="space-y-6">
      <div class="flex items-center justify-between">
        <Typography variant="h2" as="h1">
          Задания чатов
        </Typography>
        <Button size="sm" @click="openCreate">
          <Plus class="h-4 w-4 mr-1" />
          Создать
        </Button>
      </div>

      <!-- Таблица квестов -->
      <Card data-reveal>
        <CardContent class="p-0">
          <div class="overflow-x-auto">
            <table class="w-full text-sm">
              <thead>
                <tr class="border-b">
                  <th class="text-left py-3 px-4 font-medium">
                    Название
                  </th>
                  <th class="text-left py-3 px-4 font-medium hidden md:table-cell">
                    Чат
                  </th>
                  <th class="text-center py-3 px-4 font-medium">
                    Цель
                  </th>
                  <th class="text-center py-3 px-4 font-medium">
                    Награда
                  </th>
                  <th class="text-left py-3 px-4 font-medium hidden lg:table-cell">
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
                  v-for="quest in quests"
                  :key="quest.id"
                  class="border-b last:border-0 hover:bg-muted/50"
                >
                  <td class="py-3 px-4">
                    <div class="font-medium">
                      {{ quest.title }}
                    </div>
                    <div v-if="quest.description" class="text-xs text-muted-foreground mt-0.5 truncate max-w-xs">
                      {{ quest.description }}
                    </div>
                  </td>
                  <td class="py-3 px-4 hidden md:table-cell text-muted-foreground">
                    {{ quest.chatId ? chats.find(c => c.chatId === quest.chatId)?.title || `#${quest.chatId}` : 'Любой чат' }}
                  </td>
                  <td class="py-3 px-4 text-center">
                    {{ quest.targetCount }}
                  </td>
                  <td class="py-3 px-4 text-center font-bold text-yellow-500">
                    +{{ quest.pointsReward }}
                  </td>
                  <td class="py-3 px-4 hidden lg:table-cell text-xs text-muted-foreground">
                    {{ formatDate(quest.startsAt) }} — {{ formatDate(quest.endsAt) }}
                  </td>
                  <td class="py-3 px-4 text-center">
                    <span
                      class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium"
                      :class="isQuestActive(quest)
                        ? 'bg-green-500/10 text-green-500'
                        : 'bg-muted text-muted-foreground'"
                    >
                      {{ isQuestActive(quest) ? 'Активно' : 'Неактивно' }}
                    </span>
                  </td>
                  <td class="py-3 px-4 text-right">
                    <div class="flex items-center justify-end gap-1">
                      <button
                        class="p-1.5 rounded-md hover:bg-muted"
                        @click="openEdit(quest)"
                      >
                        <Pencil class="h-4 w-4" />
                      </button>
                      <button
                        class="p-1.5 rounded-md hover:bg-destructive/10 text-destructive"
                        @click="confirmDelete(quest.id)"
                      >
                        <Trash2 class="h-4 w-4" />
                      </button>
                    </div>
                  </td>
                </tr>
                <tr v-if="quests.length === 0 && !isLoading">
                  <td colspan="7" class="py-8 text-center text-muted-foreground">
                    Заданий пока нет
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </CardContent>
      </Card>

      <!-- Подтверждение удаления -->
      <Teleport to="body">
        <div
          v-if="deleteId"
          class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
          @mousedown.self="deleteId = null"
        >
          <Card class="w-full max-w-sm">
            <CardContent class="p-6 space-y-4">
              <p class="text-sm">
                Удалить задание?
              </p>
              <div class="flex justify-end gap-2">
                <Button variant="outline" size="sm" @click="deleteId = null">
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

      <!-- Модалка создания/редактирования -->
      <Teleport to="body">
        <div
          v-if="showModal"
          class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
          @mousedown.self="showModal = false"
        >
          <Card class="w-full max-w-lg">
            <CardHeader>
              <CardTitle>{{ editingQuest ? 'Редактирование' : 'Новое задание' }}</CardTitle>
            </CardHeader>
            <CardContent>
              <form class="space-y-4" @submit.prevent="handleSubmit">
                <div>
                  <label class="text-sm font-medium mb-1 block">Название</label>
                  <input
                    v-model="form.title"
                    type="text"
                    required
                    class="w-full border rounded-md px-3 py-2 text-sm bg-background"
                    placeholder="Напиши 100 сообщений"
                  >
                </div>

                <div>
                  <label class="text-sm font-medium mb-1 block">Описание</label>
                  <textarea
                    v-model="form.description"
                    rows="2"
                    class="w-full border rounded-md px-3 py-2 text-sm bg-background"
                    placeholder="Подробное описание задания..."
                  />
                </div>

                <div>
                  <label class="text-sm font-medium mb-1 block">Тип задания</label>
                  <Select v-model="form.questType">
                    <SelectTrigger>
                      <SelectValue placeholder="Тип задания" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="message_count">
                        Количество сообщений
                      </SelectItem>
                      <SelectItem value="daily_streak">
                        Дни подряд
                      </SelectItem>
                    </SelectContent>
                  </Select>
                </div>

                <div class="grid grid-cols-2 gap-4">
                  <div>
                    <label class="text-sm font-medium mb-1 block">Чат</label>
                    <Select
                      :model-value="form.chatId ? String(form.chatId) : ''"
                      @update:model-value="form.chatId = $event ? Number($event) : null"
                    >
                      <SelectTrigger>
                        <SelectValue placeholder="Любой чат" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="">
                          Любой чат
                        </SelectItem>
                        <SelectItem v-for="chat in chats" :key="chat.chatId" :value="String(chat.chatId)">
                          {{ chat.title }}
                        </SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                  <div>
                    <label class="text-sm font-medium mb-1 block">{{ form.questType === 'daily_streak' ? 'Целевое кол-во дней' : 'Целевое кол-во сообщений' }}</label>
                    <input
                      v-model.number="form.targetCount"
                      type="number"
                      min="1"
                      required
                      class="w-full border rounded-md px-3 py-2 text-sm bg-background"
                    >
                  </div>
                </div>

                <div class="grid grid-cols-2 gap-4">
                  <div>
                    <label class="text-sm font-medium mb-1 block">Награда (баллы)</label>
                    <input
                      v-model.number="form.pointsReward"
                      type="number"
                      min="1"
                      required
                      class="w-full border rounded-md px-3 py-2 text-sm bg-background"
                    >
                  </div>
                  <div class="flex items-end">
                    <label class="flex items-center gap-2 text-sm">
                      <input
                        v-model="form.isActive"
                        type="checkbox"
                        class="rounded"
                      >
                      Активно
                    </label>
                  </div>
                </div>

                <div class="grid grid-cols-2 gap-4">
                  <div>
                    <label class="text-sm font-medium mb-1 block">Начало</label>
                    <input
                      v-model="form.startsAt"
                      type="datetime-local"
                      required
                      class="w-full border rounded-md px-3 py-2 text-sm bg-background"
                    >
                  </div>
                  <div>
                    <label class="text-sm font-medium mb-1 block">Окончание</label>
                    <input
                      v-model="form.endsAt"
                      type="datetime-local"
                      required
                      class="w-full border rounded-md px-3 py-2 text-sm bg-background"
                    >
                  </div>
                </div>

                <div class="flex justify-end gap-2 pt-2">
                  <Button type="button" variant="outline" @click="showModal = false">
                    Отмена
                  </Button>
                  <Button type="submit">
                    {{ editingQuest ? 'Сохранить' : 'Создать' }}
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
