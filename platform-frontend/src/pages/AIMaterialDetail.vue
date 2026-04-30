<script setup lang="ts">
import type { AIMaterial, CreateAIMaterialRequest } from '@/models/aiMaterial'
import { ArrowLeft, ExternalLink, FileCode2, Loader2, Pencil, Sparkles, Trash2 } from 'lucide-vue-next'
import { computed, nextTick, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import AIMaterialComments from '@/components/ai-materials/AIMaterialComments.vue'
import AIMaterialContentBlock from '@/components/ai-materials/AIMaterialContentBlock.vue'
import AIMaterialEditor from '@/components/ai-materials/AIMaterialEditor.vue'
import AIMaterialReactions from '@/components/ai-materials/AIMaterialReactions.vue'
import ErrorState from '@/components/common/ErrorState.vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import { useToast } from '@/components/ui/toast'
import { Typography } from '@/components/ui/typography'
import { isUserAdmin, useUser } from '@/composables/useUser'
import { AI_MATERIAL_KIND_LABELS } from '@/models/aiMaterial'
import { aiMaterialsService } from '@/services/aiMaterials'
import { handleError } from '@/services/errorService'

const route = useRoute()
const router = useRouter()
const { toast } = useToast()
const user = useUser()
const isAdmin = isUserAdmin()

const item = ref<AIMaterial | null>(null)
const isLoading = ref(true)
const loadError = ref<string | null>(null)
const showEditDialog = ref(false)
const isSubmitting = ref(false)
const commentsRef = ref<HTMLElement | null>(null)

const isAuthor = computed(() => !!item.value && user.value?.id === item.value.authorId)
const canManage = computed(() => isAuthor.value || isAdmin.value)

const contentIcon = computed(() => {
  if (!item.value)
    return Sparkles
  switch (item.value.contentType) {
    case 'link':
      return ExternalLink
    case 'agent':
      return FileCode2
    default:
      return Sparkles
  }
})

function authorName(): string {
  const a = item.value?.author
  if (!a)
    return 'Аноним'
  const name = [a.firstName, a.lastName].filter(Boolean).join(' ')
  return name || (a.tg ? `@${a.tg}` : 'Аноним')
}

function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString('ru-RU', { day: 'numeric', month: 'long', year: 'numeric' })
}

async function fetchItem() {
  isLoading.value = true
  loadError.value = null
  try {
    const id = Number(route.params.id)
    if (!Number.isFinite(id) || id <= 0) {
      loadError.value = 'Неверный идентификатор материала'
      return
    }
    item.value = await aiMaterialsService.getById(id)
  }
  catch (error) {
    loadError.value = (await handleError(error)).message
  }
  finally {
    isLoading.value = false
  }
}

async function saveEdit(data: CreateAIMaterialRequest) {
  if (!item.value)
    return
  isSubmitting.value = true
  try {
    item.value = await aiMaterialsService.update(item.value.id, data)
    toast({ title: 'Материал обновлён' })
    showEditDialog.value = false
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isSubmitting.value = false
  }
}

async function deleteItem() {
  if (!item.value)
    return
  try {
    await aiMaterialsService.remove(item.value.id)
    toast({ title: 'Материал удалён' })
    router.push({ name: 'aiMaterials' })
  }
  catch (error) {
    handleError(error)
  }
}

function patch<K extends keyof AIMaterial>(field: K, value: AIMaterial[K]) {
  if (!item.value)
    return
  item.value = { ...item.value, [field]: value }
}

async function jumpToComments() {
  await nextTick()
  commentsRef.value?.scrollIntoView({ behavior: 'smooth', block: 'start' })
}

watch(() => route.params.id, () => {
  if (route.name === 'aiMaterialDetail')
    fetchItem()
})

onMounted(fetchItem)
</script>

<template>
  <div class="container mx-auto px-4 py-6 md:py-8 max-w-4xl">
    <button
      class="mb-4 inline-flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors"
      @click="router.push({ name: 'aiMaterials' })"
    >
      <ArrowLeft class="h-4 w-4" />
      К каталогу
    </button>

    <div v-if="isLoading" class="space-y-4">
      <div class="h-8 w-2/3 rounded bg-muted animate-pulse" />
      <div class="h-4 w-full rounded bg-muted animate-pulse" />
      <div class="h-64 w-full rounded bg-muted animate-pulse" />
    </div>

    <ErrorState
      v-else-if="loadError || !item"
      :message="loadError ?? 'Материал не найден'"
      @retry="fetchItem"
    />

    <article v-else class="space-y-6">
      <header class="space-y-3">
        <div class="flex items-center gap-2 flex-wrap text-xs text-muted-foreground">
          <span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full font-medium bg-accent/15 text-accent">
            <component :is="contentIcon" class="h-3 w-3" />
            {{ AI_MATERIAL_KIND_LABELS[item.materialKind] }}
          </span>
          <span>{{ formatDate(item.createdAt) }}</span>
          <span>·</span>
          <span>{{ authorName() }}</span>
          <span v-if="item.isHidden" class="px-2 py-0.5 rounded-full bg-yellow-500/15 text-yellow-600">
            Скрыт
          </span>
        </div>
        <Typography variant="h2" as="h1">
          {{ item.title }}
        </Typography>
        <p class="text-base text-muted-foreground">
          {{ item.summary }}
        </p>
        <div v-if="item.tags.length" class="flex flex-wrap gap-1.5">
          <span
            v-for="t in item.tags"
            :key="t"
            class="px-2 py-0.5 rounded-full text-xs font-mono bg-muted text-muted-foreground"
          >#{{ t }}</span>
        </div>
      </header>

      <AIMaterialContentBlock :item="item" />

      <div class="border-t border-border pt-4">
        <AIMaterialReactions
          :material-id="item.id"
          :liked="item.liked"
          :bookmarked="item.bookmarked"
          :likes-count="item.likesCount"
          :bookmarks-count="item.bookmarksCount"
          :comments-count="item.commentsCount"
          @update:liked="(v) => patch('liked', v)"
          @update:bookmarked="(v) => patch('bookmarked', v)"
          @update:likes-count="(v) => patch('likesCount', v)"
          @update:bookmarks-count="(v) => patch('bookmarksCount', v)"
          @jump-to-comments="jumpToComments"
        />
      </div>

      <div v-if="canManage" class="flex flex-wrap gap-2 border-t border-border pt-4">
        <button
          class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-sm text-sm font-medium border border-border hover:bg-accent hover:text-accent-foreground transition-colors"
          @click="showEditDialog = true"
        >
          <Pencil class="h-4 w-4" />
          Редактировать
        </button>
        <ConfirmDialog
          title="Удалить материал?"
          description="Материал и все связанные с ним лайки/закладки/комментарии будут удалены безвозвратно."
          confirm-label="Удалить"
          @confirm="deleteItem"
        >
          <template #trigger>
            <button
              class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-sm text-sm font-medium text-red-500 border border-red-500/30 hover:bg-red-500/10 transition-colors"
            >
              <Trash2 class="h-4 w-4" />
              Удалить
            </button>
          </template>
        </ConfirmDialog>
      </div>

      <div ref="commentsRef" class="border-t border-border pt-4">
        <AIMaterialComments
          :material-id="item.id"
          :initial-count="item.commentsCount"
          @update:count="(v) => patch('commentsCount', v)"
        />
      </div>

      <AIMaterialEditor
        v-model:open="showEditDialog"
        :initial="item"
        :is-submitting="isSubmitting"
        @submit="saveEdit"
      />
    </article>

    <div v-if="isSubmitting" class="fixed bottom-4 right-4">
      <Loader2 class="h-6 w-6 animate-spin text-primary" />
    </div>
  </div>
</template>
