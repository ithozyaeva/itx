<script setup lang="ts">
import AdminLayout from '@/components/layout/AdminLayout.vue'
import RevievOnCommunityModal from '@/components/modals/RevievOnCommunityModal.vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Pagination, PaginationEllipsis, PaginationFirst, PaginationLast, PaginationList, PaginationListItem, PaginationNext, PaginationPrev } from '@/components/ui/pagination'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { useModal } from '@/composables/useModal'
import { reviewOnCommunityService } from '@/services/reviewOnCommunityService'
import { onMounted, ref } from 'vue'

import Pencil from '~icons/lucide/pencil'
import Plus from '~icons/lucide/plus'
import Trash from '~icons/lucide/trash'

onMounted(reviewOnCommunityService.search)

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
</script>

<template>
  <AdminLayout>
    <div class="space-y-6">
      <div class="flex justify-between items-center">
        <h1 class="text-3xl font-bold">
          Отзывы на сообщество
        </h1>
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
                <TableHead>Автор</TableHead>
                <TableHead>Текст отзыва</TableHead>
                <TableHead>Дата</TableHead>
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
                <TableCell>Имя: {{ review.author.name }} <br> Telegram: {{ review.author.tg }}</TableCell>
                <TableCell>{{ review.text }}</TableCell>
                <TableCell>{{ new Date(review.date).toLocaleDateString() }}</TableCell>
                <TableCell>
                  <Button variant="ghost" size="sm" @click="selectReview(review.author.id)">
                    <Pencil class="h-4 w-4" />
                  </Button>
                  <Button variant="ghost" size="sm" @click="reviewOnCommunityService.delete(review.id)">
                    <Trash class="h-4 w-4" />
                  </Button>
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </CardContent>
      </Card>
      <div class="mt-4 flex justify-end">
        <Pagination v-slot="{ page }" :items-per-page="10" :total="reviewOnCommunityService.items.value.total" :sibling-count="1" show-edges :default-page="1">
          <PaginationList v-slot="{ items }" class="flex items-center gap-1">
            <PaginationFirst />
            <PaginationPrev />

            <template v-for="(item, index) in items">
              <PaginationListItem v-if="item.type === 'page'" :key="index" :value="item.value" as-child>
                <Button class="w-10 h-10 p-0" :variant="item.value === page ? 'default' : 'outline'" @click="reviewOnCommunityService.changePagination(item.value)">
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
    <RevievOnCommunityModal v-model:is-open="isOpen" :review-id="selectedReviewId" @saved="reviewOnCommunityService.search" />
  </AdminLayout>
</template>
