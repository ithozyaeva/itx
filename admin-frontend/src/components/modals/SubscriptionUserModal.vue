<script setup lang="ts">
import type { SubscriptionUserDetail } from '@/services/subscriptionService'
import { ref, watch } from 'vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import { Button } from '@/components/ui/button'
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Label } from '@/components/ui/label'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { subscriptionService } from '@/services/subscriptionService'

const props = defineProps<{
  isOpen: boolean
  userId: number | null
}>()

const emit = defineEmits(['update:isOpen', 'saved'])

const user = ref<SubscriptionUserDetail | null>(null)
const isLoading = ref(false)
const selectedTierSlug = ref('')

watch(() => props.isOpen, async (open) => {
  if (open && props.userId) {
    isLoading.value = true
    try {
      user.value = await subscriptionService.getUser(props.userId)
    }
    finally {
      isLoading.value = false
    }
  }
  else {
    user.value = null
    selectedTierSlug.value = ''
  }
})

function handleClose() {
  emit('update:isOpen', false)
}

async function handleSetOverride() {
  if (!user.value || !selectedTierSlug.value)
    return

  const success = await subscriptionService.setOverride(user.value.id, selectedTierSlug.value)
  if (success) {
    user.value = await subscriptionService.getUser(user.value.id)
    emit('saved')
  }
}

async function handleClearOverride() {
  if (!user.value)
    return

  const success = await subscriptionService.clearOverride(user.value.id)
  if (success) {
    user.value = await subscriptionService.getUser(user.value.id)
    emit('saved')
  }
}

async function handleRevokeAccess(chatID: number) {
  if (!user.value)
    return

  const success = await subscriptionService.revokeAccess(user.value.id, chatID)
  if (success) {
    user.value = await subscriptionService.getUser(user.value.id)
    emit('saved')
  }
}

function formatDate(dateStr?: string) {
  if (!dateStr)
    return '-'
  return new Date(dateStr).toLocaleString()
}
</script>

<template>
  <Dialog
    :open="isOpen"
    @update:open="handleClose"
  >
    <DialogContent class="sm:max-w-[600px] max-h-[90vh] overflow-auto">
      <DialogHeader>
        <DialogTitle>Информация о подписчике</DialogTitle>
      </DialogHeader>

      <div
        v-if="isLoading"
        class="py-8 text-center text-muted-foreground"
      >
        Загрузка...
      </div>

      <div
        v-else-if="user"
        class="space-y-6"
      >
        <!-- Basic Info -->
        <div class="grid grid-cols-2 gap-3 text-sm">
          <div>
            <Label class="text-muted-foreground">ID</Label>
            <div class="font-mono">
              {{ user.id }}
            </div>
          </div>
          <div>
            <Label class="text-muted-foreground">Username</Label>
            <div>{{ user.username ? `@${user.username}` : '-' }}</div>
          </div>
          <div>
            <Label class="text-muted-foreground">Имя</Label>
            <div>{{ user.fullName }}</div>
          </div>
          <div>
            <Label class="text-muted-foreground">Статус</Label>
            <div :class="user.isActive ? 'text-green-500' : 'text-red-500'">
              {{ user.isActive ? 'Активен' : 'Неактивен' }}
            </div>
          </div>
        </div>

        <!-- Tier Info -->
        <div class="space-y-2">
          <Label class="text-muted-foreground text-xs uppercase tracking-wide">Тиры</Label>
          <div class="grid grid-cols-3 gap-3 text-sm">
            <div>
              <div class="text-muted-foreground text-xs">
                Resolved
              </div>
              <div>{{ user.resolvedTierName ?? '-' }}</div>
            </div>
            <div>
              <div class="text-muted-foreground text-xs">
                Manual
              </div>
              <div class="flex items-center gap-1">
                {{ user.manualTierName ?? '-' }}
                <Button
                  v-if="user.manualTierID"
                  v-permission="'can_edit_admin_subscriptions'"
                  variant="ghost"
                  size="sm"
                  class="h-5 w-5 p-0 text-orange-500"
                  @click="handleClearOverride"
                >
                  &times;
                </Button>
              </div>
            </div>
            <div>
              <div class="text-muted-foreground text-xs">
                Effective
              </div>
              <div class="font-medium">
                {{ user.effectiveTierName ?? '-' }}
              </div>
            </div>
          </div>
        </div>

        <!-- Override -->
        <div
          v-permission="'can_edit_admin_subscriptions'"
          class="flex items-end gap-2"
        >
          <div class="flex-1">
            <Label class="text-xs">Установить ручной тир</Label>
            <select
              v-model="selectedTierSlug"
              class="w-full mt-1 h-9 rounded-md border border-input bg-background px-3 text-sm"
            >
              <option value="">
                Выберите тир
              </option>
              <option
                v-for="tier in subscriptionService.tiers.value"
                :key="tier.id"
                :value="tier.slug"
              >
                {{ tier.name }} (level {{ tier.level }})
              </option>
            </select>
          </div>
          <Button
            size="sm"
            :disabled="!selectedTierSlug"
            @click="handleSetOverride"
          >
            Установить
          </Button>
        </div>

        <!-- Dates -->
        <div class="grid grid-cols-2 gap-3 text-sm">
          <div>
            <Label class="text-muted-foreground">Последняя проверка</Label>
            <div>{{ formatDate(user.lastCheckAt) }}</div>
          </div>
          <div>
            <Label class="text-muted-foreground">Регистрация</Label>
            <div>{{ formatDate(user.createdAt) }}</div>
          </div>
        </div>

        <!-- Access -->
        <div>
          <Label class="text-muted-foreground text-xs uppercase tracking-wide">
            Доступ к чатам ({{ user.access.length }})
          </Label>
          <Table
            v-if="user.access.length > 0"
            class="mt-2"
          >
            <TableHeader>
              <TableRow>
                <TableHead>Чат</TableHead>
                <TableHead>Доступ с</TableHead>
                <TableHead />
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow
                v-for="a in user.access"
                :key="a.chatID"
              >
                <TableCell>{{ a.chatTitle ?? a.chatID }}</TableCell>
                <TableCell class="text-xs">
                  {{ formatDate(a.grantedAt) }}
                </TableCell>
                <TableCell class="text-right">
                  <ConfirmDialog
                    title="Отозвать доступ?"
                    description="Пользователь потеряет доступ к этому чату."
                    confirm-label="Отозвать"
                    @confirm="handleRevokeAccess(a.chatID)"
                  >
                    <template #trigger>
                      <Button
                        v-permission="'can_edit_admin_subscriptions'"
                        variant="ghost"
                        size="sm"
                        class="text-destructive"
                      >
                        Отозвать
                      </Button>
                    </template>
                  </ConfirmDialog>
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
          <div
            v-else
            class="mt-2 text-sm text-muted-foreground"
          >
            Нет активных доступов
          </div>
        </div>
      </div>

      <DialogFooter>
        <Button
          variant="outline"
          @click="handleClose"
        >
          Закрыть
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
