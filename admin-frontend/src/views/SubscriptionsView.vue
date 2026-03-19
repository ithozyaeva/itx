<script setup lang="ts">
import { Typography } from 'itx-ui-kit'
import { computed, onMounted, onUnmounted, ref } from 'vue'
import Eye from '~icons/lucide/eye'
import ShieldX from '~icons/lucide/shield-x'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import AdminLayout from '@/components/layout/AdminLayout.vue'
import SubscriptionUserModal from '@/components/modals/SubscriptionUserModal.vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Pagination, PaginationEllipsis, PaginationFirst, PaginationLast, PaginationList, PaginationListItem, PaginationNext, PaginationPrev } from '@/components/ui/pagination'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { subscriptionService } from '@/services/subscriptionService'

type Tab = 'overview' | 'users' | 'chats'

const activeTab = ref<Tab>('overview')
const selectedUserId = ref<number | null>(null)
const isUserModalOpen = ref(false)

const stats = computed(() => subscriptionService.stats.value)

function openUserModal(userId: number) {
  selectedUserId.value = userId
  isUserModalOpen.value = true
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
        <Card>
          <CardContent>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>ID</TableHead>
                  <TableHead>Название</TableHead>
                  <TableHead>Тип</TableHead>
                  <TableHead>Роль</TableHead>
                  <TableHead class="text-right">
                    Пользователей
                  </TableHead>
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
                <TableRow
                  v-for="chat in subscriptionService.chats.value"
                  :key="chat.id"
                >
                  <TableCell>
                    <code class="text-xs">{{ chat.id }}</code>
                  </TableCell>
                  <TableCell class="font-medium">
                    {{ chat.title }}
                  </TableCell>
                  <TableCell>{{ chat.chatType }}</TableCell>
                  <TableCell>
                    <span
                      v-if="chat.anchorTierName"
                      class="inline-flex items-center gap-1 text-blue-500 text-xs font-medium"
                    >
                      ANCHOR &rarr; {{ chat.anchorTierName }}
                    </span>
                    <span
                      v-else
                      class="text-xs text-muted-foreground"
                    >content</span>
                  </TableCell>
                  <TableCell class="text-right">
                    {{ chat.activeUsers }}
                  </TableCell>
                </TableRow>
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
  </AdminLayout>
</template>
