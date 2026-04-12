<script setup lang="ts">
import type { TaskExchange, TaskExchangeStatus } from '@/models/taskExchange'
import {
  CheckCircle,
  ClipboardList,
  Clock,
  Edit3,
  Loader2,
  Plus,
  Search,
  Trash2,
  User,
  Users,
  XCircle,
} from 'lucide-vue-next'
import { computed, onMounted, ref, watch } from 'vue'
import EmptyState from '@/components/common/EmptyState.vue'
import ErrorState from '@/components/common/ErrorState.vue'
import FormField from '@/components/common/FormField.vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import TaskCardSkeleton from '@/components/tasks/TaskCardSkeleton.vue'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { useToast } from '@/components/ui/toast'
import { Typography } from '@/components/ui/typography'
import { required, useFormValidation } from '@/composables/useFormValidation'
import { useSSE } from '@/composables/useSSE'
import { isUserAdmin, useUser } from '@/composables/useUser'
import { handleError } from '@/services/errorService'
import { taskExchangeService } from '@/services/taskExchange'

const { toast } = useToast()

const tasks = ref<TaskExchange[]>([])
const total = ref(0)
const isLoading = ref(true)
const loadError = ref<string | null>(null)
const isSubmitting = ref(false)
const actionInProgress = ref<number | null>(null)
const showCreateDialog = ref(false)
const showEditDialog = ref(false)
const newTitle = ref('')
const newDescription = ref('')
const newMaxAssignees = ref(1)
const editTask = ref<TaskExchange | null>(null)
const editTitle = ref('')
const editDescription = ref('')
const editMaxAssignees = ref(1)
const activeStatus = ref<TaskExchangeStatus | 'all' | 'active'>('active')
const searchQuery = ref('')

const user = useUser()
const isAdmin = isUserAdmin()

const { errors: createErrors, validateAll: validateCreate, clearErrors: clearCreateErrors } = useFormValidation({
  title: [required('Введите название задания')],
})

const { errors: editErrors, validateAll: validateEdit, clearErrors: clearEditErrors } = useFormValidation({
  title: [required('Введите название задания')],
})

const statusTabs: { key: TaskExchangeStatus | 'all' | 'active', label: string }[] = [
  { key: 'active', label: 'Активные' },
  { key: 'all', label: 'Все' },
  { key: 'OPEN', label: 'Открытые' },
  { key: 'IN_PROGRESS', label: 'В работе' },
  { key: 'DONE', label: 'На проверке' },
  { key: 'APPROVED', label: 'Выполненные' },
]

const activeTasks = computed(() =>
  tasks.value.filter(t => t.status === 'OPEN' || t.status === 'IN_PROGRESS' || t.status === 'DONE'),
)

const statusConfig: Record<TaskExchangeStatus, { label: string, class: string }> = {
  OPEN: { label: 'Открыто', class: 'bg-blue-500/10 text-blue-500' },
  IN_PROGRESS: { label: 'В работе', class: 'bg-yellow-500/10 text-yellow-500' },
  DONE: { label: 'На проверке', class: 'bg-purple-500/10 text-purple-500' },
  APPROVED: { label: 'Выполнено', class: 'bg-green-500/10 text-green-500' },
}

const filteredTasks = computed(() => {
  const query = searchQuery.value.toLowerCase().trim()
  let result: TaskExchange[]
  if (activeStatus.value === 'all')
    result = tasks.value
  else if (activeStatus.value === 'active')
    result = activeTasks.value
  else
    result = tasks.value.filter(t => t.status === activeStatus.value)
  if (query) {
    result = result.filter(t =>
      t.title.toLowerCase().includes(query)
      || t.description?.toLowerCase().includes(query),
    )
  }
  return result
})

async function fetchTasks() {
  isLoading.value = true
  loadError.value = null
  try {
    const res = await taskExchangeService.getAll({ limit: 100 })
    tasks.value = res.items ?? []
    total.value = res.total
  }
  catch (error) {
    const appError = await handleError(error)
    loadError.value = appError.message
  }
  finally {
    isLoading.value = false
  }
}

async function createTask() {
  if (!validateCreate({ title: newTitle.value }))
    return
  isSubmitting.value = true
  try {
    await taskExchangeService.create({
      title: newTitle.value.trim(),
      description: newDescription.value.trim(),
      maxAssignees: newMaxAssignees.value,
    })
    toast({ title: 'Задание создано' })
    showCreateDialog.value = false
    await fetchTasks()
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isSubmitting.value = false
  }
}

function openEditDialog(task: TaskExchange) {
  editTask.value = task
  editTitle.value = task.title
  editDescription.value = task.description
  editMaxAssignees.value = task.maxAssignees
  showEditDialog.value = true
}

async function updateTask() {
  if (!editTask.value || !validateEdit({ title: editTitle.value }))
    return
  isSubmitting.value = true
  try {
    await taskExchangeService.update(editTask.value.id, {
      title: editTitle.value.trim(),
      description: editDescription.value.trim(),
      maxAssignees: editMaxAssignees.value,
    })
    toast({ title: 'Задание обновлено' })
    showEditDialog.value = false
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
  if (actionInProgress.value)
    return
  actionInProgress.value = id
  // Optimistic update
  const task = tasks.value.find(t => t.id === id)
  if (task && user.value) {
    task.assignees = [...(task.assignees ?? []), user.value]
    if (task.assignees.length >= task.maxAssignees)
      task.status = 'IN_PROGRESS'
  }
  try {
    await taskExchangeService.assign(id)
    await fetchTasks()
  }
  catch (error) {
    handleError(error)
    await fetchTasks()
  }
  finally {
    actionInProgress.value = null
  }
}

async function unassignTask(id: number) {
  if (actionInProgress.value)
    return
  actionInProgress.value = id
  // Optimistic update
  const task = tasks.value.find(t => t.id === id)
  if (task && user.value) {
    task.assignees = (task.assignees ?? []).filter(a => a.id !== user.value?.id)
    task.status = 'OPEN'
  }
  try {
    await taskExchangeService.unassign(id)
    await fetchTasks()
  }
  catch (error) {
    handleError(error)
    await fetchTasks()
  }
  finally {
    actionInProgress.value = null
  }
}

async function removeAssignee(taskId: number, memberId: number) {
  // Optimistic update
  const task = tasks.value.find(t => t.id === taskId)
  if (task) {
    task.assignees = (task.assignees ?? []).filter(a => a.id !== memberId)
    if (task.assignees.length === 0)
      task.status = 'OPEN'
  }
  try {
    await taskExchangeService.removeAssignee(taskId, memberId)
    await fetchTasks()
  }
  catch (error) {
    handleError(error)
    await fetchTasks()
  }
}

async function markDone(id: number) {
  if (actionInProgress.value)
    return
  actionInProgress.value = id
  const task = tasks.value.find(t => t.id === id)
  if (task)
    task.status = 'DONE'
  try {
    await taskExchangeService.markDone(id)
    await fetchTasks()
  }
  catch (error) {
    handleError(error)
    await fetchTasks()
  }
  finally {
    actionInProgress.value = null
  }
}

async function approveTask(id: number) {
  if (actionInProgress.value)
    return
  actionInProgress.value = id
  const task = tasks.value.find(t => t.id === id)
  if (task)
    task.status = 'APPROVED'
  try {
    await taskExchangeService.approve(id)
    await fetchTasks()
  }
  catch (error) {
    handleError(error)
    await fetchTasks()
  }
  finally {
    actionInProgress.value = null
  }
}

async function rejectTask(id: number) {
  if (actionInProgress.value)
    return
  actionInProgress.value = id
  const task = tasks.value.find(t => t.id === id)
  if (task)
    task.status = 'OPEN'
  try {
    await taskExchangeService.reject(id)
    await fetchTasks()
  }
  catch (error) {
    handleError(error)
    await fetchTasks()
  }
  finally {
    actionInProgress.value = null
  }
}

async function deleteTask(id: number) {
  // Optimistic update
  const snapshot = [...tasks.value]
  tasks.value = tasks.value.filter(t => t.id !== id)
  try {
    await taskExchangeService.remove(id)
  }
  catch (error) {
    tasks.value = snapshot
    handleError(error)
  }
}

function isCreator(task: TaskExchange) {
  return user.value?.id === task.creatorId
}

function isAssignee(task: TaskExchange) {
  return task.assignees?.some(a => a.id === user.value?.id) ?? false
}

function canTakeTask(task: TaskExchange) {
  return task.status === 'OPEN'
    && !isCreator(task)
    && !isAssignee(task)
    && (task.assignees?.length ?? 0) < task.maxAssignees
}

function canEdit(task: TaskExchange) {
  return (isCreator(task) || isAdmin.value) && (task.status === 'OPEN' || task.status === 'IN_PROGRESS')
}

function canMarkDone(task: TaskExchange) {
  return task.status === 'IN_PROGRESS' && (isCreator(task) || isAdmin.value)
}

function canRemoveAssignee(task: TaskExchange) {
  return isCreator(task) || isAdmin.value
}

function displayName(member: { firstName: string, lastName: string, tg: string }) {
  const name = [member.firstName, member.lastName].filter(Boolean).join(' ')
  return name || (member.tg ? `@${member.tg}` : 'Аноним')
}

useSSE('tasks', () => fetchTasks())

onMounted(() => {
  fetchTasks()
})

watch(showCreateDialog, (open) => {
  if (!open) {
    newTitle.value = ''
    newDescription.value = ''
    newMaxAssignees.value = 1
    clearCreateErrors()
  }
})

watch(showEditDialog, (open) => {
  if (!open) {
    editTask.value = null
    editTitle.value = ''
    editDescription.value = ''
    editMaxAssignees.value = 1
    clearEditErrors()
  }
})
</script>

<template>
  <div class="container mx-auto px-4 py-6 md:py-8">
    <div class="font-mono text-[11px] text-muted-foreground/60 tracking-wider mb-2">
      ~/community/task-exchange
    </div>
    <div class="flex items-center justify-between mb-6">
      <Typography
        variant="h2"
        as="h1"
      >
        Биржа заданий
      </Typography>
      <button
        class="flex items-center gap-2 px-3 sm:px-4 py-2 rounded-sm bg-primary text-primary-foreground text-sm font-medium hover:bg-primary/90 transition-colors shrink-0"
        @click="showCreateDialog = true"
      >
        <Plus class="h-4 w-4" />
        <span class="hidden sm:inline">Предложить задание</span>
      </button>
    </div>

    <div
      v-if="isLoading"
      class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4"
    >
      <TaskCardSkeleton v-for="i in 4" :key="i" />
    </div>

    <ErrorState
      v-else-if="loadError"
      :message="loadError"
      @retry="fetchTasks"
    />

    <template v-else>
      <div class="flex flex-col sm:flex-row sm:items-center gap-4 mb-6">
        <div class="flex gap-2 flex-wrap">
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
        <div class="relative sm:ml-auto">
          <Search class="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Поиск заданий..."
            class="w-full sm:w-64 rounded-sm border border-border bg-background pl-9 pr-3 py-1.5 text-sm focus:outline-none focus:ring-2 focus:ring-primary"
          >
        </div>
      </div>

      <EmptyState
        v-if="filteredTasks.length === 0"
        :icon="ClipboardList"
        title="Заданий пока нет"
        description="Предложите задание, полезное сообществу, и заработайте баллы"
        action-label="Предложить задание"
        @action="showCreateDialog = true"
      />

      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        <div
          v-for="task in filteredTasks"
          :key="task.id"
          class="rounded-sm p-4 border bg-card border-border terminal-card"
        >
          <div class="flex items-start justify-between gap-2 mb-2">
            <h3 class="font-medium text-sm leading-tight break-words min-w-0">
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
            <User class="h-3.5 w-3.5" aria-hidden="true" />
            <span>Автор: {{ displayName(task.creator) }}</span>
          </div>

          <!-- Assignees progress -->
          <div class="flex items-center gap-1.5 text-xs text-muted-foreground mb-1">
            <Users class="h-3.5 w-3.5" aria-hidden="true" />
            <span>Исполнители: {{ task.assignees?.length ?? 0 }}/{{ task.maxAssignees }}</span>
          </div>

          <!-- Assignees list -->
          <div
            v-if="task.assignees?.length"
            class="mt-1 space-y-1"
          >
            <div
              v-for="assignee in task.assignees"
              :key="assignee.id"
              class="flex items-center gap-1.5 text-xs text-muted-foreground"
            >
              <Clock class="h-3 w-3" aria-hidden="true" />
              <span class="truncate">{{ displayName(assignee) }}</span>
              <ConfirmDialog
                v-if="canRemoveAssignee(task) && task.status !== 'DONE' && task.status !== 'APPROVED'"
                title="Удалить исполнителя?"
                description="Исполнитель будет снят с задания."
                confirm-label="Удалить"
                @confirm="removeAssignee(task.id, assignee.id)"
              >
                <template #trigger>
                  <button
                    class="ml-auto shrink-0 p-0.5 rounded hover:bg-red-500/10 text-red-500 transition-colors"
                    aria-label="Удалить исполнителя"
                  >
                    <Trash2 class="h-3 w-3" />
                  </button>
                </template>
              </ConfirmDialog>
            </div>
          </div>

          <div class="flex flex-wrap gap-2 mt-3">
            <!-- OPEN: take task (not creator, not already assigned) -->
            <button
              v-if="canTakeTask(task)"
              class="px-3 py-1.5 rounded-lg text-xs font-medium bg-primary text-primary-foreground hover:bg-primary/90 transition-colors disabled:opacity-50"
              :disabled="actionInProgress === task.id"
              @click="assignTask(task.id)"
            >
              <Loader2 v-if="actionInProgress === task.id" class="h-3 w-3 animate-spin inline mr-1" />
              Взять задание
            </button>

            <!-- IN_PROGRESS: mark done (creator/admin) -->
            <button
              v-if="canMarkDone(task)"
              class="px-3 py-1.5 rounded-lg text-xs font-medium bg-green-600 text-white hover:bg-green-700 transition-colors disabled:opacity-50"
              :disabled="actionInProgress === task.id"
              @click="markDone(task.id)"
            >
              <Loader2 v-if="actionInProgress === task.id" class="h-3 w-3 animate-spin inline mr-1" />
              Выполнено
            </button>

            <!-- Unassign self (assignee, while OPEN or IN_PROGRESS) -->
            <button
              v-if="isAssignee(task) && (task.status === 'OPEN' || task.status === 'IN_PROGRESS')"
              class="px-3 py-1.5 rounded-lg text-xs font-medium bg-muted text-muted-foreground hover:text-foreground transition-colors disabled:opacity-50"
              :disabled="actionInProgress === task.id"
              @click="unassignTask(task.id)"
            >
              Отказаться
            </button>

            <!-- Edit (creator/admin, while OPEN or IN_PROGRESS) -->
            <button
              v-if="canEdit(task)"
              class="flex items-center gap-1 px-3 py-1.5 rounded-lg text-xs font-medium text-muted-foreground hover:text-foreground transition-colors"
              aria-label="Редактировать задание"
              @click="openEditDialog(task)"
            >
              <Edit3 class="h-3.5 w-3.5" />
            </button>

            <!-- DONE: approve/reject (admin) -->
            <button
              v-if="task.status === 'DONE' && isAdmin"
              class="flex items-center gap-1 px-3 py-1.5 rounded-lg text-xs font-medium bg-green-600 text-white hover:bg-green-700 transition-colors disabled:opacity-50"
              :disabled="actionInProgress === task.id"
              @click="approveTask(task.id)"
            >
              <CheckCircle class="h-3.5 w-3.5" aria-hidden="true" />
              Одобрить
            </button>
            <button
              v-if="task.status === 'DONE' && isAdmin"
              class="flex items-center gap-1 px-3 py-1.5 rounded-lg text-xs font-medium bg-red-600 text-white hover:bg-red-700 transition-colors disabled:opacity-50"
              :disabled="actionInProgress === task.id"
              @click="rejectTask(task.id)"
            >
              <XCircle class="h-3.5 w-3.5" aria-hidden="true" />
              Отклонить
            </button>

            <!-- Delete (creator if OPEN, or admin) -->
            <ConfirmDialog
              v-if="(task.status === 'OPEN' && isCreator(task)) || isAdmin"
              title="Удалить задание?"
              description="Это действие нельзя отменить. Задание будет удалено навсегда."
              confirm-label="Удалить"
              @confirm="deleteTask(task.id)"
            >
              <template #trigger>
                <button
                  class="flex items-center gap-1 px-3 py-1.5 rounded-lg text-xs font-medium text-red-500 hover:bg-red-500/10 transition-colors ml-auto"
                  aria-label="Удалить задание"
                >
                  <Trash2 class="h-3.5 w-3.5" />
                </button>
              </template>
            </ConfirmDialog>
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
          <FormField
            label="Название"
            :error="createErrors.title"
            html-for="create-title"
            required
          >
            <input
              id="create-title"
              v-model="newTitle"
              type="text"
              class="w-full rounded-sm border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary"
              :class="{ 'border-destructive': createErrors.title }"
              placeholder="Что нужно сделать?"
            >
          </FormField>
          <FormField
            label="Описание"
            html-for="create-desc"
          >
            <textarea
              id="create-desc"
              v-model="newDescription"
              class="w-full rounded-sm border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary min-h-24 resize-none"
              placeholder="Подробности задания..."
            />
          </FormField>
          <FormField
            label="Кол-во исполнителей"
            html-for="create-assignees"
          >
            <input
              id="create-assignees"
              v-model.number="newMaxAssignees"
              type="number"
              min="1"
              max="50"
              class="w-full rounded-sm border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary"
            >
          </FormField>

          <DialogFooter>
            <button
              type="submit"
              class="px-4 py-2 rounded-sm bg-primary text-primary-foreground text-sm font-medium hover:bg-primary/90 transition-colors disabled:opacity-50"
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

    <!-- Edit task dialog -->
    <Dialog
      v-model:open="showEditDialog"
    >
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Редактировать задание</DialogTitle>
          <DialogDescription>
            Измените параметры задания
          </DialogDescription>
        </DialogHeader>

        <form
          class="space-y-4"
          @submit.prevent="updateTask"
        >
          <FormField
            label="Название"
            :error="editErrors.title"
            html-for="edit-title"
            required
          >
            <input
              id="edit-title"
              v-model="editTitle"
              type="text"
              class="w-full rounded-sm border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary"
              :class="{ 'border-destructive': editErrors.title }"
            >
          </FormField>
          <FormField
            label="Описание"
            html-for="edit-desc"
          >
            <textarea
              id="edit-desc"
              v-model="editDescription"
              class="w-full rounded-sm border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary min-h-24 resize-none"
            />
          </FormField>
          <FormField
            label="Кол-во исполнителей"
            html-for="edit-assignees"
          >
            <input
              id="edit-assignees"
              v-model.number="editMaxAssignees"
              type="number"
              min="1"
              max="50"
              class="w-full rounded-sm border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary"
            >
          </FormField>

          <DialogFooter>
            <button
              type="submit"
              class="px-4 py-2 rounded-sm bg-primary text-primary-foreground text-sm font-medium hover:bg-primary/90 transition-colors disabled:opacity-50"
              :disabled="!editTitle.trim() || isSubmitting"
            >
              <Loader2
                v-if="isSubmitting"
                class="h-4 w-4 animate-spin inline mr-1"
              />
              Сохранить
            </button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  </div>
</template>
