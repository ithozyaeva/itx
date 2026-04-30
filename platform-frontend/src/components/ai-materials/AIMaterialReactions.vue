<script setup lang="ts">
import { Bookmark, Heart, MessageCircle } from 'lucide-vue-next'
import { computed, ref } from 'vue'
import { aiMaterialsService } from '@/services/aiMaterials'
import { handleError } from '@/services/errorService'

const props = defineProps<{
  materialId: number
  liked: boolean
  bookmarked: boolean
  likesCount: number
  bookmarksCount: number
  commentsCount: number
  size?: 'sm' | 'md'
  // Когда true — клики не делают навигацию по кнопке-обёртке
  // (используется в карточке листинга, чтобы клик по ❤️ не открывал детальную).
  stopPropagation?: boolean
}>()

const emit = defineEmits<{
  'update:liked': [v: boolean]
  'update:bookmarked': [v: boolean]
  'update:likesCount': [v: number]
  'update:bookmarksCount': [v: number]
  'jumpToComments': []
}>()

const likeBusy = ref(false)
const bookmarkBusy = ref(false)

const iconClass = computed(() => props.size === 'sm' ? 'h-3 w-3' : 'h-4 w-4')

async function toggleLike(e: MouseEvent) {
  if (props.stopPropagation) {
    e.preventDefault()
    e.stopPropagation()
  }
  if (likeBusy.value)
    return
  likeBusy.value = true
  // Оптимистичный апдейт — мгновенный отклик. Откатим при ошибке.
  const prevLiked = props.liked
  const prevCount = props.likesCount
  emit('update:liked', !prevLiked)
  emit('update:likesCount', prevCount + (prevLiked ? -1 : 1))
  try {
    const res = await aiMaterialsService.toggleLike(props.materialId)
    emit('update:liked', res.liked)
    emit('update:likesCount', res.likesCount)
  }
  catch (error) {
    emit('update:liked', prevLiked)
    emit('update:likesCount', prevCount)
    handleError(error)
  }
  finally {
    likeBusy.value = false
  }
}

async function toggleBookmark(e: MouseEvent) {
  if (props.stopPropagation) {
    e.preventDefault()
    e.stopPropagation()
  }
  if (bookmarkBusy.value)
    return
  bookmarkBusy.value = true
  const prevBookmarked = props.bookmarked
  const prevCount = props.bookmarksCount
  emit('update:bookmarked', !prevBookmarked)
  emit('update:bookmarksCount', prevCount + (prevBookmarked ? -1 : 1))
  try {
    const res = await aiMaterialsService.toggleBookmark(props.materialId)
    emit('update:bookmarked', res.bookmarked)
    emit('update:bookmarksCount', res.bookmarksCount)
  }
  catch (error) {
    emit('update:bookmarked', prevBookmarked)
    emit('update:bookmarksCount', prevCount)
    handleError(error)
  }
  finally {
    bookmarkBusy.value = false
  }
}

function jumpToComments(e: MouseEvent) {
  if (props.stopPropagation) {
    e.preventDefault()
    e.stopPropagation()
  }
  emit('jumpToComments')
}
</script>

<template>
  <div class="flex items-center gap-1 text-xs text-muted-foreground">
    <button
      type="button"
      :aria-label="liked ? 'Убрать лайк' : 'Поставить лайк'"
      :aria-pressed="liked"
      :disabled="likeBusy"
      class="inline-flex items-center gap-0.5 rounded-full px-1.5 py-1 transition-colors hover:text-red-500 hover:bg-red-500/10 disabled:opacity-50"
      :class="liked ? 'text-red-500' : ''"
      @click="toggleLike"
    >
      <Heart :class="[iconClass, liked ? 'fill-red-500' : '']" />
      {{ likesCount }}
    </button>
    <button
      type="button"
      :aria-label="bookmarked ? 'Убрать из закладок' : 'В закладки'"
      :aria-pressed="bookmarked"
      :disabled="bookmarkBusy"
      class="inline-flex items-center gap-0.5 rounded-full px-1.5 py-1 transition-colors hover:text-foreground hover:bg-muted disabled:opacity-50"
      :class="bookmarked ? 'text-foreground' : ''"
      @click="toggleBookmark"
    >
      <Bookmark :class="[iconClass, bookmarked ? 'fill-current' : '']" />
      {{ bookmarksCount }}
    </button>
    <button
      type="button"
      aria-label="К комментариям"
      class="inline-flex items-center gap-0.5 rounded-full px-1.5 py-1 transition-colors hover:text-foreground hover:bg-muted"
      @click="jumpToComments"
    >
      <MessageCircle :class="iconClass" />
      {{ commentsCount }}
    </button>
  </div>
</template>
