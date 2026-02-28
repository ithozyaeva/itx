<script setup lang="ts">
import { Typography } from 'itx-ui-kit'
import { onMounted, onUnmounted, ref } from 'vue'
import Plus from '~icons/lucide/plus'
import Trash from '~icons/lucide/trash'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import AdminLayout from '@/components/layout/AdminLayout.vue'
import PointsAwardModal from '@/components/modals/PointsAwardModal.vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Pagination, PaginationEllipsis, PaginationFirst, PaginationLast, PaginationList, PaginationListItem, PaginationNext, PaginationPrev } from '@/components/ui/pagination'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { pointsService } from '@/services/pointsService'

const isAwardModalOpen = ref(false)
const memberIdFilter = ref('')

const reasonLabels: Record<string, string> = {
  event_attend: 'Участие в событии',
  event_host: 'Проведение события',
  review_community: 'Отзыв на сообщество',
  review_service: 'Отзыв на услугу',
  resume_upload: 'Загрузка резюме',
  referal_create: 'Создание реферала',
  referal_conversion: 'Конверсия реферала',
  profile_complete: 'Заполнение профиля',
  weekly_activity: 'Еженедельная активность',
  monthly_active: 'Месячная активность',
  streak_4weeks: 'Серия 4 недели',
  admin_manual: 'Ручное начисление',
}

function applyMemberFilter() {
  const id = memberIdFilter.value ? Number(memberIdFilter.value) : undefined
  pointsService.applyFilters({ memberId: id })
}

function resetFilters() {
  memberIdFilter.value = ''
  pointsService.applyFilters({})
}

onMounted(pointsService.search)
onUnmounted(pointsService.clearPagination)
</script>

<template>
  <AdminLayout>
    <div class="space-y-6">
      <div class="flex justify-between items-center">
        <Typography
          variant="h2"
          as="h1"
        >
          Баллы
        </Typography>
        <Button
          v-permission="'can_edit_admin_points'"
          @click="isAwardModalOpen = true"
        >
          <Plus class="mr-2 h-4 w-4" />
          Начислить баллы
        </Button>
      </div>

      <Card class="p-4 rounded-lg">
        <div class="flex items-center gap-3">
          <Input
            v-model="memberIdFilter"
            placeholder="ID участника"
            class="w-[200px]"
            @keydown.enter="applyMemberFilter"
          />
          <Button
            size="sm"
            @click="applyMemberFilter"
          >
            Применить
          </Button>
          <Button
            variant="outline"
            size="sm"
            @click="resetFilters"
          >
            Сбросить
          </Button>
        </div>
      </Card>

      <Card>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Участник</TableHead>
                <TableHead>Баллы</TableHead>
                <TableHead>Причина</TableHead>
                <TableHead>Описание</TableHead>
                <TableHead>Дата</TableHead>
                <TableHead />
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow
                v-if="pointsService.items.value.total === 0"
                class="h-24"
              >
                <TableCell
                  colspan="6"
                  class="text-center"
                >
                  Транзакции не найдены
                </TableCell>
              </TableRow>
              <TableRow
                v-for="tx in pointsService.items.value.items"
                :key="tx.id"
              >
                <TableCell>
                  {{ tx.memberFirstName ?? '' }} {{ tx.memberLastName ?? '' }}
                  <span class="text-muted-foreground">@{{ tx.memberUsername }}</span>
                </TableCell>
                <TableCell class="font-medium">
                  {{ tx.amount }}
                </TableCell>
                <TableCell>{{ reasonLabels[tx.reason] ?? tx.reason }}</TableCell>
                <TableCell class="max-w-[200px] truncate">
                  {{ tx.description }}
                </TableCell>
                <TableCell class="whitespace-nowrap">
                  {{ new Date(tx.createdAt).toLocaleString() }}
                </TableCell>
                <TableCell class="text-right">
                  <ConfirmDialog
                    title="Удалить транзакцию?"
                    description="Транзакция будет удалена без возможности восстановления."
                    confirm-label="Удалить"
                    @confirm="pointsService.deleteTransaction(tx.id)"
                  >
                    <template #trigger>
                      <Button
                        v-permission="'can_edit_admin_points'"
                        variant="ghost"
                        size="sm"
                        class="text-destructive"
                        :disabled="pointsService.isLoading.value"
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
        <Pagination
          v-slot="{ page }"
          :items-per-page="20"
          :total="pointsService.items.value.total"
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
                  @click="pointsService.changePagination(item.value)"
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
    </div>

    <PointsAwardModal
      v-model:is-open="isAwardModalOpen"
      @saved="pointsService.search"
    />
  </AdminLayout>
</template>
