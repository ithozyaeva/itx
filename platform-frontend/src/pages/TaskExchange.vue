<script setup lang="ts">
import type { TaskExchange, TaskExchangeStatus } from '@/models/taskExchange'
import { Typography } from 'itx-ui-kit'
import {
  CheckCircle,
  ClipboardList,
  Clock,
  Loader2,
  Plus,
  Trash2,
  User,
  XCircle,
} from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { isUserAdmin, useUser } from '@/composables/useUser'
import { handleError } from '@/services/errorService'
import { taskExchangeService } from '@/services/taskExchange'

const tasks = ref<TaskExchange[]>([])
const total = ref(0)
const isLoading = ref(true)
const isSubmitting = ref(false)
const showCreateDialog = ref(false)
const newTitle = ref('')
const newDescription = ref('')
const activeStatus = ref<TaskExchangeStatus | 'all'>('all')

const user = useUser()
const isAdmin = isUserAdmin()

const statusTabs: { key: TaskExchangeStatus | 'all', label: string }[] = [
  { key: 'all', label: 'Все' },
  { key: 'OPEN', label: 'Открытые' },
  { key: 'IN_PROGRESS', label: 'В работе' },
  { key: 'DONE', label: 'На проверке' },
  { key: 'APPROVED', label: 'Выполненные' },
]

const statusConfig: Record<TaskExchangeStatus, { label: string, class: string }> = {
  OPEN: { label: 'Открыто', class: 'bg-blue-500/10 text-blue-500' },
  IN_PROGRESS: { label: 'В работе', class: 'bg-yellow-500/10 text-yellow-500' },
  DONE: { label: 'На проверке', class: 'bg-purple-500/10 text-purple-500' },
  APPROVED: { label: 'Выполнено', class: 'bg-green-500/10 text-green-500' },
}

const filteredTasks = computed(() => {
  if (activeStatus.value === 'all')
    return tasks.value
  return tasks.value.filter(t => t.status === activeStatus.value)
})

async function fetchTasks() {
  isLoading.value = true
  try {
    const res = await taskExchangeService.getAll({ limit: 100 })
    tasks.value = res.items ?? []
    total.value = res.total
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isLoading.value = false
  }
}

async function createTask() {
  if (!newTitle.value.trim())
    return
  isSubmitting.value = true
  try {
    await taskExchangeService.create({
      title: newTitle.value.trim(),
      description: newDescription.value.trim(),
    })
    showCreateDialog.value = false
    newTitle.value = ''
    newDescription.value = ''
    await fetchTasks()
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isSubmitting.value = false
  }
}

async function assignTask(id: number) {
  try {
    await taskExchangeService.assign(id)
    await fetchTasks()
  }
  catch (error) {
    handleError(error)
  }
}

async function unassignTask(id: number) {
  try {
    await taskExchangeService.unassign(id)
    await fetchTasks()
  }
  catch (error) {
    handleError(error)
  }
}

async function markDone(id: number) {
  try {
    await taskExchangeService.markDone(id)
    await fetchTasks()
  }
  catch (error) {
    handleError(error)
  }
}

async function approveTask(id: number) {
  try {
    await taskExchangeService.approve(id)
    await fetchTasks()
  }
  catch (error) {
    handleError(error)
  }
}

async function rejectTask(id: number) {
  try {
    await taskExchangeService.reject(id)
    await fetchTasks()
  }
  catch (error) {
    handleError(error)
  }
}

async function deleteTask(id: number) {
  try {
    await taskExchangeService.remove(id)
    await fetchTasks()
  }
  catch (error) {
    handleError(error)
  }
}

function isCreator(task: TaskExchange) {
  return user.value?.id === task.creatorId
}

function isAssignee(task: TaskExchange) {
  return user.value?.id === task.assigneeId
}

function displayName(member: { firstName: string, lastName: string, tg: string }) {
  const name = [member.firstName, member.lastName].filter(Boolean).join(' ')
  return name || `@${member.tg}`
}

onMounted(() => {
  fetchTasks()
})
</script>

<template>
  <div class="container mx-auto px-4 py-6 md:py-8">
    <div class="flex items-center justify-between mb-6">
      <Typography
        variant="h2"
        as="h1"
      >
        Биржа заданий
      </Typography>
      <button
        class="flex items-center gap-2 px-4 py-2 rounded-xl bg-primary text-primary-foreground text-sm font-medium hover:bg-primary/90 transition-colors"
        @click="showCreateDialog = true"
      >
        <Plus class="h-4 w-4" />
        Предложить задание
      </button>
    </div>

    <div
      v-if="isLoading"
      class="flex justify-center py-12"
    >
      <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
    </div>

    <template v-else>
      <div class="flex gap-2 mb-6 flex-wrap">
        <button
          v-for="tab in statusTabs"
          :key="tab.key"
          class="px-3 py-1.5 rounded-full text-sm font-medium transition-colors"
          :class="activeStatus === tab.key
            ? 'bg-primary text-primary-foreground'
            : 'bg-card border border-border text-muted-foreground hover:text-foreground'"
          @click="activeStatus = tab.key"
        >
          {{ tab.label }}
        </button>
      </div>

      <div
        v-if="filteredTasks.length === 0"
        class="text-center py-12 text-muted-foreground"
      >
        <ClipboardList class="h-12 w-12 mx-auto mb-3 opacity-50" />
        <p>Заданий пока нет</p>
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        <div
          v-for="task in filteredTasks"
          :key="task.id"
          class="rounded-2xl p-4 border bg-card border-border"
        >
          <div class="flex items-start justify-between gap-2 mb-2">
            <h3 class="font-medium text-sm leading-tight">
              {{ task.title }}
            </h3>
            <span
              class="shrink-0 px-2 py-0.5 rounded-full text-xs font-medium"
              :class="statusConfig[task.status].class"
            >
              {{ statusConfig[task.status].label }}
            </span>
          </div>

          <p
            v-if="task.description"
            class="text-xs text-muted-foreground mb-3 line-clamp-3"
          >
            {{ task.description }}
          </p>

          <div class="flex items-center gap-1.5 text-xs text-muted-foreground mb-1">
            <User class="h-3.5 w-3.5" />
            <span>Автор: {{ displayName(task.creator) }}</span>
          </div>

          <div
            v-if="task.assignee"
            class="flex items-center gap-1.5 text-xs text-muted-foreground mb-1"
          >
            <Clock class="h-3.5 w-3.5" />
            <span>Исполнитель: {{ displayName(task.assignee) }}</span>
          </div>

          <div class="flex flex-wrap gap-2 mt-3">
            <!-- OPEN: take task (not creator) -->
            <button
              v-if="task.status === 'OPEN' && !isCreator(task)"
              class="px-3 py-1 rounded-lg text-xs font-medium bg-primary text-primary-foreground hover:bg-primary/90 transition-colors"
              @click="assignTask(task.id)"
            >
              Взять задание
            </button>

            <!-- IN_PROGRESS: mark done (assignee) -->
            <button
              v-if="task.status === 'IN_PROGRESS' && isAssignee(task)"
              class="px-3 py-1 rounded-lg text-xs font-medium bg-green-600 text-white hover:bg-green-700 transition-colors"
              @click="markDone(task.id)"
            >
              Выполнено
            </button>

            <!-- IN_PROGRESS: unassign (assignee) -->
            <button
              v-if="task.status === 'IN_PROGRESS' && isAssignee(task)"
              class="px-3 py-1 rounded-lg text-xs font-medium bg-muted text-muted-foreground hover:text-foreground transition-colors"
              @click="unassignTask(task.id)"
            >
              Отказаться
            </button>

            <!-- DONE: approve/reject (admin) -->
            <button
              v-if="task.status === 'DONE' && isAdmin"
              class="flex items-center gap-1 px-3 py-1 rounded-lg text-xs font-medium bg-green-600 text-white hover:bg-green-700 transition-colors"
              @click="approveTask(task.id)"
            >
              <CheckCircle class="h-3.5 w-3.5" />
              Одобрить
            </button>
            <button
              v-if="task.status === 'DONE' && isAdmin"
              class="flex items-center gap-1 px-3 py-1 rounded-lg text-xs font-medium bg-red-600 text-white hover:bg-red-700 transition-colors"
              @click="rejectTask(task.id)"
            >
              <XCircle class="h-3.5 w-3.5" />
              Отклонить
            </button>

            <!-- Delete (creator if OPEN, or admin) -->
            <button
              v-if="(task.status === 'OPEN' && isCreator(task)) || isAdmin"
              class="flex items-center gap-1 px-3 py-1 rounded-lg text-xs font-medium text-red-500 hover:bg-red-500/10 transition-colors ml-auto"
              @click="deleteTask(task.id)"
            >
              <Trash2 class="h-3.5 w-3.5" />
            </button>
          </div>
        </div>
      </div>
    </template>

    <!-- Create task dialog -->
    <Dialog
      v-model:open="showCreateDialog"
    >
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Предложить задание</DialogTitle>
          <DialogDescription>
            Опишите задание, полезное сообществу
          </DialogDescription>
        </DialogHeader>

        <form
          class="space-y-4"
          @submit.prevent="createTask"
        >
          <div>
            <label class="block text-sm font-medium mb-1">Название</label>
            <input
              v-model="newTitle"
              type="text"
              class="w-full rounded-xl border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary"
              placeholder="Что нужно сделать?"
            >
          </div>
          <div>
            <label class="block text-sm font-medium mb-1">Описание</label>
            <textarea
              v-model="newDescription"
              class="w-full rounded-xl border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary min-h-24 resize-none"
              placeholder="Подробности задания..."
            />
          </div>

          <DialogFooter>
            <button
              type="submit"
              class="px-4 py-2 rounded-xl bg-primary text-primary-foreground text-sm font-medium hover:bg-primary/90 transition-colors disabled:opacity-50"
              :disabled="!newTitle.trim() || isSubmitting"
            >
              <Loader2
                v-if="isSubmitting"
                class="h-4 w-4 animate-spin inline mr-1"
              />
              Создать
            </button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  </div>
</template>
