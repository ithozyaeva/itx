<script setup lang="ts">
import type { SubscriptionChatDetail } from '@/services/subscriptionService'
import { computed, onMounted, onUnmounted, ref } from 'vue'
import Anchor from '~icons/lucide/anchor'
import ChevronRight from '~icons/lucide/chevron-right'
import Eye from '~icons/lucide/eye'
import Loader2 from '~icons/lucide/loader-2'
import Plus from '~icons/lucide/plus'
import ShieldX from '~icons/lucide/shield-x'
import Trash2 from '~icons/lucide/trash-2'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import AdminLayout from '@/components/layout/AdminLayout.vue'
import SubscriptionChatModal from '@/components/modals/SubscriptionChatModal.vue'
import SubscriptionUserModal from '@/components/modals/SubscriptionUserModal.vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Pagination, PaginationEllipsis, PaginationFirst, PaginationLast, PaginationList, PaginationListItem, PaginationNext, PaginationPrev } from '@/components/ui/pagination'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Typography } from '@/components/ui/typography'
import { subscriptionService } from '@/services/subscriptionService'

type Tab = 'overview' | 'users' | 'chats'

const activeTab = ref<Tab>('overview')
const selectedUserId = ref<number | null>(null)
const isUserModalOpen = ref(false)
const selectedChatId = ref<number | null>(null)
const isChatModalOpen = ref(false)

const expandedChatId = ref<number | null>(null)
const expandedChatDetail = ref<SubscriptionChatDetail | null>(null)
const expandedLoading = ref(false)
const savingAction = ref<string | null>(null)

const stats = computed(() => subscriptionService.stats.value)

function openUserModal(userId: number) {
  selectedUserId.value = userId
  isUserModalOpen.value = true
}

function openChatModal(chatId: number | null) {
  selectedChatId.value = chatId
  isChatModalOpen.value = true
}

async function toggleChatExpand(chatId: number) {
  if (expandedChatId.value === chatId) {
    expandedChatId.value = null
    expandedChatDetail.value = null
    return
  }

  expandedChatId.value = chatId
  expandedLoading.value = true
  try {
    expandedChatDetail.value = await subscriptionService.getChatDetail(chatId)
  }
  finally {
    expandedLoading.value = false
  }
}

async function toggleAnchor(chatId: number, tierID: number) {
  if (!expandedChatDetail.value)
    return
  const isActive = expandedChatDetail.value.anchorForTierID === tierID
  const newTierID = isActive ? null : tierID

  savingAction.value = `anchor-${tierID}`
  try {
    const success = await subscriptionService.updateChat(chatId, {
      anchorForTierID: newTierID ?? undefined,
      clearAnchor: newTierID === null,
      tierIDs: newTierID !== null ? [] : expandedChatDetail.value.tierIDs,
    })
    if (success) {
      expandedChatDetail.value = {
        ...expandedChatDetail.value,
        anchorForTierID: newTierID ?? undefined,
        tierIDs: newTierID !== null ? [] : expandedChatDetail.value.tierIDs,
      }
      await subscriptionService.fetchChats()
      await subscriptionService.fetchStats()
    }
  }
  finally {
    savingAction.value = null
  }
}

async function toggleContentTier(chatId: number, tierId: number) {
  if (!expandedChatDetail.value)
    return

  const currentTiers = [...(expandedChatDetail.value.tierIDs ?? [])]
  const idx = currentTiers.indexOf(tierId)
  if (idx >= 0)
    currentTiers.splice(idx, 1)
  else
    currentTiers.push(tierId)

  savingAction.value = `tier-${tierId}`
  try {
    const success = await subscriptionService.updateChat(chatId, {
      clearAnchor: true,
      tierIDs: currentTiers,
    })
    if (success) {
      expandedChatDetail.value = {
        ...expandedChatDetail.value,
        anchorForTierID: undefined,
        tierIDs: currentTiers,
      }
      await subscriptionService.fetchChats()
    }
  }
  finally {
    savingAction.value = null
  }
}

async function handleDeleteChat(chatId: number) {
  const success = await subscriptionService.deleteChat(chatId)
  if (success) {
    if (expandedChatId.value === chatId)
      expandedChatId.value = null
    await subscriptionService.fetchChats()
    await subscriptionService.fetchStats()
  }
}

async function handleChatSaved() {
  await subscriptionService.fetchChats()
  await subscriptionService.fetchStats()
}

async function handleClearOverride(userId: number) {
  await subscriptionService.clearOverride(userId)
  await subscriptionService.searchUsers()
  await subscriptionService.fetchStats()
}

function switchTab(tab: Tab) {
  activeTab.value = tab
  if (tab === 'users') {
    subscriptionService.searchUsers()
  }
  else if (tab === 'chats') {
    subscriptionService.fetchChats()
  }
}

onMounted(() => {
  subscriptionService.fetchStats()
  subscriptionService.fetchTiers()
})
onUnmounted(subscriptionService.clearPagination)
</script>

<template>
  <AdminLayout>
    <div class="space-y-6">
      <Typography
        variant="h2"
        as="h1"
      >
        Подписки
      </Typography>

      <!-- Tabs -->
      <div class="flex gap-2 border-b pb-2">
        <Button
          :variant="activeTab === 'overview' ? 'default' : 'ghost'"
          size="sm"
          @click="switchTab('overview')"
        >
          Обзор
        </Button>
        <Button
          :variant="activeTab === 'users' ? 'default' : 'ghost'"
          size="sm"
          @click="switchTab('users')"
        >
          Пользователи
        </Button>
        <Button
          :variant="activeTab === 'chats' ? 'default' : 'ghost'"
          size="sm"
          @click="switchTab('chats')"
        >
          Чаты
        </Button>
      </div>

      <!-- Overview Tab -->
      <template v-if="activeTab === 'overview'">
        <div
          v-if="stats"
          class="grid grid-cols-2 md:grid-cols-4 gap-4"
        >
          <Card>
            <CardContent class="p-4">
              <div class="text-sm text-muted-foreground">
                Всего пользователей
              </div>
              <div class="text-2xl font-bold">
                {{ stats.totalUsers }}
              </div>
            </CardContent>
          </Card>
          <Card>
            <CardContent class="p-4">
              <div class="text-sm text-muted-foreground">
                Всего чатов
              </div>
              <div class="text-2xl font-bold">
                {{ stats.totalChats }}
              </div>
            </CardContent>
          </Card>
          <Card>
            <CardContent class="p-4">
              <div class="text-sm text-muted-foreground">
                Anchor чатов
              </div>
              <div class="text-2xl font-bold">
                {{ stats.anchorChats }}
              </div>
            </CardContent>
          </Card>
          <Card>
            <CardContent class="p-4">
              <div class="text-sm text-muted-foreground">
                Content чатов
              </div>
              <div class="text-2xl font-bold">
                {{ stats.contentChats }}
              </div>
            </CardContent>
          </Card>
        </div>

        <Card v-if="stats">
          <CardContent>
            <Typography
              variant="h4"
              as="h3"
              class="mb-4 pt-4"
            >
              Тиры подписок
            </Typography>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Level</TableHead>
                  <TableHead>Название</TableHead>
                  <TableHead>Slug</TableHead>
                  <TableHead class="text-right">
                    Пользователей
                  </TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                <TableRow
                  v-for="tier in stats.tiers"
                  :key="tier.id"
                >
                  <TableCell>{{ tier.level }}</TableCell>
                  <TableCell class="font-medium">
                    {{ tier.name }}
                  </TableCell>
                  <TableCell>
                    <code class="text-xs bg-muted px-1.5 py-0.5 rounded">{{ tier.slug }}</code>
                  </TableCell>
                  <TableCell class="text-right">
                    {{ tier.users }}
                  </TableCell>
                </TableRow>
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      </template>

      <!-- Users Tab -->
      <template v-if="activeTab === 'users'">
        <Card>
          <CardContent>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>ID</TableHead>
                  <TableHead>Имя</TableHead>
                  <TableHead>Username</TableHead>
                  <TableHead>Тир</TableHead>
                  <TableHead>Чатов</TableHead>
                  <TableHead>Статус</TableHead>
                  <TableHead />
                </TableRow>
              </TableHeader>
              <TableBody>
                <TableRow
                  v-if="subscriptionService.users.value.total === 0"
                  class="h-24"
                >
                  <TableCell
                    colspan="7"
                    class="text-center"
                  >
                    Пользователи не найдены
                  </TableCell>
                </TableRow>
                <TableRow
                  v-for="user in subscriptionService.users.value.items"
                  :key="user.id"
                >
                  <TableCell>
                    <code class="text-xs">{{ user.id }}</code>
                  </TableCell>
                  <TableCell>{{ user.fullName }}</TableCell>
                  <TableCell>
                    <span
                      v-if="user.username"
                      class="text-muted-foreground"
                    >@{{ user.username }}</span>
                    <span
                      v-else
                      class="text-muted-foreground"
                    >-</span>
                  </TableCell>
                  <TableCell>
                    <span
                      v-if="user.tierName"
                      class="inline-flex items-center gap-1"
                    >
                      {{ user.tierName }}
                      <span
                        v-if="user.manualTierID"
                        class="text-xs text-orange-500"
                        title="Ручной тир"
                      >M</span>
                    </span>
                    <span
                      v-else
                      class="text-muted-foreground"
                    >-</span>
                  </TableCell>
                  <TableCell>{{ user.activeChats }}</TableCell>
                  <TableCell>
                    <span
                      :class="user.isActive ? 'text-green-500' : 'text-red-500'"
                      class="text-xs font-medium"
                    >{{ user.isActive ? 'Активен' : 'Неактивен' }}</span>
                  </TableCell>
                  <TableCell class="text-right space-x-1">
                    <Button
                      variant="ghost"
                      size="sm"
                      @click="openUserModal(user.id)"
                    >
                      <Eye class="h-4 w-4" />
                    </Button>
                    <ConfirmDialog
                      v-if="user.manualTierID"
                      title="Снять ручной тир?"
                      description="Пользователь вернётся к автоматически определённому тиру."
                      confirm-label="Снять"
                      @confirm="handleClearOverride(user.id)"
                    >
                      <template #trigger>
                        <Button
                          v-permission="'can_edit_admin_subscriptions'"
                          variant="ghost"
                          size="sm"
                          class="text-orange-500"
                        >
                          <ShieldX class="h-4 w-4" />
                        </Button>
                      </template>
                    </ConfirmDialog>
                  </TableCell>
                </TableRow>
              </TableBody>
            </Table>
          </CardContent>
        </Card>

        <div class="mt-4 flex justify-end">
          <Pagination
            v-slot="{ page }"
            :items-per-page="20"
            :total="subscriptionService.users.value.total"
            :sibling-count="1"
            show-edges
            :default-page="1"
          >
            <PaginationList
              v-slot="{ items }"
              class="flex items-center gap-1"
            >
              <PaginationFirst />
              <PaginationPrev />

              <template v-for="(item, index) in items">
                <PaginationListItem
                  v-if="item.type === 'page'"
                  :key="index"
                  :value="item.value"
                  as-child
                >
                  <Button
                    class="w-10 h-10 p-0"
                    :variant="item.value === page ? 'default' : 'outline'"
                    @click="subscriptionService.changePagination(item.value)"
                  >
                    {{ item.value }}
                  </Button>
                </PaginationListItem>
                <PaginationEllipsis
                  v-else
                  :key="item.type"
                  :index="index"
                />
              </template>

              <PaginationNext />
              <PaginationLast />
            </PaginationList>
          </Pagination>
        </div>
      </template>

      <!-- Chats Tab -->
      <template v-if="activeTab === 'chats'">
        <div class="flex justify-end">
          <Button
            size="sm"
            @click="openChatModal(null)"
          >
            <Plus class="h-4 w-4 mr-1" />
            Добавить чат
          </Button>
        </div>

        <div class="space-y-2">
          <div
            v-if="subscriptionService.chats.value.length === 0"
            class="text-center py-12 text-muted-foreground"
          >
            Чаты не найдены
          </div>

          <Card
            v-for="chat in subscriptionService.chats.value"
            :key="chat.id"
            class="transition-colors"
            :class="{ 'ring-1 ring-primary/20': expandedChatId === chat.id }"
          >
            <div
              class="flex items-center gap-3 px-4 py-3 cursor-pointer select-none hover:bg-muted/40 transition-colors"
              @click="toggleChatExpand(chat.id)"
            >
              <ChevronRight
                class="h-4 w-4 text-muted-foreground shrink-0 transition-transform duration-200"
                :class="{ 'rotate-90': expandedChatId === chat.id }"
              />

              <div class="flex-1 min-w-0">
                <div class="font-medium truncate">
                  {{ chat.title }}
                </div>
                <div class="text-xs text-muted-foreground">
                  {{ chat.chatType }} &middot; {{ chat.id }}
                </div>
              </div>

              <div class="flex flex-wrap gap-1 shrink-0">
                <span
                  v-if="chat.anchorTierName"
                  class="inline-flex items-center gap-1 text-blue-500 text-xs font-medium bg-blue-500/10 px-2 py-0.5 rounded-full"
                >
                  <Anchor class="h-3 w-3" />
                  {{ chat.anchorTierName }}
                </span>
                <span
                  v-for="name in (chat.tierNames || [])"
                  :key="name"
                  class="text-xs bg-muted px-2 py-0.5 rounded-full"
                >
                  {{ name }}
                </span>
                <span
                  v-if="!chat.anchorTierName && (!chat.tierNames || chat.tierNames.length === 0)"
                  class="text-xs text-muted-foreground/50 italic"
                >не привязан</span>
              </div>

              <div class="text-xs text-muted-foreground tabular-nums shrink-0 w-8 text-right">
                {{ chat.activeUsers }}
              </div>

              <div
                class="shrink-0"
                @click.stop
              >
                <ConfirmDialog
                  title="Удалить чат?"
                  description="Чат будет удалён из системы подписок. Все привязки и доступы будут удалены."
                  confirm-label="Удалить"
                  @confirm="handleDeleteChat(chat.id)"
                >
                  <template #trigger>
                    <Button
                      variant="ghost"
                      size="sm"
                      class="h-8 w-8 p-0 text-muted-foreground hover:text-destructive"
                    >
                      <Trash2 class="h-3.5 w-3.5" />
                    </Button>
                  </template>
                </ConfirmDialog>
              </div>
            </div>

            <!-- Expanded panel -->
            <div
              v-if="expandedChatId === chat.id"
              class="border-t px-4 py-4 bg-muted/20"
            >
              <div
                v-if="expandedLoading"
                class="flex items-center gap-2 text-sm text-muted-foreground py-2"
              >
                <Loader2 class="h-4 w-4 animate-spin" />
                Загрузка...
              </div>

              <div
                v-else-if="expandedChatDetail"
                class="space-y-4"
              >
                <!-- Anchor -->
                <div>
                  <div class="flex items-center gap-2 mb-2">
                    <Anchor class="h-3.5 w-3.5 text-blue-500" />
                    <span class="text-xs font-medium uppercase tracking-wide text-muted-foreground">Якорь</span>
                    <span class="text-xs text-muted-foreground/60">&mdash; членство в чате = уровень подписки</span>
                  </div>
                  <div class="flex flex-wrap gap-2">
                    <button
                      v-for="tier in subscriptionService.tiers.value"
                      :key="tier.id"
                      class="inline-flex items-center gap-1.5 text-sm px-3 py-1.5 rounded-lg border transition-all duration-150 cursor-pointer"
                      :class="expandedChatDetail.anchorForTierID === tier.id
                        ? 'bg-blue-500 text-white border-blue-500 shadow-sm'
                        : 'bg-background border-border hover:border-blue-300 hover:bg-blue-50 dark:hover:bg-blue-950/30'"
                      :disabled="savingAction !== null"
                      @click="toggleAnchor(chat.id, tier.id)"
                    >
                      <Loader2
                        v-if="savingAction === `anchor-${tier.id}`"
                        class="h-3.5 w-3.5 animate-spin"
                      />
                      <Anchor
                        v-else-if="expandedChatDetail.anchorForTierID === tier.id"
                        class="h-3.5 w-3.5"
                      />
                      {{ tier.name }}
                    </button>
                  </div>
                </div>

                <!-- Content tiers -->
                <div v-if="!expandedChatDetail.anchorForTierID">
                  <div class="flex items-center gap-2 mb-2">
                    <span class="text-xs font-medium uppercase tracking-wide text-muted-foreground">Доступ по тирам</span>
                    <span class="text-xs text-muted-foreground/60">&mdash; какие уровни получат доступ к этому чату</span>
                  </div>
                  <div class="flex flex-wrap gap-2">
                    <button
                      v-for="tier in subscriptionService.tiers.value"
                      :key="tier.id"
                      class="inline-flex items-center gap-1.5 text-sm px-3 py-1.5 rounded-lg border transition-all duration-150 cursor-pointer"
                      :class="(expandedChatDetail.tierIDs || []).includes(tier.id)
                        ? 'bg-primary text-primary-foreground border-primary shadow-sm'
                        : 'bg-background border-border hover:border-primary/30 hover:bg-accent'"
                      :disabled="savingAction !== null"
                      @click="toggleContentTier(chat.id, tier.id)"
                    >
                      <Loader2
                        v-if="savingAction === `tier-${tier.id}`"
                        class="h-3.5 w-3.5 animate-spin"
                      />
                      {{ tier.name }}
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </Card>
        </div>
      </template>
    </div>

    <SubscriptionUserModal
      v-model:is-open="isUserModalOpen"
      :user-id="selectedUserId"
      @saved="subscriptionService.searchUsers(); subscriptionService.fetchStats()"
    />

    <SubscriptionChatModal
      v-model:is-open="isChatModalOpen"
      :chat-id="selectedChatId"
      @saved="handleChatSaved"
    />
  </AdminLayout>
</template>
