<script setup lang="ts">
import type { AIMaterialComment } from '@/models/aiMaterial'
import { Loader2, MessageCircle, Pencil, Send, Trash2, X } from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import { useToast } from '@/components/ui/toast'
import { isUserAdmin, useUser } from '@/composables/useUser'
import { aiMaterialsService } from '@/services/aiMaterials'
import { handleError } from '@/services/errorService'

const props = defineProps<{
  materialId: number
  initialCount: number
}>()

const emit = defineEmits<{
  'update:count': [v: number]
}>()

const user = useUser()
const isAdmin = isUserAdmin()
const { toast } = useToast()

const comments = ref<AIMaterialComment[]>([])
const isLoading = ref(true)
const loadError = ref<string | null>(null)

const newBody = ref('')
const submitting = ref(false)

const editingId = ref<number | null>(null)
const editingBody = ref('')
const editing = ref(false)

const COMMENT_MAX = 4_000

function emitCount() {
  emit('update:count', comments.value.length)
}

async function fetchComments() {
  isLoading.value = true
  loadError.value = null
  try {
    const res = await aiMaterialsService.listComments(props.materialId)
    comments.value = res.items ?? []
    emitCount()
  }
  catch (error) {
    loadError.value = (await handleError(error)).message
  }
  finally {
    isLoading.value = false
  }
}

async function postComment() {
  const body = newBody.value.trim()
  if (!body || submitting.value)
    return
  submitting.value = true
  try {
    const created = await aiMaterialsService.createComment(props.materialId, body)
    comments.value.push(created)
    newBody.value = ''
    emitCount()
  }
  catch (error) {
    handleError(error)
  }
  finally {
    submitting.value = false
  }
}

function startEdit(c: AIMaterialComment) {
  editingId.value = c.id
  editingBody.value = c.body
}

function cancelEdit() {
  editingId.value = null
  editingBody.value = ''
}

async function saveEdit() {
  if (editingId.value == null || editing.value)
    return
  const body = editingBody.value.trim()
  if (!body)
    return
  editing.value = true
  try {
    const updated = await aiMaterialsService.updateComment(editingId.value, body)
    const idx = comments.value.findIndex(c => c.id === updated.id)
    if (idx !== -1)
      comments.value[idx] = updated
    cancelEdit()
  }
  catch (error) {
    handleError(error)
  }
  finally {
    editing.value = false
  }
}

async function deleteComment(c: AIMaterialComment) {
  try {
    await aiMaterialsService.deleteComment(c.id)
    comments.value = comments.value.filter(it => it.id !== c.id)
    toast({ title: 'Комментарий удалён' })
    emitCount()
  }
  catch (error) {
    handleError(error)
  }
}

function canManage(c: AIMaterialComment) {
  return isAdmin.value || user.value?.id === c.authorId
}

function authorName(c: AIMaterialComment): string {
  const a = c.author
  if (!a)
    return 'Аноним'
  const name = [a.firstName, a.lastName].filter(Boolean).join(' ')
  return name || (a.tg ? `@${a.tg}` : 'Аноним')
}

function formatDate(iso: string): string {
  const d = new Date(iso)
  const date = d.toLocaleDateString('ru-RU', { day: 'numeric', month: 'long', year: 'numeric' })
  const time = d.toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit' })
  return `${date}, ${time}`
}

const visibleCount = computed(() => comments.value.length || props.initialCount)

onMounted(fetchComments)

defineExpose({ refresh: fetchComments })
</script>

<template>
  <section class="space-y-4">
    <header class="flex items-center gap-2">
      <MessageCircle class="h-5 w-5 text-muted-foreground" />
      <h2 class="font-medium">
        Комментарии
        <span class="text-muted-foreground font-normal">({{ visibleCount }})</span>
      </h2>
    </header>

    <form
      class="rounded-sm border border-border bg-card p-3 space-y-2"
      @submit.prevent="postComment"
    >
      <textarea
        v-model="newBody"
        :maxlength="COMMENT_MAX"
        class="w-full bg-transparent text-sm focus:outline-none resize-none min-h-16"
        placeholder="Поделитесь впечатлением, спросите автора, оставьте свой опыт..."
      />
      <div class="flex items-center justify-between">
        <span class="text-xs text-muted-foreground">{{ newBody.length }} / {{ COMMENT_MAX }}</span>
        <button
          type="submit"
          class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-sm bg-primary text-primary-foreground text-sm font-medium hover:bg-primary/90 transition-colors disabled:opacity-50"
          :disabled="!newBody.trim() || submitting"
        >
          <Loader2 v-if="submitting" class="h-3.5 w-3.5 animate-spin" />
          <Send v-else class="h-3.5 w-3.5" />
          Отправить
        </button>
      </div>
    </form>

    <div v-if="isLoading" class="space-y-3">
      <div v-for="i in 2" :key="i" class="rounded-sm border border-border bg-card p-3">
        <div class="h-3 w-1/3 rounded bg-muted animate-pulse mb-2" />
        <div class="h-3 w-full rounded bg-muted animate-pulse mb-1" />
        <div class="h-3 w-2/3 rounded bg-muted animate-pulse" />
      </div>
    </div>

    <p v-else-if="loadError" class="text-sm text-destructive">
      {{ loadError }}
    </p>

    <p v-else-if="comments.length === 0" class="text-sm text-muted-foreground">
      Пока нет комментариев. Будьте первым.
    </p>

    <ul v-else class="space-y-3">
      <li
        v-for="c in comments"
        :key="c.id"
        class="rounded-sm border border-border bg-card p-3"
      >
        <header class="flex items-center justify-between gap-2 mb-1.5 text-xs text-muted-foreground">
          <div class="flex items-center gap-2 min-w-0">
            <span class="font-medium text-foreground truncate">{{ authorName(c) }}</span>
            <span class="truncate">{{ formatDate(c.createdAt) }}</span>
            <span
              v-if="c.updatedAt && c.updatedAt !== c.createdAt"
              class="italic"
            >· изменено</span>
            <span
              v-if="c.isHidden"
              class="px-1.5 py-0.5 rounded-full bg-yellow-500/15 text-yellow-600"
            >Скрыт</span>
          </div>
          <div v-if="canManage(c) && editingId !== c.id" class="flex items-center gap-1 shrink-0">
            <button
              type="button"
              aria-label="Редактировать"
              class="p-1 rounded hover:bg-muted text-muted-foreground hover:text-foreground"
              @click="startEdit(c)"
            >
              <Pencil class="h-3.5 w-3.5" />
            </button>
            <ConfirmDialog
              title="Удалить комментарий?"
              description="Удаление безвозвратно."
              confirm-label="Удалить"
              @confirm="deleteComment(c)"
            >
              <template #trigger>
                <button
                  type="button"
                  aria-label="Удалить"
                  class="p-1 rounded hover:bg-red-500/10 text-muted-foreground hover:text-red-500"
                >
                  <Trash2 class="h-3.5 w-3.5" />
                </button>
              </template>
            </ConfirmDialog>
          </div>
        </header>

        <div v-if="editingId === c.id" class="space-y-2">
          <textarea
            v-model="editingBody"
            :maxlength="COMMENT_MAX"
            class="w-full rounded-sm border border-border bg-background px-2 py-1.5 text-sm focus:outline-none focus:ring-2 focus:ring-primary min-h-16 resize-none"
          />
          <div class="flex items-center gap-2">
            <button
              type="button"
              class="inline-flex items-center gap-1 px-2.5 py-1 rounded-sm bg-primary text-primary-foreground text-xs font-medium hover:bg-primary/90 transition-colors disabled:opacity-50"
              :disabled="!editingBody.trim() || editing"
              @click="saveEdit"
            >
              <Loader2 v-if="editing" class="h-3 w-3 animate-spin" />
              Сохранить
            </button>
            <button
              type="button"
              class="inline-flex items-center gap-1 px-2.5 py-1 rounded-sm border border-border text-xs text-muted-foreground hover:text-foreground transition-colors"
              @click="cancelEdit"
            >
              <X class="h-3 w-3" />
              Отмена
            </button>
          </div>
        </div>
        <p
          v-else
          class="text-sm whitespace-pre-wrap break-words"
        >
          {{ c.body }}
        </p>
      </li>
    </ul>
  </section>
</template>
