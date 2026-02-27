<script setup lang="ts">
import type { ReviewStatus } from '@/models/reviewOnCommunity'
import { Typography } from 'itx-ui-kit'
import { onMounted, onUnmounted, ref } from 'vue'
import Pencil from '~icons/lucide/pencil'
import Plus from '~icons/lucide/plus'
import Trash from '~icons/lucide/trash'
import BulkActionBar from '@/components/BulkActionBar.vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import AdminLayout from '@/components/layout/AdminLayout.vue'
import RevievOnCommunityModal from '@/components/modals/RevievOnCommunityModal.vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Pagination, PaginationEllipsis, PaginationFirst, PaginationLast, PaginationList, PaginationListItem, PaginationNext, PaginationPrev } from '@/components/ui/pagination'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'

import { useBulkSelection } from '@/composables/useBulkSelection'
import { useDictionary } from '@/composables/useDictionary'
import { useModal } from '@/composables/useModal'
import { bulkService } from '@/services/bulkService'
import { reviewOnCommunityService } from '@/services/reviewOnCommunityService'

onMounted(reviewOnCommunityService.search)
onUnmounted(reviewOnCommunityService.clearPagination)

const selectedReviewId = ref<number>()
const { open, isOpen } = useModal()

/**
 * Выбор отзыва для редактирования.
 *
 * @param reviewId - ID отзыва.
 */
function selectReview(reviewId: number) {
  selectedReviewId.value = reviewId
  open()
}

const { reviewStatuses } = useDictionary<ReviewStatus>(['reviewStatuses'])
const bulk = useBulkSelection()

async function handleBulkDelete() {
  await bulkService.deleteReviews(bulk.ids.value)
  bulk.clearSelection()
  reviewOnCommunityService.search()
}

async function handleBulkApprove() {
  await bulkService.approveReviews(bulk.ids.value)
  bulk.clearSelection()
  reviewOnCommunityService.search()
}
</script>

<template>
  <AdminLayout>
    <div class="space-y-6">
      <div class="flex justify-between items-center">
        <Typography variant="h2" as="h1">
          Отзывы на сообщество
        </Typography>
        <Button @click="open">
          <Plus class="mr-2 h-4 w-4" />
          Добавить отзыв
        </Button>
      </div>

      <Card>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead class="w-10">
                  <input type="checkbox" :checked="bulk.count.value === reviewOnCommunityService.items.value.items.length && bulk.count.value > 0" @change="bulk.toggleAll(reviewOnCommunityService.items.value.items.map(r => r.id))">
                </TableHead>
                <TableHead>Автор</TableHead>
                <TableHead>Текст отзыва</TableHead>
                <TableHead>Дата</TableHead>
                <TableHead>Статус</TableHead>
                <TableHead />
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow v-if="reviewOnCommunityService.items.value.total === 0" class="h-24">
                <TableCell colspan="7" class="text-center">
                  Отзывы не найдены
                </TableCell>
              </TableRow>
              <TableRow v-for="review in reviewOnCommunityService.items.value.items" :key="review.id">
                <TableCell>
                  <input type="checkbox" :checked="bulk.isSelected(review.id)" @change="bulk.toggleItem(review.id)">
                </TableCell>
                <TableCell>
                  Имя: {{ review.author.firstName ?? "" }} {{ review.author.lastName ?? "" }} <br> Telegram: {{
                    review.author.tg }}
                </TableCell>
                <TableCell>{{ review.text }}</TableCell>
                <TableCell>{{ new Date(review.date).toLocaleDateString() }}</TableCell>
                <TableCell>{{ reviewStatuses.find((status) => status.value === review.status)?.label }}</TableCell>
                <TableCell>
                  <div class="flex items-center justify-end">
                    <ConfirmDialog
                      v-if="review.status !== 'APPROVED'"
                      title="Опубликовать отзыв?"
                      description="Отзыв станет виден на сайте."
                      confirm-label="Опубликовать"
                      variant="default"
                      @confirm="reviewOnCommunityService.approve(review.id)"
                    >
                      <template #trigger>
                        <Button :disabled="reviewOnCommunityService.isLoading.value">
                          Опубликовать
                        </Button>
                      </template>
                    </ConfirmDialog>
                    <Button variant="ghost" size="sm" @click="selectReview(review.id)">
                      <Pencil class="h-4 w-4" />
                    </Button>
                    <ConfirmDialog
                      title="Удалить отзыв?"
                      description="Отзыв будет удалён без возможности восстановления."
                      confirm-label="Удалить"
                      @confirm="reviewOnCommunityService.delete(review.id)"
                    >
                      <template #trigger>
                        <Button variant="ghost" size="sm" :disabled="reviewOnCommunityService.isLoading.value">
                          <Trash class="h-4 w-4" />
                        </Button>
                      </template>
                    </ConfirmDialog>
                  </div>
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </CardContent>
      </Card>
      <div class="mt-4 flex justify-end">
        <Pagination
          v-slot="{ page }" :items-per-page="10" :total="reviewOnCommunityService.items.value.total"
          :sibling-count="1" show-edges :default-page="1"
        >
          <PaginationList v-slot="{ items }" class="flex items-center gap-1">
            <PaginationFirst />
            <PaginationPrev />

            <template v-for="(item, index) in items">
              <PaginationListItem v-if="item.type === 'page'" :key="index" :value="item.value" as-child>
                <Button
                  class="w-10 h-10 p-0" :variant="item.value === page ? 'default' : 'outline'"
                  @click="reviewOnCommunityService.changePagination(item.value)"
                >
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
    <RevievOnCommunityModal
      v-model:is-open="isOpen" :review-id="selectedReviewId"
      @saved="reviewOnCommunityService.search"
    />
    <BulkActionBar
      :count="bulk.count.value"
      :actions="[
        { label: 'Опубликовать', variant: 'default', handler: handleBulkApprove },
        { label: 'Удалить', handler: handleBulkDelete },
      ]"
      @clear="bulk.clearSelection"
    />
  </AdminLayout>
</template>
