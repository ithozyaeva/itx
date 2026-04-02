<script setup lang="ts">
import type { SubscriptionChatDetail } from '@/services/subscriptionService'
import { Typography } from 'itx-ui-kit'
import { computed, onMounted, onUnmounted, ref } from 'vue'
import Anchor from '~icons/lucide/anchor'
import ChevronDown from '~icons/lucide/chevron-down'
import Eye from '~icons/lucide/eye'
import Plus from '~icons/lucide/plus'
import ShieldX from '~icons/lucide/shield-x'
import Trash2 from '~icons/lucide/trash-2'
import XIcon from '~icons/lucide/x'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import AdminLayout from '@/components/layout/AdminLayout.vue'
import SubscriptionChatModal from '@/components/modals/SubscriptionChatModal.vue'
import SubscriptionUserModal from '@/components/modals/SubscriptionUserModal.vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Checkbox } from '@/components/ui/checkbox'
import { Pagination, PaginationEllipsis, PaginationFirst, PaginationLast, PaginationList, PaginationListItem, PaginationNext, PaginationPrev } from '@/components/ui/pagination'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
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
const expandedSaving = ref(false)
const editAnchorTierID = ref<number | null>(null)
const editTierIDs = ref<number[]>([])

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
    const detail = await subscriptionService.getChatDetail(chatId)
    expandedChatDetail.value = detail
    if (detail) {
      editAnchorTierID.value = detail.anchorForTierID ?? null
      editTierIDs.value = detail.tierIDs ?? []
    }
  }
  finally {
    expandedLoading.value = false
  }
}

async function saveAnchor(chatId: number, tierID: number | null) {
  expandedSaving.value = true
  try {
    const success = await subscriptionService.updateChat(chatId, {
      anchorForTierID: tierID ?? undefined,
      clearAnchor: tierID === null,
      tierIDs: tierID !== null ? [] : undefined,
    })
    if (success) {
      editAnchorTierID.value = tierID
      if (tierID !== null)
        editTierIDs.value = []
      await subscriptionService.fetchChats()
      await subscriptionService.fetchStats()
    }
  }
  finally {
    expandedSaving.value = false
  }
}

function toggleContentTier(tierId: number) {
  const idx = editTierIDs.value.indexOf(tierId)
  if (idx >= 0)
    editTierIDs.value.splice(idx, 1)
  else
    editTierIDs.value.push(tierId)
}

async function saveContentTiers(chatId: number) {
  expandedSaving.value = true
  try {
    const success = await subscriptionService.updateChat(chatId, {
      clearAnchor: true,
      tierIDs: editTierIDs.value,
    })
    if (success) {
      editAnchorTierID.value = null
      await subscriptionService.fetchChats()
    }
  }
  finally {
    expandedSaving.value = false
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

        <Card>
          <CardContent>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead class="w-8" />
                  <TableHead>Название</TableHead>
                  <TableHead>Привязка</TableHead>
                  <TableHead class="text-right">
                    Участников
                  </TableHead>
                  <TableHead class="w-10" />
                </TableRow>
              </TableHeader>
              <TableBody>
                <TableRow
                  v-if="subscriptionService.chats.value.length === 0"
                  class="h-24"
                >
                  <TableCell
                    colspan="5"
                    class="text-center"
                  >
                    Чаты не найдены
                  </TableCell>
                </TableRow>
                <template
                  v-for="chat in subscriptionService.chats.value"
                  :key="chat.id"
                >
                  <TableRow
                    class="cursor-pointer hover:bg-muted/50"
                    @click="toggleChatExpand(chat.id)"
                  >
                    <TableCell class="w-8 pr-0">
                      <ChevronDown
                        class="h-4 w-4 text-muted-foreground transition-transform"
                        :class="{ 'rotate-180': expandedChatId === chat.id }"
                      />
                    </TableCell>
                    <TableCell>
                      <div class="font-medium">
                        {{ chat.title }}
                      </div>
                      <div class="text-xs text-muted-foreground">
                        {{ chat.chatType }} &middot; {{ chat.id }}
                      </div>
                    </TableCell>
                    <TableCell>
                      <div class="flex flex-wrap gap-1">
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
                          class="text-xs text-muted-foreground italic"
                        >не привязан</span>
                      </div>
                    </TableCell>
                    <TableCell class="text-right">
                      {{ chat.activeUsers }}
                    </TableCell>
                    <TableCell
                      class="w-10"
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
                            class="text-destructive"
                          >
                            <Trash2 class="h-4 w-4" />
                          </Button>
                        </template>
                      </ConfirmDialog>
                    </TableCell>
                  </TableRow>

                  <!-- Expanded panel -->
                  <TableRow v-if="expandedChatId === chat.id">
                    <TableCell
                      colspan="5"
                      class="bg-muted/30 p-0"
                    >
                      <div
                        v-if="expandedLoading"
                        class="px-6 py-4 text-sm text-muted-foreground"
                      >
                        Загрузка...
                      </div>
                      <div
                        v-else-if="expandedChatDetail"
                        class="px-6 py-4 space-y-4"
                      >
                        <!-- Anchor section -->
                        <div>
                          <div class="text-xs font-medium text-muted-foreground uppercase tracking-wide mb-2">
                            Якорь для тира
                          </div>
                          <div class="flex flex-wrap gap-2">
                            <Button
                              v-for="tier in subscriptionService.tiers.value"
                              :key="tier.id"
                              size="sm"
                              :variant="editAnchorTierID === tier.id ? 'default' : 'outline'"
                              :disabled="expandedSaving"
                              @click="saveAnchor(chat.id, editAnchorTierID === tier.id ? null : tier.id)"
                            >
                              <Anchor
                                v-if="editAnchorTierID === tier.id"
                                class="h-3.5 w-3.5 mr-1.5"
                              />
                              {{ tier.name }}
                              <span class="text-xs text-muted-foreground ml-1">(lvl {{ tier.level }})</span>
                            </Button>
                            <Button
                              v-if="editAnchorTierID !== null"
                              size="sm"
                              variant="ghost"
                              class="text-muted-foreground"
                              :disabled="expandedSaving"
                              @click="saveAnchor(chat.id, null)"
                            >
                              <XIcon class="h-3.5 w-3.5 mr-1" />
                              Снять якорь
                            </Button>
                          </div>
                        </div>

                        <!-- Content tiers section -->
                        <div v-if="editAnchorTierID === null">
                          <div class="text-xs font-medium text-muted-foreground uppercase tracking-wide mb-2">
                            Доступен для тиров (content)
                          </div>
                          <div class="flex flex-wrap items-center gap-3">
                            <label
                              v-for="tier in subscriptionService.tiers.value"
                              :key="tier.id"
                              class="flex items-center gap-2 cursor-pointer"
                            >
                              <Checkbox
                                :checked="editTierIDs.includes(tier.id)"
                                :disabled="expandedSaving"
                                @update:checked="toggleContentTier(tier.id)"
                              />
                              <span class="text-sm">{{ tier.name }}</span>
                            </label>
                            <Button
                              size="sm"
                              :disabled="expandedSaving"
                              @click="saveContentTiers(chat.id)"
                            >
                              {{ expandedSaving ? 'Сохранение...' : 'Сохранить' }}
                            </Button>
                          </div>
                        </div>
                      </div>
                    </TableCell>
                  </TableRow>
                </template>
              </TableBody>
            </Table>
          </CardContent>
        </Card>
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
