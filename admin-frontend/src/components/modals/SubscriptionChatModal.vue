<script setup lang="ts">
import type { SubscriptionChatDetail } from '@/services/subscriptionService'
import { onUnmounted, ref, watch } from 'vue'
import Loader2 from '~icons/lucide/loader-2'
import { Button } from '@/components/ui/button'
import { Checkbox } from '@/components/ui/checkbox'
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { subscriptionService } from '@/services/subscriptionService'

const props = defineProps<{
  isOpen: boolean
  chatId: number | null
}>()

const emit = defineEmits(['update:isOpen', 'saved'])

const isLoading = ref(false)
const isSaving = ref(false)
const isCreateMode = ref(false)
const isResolving = ref(false)
const resolveError = ref('')

const formId = ref('')
const formTitle = ref('')
const formChatType = ref('supergroup')

let resolveTimeout: ReturnType<typeof setTimeout> | null = null

function onIdInput() {
  resolveError.value = ''
  if (resolveTimeout)
    clearTimeout(resolveTimeout)

  const id = Number(formId.value)
  if (!id || !isCreateMode.value)
    return

  resolveTimeout = setTimeout(async () => {
    isResolving.value = true
    resolveError.value = ''
    try {
      const result = await subscriptionService.resolveChat(id)
      if (result) {
        formTitle.value = result.title
        formChatType.value = result.chatType
      }
      else {
        resolveError.value = 'Бот не в чате или чат не найден'
      }
    }
    finally {
      isResolving.value = false
    }
  }, 500)
}
const formRole = ref<'anchor' | 'content'>('content')
const formAnchorTierID = ref<number | null>(null)
const formTierIDs = ref<number[]>([])
const formCategory = ref('')
const formEmoji = ref('')

watch(() => props.isOpen, async (open) => {
  if (!open) {
    resetForm()
    return
  }

  if (!props.chatId) {
    isCreateMode.value = true
    return
  }

  isCreateMode.value = false
  isLoading.value = true
  try {
    const detail = await subscriptionService.getChatDetail(props.chatId)
    if (detail) {
      fillForm(detail)
    }
  }
  finally {
    isLoading.value = false
  }
})

function resetForm() {
  formId.value = ''
  formTitle.value = ''
  formChatType.value = 'supergroup'
  formRole.value = 'content'
  formAnchorTierID.value = null
  formTierIDs.value = []
  formCategory.value = ''
  formEmoji.value = ''
}

function fillForm(detail: SubscriptionChatDetail) {
  formId.value = String(detail.id)
  formTitle.value = detail.title
  formChatType.value = detail.chatType
  if (detail.anchorForTierID) {
    formRole.value = 'anchor'
    formAnchorTierID.value = detail.anchorForTierID
  }
  else {
    formRole.value = 'content'
    formAnchorTierID.value = null
  }
  formTierIDs.value = detail.tierIDs ?? []
  formCategory.value = detail.category ?? ''
  formEmoji.value = detail.emoji ?? ''
}

function handleClose() {
  emit('update:isOpen', false)
}

// Чистим pending debounce при unmount — иначе закрытие модалки за 500ms
// до выстрела resolveChat-таймера всё равно отправляет запрос на сервер.
onUnmounted(() => {
  if (resolveTimeout) {
    clearTimeout(resolveTimeout)
    resolveTimeout = null
  }
})

function toggleTier(tierId: number) {
  const idx = formTierIDs.value.indexOf(tierId)
  if (idx >= 0) {
    formTierIDs.value.splice(idx, 1)
  }
  else {
    formTierIDs.value.push(tierId)
  }
}

async function handleSave() {
  isSaving.value = true
  try {
    if (isCreateMode.value) {
      const chatId = Number(formId.value)
      if (!chatId || !formTitle.value)
        return

      const success = await subscriptionService.createChat({
        id: chatId,
        title: formTitle.value,
        chatType: formChatType.value,
        anchorForTierID: formRole.value === 'anchor' && formAnchorTierID.value ? formAnchorTierID.value : undefined,
        tierIDs: formRole.value === 'content' ? formTierIDs.value : undefined,
        category: formCategory.value || null,
        emoji: formEmoji.value || null,
      })
      if (success) {
        emit('saved')
        handleClose()
      }
    }
    else {
      const chatId = Number(formId.value)
      const success = await subscriptionService.updateChat(chatId, {
        title: formTitle.value,
        anchorForTierID: formRole.value === 'anchor' && formAnchorTierID.value ? formAnchorTierID.value : undefined,
        clearAnchor: formRole.value === 'content',
        tierIDs: formRole.value === 'content' ? formTierIDs.value : [],
        category: formCategory.value || null,
        emoji: formEmoji.value || null,
        clearCategory: !formCategory.value && !formEmoji.value,
      })
      if (success) {
        emit('saved')
        handleClose()
      }
    }
  }
  finally {
    isSaving.value = false
  }
}
</script>

<template>
  <Dialog
    :open="isOpen"
    @update:open="handleClose"
  >
    <DialogContent class="sm:max-w-[500px]">
      <DialogHeader>
        <DialogTitle>{{ isCreateMode ? 'Добавить чат' : 'Редактировать чат' }}</DialogTitle>
      </DialogHeader>

      <div
        v-if="isLoading"
        class="py-8 text-center text-muted-foreground"
      >
        Загрузка...
      </div>

      <div
        v-else
        class="space-y-4"
      >
        <div class="space-y-2">
          <Label>Telegram Chat ID</Label>
          <div class="relative">
            <Input
              v-model="formId"
              :disabled="!isCreateMode"
              placeholder="-1001234567890"
              @input="onIdInput"
            />
            <Loader2
              v-if="isResolving"
              class="absolute right-3 top-1/2 -translate-y-1/2 h-4 w-4 animate-spin text-muted-foreground"
            />
          </div>
          <p
            v-if="resolveError"
            class="text-xs text-destructive"
          >
            {{ resolveError }}
          </p>
        </div>

        <div class="space-y-2">
          <Label>Название</Label>
          <Input
            v-model="formTitle"
            placeholder="Название чата"
          />
        </div>

        <div class="space-y-2">
          <Label>Тип</Label>
          <select
            v-model="formChatType"
            class="w-full h-9 rounded-sm border border-input bg-background px-3 text-sm font-mono text-xs focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
          >
            <option value="supergroup">
              supergroup
            </option>
            <option value="group">
              group
            </option>
            <option value="channel">
              channel
            </option>
          </select>
        </div>

        <div class="space-y-2">
          <Label>Роль</Label>
          <div class="flex gap-4">
            <label class="flex items-center gap-2 cursor-pointer">
              <input
                v-model="formRole"
                type="radio"
                value="content"
                class="accent-primary"
              >
              <span class="text-sm">Content</span>
            </label>
            <label class="flex items-center gap-2 cursor-pointer">
              <input
                v-model="formRole"
                type="radio"
                value="anchor"
                class="accent-primary"
              >
              <span class="text-sm">Anchor</span>
            </label>
          </div>
        </div>

        <div
          v-if="formRole === 'anchor'"
          class="space-y-2"
        >
          <Label>Anchor для тира</Label>
          <select
            v-model.number="formAnchorTierID"
            class="w-full h-9 rounded-sm border border-input bg-background px-3 text-sm font-mono text-xs focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
          >
            <option :value="null">
              Выберите тир
            </option>
            <option
              v-for="tier in subscriptionService.tiers.value"
              :key="tier.id"
              :value="tier.id"
            >
              {{ tier.name }} (level {{ tier.level }})
            </option>
          </select>
        </div>

        <div
          v-if="formRole === 'content'"
          class="space-y-2"
        >
          <Label>Доступен для тиров</Label>
          <div class="space-y-2">
            <label
              v-for="tier in subscriptionService.tiers.value"
              :key="tier.id"
              class="flex items-center gap-2 cursor-pointer"
            >
              <Checkbox
                :checked="formTierIDs.includes(tier.id)"
                @update:checked="toggleTier(tier.id)"
              />
              <span class="text-sm">{{ tier.name }} (level {{ tier.level }})</span>
            </label>
          </div>
        </div>

        <div class="grid grid-cols-[64px_1fr] gap-3 items-end">
          <div class="space-y-2">
            <Label>Эмодзи</Label>
            <Input
              v-model="formEmoji"
              placeholder="💬"
              class="text-center text-lg"
              maxlength="8"
            />
          </div>
          <div class="space-y-2">
            <Label>Категория</Label>
            <Input
              v-model="formCategory"
              placeholder="Например: По интересам"
              maxlength="100"
            />
          </div>
        </div>
      </div>

      <DialogFooter>
        <Button
          variant="outline"
          @click="handleClose"
        >
          Отмена
        </Button>
        <Button
          :disabled="isSaving || isLoading"
          @click="handleSave"
        >
          {{ isSaving ? 'Сохранение...' : 'Сохранить' }}
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
