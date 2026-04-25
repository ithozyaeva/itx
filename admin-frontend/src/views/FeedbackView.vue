<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import AdminLayout from '@/components/layout/AdminLayout.vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Pagination, PaginationEllipsis, PaginationFirst, PaginationLast, PaginationList, PaginationListItem, PaginationNext, PaginationPrev } from '@/components/ui/pagination'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Typography } from '@/components/ui/typography'
import { feedbackService } from '@/services/feedbackService'

onMounted(feedbackService.search)
onUnmounted(feedbackService.clearPagination)

function userName(item: { userFirstName: string | null, userLastName: string | null, userUsername: string | null }): string {
  const name = [item.userFirstName, item.userLastName].filter(Boolean).join(' ').trim()
  if (name && item.userUsername)
    return `${name} (@${item.userUsername})`
  if (name)
    return name
  if (item.userUsername)
    return `@${item.userUsername}`
  return 'Аноним'
}

function scoreColor(score: number): string {
  if (score >= 9)
    return 'text-green-600'
  if (score >= 7)
    return 'text-yellow-600'
  return 'text-red-600'
}
</script>

<template>
  <AdminLayout>
    <div class="space-y-6">
      <Typography variant="h2" as="h1">
        Обратная связь
      </Typography>

      <Card>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Дата</TableHead>
                <TableHead>Пользователь</TableHead>
                <TableHead class="w-20 text-center">
                  Оценка
                </TableHead>
                <TableHead>Комментарий</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow v-if="feedbackService.items.value.total === 0" class="h-24">
                <TableCell colspan="4" class="text-center">
                  Отзывов пока нет
                </TableCell>
              </TableRow>
              <TableRow v-for="item in feedbackService.items.value.items" :key="item.id">
                <TableCell class="whitespace-nowrap">
                  {{ new Date(item.createdAt).toLocaleString() }}
                </TableCell>
                <TableCell class="whitespace-nowrap">
                  {{ userName(item) }}
                </TableCell>
                <TableCell class="text-center font-semibold" :class="scoreColor(item.score)">
                  {{ item.score }}
                </TableCell>
                <TableCell class="whitespace-pre-wrap">
                  {{ item.comment || '—' }}
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      <div class="mt-4 flex justify-end">
        <Pagination v-slot="{ page }" :items-per-page="10" :total="feedbackService.items.value.total" :sibling-count="1" show-edges :default-page="1">
          <PaginationList v-slot="{ items }" class="flex items-center gap-1">
            <PaginationFirst />
            <PaginationPrev />

            <template v-for="(item, index) in items">
              <PaginationListItem v-if="item.type === 'page'" :key="index" :value="item.value" as-child>
                <Button class="w-10 h-10 p-0" :variant="item.value === page ? 'default' : 'outline'" @click="feedbackService.changePagination(item.value)">
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
  </AdminLayout>
</template>
