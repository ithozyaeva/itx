<script setup lang="ts">
import type { KudosItem } from '@/models/kudos'
import type { TelegramUser } from '@/models/profile'
import { Heart, Loader2, Send } from 'lucide-vue-next'
import { onMounted, ref } from 'vue'
import EmptyState from '@/components/common/EmptyState.vue'
import ErrorState from '@/components/common/ErrorState.vue'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogFooter,
  DialogHeader,
  DialogScrollContent,
  DialogTitle,
} from '@/components/ui/dialog'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { useToast } from '@/components/ui/toast'
import { Typography } from '@/components/ui/typography'
import { useSSE } from '@/composables/useSSE'
import { useUser } from '@/composables/useUser'
import { displayName } from '@/lib/utils'
import { apiClient } from '@/services/api'
import { handleError } from '@/services/errorService'
import { kudosService } from '@/services/kudos'

const { toast } = useToast()

const PAGE_SIZE = 20
const items = ref<KudosItem[]>([])
const total = ref(0)
const isLoading = ref(true)
const isLoadingMore = ref(false)
const isSubmitting = ref(false)
const loadError = ref<string | null>(null)
const showDialog = ref(false)
const members = ref<TelegramUser[]>([])
const selectedMemberId = ref<number | null>(null)
const message = ref('')
const user = useUser()

async function fetchKudos() {
  isLoading.value = true
  loadError.value = null
  try {
    const res = await kudosService.getRecent(PAGE_SIZE, 0)
    items.value = res.items ?? []
    total.value = res.total
  }
  catch (error) {
    loadError.value = (await handleError(error)).message
  }
  finally {
    isLoading.value = false
  }
}

async function loadMore() {
  isLoadingMore.value = true
  try {
    const res = await kudosService.getRecent(PAGE_SIZE, items.value.length)
    items.value.push(...(res.items ?? []))
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isLoadingMore.value = false
  }
}

async function fetchMembers() {
  try {
    const res = await apiClient.get('members', { prefixUrl: '/api/' }).json<{ items: TelegramUser[] }>()
    members.value = (res.items ?? []).filter(m => m.id !== user.value?.id)
  }
  catch (error) {
    handleError(error)
  }
}

async function sendKudos() {
  if (!selectedMemberId.value || !message.value.trim())
    return
  isSubmitting.value = true
  try {
    await kudosService.send(selectedMemberId.value, message.value.trim())
    toast({ title: 'Благодарность отправлена' })
    showDialog.value = false
    selectedMemberId.value = null
    message.value = ''
    await fetchKudos()
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isSubmitting.value = false
  }
}

function timeAgo(dateStr: string) {
  const diff = Date.now() - new Date(dateStr).getTime()
  const mins = Math.floor(diff / 60000)
  if (mins < 60)
    return `${mins} мин. назад`
  const hours = Math.floor(mins / 60)
  if (hours < 24)
    return `${hours} ч. назад`
  const days = Math.floor(hours / 24)
  return `${days} дн. назад`
}

function openDialog() {
  fetchMembers()
  showDialog.value = true
}

useSSE('kudos', () => fetchKudos())

onMounted(() => {
  fetchKudos()
})
</script>

<template>
  <div class="container mx-auto px-4 py-6 md:py-8">
    <div class="font-mono text-[11px] text-muted-foreground/60 tracking-wider mb-2">
      ~/community/kudos
    </div>
    <div class="flex items-center justify-between mb-6">
      <Typography
        variant="h2"
        as="h1"
      >
        Стена благодарностей
      </Typography>
      <button
        class="flex items-center gap-2 px-3 sm:px-4 py-2 rounded-xl bg-primary text-primary-foreground text-sm font-medium hover:bg-primary/90 transition-colors shrink-0"
        @click="openDialog"
      >
        <Send class="h-4 w-4" />
        <span class="hidden sm:inline">Поблагодарить</span>
      </button>
    </div>

    <div
      v-if="isLoading"
      class="flex justify-center py-12"
    >
      <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
    </div>

    <ErrorState
      v-else-if="loadError"
      :message="loadError"
      @retry="fetchKudos"
    />

    <template v-else>
      <EmptyState
        v-if="items.length === 0"
        :icon="Heart"
        title="Благодарностей пока нет"
        description="Отправьте благодарность участнику сообщества"
        action-label="Поблагодарить"
        @action="openDialog"
      />

      <div class="space-y-3">
        <div
          v-for="item in items"
          :key="item.id"
          class="rounded-sm border bg-card border-border p-4"
        >
          <div class="flex items-start gap-3">
            <img
              :src="item.fromAvatarUrl || `https://ui-avatars.com/api/?name=${encodeURIComponent(item.fromFirstName || '?')}&background=random`"
              :alt="displayName(item.fromFirstName, item.fromLastName)"
              class="h-10 w-10 rounded-full shrink-0 object-cover"
            >
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-1.5 flex-wrap text-sm">
                <router-link
                  :to="`/members/${item.fromId}`"
                  class="font-medium hover:underline"
                >
                  {{ displayName(item.fromFirstName, item.fromLastName) }}
                </router-link>
                <Heart class="h-3.5 w-3.5 text-red-500 fill-red-500" />
                <router-link
                  :to="`/members/${item.toId}`"
                  class="font-medium hover:underline"
                >
                  {{ displayName(item.toFirstName, item.toLastName) }}
                </router-link>
              </div>
              <p class="text-sm mt-1">
                {{ item.message }}
              </p>
              <p class="text-xs text-muted-foreground mt-1">
                {{ timeAgo(item.createdAt) }}
              </p>
            </div>
          </div>
        </div>
      </div>

      <div
        v-if="items.length < total"
        class="mt-4 flex justify-center"
      >
        <Button
          variant="outline"
          :disabled="isLoadingMore"
          @click="loadMore"
        >
          <Loader2
            v-if="isLoadingMore"
            class="mr-2 h-4 w-4 animate-spin"
          />
          Показать ещё
        </Button>
      </div>
    </template>

    <Dialog v-model:open="showDialog">
      <DialogScrollContent>
        <DialogHeader>
          <DialogTitle>Поблагодарить участника</DialogTitle>
        </DialogHeader>
        <form
          class="space-y-4"
          @submit.prevent="sendKudos"
        >
          <div>
            <label class="block text-sm font-medium mb-1">Кому</label>
            <Select
              :model-value="selectedMemberId ? String(selectedMemberId) : ''"
              @update:model-value="selectedMemberId = $event ? Number($event) : null"
            >
              <SelectTrigger>
                <SelectValue placeholder="Выберите участника" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem
                  v-for="m in members"
                  :key="m.id"
                  :value="String(m.id)"
                >
                  {{ m.firstName }} {{ m.lastName }}{{ m.tg ? ` (@${m.tg})` : '' }}
                </SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div>
            <label class="block text-sm font-medium mb-1">Сообщение</label>
            <textarea
              v-model="message"
              required
              class="w-full rounded-xl border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary min-h-20 resize-none"
              placeholder="За что вы хотите поблагодарить?"
            />
          </div>
          <p class="text-xs text-muted-foreground">
            Получатель получит +5 баллов. Лимит: 3 благодарности в день.
          </p>
          <DialogFooter>
            <button
              type="submit"
              class="px-4 py-2 rounded-xl bg-primary text-primary-foreground text-sm font-medium hover:bg-primary/90 transition-colors disabled:opacity-50"
              :disabled="!selectedMemberId || !message.trim() || isSubmitting"
            >
              <Loader2
                v-if="isSubmitting"
                class="h-4 w-4 animate-spin inline mr-1"
              />
              Отправить
            </button>
          </DialogFooter>
        </form>
      </DialogScrollContent>
    </Dialog>
  </div>
</template>
