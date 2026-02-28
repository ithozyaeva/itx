<script setup lang="ts">
import { Label, Typography } from 'itx-ui-kit'
import { onMounted, onUnmounted, ref } from 'vue'
import Pencil from '~icons/lucide/pencil'
import Plus from '~icons/lucide/plus'
import Trash from '~icons/lucide/trash'
import BulkActionBar from '@/components/BulkActionBar.vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import AdminLayout from '@/components/layout/AdminLayout.vue'
import EventsModal from '@/components/modals/EventsModal.vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Pagination, PaginationEllipsis, PaginationFirst, PaginationLast, PaginationList, PaginationListItem, PaginationNext, PaginationPrev } from '@/components/ui/pagination'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { useBulkSelection } from '@/composables/useBulkSelection'
import { useModal } from '@/composables/useModal'
import { bulkService } from '@/services/bulkService'
import { eventsService } from '@/services/eventsService'

onMounted(eventsService.search)
onUnmounted(eventsService.clearPagination)

const selectedEventId = ref<number>()
const { open, isOpen } = useModal()
const bulk = useBulkSelection()

async function handleBulkDelete() {
  await bulkService.deleteEvents(bulk.ids.value)
  bulk.clearSelection()
  eventsService.search()
}

/**
 * Выбор сущности для редактирования.
 *
 * @param entityId - ID сущности.
 */
function selectEvent(entityId: number) {
  selectedEventId.value = entityId
  open()
}
</script>

<template>
  <AdminLayout>
    <div class="space-y-6">
      <div class="flex justify-between items-center">
        <Typography variant="h2" as="h1">
          События сообщества
        </Typography>
        <Button @click="open">
          <Plus class="mr-2 h-4 w-4" />
          Добавить событие
        </Button>
      </div>

      <Card>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead class="w-10">
                  <input
                    type="checkbox"
                    :checked="bulk.count.value === eventsService.items.value.items.length && bulk.count.value > 0"
                    @change="bulk.toggleAll(eventsService.items.value.items.map(e => e.id))"
                  >
                </TableHead>
                <TableHead>Дата</TableHead>
                <TableHead>Название</TableHead>
                <TableHead>Темы</TableHead>
                <TableHead>Тип</TableHead>
                <TableHead>Место</TableHead>
                <TableHead>Спикеры</TableHead>
                <TableHead />
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow v-if="eventsService.items.value.total === 0" class="h-24">
                <TableCell colspan="7" class="text-center">
                  События не найдены
                </TableCell>
              </TableRow>
              <TableRow v-for="event in eventsService.items.value.items" :key="event.id">
                <TableCell>
                  <input type="checkbox" :checked="bulk.isSelected(event.id)" @change="bulk.toggleItem(event.id)">
                </TableCell>
                <TableCell>
                  <div class="flex flex-col">
                    <span>{{ new Date(event.date).toLocaleString() }} ({{ event.timezone || 'UTC' }})</span>
                    <span v-if="event.isRepeating && event.repeatPeriod" class="text-xs text-muted-foreground italic">
                      Повторяется
                    </span>
                  </div>
                </TableCell>
                <TableCell class="max-w-[200px] truncate">
                  {{ event.title }}
                </TableCell>
                <TableCell>
                  <div v-if="event.eventTags.length > 0" class="flex flex-wrap items-center gap-1">
                    <Label v-for="tag in event.eventTags" :key="tag.id" :title="tag.name">{{ tag.name.length > 30 ? `${tag.name.slice(0, 24)}...` : tag.name }}</Label>
                  </div>
                </TableCell>
                <TableCell>{{ event.eventType }}</TableCell>
                <TableCell class="max-w-[200px] truncate">
                  {{ event.place }}
                </TableCell>
                <TableCell class="max-w-[200px] truncate">
                  {{ event.hosts?.map((host) => `${host.firstName} ${host.lastName}`).join(', ') }}
                </TableCell>
                <TableCell>
                  <Button variant="ghost" size="sm" @click="selectEvent(event.id)">
                    <Pencil class="h-4 w-4" />
                  </Button>
                  <ConfirmDialog
                    title="Удалить событие?"
                    description="Событие будет удалено без возможности восстановления."
                    confirm-label="Удалить"
                    @confirm="eventsService.delete(event.id)"
                  >
                    <template #trigger>
                      <Button variant="ghost" size="sm" :disabled="eventsService.isLoading.value">
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
        <Pagination v-slot="{ page }" :items-per-page="10" :total="eventsService.items.value.total" :sibling-count="1" show-edges :default-page="1">
          <PaginationList v-slot="{ items }" class="flex items-center gap-1">
            <PaginationFirst />
            <PaginationPrev />

            <template v-for="(item, index) in items">
              <PaginationListItem v-if="item.type === 'page'" :key="index" :value="item.value" as-child>
                <Button class="w-10 h-10 p-0" :variant="item.value === page ? 'default' : 'outline'" @click="eventsService.changePagination(item.value)">
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
    </div>
    <EventsModal v-model:is-open="isOpen" :event-id="selectedEventId" @saved="eventsService.search" />
    <BulkActionBar
      :count="bulk.count.value"
      :actions="[{ label: 'Удалить', handler: handleBulkDelete }]"
      @clear="bulk.clearSelection"
    />
  </AdminLayout>
</template>
