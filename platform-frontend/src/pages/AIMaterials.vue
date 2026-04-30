<script setup lang="ts">
import type { AIMaterial, CreateAIMaterialRequest, AIMaterialFilters as Filters } from '@/models/aiMaterial'
import { Loader2, Plus, Sparkles } from 'lucide-vue-next'
import { computed, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import AIMaterialCard from '@/components/ai-materials/AIMaterialCard.vue'
import AIMaterialEditor from '@/components/ai-materials/AIMaterialEditor.vue'
import AIMaterialFilters from '@/components/ai-materials/AIMaterialFilters.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import ErrorState from '@/components/common/ErrorState.vue'
import { useToast } from '@/components/ui/toast'
import { Typography } from '@/components/ui/typography'
import { aiMaterialsService } from '@/services/aiMaterials'
import { handleError } from '@/services/errorService'

const router = useRouter()
const { toast } = useToast()

const PAGE_SIZE = 24
const items = ref<AIMaterial[]>([])
const total = ref(0)
const isLoading = ref(true)
const isLoadingMore = ref(false)
const loadError = ref<string | null>(null)
const showCreateDialog = ref(false)
const isSubmitting = ref(false)

const filters = ref<Filters>({
  kind: '',
  q: '',
  sort: 'new',
  limit: PAGE_SIZE,
  offset: 0,
})

const hasMore = computed(() => items.value.length < total.value)

let queryDebounce: ReturnType<typeof setTimeout> | null = null

watch(
  () => filters.value,
  (next, prev) => {
    if (!prev)
      return
    // Offset — это курсор пагинации, не фильтр. Изменение только offset
    // (через loadMore) не должно триггерить reset, иначе append-запрос
    // отменяется reset-запросом и список схлопывается обратно к page 1.
    if (
      next.offset !== prev.offset
      && next.q === prev.q
      && next.kind === prev.kind
      && next.sort === prev.sort
      && next.tag === prev.tag
    ) {
      return
    }
    // Дебаунсим только query — kind/sort должны срабатывать сразу.
    if (next.q !== prev.q && next.kind === prev.kind && next.sort === prev.sort) {
      if (queryDebounce)
        clearTimeout(queryDebounce)
      queryDebounce = setTimeout(fetchItems, 250, true)
      return
    }
    fetchItems(true)
  },
  { deep: true },
)

async function fetchItems(reset: boolean) {
  if (reset) {
    isLoading.value = true
    loadError.value = null
    filters.value.offset = 0
  }
  else {
    isLoadingMore.value = true
  }
  try {
    const res = await aiMaterialsService.search(filters.value)
    if (reset)
      items.value = res.items ?? []
    else
      items.value.push(...(res.items ?? []))
    total.value = res.total
  }
  catch (error) {
    if (reset)
      loadError.value = (await handleError(error)).message
    else
      handleError(error)
  }
  finally {
    isLoading.value = false
    isLoadingMore.value = false
  }
}

function loadMore() {
  if (isLoadingMore.value || !hasMore.value)
    return
  filters.value = { ...filters.value, offset: items.value.length }
  fetchItems(false)
}

async function createMaterial(data: CreateAIMaterialRequest) {
  isSubmitting.value = true
  try {
    const created = await aiMaterialsService.create(data)
    toast({ title: 'Материал опубликован' })
    showCreateDialog.value = false
    router.push({ name: 'aiMaterialDetail', params: { id: created.id } })
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isSubmitting.value = false
  }
}

function openMaterial(item: AIMaterial) {
  router.push({ name: 'aiMaterialDetail', params: { id: item.id } })
}

onMounted(() => fetchItems(true))
</script>

<template>
  <div class="container mx-auto px-4 py-6 md:py-8">
    <div class="font-mono text-[11px] text-muted-foreground/60 tracking-wider mb-2">
      ~/knowledge/ai-materials
    </div>
    <div class="flex items-center justify-between mb-6 gap-3">
      <Typography variant="h2" as="h1">
        AI-материалы
      </Typography>
      <button
        class="flex items-center gap-2 px-3 sm:px-4 py-2 rounded-sm bg-primary text-primary-foreground text-sm font-medium hover:bg-primary/90 transition-colors shrink-0"
        @click="showCreateDialog = true"
      >
        <Plus class="h-4 w-4" />
        <span class="hidden sm:inline">Новый материал</span>
      </button>
    </div>

    <p class="text-sm text-muted-foreground mb-6 max-w-2xl">
      Каталог промтов, скиллов, библиотек и AI-агентов от участников. Делитесь своими наработками — и забирайте чужие в один клик.
    </p>

    <AIMaterialFilters v-model="filters" />

    <div
      v-if="isLoading"
      class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4"
    >
      <div v-for="i in 6" :key="i" class="rounded-sm border border-border bg-card p-4 animate-pulse h-44" />
    </div>

    <ErrorState
      v-else-if="loadError"
      :message="loadError"
      @retry="fetchItems(true)"
    />

    <template v-else>
      <EmptyState
        v-if="items.length === 0"
        :icon="Sparkles"
        title="Здесь пока пусто"
        description="Будьте первым: добавьте промт, скилл, библиотеку или конфиг агента"
        action-label="Добавить материал"
        @action="showCreateDialog = true"
      />

      <div
        v-else
        class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4"
      >
        <!--
          role="button" + tabindex/keydown — нельзя обернуть карточку в <button>,
          потому что внутри AIMaterialReactions есть свои <button>'ы
          (лайк/закладка/комментарии). Вложенные buttons — невалидный HTML,
          парсер закрывает внешний button раньше → клик по нижней части
          карточки перестаёт навигировать.
        -->
        <div
          v-for="(item, idx) in items"
          :key="item.id"
          role="button"
          tabindex="0"
          class="text-left cursor-pointer focus:outline-none focus-visible:ring-2 focus-visible:ring-primary rounded-sm"
          @click="openMaterial(item)"
          @keydown.enter.prevent="openMaterial(item)"
          @keydown.space.prevent="openMaterial(item)"
        >
          <AIMaterialCard
            :item="item"
            @update:item="(v) => (items[idx] = v)"
          />
        </div>
      </div>

      <div v-if="hasMore" class="mt-6 flex justify-center">
        <button
          class="px-4 py-2 rounded-sm border border-border text-sm font-medium hover:bg-accent hover:text-accent-foreground transition-colors disabled:opacity-50"
          :disabled="isLoadingMore"
          @click="loadMore"
        >
          <Loader2 v-if="isLoadingMore" class="h-4 w-4 animate-spin inline mr-1" />
          Показать ещё
        </button>
      </div>
    </template>

    <AIMaterialEditor
      v-model:open="showCreateDialog"
      :is-submitting="isSubmitting"
      @submit="createMaterial"
    />
  </div>
</template>
