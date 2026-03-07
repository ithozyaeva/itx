<script setup lang="ts">
import type { KudosItem } from '@/models/kudos'
import type { TelegramUser } from '@/models/profile'
import { Typography } from 'itx-ui-kit'
import { Heart, Loader2, Send } from 'lucide-vue-next'
import { onMounted, ref } from 'vue'
import {
  Dialog,
  DialogFooter,
  DialogHeader,
  DialogScrollContent,
  DialogTitle,
} from '@/components/ui/dialog'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { useSSE } from '@/composables/useSSE'
import { useUser } from '@/composables/useUser'
import { apiClient } from '@/services/api'
import { handleError } from '@/services/errorService'
import { kudosService } from '@/services/kudos'

const items = ref<KudosItem[]>([])
const total = ref(0)
const isLoading = ref(true)
const isSubmitting = ref(false)
const showDialog = ref(false)
const members = ref<TelegramUser[]>([])
const selectedMemberId = ref<number | null>(null)
const message = ref('')
const user = useUser()

async function fetchKudos() {
  isLoading.value = true
  try {
    const res = await kudosService.getRecent(50)
    items.value = res.items ?? []
    total.value = res.total
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isLoading.value = false
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

function displayName(firstName: string, lastName: string) {
  return [firstName, lastName].filter(Boolean).join(' ')
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
    <div class="flex items-center justify-between mb-6">
      <Typography
        variant="h2"
        as="h1"
      >
        Стена благодарностей
      </Typography>
      <button
        class="flex items-center gap-2 px-4 py-2 rounded-xl bg-primary text-primary-foreground text-sm font-medium hover:bg-primary/90 transition-colors"
        @click="openDialog"
      >
        <Send class="h-4 w-4" />
        Поблагодарить
      </button>
    </div>

    <div
      v-if="isLoading"
      class="flex justify-center py-12"
    >
      <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
    </div>

    <template v-else>
      <div
        v-if="items.length === 0"
        class="text-center py-12 text-muted-foreground"
      >
        <Heart class="h-12 w-12 mx-auto mb-3 opacity-50" />
        <p>Пока нет благодарностей. Будьте первым!</p>
      </div>

      <div class="space-y-3">
        <div
          v-for="item in items"
          :key="item.id"
          class="rounded-2xl border bg-card border-border p-4"
        >
          <div class="flex items-start gap-3">
            <img
              :src="item.fromAvatarUrl || `https://ui-avatars.com/api/?name=${item.fromFirstName}`"
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
                  {{ m.firstName }} {{ m.lastName }} (@{{ m.tg }})
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
