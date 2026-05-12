<script setup lang="ts">
import type { Comment, CommentEntityType } from '@/models/comment'
import { Heart, Loader2, MessageCircle, Pencil, Send, Trash2, X } from 'lucide-vue-next'
import { computed, onMounted, ref, watch } from 'vue'
import { RouterLink } from 'vue-router'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import { useToast } from '@/components/ui/toast'
import { isUserAdmin, useUser } from '@/composables/useUser'
import { COMMENT_MAX_LEN } from '@/models/comment'
import { commentsService } from '@/services/comments'
import { handleError } from '@/services/errorService'

// withDefaults обязателен для autoLoad: type Boolean без явного default
// в Vue 3 трактуется как false, а не undefined. Без него `props.autoLoad`
// всегда `false`, и fetchComments никогда не запускался — на проде
// Comments монтировался, но GET к /comments не уходил, после refresh
// список оставался пустым (только что отправленный коммент исчезал).
const props = withDefaults(defineProps<{
  entityType: CommentEntityType
  entityId: number
  initialCount: number
  // autoLoad=false подавляет fetch на mount — для раскрывающихся
  // блоков типа EventCard accordion, чтобы не дёргать API на свёрнутых.
  autoLoad?: boolean
}>(), {
  autoLoad: true,
})

const emit = defineEmits<{
  'update:count': [v: number]
}>()

const user = useUser()
const isAdmin = isUserAdmin()
const { toast } = useToast()

const PAGE_SIZE = 20

const comments = ref<Comment[]>([])
const total = ref(0)
const loaded = ref(false)
const isLoading = ref(false)
const isLoadingMore = ref(false)
const loadError = ref<string | null>(null)

const newBody = ref('')
const submitting = ref(false)

const editingId = ref<number | null>(null)
const editingBody = ref('')
const editing = ref(false)

const likeBusy = ref<Set<number>>(new Set())

const hasMore = computed(() => comments.value.length < total.value)
// До первой загрузки показываем initialCount (с карточки), после — actual total.
// Раньше fallback `total || initialCount` врал когда total реально 0.
const visibleCount = computed(() => loaded.value ? total.value : props.initialCount)

function syncCountUp() {
  emit('update:count', total.value)
}

async function fetchComments(append = false) {
  if (append) {
    if (isLoadingMore.value || !hasMore.value)
      return
    isLoadingMore.value = true
  }
  else {
    isLoading.value = true
    loadError.value = null
  }
  try {
    const offset = append ? comments.value.length : 0
    const res = await commentsService.list(props.entityType, props.entityId, PAGE_SIZE, offset)
    if (append)
      comments.value.push(...(res.items ?? []))
    else
      comments.value = res.items ?? []
    total.value = res.total
    loaded.value = true
    syncCountUp()
  }
  catch (error) {
    if (append)
      handleError(error)
    else
      loadError.value = (await handleError(error)).message
  }
  finally {
    isLoading.value = false
    isLoadingMore.value = false
  }
}

async function postComment() {
  const body = newBody.value.trim()
  if (!body || submitting.value)
    return
  submitting.value = true
  try {
    const created = await commentsService.create(props.entityType, props.entityId, body)
    comments.value.push(created)
    total.value += 1
    newBody.value = ''
    syncCountUp()
  }
  catch (error) {
    handleError(error)
  }
  finally {
    submitting.value = false
  }
}

function startEdit(c: Comment) {
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
    const updated = await commentsService.update(editingId.value, body)
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

async function deleteComment(c: Comment) {
  try {
    await commentsService.remove(c.id)
    comments.value = comments.value.filter(it => it.id !== c.id)
    total.value = Math.max(0, total.value - 1)
    toast({ title: 'Комментарий удалён' })
    syncCountUp()
  }
  catch (error) {
    handleError(error)
  }
}

async function toggleLike(c: Comment) {
  if (likeBusy.value.has(c.id))
    return
  likeBusy.value.add(c.id)
  const prevLiked = c.liked
  const prevCount = c.likesCount
  c.liked = !prevLiked
  c.likesCount = prevCount + (prevLiked ? -1 : 1)
  try {
    const res = await commentsService.toggleLike(c.id)
    c.liked = res.liked
    c.likesCount = res.likesCount
  }
  catch (error) {
    c.liked = prevLiked
    c.likesCount = prevCount
    handleError(error)
  }
  finally {
    likeBusy.value.delete(c.id)
  }
}

function canManage(c: Comment) {
  return isAdmin.value || user.value?.id === c.authorId
}

function authorName(c: Comment): string {
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

watch(() => [props.entityType, props.entityId], () => {
  comments.value = []
  total.value = 0
  loaded.value = false
  if (props.autoLoad !== false)
    fetchComments(false)
})

onMounted(() => {
  if (props.autoLoad !== false)
    fetchComments(false)
})

defineExpose({ refresh: () => fetchComments(false) })
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
        :maxlength="COMMENT_MAX_LEN"
        class="w-full bg-transparent text-base sm:text-sm focus:outline-none resize-none min-h-16"
        placeholder="Поделитесь мнением, спросите автора, оставьте свой опыт..."
      />
      <div class="flex items-center justify-between">
        <span class="text-xs text-muted-foreground">{{ newBody.length }} / {{ COMMENT_MAX_LEN }}</span>
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

    <p v-else-if="loaded && total === 0" class="text-sm text-muted-foreground">
      Пока нет комментариев. Будьте первым.
    </p>

    <!--
      Рассинхрон: счётчик на parent показывает N, но API вернул items=[].
      Может случиться при race в денормализации; даём явный retry, не молчим.
    -->
    <div v-else-if="loaded && total > 0 && comments.length === 0" class="text-sm text-muted-foreground space-y-2">
      <p>Комментарии не подгрузились ({{ total }} в счётчике, 0 в ответе).</p>
      <button
        type="button"
        class="px-2 py-1 rounded-sm border border-border text-xs hover:text-foreground"
        @click="fetchComments(false)"
      >
        Перезагрузить
      </button>
    </div>

    <ul v-else-if="comments.length" class="space-y-3">
      <li
        v-for="c in comments"
        :key="c.id"
        class="rounded-sm border border-border bg-card p-3"
      >
        <header class="flex items-center justify-between gap-2 mb-1.5 text-xs text-muted-foreground">
          <div class="flex items-center gap-2 min-w-0">
            <RouterLink
              v-if="c.author"
              :to="{ name: 'memberProfile', params: { id: c.authorId } }"
              class="font-medium text-foreground truncate hover:text-accent hover:underline"
            >
              {{ authorName(c) }}
            </RouterLink>
            <span v-else class="font-medium text-foreground truncate">{{ authorName(c) }}</span>
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
          <div class="flex items-center gap-1 shrink-0">
            <button
              type="button"
              :aria-label="c.liked ? 'Убрать лайк' : 'Поставить лайк'"
              :aria-pressed="c.liked"
              :disabled="likeBusy.has(c.id)"
              class="inline-flex items-center gap-0.5 rounded-full px-1.5 py-1 transition-colors hover:text-red-500 hover:bg-red-500/10 disabled:opacity-50"
              :class="c.liked ? 'text-red-500' : ''"
              @click="toggleLike(c)"
            >
              <Heart class="h-3 w-3" :class="c.liked ? 'fill-red-500' : ''" />
              {{ c.likesCount }}
            </button>
            <template v-if="canManage(c) && editingId !== c.id">
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
            </template>
          </div>
        </header>

        <div v-if="editingId === c.id" class="space-y-2">
          <textarea
            v-model="editingBody"
            :maxlength="COMMENT_MAX_LEN"
            class="w-full rounded-sm border border-border bg-background px-2 py-1.5 text-base sm:text-sm focus:outline-none focus:ring-2 focus:ring-primary min-h-16 resize-none"
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

    <div v-if="hasMore" class="flex justify-center">
      <button
        type="button"
        class="px-3 py-1.5 rounded-sm border border-border text-xs font-medium text-muted-foreground hover:text-foreground hover:bg-muted transition-colors disabled:opacity-50"
        :disabled="isLoadingMore"
        @click="fetchComments(true)"
      >
        <Loader2 v-if="isLoadingMore" class="h-3 w-3 animate-spin inline mr-1" />
        Показать ещё ({{ total - comments.length }})
      </button>
    </div>
  </section>
</template>
