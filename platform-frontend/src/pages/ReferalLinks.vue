<script setup lang="ts">
import type { ReferalLink } from '@/models/referals'
import type { ReferalSearchFilters } from '@/services/referals'
import { Typography } from 'itx-ui-kit'
import { Loader2, Share2 } from 'lucide-vue-next'
import { onMounted, ref } from 'vue'
import EmptyState from '@/components/common/EmptyState.vue'
import ErrorState from '@/components/common/ErrorState.vue'
import ReferalFilters from '@/components/referals/ReferalFilters.vue'
import ReferalLinkCard from '@/components/referals/ReferalLinkCard.vue'
import ReferalLinkForm from '@/components/referals/ReferalLinkForm.vue'
import { useCardReveal } from '@/composables/useCardReveal'
import { handleError } from '@/services/errorService'
import { referalLinkService } from '@/services/referals'

const containerRef = ref<HTMLElement | null>(null)
useCardReveal(containerRef)

const referalLinks = ref<ReferalLink[]>([])
const showAddForm = ref(false)
const isSaving = ref(false)
const totalLinks = ref(0)
const currentOffset = ref(0)
const ITEMS_PER_PAGE = 10
const currentFilters = ref<ReferalSearchFilters>({})
const isLoading = ref(false)
const isLoadingMore = ref(false)
const loadError = ref<string | null>(null)

async function fetchReferalLinks(filters?: ReferalSearchFilters) {
  if (filters) {
    currentFilters.value = filters
    currentOffset.value = 0
  }
  isLoading.value = currentOffset.value === 0
  loadError.value = null
  try {
    const response = await referalLinkService.search(ITEMS_PER_PAGE, currentOffset.value, currentFilters.value)
    if (currentOffset.value === 0) {
      referalLinks.value = response.items ?? []
    }
    else {
      referalLinks.value = [...referalLinks.value, ...(response.items ?? [])]
    }
    totalLinks.value = response.total
  }
  catch (error) {
    loadError.value = (await handleError(error)).message
  }
  finally {
    isLoading.value = false
  }
}

async function loadMore() {
  if (isLoadingMore.value)
    return
  isLoadingMore.value = true
  try {
    currentOffset.value += ITEMS_PER_PAGE
    await fetchReferalLinks()
  }
  finally {
    isLoadingMore.value = false
  }
}

onMounted(() => {
  fetchReferalLinks()
})

function toggleAddForm() {
  showAddForm.value = !showAddForm.value
}

async function saveNewLink(newLink: Partial<ReferalLink>) {
  isSaving.value = true
  try {
    const addedLink = await referalLinkService.addLink(newLink)
    referalLinks.value.unshift(addedLink)
    showAddForm.value = false
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isSaving.value = false
  }
}

function cancelAdd() {
  showAddForm.value = false
}

function handleLinkUpdated(updatedLink: ReferalLink) {
  const index = referalLinks.value.findIndex(link => link.id === updatedLink.id)
  if (index !== -1) {
    referalLinks.value[index] = updatedLink
  }
}

function handleLinkDeleted(deletedLinkId: number) {
  const index = referalLinks.value.findIndex(link => link.id === deletedLinkId)
  if (index !== -1) {
    referalLinks.value.splice(index, 1)
  }
}
</script>

<template>
  <div ref="containerRef" class="container mx-auto px-4 py-6 md:py-8">
    <Typography variant="h2" as="h1" class="mb-4">
      Реферальные ссылки
    </Typography>

    <ReferalFilters class="mb-6" @change="fetchReferalLinks" />

    <div v-if="isLoading" class="flex justify-center py-12">
      <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
    </div>

    <ErrorState
      v-else-if="loadError"
      :message="loadError"
      @retry="fetchReferalLinks"
    />

    <EmptyState
      v-else-if="referalLinks.length === 0 && !showAddForm"
      :icon="Share2"
      title="Нет реферальных ссылок"
      description="Создайте реферальную ссылку и приглашайте новых участников"
      action-label="Добавить ссылку"
      @action="toggleAddForm"
    />

    <div v-else class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
      <div class="bg-card rounded-3xl border p-4 shadow-md ">
        <button
          v-if="!showAddForm"
          type="button"
          class="transition-shadow flex flex-col items-center justify-center gap-2 text-center cursor-pointer min-h-[100px] sm:min-h-[150px] w-full bg-transparent border-none"
          @click="toggleAddForm"
        >
          <span class="text-4xl">+</span>
          <span class="text-lg font-semibold">Добавить ссылку</span>
        </button>
        <ReferalLinkForm
          v-if="showAddForm"
          :is-saving="isSaving"
          @save="saveNewLink"
          @cancel="cancelAdd"
        />
      </div>

      <ReferalLinkCard
        v-for="link in referalLinks"
        :key="link.id"
        :link="link"
        @updated="handleLinkUpdated"
        @deleted="handleLinkDeleted"
      />

      <button v-if="referalLinks.length < totalLinks" type="button" class="bg-card rounded-3xl border p-4 hover:shadow-md flex justify-center items-center cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed" :disabled="isLoadingMore" @click="loadMore">
        <Loader2 v-if="isLoadingMore" class="h-5 w-5 animate-spin" />
        <span v-else class="m-auto">
          Показать ещё
        </span>
      </button>
    </div>
  </div>
</template>
