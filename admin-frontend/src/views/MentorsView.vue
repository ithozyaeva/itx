<script setup lang="ts">
import type { Mentor } from '@/models/mentors'
import { Tag, Typography } from 'itx-ui-kit'
import { onMounted, onUnmounted, ref } from 'vue'
import Download from '~icons/lucide/download'
import Pencil from '~icons/lucide/pencil'
import Trash from '~icons/lucide/trash'
import BulkActionBar from '@/components/BulkActionBar.vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import AdminLayout from '@/components/layout/AdminLayout.vue'
import MentorModal from '@/components/modals/MentorModal.vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Pagination, PaginationEllipsis, PaginationFirst, PaginationLast, PaginationList, PaginationListItem, PaginationNext, PaginationPrev } from '@/components/ui/pagination'

import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { useBulkSelection } from '@/composables/useBulkSelection'
import { downloadFile } from '@/lib/utils'
import { bulkService } from '@/services/bulkService'
import { mentorService } from '@/services/mentorService'

const isModalOpen = ref(false)
const selectedMentor = ref<Mentor | null>(null)

/**
 * Обработчик изменения ментора.
 *
 * @param mentor - Ментор.
 */
function handleEditMentor(mentor: Mentor) {
  selectedMentor.value = mentor
  isModalOpen.value = true
}

const bulk = useBulkSelection()

async function handleBulkDelete() {
  await bulkService.deleteMentors(bulk.ids.value)
  bulk.clearSelection()
  mentorService.search()
}

onMounted(mentorService.search)
onUnmounted(mentorService.clearPagination)
</script>

<template>
  <AdminLayout>
    <div class="space-y-6">
      <div class="flex justify-between items-center">
        <Typography variant="h2" as="h1">
          Менторы
        </Typography>
        <Button variant="outline" @click="downloadFile('mentors/export/csv', 'mentors.csv')">
          <Download class="mr-2 h-4 w-4" />
          Экспорт CSV
        </Button>
      </div>
      <Card>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead class="w-10">
                  <input type="checkbox" :checked="bulk.count.value === mentorService.items.value.items.length && bulk.count.value > 0" @change="bulk.toggleAll(mentorService.items.value.items.map(m => m.id))">
                </TableHead>
                <TableHead>Имя</TableHead>
                <TableHead>Telegram</TableHead>
                <TableHead>Должность</TableHead>
                <TableHead>Опыт</TableHead>
                <TableHead>Теги</TableHead>
                <TableHead class="text-right" />
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow v-if="mentorService.items.value.total === 0" class="h-24">
                <TableCell colspan="7" class="text-center">
                  Менторы не найдены
                </TableCell>
              </TableRow>
              <TableRow v-for="mentor in mentorService.items.value.items" :key="mentor.id">
                <TableCell>
                  <input type="checkbox" :checked="bulk.isSelected(mentor.id)" @change="bulk.toggleItem(mentor.id)">
                </TableCell>
                <TableCell>{{ mentor.firstName ?? "" }} {{ mentor.lastName ?? "" }}</TableCell>
                <TableCell>{{ mentor.tg }}</TableCell>
                <TableCell>{{ mentor.occupation }}</TableCell>
                <TableCell>{{ mentor.experience }}</TableCell>
                <TableCell>
                  <div class="flex flex-wrap gap-1">
                    <Tag
                      v-for="tag in mentor.profTags"
                      :key="tag.id"
                    >
                      {{ tag.title }}
                    </Tag>
                  </div>
                </TableCell>
                <TableCell class="text-right">
                  <Button
                    v-permission="'can_edit_admin_mentors'"
                    variant="ghost"
                    size="sm"
                    @click="handleEditMentor(mentor)"
                  >
                    <Pencil class="h-4 w-4" />
                  </Button>
                  <ConfirmDialog
                    title="Удалить ментора?"
                    description="Ментор будет удалён без возможности восстановления."
                    confirm-label="Удалить"
                    @confirm="mentorService.delete(mentor.id)"
                  >
                    <template #trigger>
                      <Button
                        v-permission="'can_edit_admin_mentors'"
                        variant="ghost"
                        size="sm"
                        class="text-destructive"
                        :disabled="mentorService.isLoading.value"
                      >
                        <Trash class="h-4 w-4" />
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
        <Pagination v-slot="{ page }" :items-per-page="10" :total="mentorService.items.value.total" :sibling-count="1" show-edges :default-page="1">
          <PaginationList v-slot="{ items }" class="flex items-center gap-1">
            <PaginationFirst />
            <PaginationPrev />

            <template v-for="(item, index) in items">
              <PaginationListItem v-if="item.type === 'page'" :key="index" :value="item.value" as-child>
                <Button class="w-10 h-10 p-0" :variant="item.value === page ? 'default' : 'outline'" @click="mentorService.changePagination(item.value)">
                  {{ item.value }}
                </Button>
              </PaginationListItem>
              <PaginationEllipsis v-else :key="item.type" :index="index" />
            </template>

            <PaginationNext />
            <PaginationLast />
          </PaginationList>
        </Pagination>
      </div>

      <MentorModal
        v-model:is-open="isModalOpen"
        :mentor-id="selectedMentor?.id!"
        :member-id="selectedMentor?.memberId!"
        @saved="mentorService.search"
      />
    </div>
    <BulkActionBar
      :count="bulk.count.value"
      :actions="[{ label: 'Удалить', handler: handleBulkDelete }]"
      @clear="bulk.clearSelection"
    />
  </AdminLayout>
</template>
