<script setup lang="ts">
import { Typography } from 'itx-ui-kit'
import { onMounted, onUnmounted, ref } from 'vue'
import Download from '~icons/lucide/download'
import Pencil from '~icons/lucide/pencil'
import Plus from '~icons/lucide/plus'
import Trash from '~icons/lucide/trash'
import BulkActionBar from '@/components/BulkActionBar.vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import AdminLayout from '@/components/layout/AdminLayout.vue'
import MemberSearchFilters from '@/components/MemberSearchFilters.vue'
import MemberModal from '@/components/modals/MemberModal.vue'
import MentorModal from '@/components/modals/MentorModal.vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'

import { Pagination, PaginationEllipsis, PaginationFirst, PaginationList, PaginationListItem, PaginationNext, PaginationPrev } from '@/components/ui/pagination'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { useBulkSelection } from '@/composables/useBulkSelection'
import { useDictionary } from '@/composables/useDictionary'
import { downloadFile } from '@/lib/utils'
import { bulkService } from '@/services/bulkService'
import { memberService } from '@/services/memberService'

const isModalOpen = ref(false)
const currentMemberId = ref<number | null>(null)

const isMentorModalOpen = ref(false)
const selectedMemberId = ref<number | null>(null)

/**
 * Обработчик открытия модального окна для одобавления участника.
 */
function openAddModal() {
  currentMemberId.value = null
  isModalOpen.value = true
}

/**
 * Обработчик открытия модального окна для редактирования участника.
 *
 * @param id - ID участника.
 */
function openEditModal(id: number) {
  currentMemberId.value = id
  isModalOpen.value = true
}

/**
 * Обработчик открытия модального окна для создания ментора из участинка.
 *
 * @param memberId - ID участника.
 */
function handleMakeMentor(memberId: number) {
  selectedMemberId.value = memberId
  isMentorModalOpen.value = true
}

const { memberRolesObject } = useDictionary(['memberRoles'])
const bulk = useBulkSelection()

async function handleBulkDelete() {
  await bulkService.deleteMembers(bulk.ids.value)
  bulk.clearSelection()
  memberService.search()
}

onMounted(memberService.search)
onUnmounted(memberService.clearPagination)
</script>

<template>
  <AdminLayout>
    <div class="space-y-6">
      <div class="flex justify-between items-center">
        <Typography variant="h2" as="h1">
          Участники сообщества
        </Typography>
        <div class="flex gap-2">
          <Button variant="outline" @click="downloadFile('members/export/csv', 'members.csv')">
            <Download class="mr-2 h-4 w-4" />
            Экспорт CSV
          </Button>
          <Button
            v-permission="'can_edit_admin_members'"
            @click="openAddModal"
          >
            <Plus class="mr-2 h-4 w-4" />
            Добавить участника
          </Button>
        </div>
      </div>
      <MemberSearchFilters @apply="memberService.search" />
      <Card>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead class="w-10">
                  <input type="checkbox" :checked="bulk.count.value === memberService.items.value.items.length && bulk.count.value > 0" @change="bulk.toggleAll(memberService.items.value.items.map(m => m.id))">
                </TableHead>
                <TableHead>Имя</TableHead>
                <TableHead>Telegram</TableHead>
                <TableHead>Роли</TableHead>
                <TableHead />
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow v-if="memberService.items.value.total === 0" class="h-24">
                <TableCell colspan="7" class="text-center">
                  Участники не найдены
                </TableCell>
              </TableRow>
              <TableRow v-for="member in memberService.items.value.items" :key="member.id">
                <TableCell>
                  <input type="checkbox" :checked="bulk.isSelected(member.id)" @change="bulk.toggleItem(member.id)">
                </TableCell>
                <TableCell>{{ member.firstName ?? "" }} {{ member.lastName ?? "" }}</TableCell>
                <TableCell>{{ member.tg }}</TableCell>
                <TableCell>{{ member.roles.map(item => memberRolesObject[item])?.join(', ') || '' }}</TableCell>
                <TableCell class="text-right">
                  <Button
                    v-permission="'can_edit_admin_members'"
                    variant="ghost"
                    size="sm"
                    @click="openEditModal(member.id)"
                  >
                    <Pencil class="h-4 w-4" />
                  </Button>
                  <ConfirmDialog
                    title="Удалить участника?"
                    description="Участник будет удалён без возможности восстановления."
                    confirm-label="Удалить"
                    @confirm="memberService.delete(member.id)"
                  >
                    <template #trigger>
                      <Button
                        v-permission="'can_edit_admin_members'"
                        variant="ghost"
                        size="sm"
                        class="text-destructive"
                        :disabled="memberService.isLoading.value"
                      >
                        <Trash class="h-4 w-4" />
                      </Button>
                    </template>
                  </ConfirmDialog>
                  <Button
                    v-if="!member.roles?.includes('MENTOR')"
                    v-permission="'can_edit_admin_mentors'"
                    variant="ghost"
                    size="sm"
                    class="text-primary"
                    @click="handleMakeMentor(member.id)"
                  >
                    Сделать ментором
                  </Button>
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </CardContent>
      </Card>
      <div class="mt-4 flex justify-end">
        <Pagination v-slot="{ page }" :items-per-page="10" :total="memberService.items.value.total" :sibling-count="1" show-edges :default-page="1">
          <PaginationList v-slot="{ items }" class="flex items-center gap-1">
            <PaginationFirst @click="memberService.changePagination(1)" />
            <PaginationPrev @click="memberService.changePagination(page - 1)" />

            <template v-for="(item, index) in items">
              <PaginationListItem v-if="item.type === 'page'" :key="index" :value="item.value" as-child>
                <Button class="w-10 h-10 p-0" :variant="item.value === page ? 'default' : 'outline'" @click="memberService.changePagination(item.value)">
                  {{ item.value }}
                </Button>
              </PaginationListItem>
              <PaginationEllipsis v-else :key="item.type" :index="index" />
            </template>

            <PaginationNext @click="memberService.changePagination(page + 1)" />
          </PaginationList>
        </Pagination>
      </div>
    </div>

    <MemberModal
      v-model:is-open="isModalOpen"
      :member-id="currentMemberId"
      @saved="memberService.search"
    />

    <!-- Модальное окно для создания ментора -->
    <MentorModal
      v-model:is-open="isMentorModalOpen"
      :member-id="selectedMemberId!"
      @saved="memberService.search"
    />
    <BulkActionBar
      :count="bulk.count.value"
      :actions="[{ label: 'Удалить', handler: handleBulkDelete }]"
      @clear="bulk.clearSelection"
    />
  </AdminLayout>
</template>
