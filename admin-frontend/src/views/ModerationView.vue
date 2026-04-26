<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import Globe from '~icons/lucide/globe'
import ShieldX from '~icons/lucide/shield-x'
import Vote from '~icons/lucide/vote'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import AdminLayout from '@/components/layout/AdminLayout.vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Typography } from '@/components/ui/typography'
import { usePermissions } from '@/composables/usePermissions'
import { moderationService } from '@/services/moderationService'

const { hasPermission } = usePermissions()
const canEdit = computed(() => hasPermission.value('can_edit_admin_moderation'))

type Tab = 'sanctions' | 'global' | 'votings'
const activeTab = ref<Tab>('sanctions')

function switchTab(tab: Tab) {
  activeTab.value = tab
  loadActiveTab()
}

function loadActiveTab() {
  if (activeTab.value === 'sanctions')
    moderationService.fetchSanctions()
  else if (activeTab.value === 'global')
    moderationService.fetchGlobalBans()
  else
    moderationService.fetchVotebans()
}

onMounted(loadActiveTab)

const ACTION_LABEL: Record<string, string> = {
  ban: 'Бан',
  mute: 'Мут',
  voteban_mute: 'Voteban (мут)',
  voteban_kick: 'Voteban (кик)',
}

function targetDisplay(target: { targetUsername: string, targetFirstName: string, targetUserId: number }): string {
  if (target.targetUsername)
    return `@${target.targetUsername}`
  if (target.targetFirstName)
    return target.targetFirstName
  return `id=${target.targetUserId}`
}

function chatDisplay(row: { chatTitle: string, chatId: number }): string {
  return row.chatTitle || `chat=${row.chatId}`
}

function fmtDateTime(s: string | null): string {
  if (!s)
    return '—'
  const d = new Date(s)
  if (Number.isNaN(d.getTime()))
    return s
  return d.toLocaleString('ru-RU', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}

function fmtExpires(s: string | null): string {
  if (!s)
    return 'permanent'
  return fmtDateTime(s)
}

function fmtDuration(seconds: number | null): string {
  if (seconds == null || seconds <= 0)
    return '∞'
  if (seconds % 86400 === 0)
    return `${seconds / 86400}д`
  if (seconds % 3600 === 0)
    return `${seconds / 3600}ч`
  if (seconds % 60 === 0)
    return `${seconds / 60}м`
  return `${seconds}с`
}
</script>

<template>
  <AdminLayout>
    <div class="space-y-6">
      <Typography variant="h2" as="h1">
        Модерация
      </Typography>

      <!-- Tabs -->
      <div class="flex gap-2 border-b pb-2">
        <Button
          :variant="activeTab === 'sanctions' ? 'default' : 'ghost'"
          size="sm"
          @click="switchTab('sanctions')"
        >
          <ShieldX class="mr-2 h-4 w-4" />
          Активные санкции
          <span v-if="activeTab === 'sanctions' && moderationService.sanctions.value.length" class="ml-2 text-xs opacity-70">
            ({{ moderationService.sanctions.value.length }})
          </span>
        </Button>
        <Button
          :variant="activeTab === 'global' ? 'default' : 'ghost'"
          size="sm"
          @click="switchTab('global')"
        >
          <Globe class="mr-2 h-4 w-4" />
          Глобальные баны
          <span v-if="activeTab === 'global' && moderationService.globalBans.value.length" class="ml-2 text-xs opacity-70">
            ({{ moderationService.globalBans.value.length }})
          </span>
        </Button>
        <Button
          :variant="activeTab === 'votings' ? 'default' : 'ghost'"
          size="sm"
          @click="switchTab('votings')"
        >
          <Vote class="mr-2 h-4 w-4" />
          Активные голосования
          <span v-if="activeTab === 'votings' && moderationService.votebans.value.length" class="ml-2 text-xs opacity-70">
            ({{ moderationService.votebans.value.length }})
          </span>
        </Button>

        <Button
          variant="ghost"
          size="sm"
          class="ml-auto"
          :disabled="moderationService.isLoading.value"
          @click="loadActiveTab"
        >
          Обновить
        </Button>
      </div>

      <!-- Активные санкции -->
      <template v-if="activeTab === 'sanctions'">
        <Card>
          <CardContent class="p-0">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Чат</TableHead>
                  <TableHead>Юзер</TableHead>
                  <TableHead>Действие</TableHead>
                  <TableHead>Длительность</TableHead>
                  <TableHead>До</TableHead>
                  <TableHead>Когда</TableHead>
                  <TableHead>Причина</TableHead>
                  <TableHead class="w-32 text-right">
                    Действия
                  </TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                <TableRow v-if="!moderationService.sanctions.value.length" class="h-24">
                  <TableCell colspan="8" class="text-center text-muted-foreground">
                    Активных санкций нет
                  </TableCell>
                </TableRow>
                <TableRow v-for="row in moderationService.sanctions.value" :key="row.id">
                  <TableCell>{{ chatDisplay(row) }}</TableCell>
                  <TableCell>
                    <a
                      :href="`tg://user?id=${row.targetUserId}`"
                      class="text-primary hover:underline"
                    >{{ targetDisplay(row) }}</a>
                    <div class="text-xs text-muted-foreground">
                      id={{ row.targetUserId }}
                    </div>
                  </TableCell>
                  <TableCell>{{ ACTION_LABEL[row.action] ?? row.action }}</TableCell>
                  <TableCell>{{ fmtDuration(row.durationSeconds) }}</TableCell>
                  <TableCell>{{ fmtExpires(row.expiresAt) }}</TableCell>
                  <TableCell>{{ fmtDateTime(row.createdAt) }}</TableCell>
                  <TableCell class="max-w-[200px] truncate">
                    {{ row.reason || '—' }}
                  </TableCell>
                  <TableCell class="text-right">
                    <ConfirmDialog
                      v-if="canEdit"
                      title="Снять санкцию?"
                      description="Бот разблокирует пользователя в этом чате через несколько секунд."
                      confirm-label="Снять"
                      @confirm="moderationService.revokeSanction(row.id)"
                    >
                      <template #trigger>
                        <Button
                          variant="ghost"
                          size="sm"
                          class="text-destructive"
                          :disabled="moderationService.isLoading.value"
                        >
                          Снять
                        </Button>
                      </template>
                    </ConfirmDialog>
                    <span v-else class="text-xs text-muted-foreground">view-only</span>
                  </TableCell>
                </TableRow>
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      </template>

      <!-- Глобальные баны -->
      <template v-else-if="activeTab === 'global'">
        <Card>
          <CardContent class="p-0">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>User ID</TableHead>
                  <TableHead>Кем выдан</TableHead>
                  <TableHead>Когда</TableHead>
                  <TableHead>До</TableHead>
                  <TableHead>Причина</TableHead>
                  <TableHead class="w-32 text-right">
                    Действия
                  </TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                <TableRow v-if="!moderationService.globalBans.value.length" class="h-24">
                  <TableCell colspan="6" class="text-center text-muted-foreground">
                    Активных глобальных банов нет
                  </TableCell>
                </TableRow>
                <TableRow v-for="b in moderationService.globalBans.value" :key="b.userId">
                  <TableCell>
                    <a
                      :href="`tg://user?id=${b.userId}`"
                      class="text-primary hover:underline"
                    >id={{ b.userId }}</a>
                  </TableCell>
                  <TableCell>{{ b.bannedBy }}</TableCell>
                  <TableCell>{{ fmtDateTime(b.createdAt) }}</TableCell>
                  <TableCell>{{ fmtExpires(b.expiresAt) }}</TableCell>
                  <TableCell class="max-w-[300px] truncate">
                    {{ b.reason || '—' }}
                  </TableCell>
                  <TableCell class="text-right">
                    <ConfirmDialog
                      v-if="canEdit"
                      title="Снять глобальный бан?"
                      description="Бот разблокирует пользователя во всех известных чатах."
                      confirm-label="Снять"
                      @confirm="moderationService.revokeGlobalBan(b.userId)"
                    >
                      <template #trigger>
                        <Button
                          variant="ghost"
                          size="sm"
                          class="text-destructive"
                          :disabled="moderationService.isLoading.value"
                        >
                          Снять
                        </Button>
                      </template>
                    </ConfirmDialog>
                    <span v-else class="text-xs text-muted-foreground">view-only</span>
                  </TableCell>
                </TableRow>
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      </template>

      <!-- Голосования -->
      <template v-else>
        <Card>
          <CardContent class="p-0">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Чат</TableHead>
                  <TableHead>Цель</TableHead>
                  <TableHead>Инициатор</TableHead>
                  <TableHead>Голоса</TableHead>
                  <TableHead>Окно</TableHead>
                  <TableHead>До</TableHead>
                  <TableHead class="w-32 text-right">
                    Действия
                  </TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                <TableRow v-if="!moderationService.votebans.value.length" class="h-24">
                  <TableCell colspan="7" class="text-center text-muted-foreground">
                    Активных голосований нет
                  </TableCell>
                </TableRow>
                <TableRow v-for="v in moderationService.votebans.value" :key="v.id">
                  <TableCell>{{ v.chatTitle || `chat=${v.chatId}` }}</TableCell>
                  <TableCell>
                    <a
                      :href="`tg://user?id=${v.targetUserId}`"
                      class="text-primary hover:underline"
                    >{{ v.targetUsername ? `@${v.targetUsername}` : v.targetFirstName || `id=${v.targetUserId}` }}</a>
                  </TableCell>
                  <TableCell>id={{ v.initiatorUserId }}</TableCell>
                  <TableCell>
                    <span class="text-emerald-500">✅ {{ v.votesFor }}</span>
                    <span class="mx-2 text-muted-foreground">/ {{ v.requiredVotes }}</span>
                    <span class="text-rose-500">❌ {{ v.votesAgainst }}</span>
                  </TableCell>
                  <TableCell>{{ fmtDuration(v.muteSeconds) }}</TableCell>
                  <TableCell>{{ fmtDateTime(v.expiresAt) }}</TableCell>
                  <TableCell class="text-right">
                    <ConfirmDialog
                      v-if="canEdit"
                      title="Отменить голосование?"
                      description="Голосование закроется без санкций. Цель остаётся в чате."
                      confirm-label="Отменить"
                      @confirm="moderationService.cancelVoteban(v.id)"
                    >
                      <template #trigger>
                        <Button
                          variant="ghost"
                          size="sm"
                          class="text-destructive"
                          :disabled="moderationService.isLoading.value"
                        >
                          Отменить
                        </Button>
                      </template>
                    </ConfirmDialog>
                    <span v-else class="text-xs text-muted-foreground">view-only</span>
                  </TableCell>
                </TableRow>
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      </template>
    </div>
  </AdminLayout>
</template>
