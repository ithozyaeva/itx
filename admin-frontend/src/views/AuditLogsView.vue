<script setup lang="ts">
import type { AuditLogFilters } from '@/services/auditLogService'
import { Typography } from 'itx-ui-kit'
import { onMounted, onUnmounted, ref } from 'vue'
import AdminLayout from '@/components/layout/AdminLayout.vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Pagination, PaginationEllipsis, PaginationFirst, PaginationLast, PaginationList, PaginationListItem, PaginationNext, PaginationPrev } from '@/components/ui/pagination'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { auditLogService } from '@/services/auditLogService'

const filters = ref<AuditLogFilters>({})

const actionLabels: Record<string, string> = {
  create: 'Создание',
  update: 'Обновление',
  delete: 'Удаление',
  approve: 'Одобрение',
}

const entityTypeLabels: Record<string, string> = {
  event: 'Событие',
  mentor: 'Ментор',
  member: 'Участник',
  review_on_community: 'Отзыв на сообщество',
  review_on_service: 'Отзыв на услугу',
  referal_link: 'Реферальная ссылка',
  resume: 'Резюме',
}

const actorTypeLabels: Record<string, string> = {
  admin: 'Админ',
  platform: 'Платформа',
}

function applyFilters() {
  auditLogService.applyFilters(filters.value)
}

function resetFilters() {
  filters.value = {}
  auditLogService.applyFilters({})
}

onMounted(auditLogService.search)
onUnmounted(auditLogService.clearPagination)
</script>

<template>
  <AdminLayout>
    <div class="space-y-6">
      <Typography variant="h2" as="h1">
        Журнал действий
      </Typography>

      <Card>
        <CardContent>
          <div class="flex flex-wrap items-center gap-3 mb-4">
            <Select v-model="filters.actorType" @update:model-value="applyFilters">
              <SelectTrigger class="w-full sm:w-[180px]">
                <SelectValue placeholder="Тип актора" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="admin">
                  Админ
                </SelectItem>
                <SelectItem value="platform">
                  Платформа
                </SelectItem>
              </SelectContent>
            </Select>

            <Select v-model="filters.action" @update:model-value="applyFilters">
              <SelectTrigger class="w-full sm:w-[180px]">
                <SelectValue placeholder="Действие" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="create">
                  Создание
                </SelectItem>
                <SelectItem value="update">
                  Обновление
                </SelectItem>
                <SelectItem value="delete">
                  Удаление
                </SelectItem>
                <SelectItem value="approve">
                  Одобрение
                </SelectItem>
              </SelectContent>
            </Select>

            <Select v-model="filters.entityType" @update:model-value="applyFilters">
              <SelectTrigger class="w-full sm:w-[180px]">
                <SelectValue placeholder="Тип сущности" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="event">
                  Событие
                </SelectItem>
                <SelectItem value="mentor">
                  Ментор
                </SelectItem>
                <SelectItem value="member">
                  Участник
                </SelectItem>
                <SelectItem value="review_on_community">
                  Отзыв на сообщество
                </SelectItem>
                <SelectItem value="review_on_service">
                  Отзыв на услугу
                </SelectItem>
                <SelectItem value="referal_link">
                  Реферальная ссылка
                </SelectItem>
                <SelectItem value="resume">
                  Резюме
                </SelectItem>
              </SelectContent>
            </Select>

            <Button variant="outline" @click="resetFilters">
              Сбросить
            </Button>
          </div>

          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Дата</TableHead>
                <TableHead>Актор</TableHead>
                <TableHead>Тип</TableHead>
                <TableHead>Действие</TableHead>
                <TableHead>Сущность</TableHead>
                <TableHead>ID</TableHead>
                <TableHead>Название</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow v-if="auditLogService.items.value.total === 0" class="h-24">
                <TableCell colspan="7" class="text-center">
                  Записи не найдены
                </TableCell>
              </TableRow>
              <TableRow v-for="log in auditLogService.items.value.items" :key="log.id">
                <TableCell class="whitespace-nowrap">
                  {{ new Date(log.createdAt).toLocaleString() }}
                </TableCell>
                <TableCell>{{ log.actorName }}</TableCell>
                <TableCell>{{ actorTypeLabels[log.actorType] ?? log.actorType }}</TableCell>
                <TableCell>{{ actionLabels[log.action] ?? log.action }}</TableCell>
                <TableCell>{{ entityTypeLabels[log.entityType] ?? log.entityType }}</TableCell>
                <TableCell>{{ log.entityId }}</TableCell>
                <TableCell>{{ log.entityName }}</TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      <div class="mt-4 flex justify-end">
        <Pagination
          v-slot="{ page }" :items-per-page="20" :total="auditLogService.items.value.total"
          :sibling-count="1" show-edges :default-page="1"
        >
          <PaginationList v-slot="{ items }" class="flex items-center gap-1">
            <PaginationFirst />
            <PaginationPrev />

            <template v-for="(item, index) in items">
              <PaginationListItem v-if="item.type === 'page'" :key="index" :value="item.value" as-child>
                <Button
                  class="w-10 h-10 p-0" :variant="item.value === page ? 'default' : 'outline'"
                  @click="auditLogService.changePagination(item.value)"
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
  </AdminLayout>
</template>
