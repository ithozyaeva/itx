<script setup lang="ts">
import { Tag, Typography } from 'itx-ui-kit'
import { onMounted, onUnmounted } from 'vue'
import Trash from '~icons/lucide/trash'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import AdminLayout from '@/components/layout/AdminLayout.vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Pagination, PaginationEllipsis, PaginationFirst, PaginationLast, PaginationList, PaginationListItem, PaginationNext, PaginationPrev } from '@/components/ui/pagination'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { referralsService } from '@/services/referralsService'

onMounted(referralsService.search)
onUnmounted(referralsService.clearPagination)
</script>

<template>
  <AdminLayout>
    <div class="space-y-6">
      <div class="flex justify-between items-center">
        <Typography variant="h2" as="h1">
          Реферальные ссылки
        </Typography>
      </div>

      <Card>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Автор</TableHead>
                <TableHead>Компания</TableHead>
                <TableHead>Грейд</TableHead>
                <TableHead>Теги</TableHead>
                <TableHead>Статус</TableHead>
                <TableHead>Вакансии</TableHead>
                <TableHead>Конверсии</TableHead>
                <TableHead />
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow v-if="referralsService.items.value.total === 0" class="h-24">
                <TableCell colspan="8" class="text-center">
                  Реферальные ссылки не найдены
                </TableCell>
              </TableRow>
              <TableRow v-for="link in referralsService.items.value.items" :key="link.id">
                <TableCell>
                  {{ link.author?.firstName }} {{ link.author?.lastName }}
                </TableCell>
                <TableCell>{{ link.company }}</TableCell>
                <TableCell>{{ link.grade }}</TableCell>
                <TableCell>
                  <div class="flex flex-wrap gap-1">
                    <Tag v-for="tag in link.profTags" :key="tag.id">
                      {{ tag.title }}
                    </Tag>
                  </div>
                </TableCell>
                <TableCell>
                  <Tag :class="link.status === 'active' ? 'bg-green-500/10 text-green-600' : 'bg-gray-500/10 text-gray-600'">
                    {{ link.status === 'active' ? 'Активна' : 'Заморожена' }}
                  </Tag>
                </TableCell>
                <TableCell>{{ link.vacationsCount }}</TableCell>
                <TableCell>{{ link.conversionsCount }}</TableCell>
                <TableCell>
                  <ConfirmDialog
                    title="Удалить реферальную ссылку?"
                    description="Ссылка будет удалена без возможности восстановления."
                    confirm-label="Удалить"
                    @confirm="referralsService.delete(link.id)"
                  >
                    <template #trigger>
                      <Button
                        variant="ghost"
                        size="sm"
                        class="text-destructive"
                        :disabled="referralsService.isLoading.value"
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
        <Pagination v-slot="{ page }" :items-per-page="10" :total="referralsService.items.value.total" :sibling-count="1" show-edges :default-page="1">
          <PaginationList v-slot="{ items }" class="flex items-center gap-1">
            <PaginationFirst />
            <PaginationPrev />

            <template v-for="(item, index) in items">
              <PaginationListItem v-if="item.type === 'page'" :key="index" :value="item.value" as-child>
                <Button class="w-10 h-10 p-0" :variant="item.value === page ? 'default' : 'outline'" @click="referralsService.changePagination(item.value)">
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
